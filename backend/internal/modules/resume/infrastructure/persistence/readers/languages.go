package readers

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type LanguagesReader struct{ db *gorm.DB }

func NewLanguagesReader(db *gorm.DB) *LanguagesReader { return &LanguagesReader{db: db} }

func (r *LanguagesReader) ListLanguages(ctx context.Context, q queries.ListLanguagesQuery) ([]queries.LanguageView, error) {
	db, err := resolveDB(ctx, r.db, "list languages")
	if err != nil {
		return nil, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("list languages: invalid query")
	}
	var models []gormmodel.ResumeLanguage
	if err := db.Joins("JOIN resumes ON resumes.id = resume_languages.resume_id").Where("resume_languages.resume_id = ? AND resumes.user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).Order("position").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list languages: %w", err)
	}
	out := make([]queries.LanguageView, 0, len(models))
	for _, model := range models {
		out = append(out, queries.LanguageView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Language: model.Language, Grade: model.Grade, Position: model.Position})
	}
	return out, nil
}

func (r *LanguagesReader) GetLanguage(ctx context.Context, q queries.GetLanguageQuery) (queries.LanguageView, error) {
	db, err := resolveDB(ctx, r.db, "get language")
	if err != nil {
		return queries.LanguageView{}, err
	}
	if q.ResumeID.IsZero() || q.LanguageID.IsZero() || q.UserID.IsZero() {
		return queries.LanguageView{}, fmt.Errorf("get language: invalid query")
	}
	var model gormmodel.ResumeLanguage
	err = db.Joins("JOIN resumes ON resumes.id = resume_languages.resume_id").Where("resume_languages.id = ? AND resume_languages.resume_id = ? AND resumes.user_id = ?", q.LanguageID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.LanguageView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.LanguageView{}, fmt.Errorf("get language: %w", err)
	}
	return queries.LanguageView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Language: model.Language, Grade: model.Grade, Position: model.Position}, nil
}
