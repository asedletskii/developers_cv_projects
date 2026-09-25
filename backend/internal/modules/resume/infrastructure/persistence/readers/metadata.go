package readers

import (
	"context"
	"errors"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/queries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResumeMetadataReader struct{ pool *pgxpool.Pool }

func NewResumeMetadataReader(pool *pgxpool.Pool) *ResumeMetadataReader {
	return &ResumeMetadataReader{pool: pool}
}

func (r *ResumeMetadataReader) GetResumeMetadata(ctx context.Context, q queries.GetResumeMetadataQuery) (queries.ResumeMetadataView, error) {
	if r == nil || r.pool == nil {
		return queries.ResumeMetadataView{}, fmt.Errorf("get resume metadata: reader dependencies are nil")
	}
	var v queries.ResumeMetadataView
	err := r.pool.QueryRow(ctx, `SELECT id::text,title,language,created_at,updated_at FROM resumes WHERE id=$1 AND user_id=$2`, q.ResumeID.UUID(), q.UserID.UUID()).Scan(&v.ID, &v.Title, &v.Language, &v.CreatedAt, &v.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return queries.ResumeMetadataView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.ResumeMetadataView{}, fmt.Errorf("get resume metadata: %w", err)
	}
	return v, nil
}
