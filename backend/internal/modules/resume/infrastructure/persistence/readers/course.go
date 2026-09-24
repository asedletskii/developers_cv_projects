package readers

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourseReader struct{ pool *pgxpool.Pool }

func NewCourseReader(pool *pgxpool.Pool) *CourseReader { return &CourseReader{pool: pool} }
func (r *CourseReader) ListCourses(ctx context.Context, q queries.ListCoursesQuery) ([]queries.CourseView, error) {
	if r == nil || r.pool == nil {
		return nil, fmt.Errorf("list courses: reader dependencies are nil")
	}
	rows, err := r.pool.Query(ctx, `SELECT c.id::text,c.resume_id::text,c.name,c.organization,to_char(c.end_date,'YYYY-MM'),c.hours,c.certificate_link FROM resume_courses c JOIN resume_roots rr ON rr.id=c.resume_id WHERE c.resume_id=$1 AND rr.user_id=$2 ORDER BY c.sort_order`, q.ResumeID.UUID(), q.UserID.UUID())
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	defer rows.Close()
	out := make([]queries.CourseView, 0)
	for rows.Next() {
		var v queries.CourseView
		if err := rows.Scan(&v.ID, &v.ResumeID, &v.Name, &v.Organization, &v.EndDate, &v.Hours, &v.CertificateLink); err != nil {
			return nil, fmt.Errorf("scan course: %w", err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate courses: %w", err)
	}
	return out, nil
}
func (r *CourseReader) GetCourse(ctx context.Context, q queries.GetCourseQuery) (queries.CourseView, error) {
	if r == nil || r.pool == nil {
		return queries.CourseView{}, fmt.Errorf("get course: reader dependencies are nil")
	}
	var v queries.CourseView
	err := r.pool.QueryRow(ctx, `SELECT c.id::text,c.resume_id::text,c.name,c.organization,to_char(c.end_date,'YYYY-MM'),c.hours,c.certificate_link FROM resume_courses c JOIN resume_roots rr ON rr.id=c.resume_id WHERE c.id=$1 AND c.resume_id=$2 AND rr.user_id=$3`, q.CourseID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).Scan(&v.ID, &v.ResumeID, &v.Name, &v.Organization, &v.EndDate, &v.Hours, &v.CertificateLink)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, queries.ErrResumeNotFound
	}
	if err != nil {
		return v, fmt.Errorf("get course: %w", err)
	}
	return v, nil
}
