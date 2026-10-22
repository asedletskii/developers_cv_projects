package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistLanguages(tx *gorm.DB, resume *domain.Resume) error {
	languages := resume.SpokenLanguages()
	models := mappers.ResumeSpokenLanguagesToDBModels(languages, resume.ID())
	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeSpokenLanguage{}).Error; err != nil {
		return fmt.Errorf("replace resume languages: %w", mapDBError(err))
	}
	if len(models) == 0 {
		return nil
	}
	if err := tx.Create(&models).Error; err != nil {
		return fmt.Errorf("create resume languages: %w", mapDBError(err))
	}
	return nil
}

func loadLanguages(tx *gorm.DB, resumeID sharedvalueobjects.ID) (*domain.ResumeSpokenLanguages, error) {
	var rows []gormmodel.ResumeSpokenLanguage
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume languages: %w", mapDBError(err))
	}
	languages, err := mappers.LanguagesToDomain(rows)
	if err != nil {
		return nil, fmt.Errorf("map resume languages: %w", err)
	}
	return languages, nil
}
