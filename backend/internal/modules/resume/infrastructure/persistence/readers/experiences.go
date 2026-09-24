package readers

import (
	"context"
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

func (r *ExperiencesReader) ListExperiences(ctx context.Context, q queries.ListExperiencesQuery) ([]queries.ExperienceView, error) {
	if r == nil || r.pool == nil {
		return nil, fmt.Errorf("list experiences: reader dependencies are nil")
	}
	rows, err := r.pool.Query(ctx, `SELECT e.id::text,e.resume_id::text,e.company,e.position,e.location,to_char(e.start_date,'YYYY-MM'),to_char(e.end_date,'YYYY-MM'),e.description FROM resume_experiences e JOIN resume_roots rr ON rr.id=e.resume_id WHERE e.resume_id=$1 AND rr.user_id=$2 ORDER BY e.sort_order`, q.ResumeID.UUID(), q.UserID.UUID())
	if err != nil {
		return nil, fmt.Errorf("list experiences: %w", err)
	}
	defer rows.Close()
	out := make([]queries.ExperienceView, 0)
	for rows.Next() {
		var v queries.ExperienceView
		if err := rows.Scan(&v.ID, &v.ResumeID, &v.Company, &v.Position, &v.Location, &v.StartDate, &v.EndDate, &v.Description); err != nil {
			return nil, fmt.Errorf("scan experience: %w", err)
		}
		if err := r.loadChildren(ctx, v.ID, &v); err != nil {
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
	err := r.pool.QueryRow(ctx, `SELECT e.id::text,e.resume_id::text,e.company,e.position,e.location,to_char(e.start_date,'YYYY-MM'),to_char(e.end_date,'YYYY-MM'),e.description FROM resume_experiences e JOIN resume_roots rr ON rr.id=e.resume_id WHERE e.id=$1 AND e.resume_id=$2 AND rr.user_id=$3`, q.ExperienceID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).Scan(&v.ID, &v.ResumeID, &v.Company, &v.Position, &v.Location, &v.StartDate, &v.EndDate, &v.Description)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, queries.ErrResumeNotFound
	}
	if err != nil {
		return v, fmt.Errorf("get experience: %w", err)
	}
	if err := r.loadChildren(ctx, v.ID, &v); err != nil {
		return v, err
	}
	return v, nil
}

func (r *ExperiencesReader) loadChildren(ctx context.Context, experienceID string, v *queries.ExperienceView) error {
	rows, err := r.pool.Query(ctx, `SELECT value FROM resume_experience_achievements WHERE experience_id=$1::uuid ORDER BY position`, experienceID)
	if err != nil {
		return fmt.Errorf("load experience achievements: %w", err)
	}
	defer rows.Close()
	v.Achievements = make([]string, 0)
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return fmt.Errorf("scan experience achievement: %w", err)
		}
		v.Achievements = append(v.Achievements, value)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate experience achievements: %w", err)
	}

	techRows, err := r.pool.Query(ctx, `SELECT name,grade FROM resume_experience_technologies WHERE experience_id=$1::uuid ORDER BY position`, experienceID)
	if err != nil {
		return fmt.Errorf("load experience technologies: %w", err)
	}
	defer techRows.Close()
	v.Technologies = make([]queries.ExperienceTechnologyView, 0)
	for techRows.Next() {
		var t queries.ExperienceTechnologyView
		if err := techRows.Scan(&t.Name, &t.Grade); err != nil {
			return fmt.Errorf("scan experience technology: %w", err)
		}
		v.Technologies = append(v.Technologies, t)
	}
	if err := techRows.Err(); err != nil {
		return fmt.Errorf("iterate experience technologies: %w", err)
	}
	return nil
}
