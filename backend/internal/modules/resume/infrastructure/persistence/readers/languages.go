package readers

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type LanguagesReader struct{ pool *pgxpool.Pool }

func NewLanguagesReader(pool *pgxpool.Pool) *LanguagesReader { return &LanguagesReader{pool: pool} }

func (r *LanguagesReader) ListLanguages(ctx context.Context, q queries.ListLanguagesQuery) ([]queries.LanguageView, error) {
	if r == nil || r.pool == nil {
		return nil, fmt.Errorf("list languages: reader dependencies are nil")
	}
	rows, err := r.pool.Query(ctx, `SELECT l.id::text,l.resume_id::text,l.language,l.grade,l.position FROM resume_languages l JOIN resumes rr ON rr.id=l.resume_id WHERE l.resume_id=$1 AND rr.user_id=$2 ORDER BY l.position`, q.ResumeID.UUID(), q.UserID.UUID())
	if err != nil {
		return nil, fmt.Errorf("list languages: %w", err)
	}
	defer rows.Close()
	out := make([]queries.LanguageView, 0)
	for rows.Next() {
		var v queries.LanguageView
		if err := rows.Scan(&v.ID, &v.ResumeID, &v.Language, &v.Grade, &v.Position); err != nil {
			return nil, fmt.Errorf("scan language: %w", err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate languages: %w", err)
	}
	return out, nil
}

func (r *LanguagesReader) GetLanguage(ctx context.Context, q queries.GetLanguageQuery) (queries.LanguageView, error) {
	if r == nil || r.pool == nil {
		return queries.LanguageView{}, fmt.Errorf("get language: reader dependencies are nil")
	}
	var v queries.LanguageView
	err := r.pool.QueryRow(ctx, `SELECT l.id::text,l.resume_id::text,l.language,l.grade,l.position FROM resume_languages l JOIN resumes rr ON rr.id=l.resume_id WHERE l.id=$1 AND l.resume_id=$2 AND rr.user_id=$3`, q.LanguageID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).Scan(&v.ID, &v.ResumeID, &v.Language, &v.Grade, &v.Position)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, queries.ErrResumeNotFound
	}
	if err != nil {
		return v, fmt.Errorf("get language: %w", err)
	}
	return v, nil
}
