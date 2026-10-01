package readers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func (r *ResumeReader) loadActivities(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var models []gormmodel.ResumeActivity
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("sort_order").Find(&models).Error; err != nil {
		return fmt.Errorf("list activities: %w", err)
	}
	result := make([]queries.ActivityView, 0, len(models))
	for _, model := range models {
		result = append(result, queries.ActivityView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Type: model.Type, Title: model.Title, Description: model.Description, Date: optionalYearMonth(model.Date), Link: model.Link})
	}
	view.Activities = result
	return nil
}
