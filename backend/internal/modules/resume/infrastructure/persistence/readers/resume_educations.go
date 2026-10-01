package readers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func (r *ResumeReader) loadEducations(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var models []gormmodel.ResumeEducation
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("sort_order").Find(&models).Error; err != nil {
		return fmt.Errorf("list educations: %w", err)
	}
	result := make([]queries.EducationView, 0, len(models))
	for _, model := range models {
		result = append(result, queries.EducationView{ID: model.ID.String(), Location: model.Location, University: model.University, Grade: model.Grade, ProgramName: model.ProgramName, ProfileName: model.ProfileName, Faculty: model.Faculty, Course: model.Course, EndDate: optionalYear(model.EndDate)})
	}
	view.Educations = result
	return nil
}
