package readers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func (r *ResumeReader) loadTechnologies(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var models []gormmodel.ResumeTechnology
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("position").Find(&models).Error; err != nil {
		return fmt.Errorf("list technologies: %w", err)
	}
	result := make([]queries.TechnologyView, 0, len(models))
	for _, model := range models {
		result = append(result, queries.TechnologyView{ID: model.ID.String(), Name: model.Name, Grade: model.Grade, Position: model.Position})
	}
	view.Technologies = result
	return nil
}
