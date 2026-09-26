package readers

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
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
	readDB := db.Where("id = ?", q.ResumeID.UUID())
	readDB = readDB.Where("user_id = ?", q.UserID.UUID())
	err = readDB.First(&model).Error
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

func (r *ResumeMetadataReader) ReadResumeMetadata(ctx context.Context, resumeID valueobjects.ID) (queries.ResumeMetadataView, error) {
	db, err := resolveDB(ctx, r.db, "get resume metadata")
	if err != nil {
		return queries.ResumeMetadataView{}, err
	}
	if resumeID.IsZero() {
		return queries.ResumeMetadataView{}, fmt.Errorf("get resume metadata: invalid query")
	}

	var model gormmodel.Resume
	readDB := db.Where("id = ?", resumeID.UUID())
	err = readDB.First(&model).Error
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
