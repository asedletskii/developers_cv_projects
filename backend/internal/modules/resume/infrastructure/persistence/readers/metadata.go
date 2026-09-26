package readers

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type ResumeMetadataReader struct{ db *gorm.DB }

func NewResumeMetadataReader(db *gorm.DB) *ResumeMetadataReader {
	return &ResumeMetadataReader{db: db}
}

func (r *ResumeMetadataReader) GetResumeMetadata(ctx context.Context, q queries.GetResumeMetadataQuery) (queries.ResumeMetadataView, error) {
	db, err := resolveDB(ctx, r.db, "get resume metadata")
	if err != nil {
		return queries.ResumeMetadataView{}, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return queries.ResumeMetadataView{}, fmt.Errorf("get resume metadata: invalid query")
	}

	var model gormmodel.Resume
	err = db.Where("id = ? AND user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.ResumeMetadataView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.ResumeMetadataView{}, fmt.Errorf("get resume metadata: %w", err)
	}

	return queries.ResumeMetadataView{
		ID: model.ID.String(), Title: model.Title, Language: model.Language,
		CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt,
	}, nil
}
