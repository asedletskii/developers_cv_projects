package readers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func (r *ResumeReader) loadCourses(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var models []gormmodel.ResumeCourse
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("sort_order").Find(&models).Error; err != nil {
		return fmt.Errorf("list courses: %w", err)
	}
	result := make([]queries.CourseView, 0, len(models))
	for _, model := range models {
		result = append(result, queries.CourseView{ID: model.ID.String(), Name: model.Name, Organization: model.Organization, EndDate: yearMonth(model.EndDate), Hours: model.Hours, CertificateLink: model.CertificateLink})
	}
	view.Courses = result
	return nil
}
