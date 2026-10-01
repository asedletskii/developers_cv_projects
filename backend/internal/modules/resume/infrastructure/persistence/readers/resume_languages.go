package readers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func (r *ResumeReader) loadLanguages(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var models []gormmodel.ResumeLanguage
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("position").Find(&models).Error; err != nil {
		return fmt.Errorf("list languages: %w", err)
	}
	result := make([]queries.LanguageView, 0, len(models))
	for _, model := range models {
		result = append(result, queries.LanguageView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Language: model.Language, Grade: model.Grade, Position: model.Position})
	}
	view.Languages = result
	return nil
}
