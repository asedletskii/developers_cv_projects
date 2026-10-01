package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistLanguages(tx *gorm.DB, resume *domain.Resume) error {
	languages := resume.Languages()
	if languages == nil {
		if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeLanguage{}).Error; err != nil {
			return fmt.Errorf("delete resume languages: %w", mapDBError(err))
		}
		return nil
	}

	entries := languages.Entries()
	models := make([]gormmodel.ResumeLanguage, 0, len(entries))
	seenIDs := make(map[identity.ID]struct{}, len(entries))
	seenLanguages := make(map[domain.SpokenLanguage]struct{}, len(entries))
	for _, entry := range entries {
		id := entry.ID()
		if id.IsZero() {
			return fmt.Errorf("%w: language id is required", domain.ErrResumeInvariant)
		}
		if _, ok := seenIDs[id]; ok {
			return fmt.Errorf("%w: duplicate language id", domain.ErrInvalidEntity)
		}
		seenIDs[id] = struct{}{}
		language := entry.Language()
		if _, ok := seenLanguages[language]; ok {
			return fmt.Errorf("%w: duplicate language", domain.ErrInvalidEntity)
		}
		seenLanguages[language] = struct{}{}
		model, err := mappers.LanguageToDB(entry, resume.ID())
		if err != nil {
			return fmt.Errorf("map resume language: %w", err)
		}
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeLanguage{}).Error; err != nil {
		return fmt.Errorf("replace resume languages: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume languages: %w", mapDBError(err))
		}
	}
	return nil
}

func loadLanguages(tx *gorm.DB, resumeID identity.ID) (*domain.ResumeLanguages, error) {
	var rows []gormmodel.ResumeLanguage
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume languages: %w", mapDBError(err))
	}
	return mappers.LanguagesToDomain(rows, resumeID)
}
