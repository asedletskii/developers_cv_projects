package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistEducations(tx *gorm.DB, resume *domain.Resume) error {
	models := mappers.EducationsToDBModels(resume.Educations())
	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeEducation{}).Error; err != nil {
		return fmt.Errorf("replace resume educations: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume educations: %w", mapDBError(err))
		}
	}
	return nil
}

func loadEducations(tx *gorm.DB, resumeID sharedvalueobjects.ID) ([]*domain.Education, error) {
	var rows []gormmodel.ResumeEducation
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume educations: %w", mapDBError(err))
	}
	return mappers.EducationsToDomain(rows)
}
