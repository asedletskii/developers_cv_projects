package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistTechnologies(tx *gorm.DB, resume *domain.Resume) error {
	technologies := resume.Technologies()
	models := mappers.ResumeTechnologiesToDBModels(technologies, resume.ID())
	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeTechnology{}).Error; err != nil {
		return fmt.Errorf("replace resume technologies: %w", mapDBError(err))
	}
	if len(models) == 0 {
		return nil
	}
	if err := tx.Create(&models).Error; err != nil {
		return fmt.Errorf("create resume technologies: %w", mapDBError(err))
	}
	return nil
}

func loadTechnologies(tx *gorm.DB, resumeID sharedvalueobjects.ID) (*domain.ResumeTechnologies, error) {
	var rows []gormmodel.ResumeTechnology
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume technologies: %w", mapDBError(err))
	}
	technologies, err := mappers.TechnologiesToDomain(rows)
	if err != nil {
		return nil, fmt.Errorf("map resume technologies: %w", err)
	}
	return technologies, nil
}
