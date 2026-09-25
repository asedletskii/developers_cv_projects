package readers

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type TechnologiesReader struct{ pool *pgxpool.Pool }

func NewTechnologiesReader(pool *pgxpool.Pool) *TechnologiesReader {
	return &TechnologiesReader{pool: pool}
}

func (r *TechnologiesReader) ListTechnologies(ctx context.Context, q queries.ListTechnologiesQuery) ([]queries.TechnologyView, error) {
	if r == nil || r.pool == nil {
		return nil, fmt.Errorf("list technologies: reader dependencies are nil")
	}
	rows, err := r.pool.Query(ctx, `SELECT t.id::text,t.resume_id::text,t.name,t.grade,t.position FROM resume_technologies t JOIN resume_roots rr ON rr.id=t.resume_id WHERE t.resume_id=$1 AND rr.user_id=$2 ORDER BY t.position`, q.ResumeID.UUID(), q.UserID.UUID())
	if err != nil {
		return nil, fmt.Errorf("list technologies: %w", err)
	}
	defer rows.Close()
	out := make([]queries.TechnologyView, 0)
	for rows.Next() {
		var v queries.TechnologyView
		if err := rows.Scan(&v.ID, &v.ResumeID, &v.Name, &v.Grade, &v.Position); err != nil {
			return nil, fmt.Errorf("scan technology: %w", err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate technologies: %w", err)
	}
	return out, nil
}

func (r *TechnologiesReader) GetTechnology(ctx context.Context, q queries.GetTechnologyQuery) (queries.TechnologyView, error) {
	if r == nil || r.pool == nil {
		return queries.TechnologyView{}, fmt.Errorf("get technology: reader dependencies are nil")
	}
	var v queries.TechnologyView
	err := r.pool.QueryRow(ctx, `SELECT t.id::text,t.resume_id::text,t.name,t.grade,t.position FROM resume_technologies t JOIN resume_roots rr ON rr.id=t.resume_id WHERE t.id=$1 AND t.resume_id=$2 AND rr.user_id=$3`, q.TechnologyID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).Scan(&v.ID, &v.ResumeID, &v.Name, &v.Grade, &v.Position)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, queries.ErrResumeNotFound
	}
	if err != nil {
		return v, fmt.Errorf("get technology: %w", err)
	}
	return v, nil
}
