package readers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type ExperiencesReader struct{ pool *pgxpool.Pool }

func NewExperiencesReader(pool *pgxpool.Pool) *ExperiencesReader {
	return &ExperiencesReader{pool: pool}
}

const experienceSelect = `
SELECT
 e.id::text,
 e.resume_id::text,
 e.company,
 e.position,
 e.location,
 to_char(e.start_date,'YYYY-MM'),
 to_char(e.end_date,'YYYY-MM'),
 e.description,
 COALESCE((SELECT jsonb_agg(a.value ORDER BY a.position) FROM resume_experience_achievements a WHERE a.experience_id=e.id), '[]'::jsonb),
 COALESCE((SELECT jsonb_agg(jsonb_build_object('name', t.name, 'grade', t.grade) ORDER BY t.position) FROM resume_experience_technologies t WHERE t.experience_id=e.id), '[]'::jsonb)
FROM resume_experiences e
JOIN resumes rr ON rr.id=e.resume_id
`

func (r *ExperiencesReader) ListExperiences(ctx context.Context, q queries.ListExperiencesQuery) ([]queries.ExperienceView, error) {
	if r == nil || r.pool == nil {
		return nil, fmt.Errorf("list experiences: reader dependencies are nil")
	}
	rows, err := r.pool.Query(ctx, experienceSelect+`WHERE e.resume_id=$1 AND rr.user_id=$2 ORDER BY e.sort_order`, q.ResumeID.UUID(), q.UserID.UUID())
	if err != nil {
		return nil, fmt.Errorf("list experiences: %w", err)
	}
	defer rows.Close()
	out := make([]queries.ExperienceView, 0)
	for rows.Next() {
		var v queries.ExperienceView
		var achievements, technologies []byte
		if err := rows.Scan(&v.ID, &v.ResumeID, &v.Company, &v.Position, &v.Location, &v.StartDate, &v.EndDate, &v.Description, &achievements, &technologies); err != nil {
			return nil, fmt.Errorf("scan experience: %w", err)
		}
		if err := decodeExperienceChildren(achievements, technologies, &v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate experiences: %w", err)
	}
	return out, nil
}

func (r *ExperiencesReader) GetExperience(ctx context.Context, q queries.GetExperienceQuery) (queries.ExperienceView, error) {
	if r == nil || r.pool == nil {
		return queries.ExperienceView{}, fmt.Errorf("get experience: reader dependencies are nil")
	}
	var v queries.ExperienceView
	var achievements, technologies []byte
	err := r.pool.QueryRow(ctx, experienceSelect+`WHERE e.id=$1 AND e.resume_id=$2 AND rr.user_id=$3`, q.ExperienceID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).
		Scan(&v.ID, &v.ResumeID, &v.Company, &v.Position, &v.Location, &v.StartDate, &v.EndDate, &v.Description, &achievements, &technologies)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, queries.ErrResumeNotFound
	}
	if err != nil {
		return v, fmt.Errorf("get experience: %w", err)
	}
	if err := decodeExperienceChildren(achievements, technologies, &v); err != nil {
		return queries.ExperienceView{}, err
	}
	return v, nil
}

func decodeExperienceChildren(achievements, technologies []byte, view *queries.ExperienceView) error {
	if err := json.Unmarshal(achievements, &view.Achievements); err != nil {
		return fmt.Errorf("decode experience achievements: %w", err)
	}
	if err := json.Unmarshal(technologies, &view.Technologies); err != nil {
		return fmt.Errorf("decode experience technologies: %w", err)
	}
	if view.Achievements == nil {
		view.Achievements = []string{}
	}
	if view.Technologies == nil {
		view.Technologies = []queries.ExperienceTechnologyView{}
	}
	return nil
}
