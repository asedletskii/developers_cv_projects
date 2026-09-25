package readers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type ListResumesReader struct{ pool *pgxpool.Pool }

func NewListResumesReader(pool *pgxpool.Pool) *ListResumesReader {
	return &ListResumesReader{pool: pool}
}

func (r *ListResumesReader) ListResumes(
	ctx context.Context,
	query queries.ListResumesQuery,
) ([]queries.ResumeListItem, error) {
	if r == nil || r.pool == nil {
		return nil, fmt.Errorf("list resumes: reader dependencies are nil")
	}
	if query.UserID.IsZero() {
		return nil, fmt.Errorf("invalid user id")
	}

	rows, err := r.pool.Query(
		ctx,
		`SELECT id::text, title, language FROM resume_roots WHERE user_id = $1 ORDER BY updated_at DESC, id`,
		query.UserID.UUID(),
	)
	if err != nil {
		return nil, fmt.Errorf("list resumes: %w", err)
	}
	defer rows.Close()

	items := make([]queries.ResumeListItem, 0)
	for rows.Next() {
		var item queries.ResumeListItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Language); err != nil {
			return nil, fmt.Errorf("scan resume list: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate resume list: %w", err)
	}

	return items, nil
}
