package readers

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type ListResumesReader struct{ db *gorm.DB }

func NewListResumesReader(db *gorm.DB) *ListResumesReader {
	return &ListResumesReader{db: db}
}

func (r *ListResumesReader) ListResumes(ctx context.Context, query queries.ListResumesQuery) ([]queries.ResumeListItem, error) {
	db, err := resolveDB(ctx, r.db, "list resumes")
	if err != nil {
		return nil, err
	}
	if query.UserID.IsZero() {
		return nil, fmt.Errorf("list resumes: invalid user id")
	}

	var models []gormmodel.Resume
	if err := db.Where("user_id = ?", query.UserID.UUID()).Order("updated_at DESC, id").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list resumes: %w", err)
	}

	items := make([]queries.ResumeListItem, 0, len(models))
	for _, model := range models {
		items = append(items, queries.ResumeListItem{ID: model.ID.String(), Title: model.Title, Language: model.Language})
	}
	return items, nil
}
