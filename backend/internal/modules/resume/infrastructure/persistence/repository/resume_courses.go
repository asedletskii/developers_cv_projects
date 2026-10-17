package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistCourses(tx *gorm.DB, resume *domain.Resume) error {
	models := mappers.CoursesToDBModels(resume.Courses())
	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeCourse{}).Error; err != nil {
		return fmt.Errorf("replace resume courses: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume courses: %w", mapDBError(err))
		}
	}
	return nil
}

func loadCourses(tx *gorm.DB, resumeID sharedvalueobjects.ID) ([]*domain.Course, error) {
	var rows []gormmodel.ResumeCourse
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume courses: %w", mapDBError(err))
	}
	return mappers.CoursesToDomain(rows)
}
