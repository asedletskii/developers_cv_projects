package readers

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EducationReader struct{ pool *pgxpool.Pool }

func NewEducationReader(pool *pgxpool.Pool) *EducationReader { return &EducationReader{pool: pool} }
func (r *EducationReader) ListEducations(ctx context.Context, q queries.ListEducationsQuery) ([]queries.EducationView, error) {
	if r == nil || r.pool == nil {
		return nil, fmt.Errorf("list educations: reader dependencies are nil")
	}
	rows, err := r.pool.Query(ctx, `SELECT e.id::text,e.resume_id::text,e.location,e.university,e.grade,e.program_name,e.profile_name,e.faculty,e.course,to_char(e.end_date,'YYYY') FROM resume_educations e JOIN resume_roots rr ON rr.id=e.resume_id WHERE e.resume_id=$1 AND rr.user_id=$2 ORDER BY e.sort_order`, q.ResumeID.UUID(), q.UserID.UUID())
	if err != nil {
		return nil, fmt.Errorf("list educations: %w", err)
	}
	defer rows.Close()
	out := make([]queries.EducationView, 0)
	for rows.Next() {
		var v queries.EducationView
		if err := rows.Scan(&v.ID, &v.ResumeID, &v.Location, &v.University, &v.Grade, &v.ProgramName, &v.ProfileName, &v.Faculty, &v.Course, &v.EndDate); err != nil {
			return nil, fmt.Errorf("scan education: %w", err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate educations: %w", err)
	}
	return out, nil
}
func (r *EducationReader) GetEducation(ctx context.Context, q queries.GetEducationQuery) (queries.EducationView, error) {
	if r == nil || r.pool == nil {
		return queries.EducationView{}, fmt.Errorf("get education: reader dependencies are nil")
	}
	var v queries.EducationView
	err := r.pool.QueryRow(ctx, `SELECT e.id::text,e.resume_id::text,e.location,e.university,e.grade,e.program_name,e.profile_name,e.faculty,e.course,to_char(e.end_date,'YYYY') FROM resume_educations e JOIN resume_roots rr ON rr.id=e.resume_id WHERE e.id=$1 AND e.resume_id=$2 AND rr.user_id=$3`, q.EducationID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).Scan(&v.ID, &v.ResumeID, &v.Location, &v.University, &v.Grade, &v.ProgramName, &v.ProfileName, &v.Faculty, &v.Course, &v.EndDate)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, queries.ErrResumeNotFound
	}
	if err != nil {
		return v, fmt.Errorf("get education: %w", err)
	}
	return v, nil
}
