package readers

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityReader struct{ pool *pgxpool.Pool }

func NewActivityReader(pool *pgxpool.Pool) *ActivityReader { return &ActivityReader{pool: pool} }
func (r *ActivityReader) ListActivities(ctx context.Context, q queries.ListActivitiesQuery) ([]queries.ActivityView, error) {
	if r == nil || r.pool == nil {
		return nil, fmt.Errorf("list activities: reader dependencies are nil")
	}
	rows, err := r.pool.Query(ctx, `SELECT a.id::text,a.resume_id::text,a.type,a.title,a.description,to_char(a.date,'YYYY-MM'),a.link FROM resume_activities a JOIN resume_roots rr ON rr.id=a.resume_id WHERE a.resume_id=$1 AND rr.user_id=$2 ORDER BY a.sort_order`, q.ResumeID.UUID(), q.UserID.UUID())
	if err != nil {
		return nil, fmt.Errorf("list activities: %w", err)
	}
	defer rows.Close()
	out := make([]queries.ActivityView, 0)
	for rows.Next() {
		var v queries.ActivityView
		if err := rows.Scan(&v.ID, &v.ResumeID, &v.Type, &v.Title, &v.Description, &v.Date, &v.Link); err != nil {
			return nil, fmt.Errorf("scan activity: %w", err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate activities: %w", err)
	}
	return out, nil
}
func (r *ActivityReader) GetActivity(ctx context.Context, q queries.GetActivityQuery) (queries.ActivityView, error) {
	if r == nil || r.pool == nil {
		return queries.ActivityView{}, fmt.Errorf("get activity: reader dependencies are nil")
	}
	var v queries.ActivityView
	err := r.pool.QueryRow(ctx, `SELECT a.id::text,a.resume_id::text,a.type,a.title,a.description,to_char(a.date,'YYYY-MM'),a.link FROM resume_activities a JOIN resume_roots rr ON rr.id=a.resume_id WHERE a.id=$1 AND a.resume_id=$2 AND rr.user_id=$3`, q.ActivityID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).Scan(&v.ID, &v.ResumeID, &v.Type, &v.Title, &v.Description, &v.Date, &v.Link)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, queries.ErrResumeNotFound
	}
	if err != nil {
		return v, fmt.Errorf("get activity: %w", err)
	}
	return v, nil
}
