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
	if technologies == nil {
		if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeTechnology{}).Error; err != nil {
			return fmt.Errorf("delete resume technologies: %w", mapDBError(err))
		}
		return nil
	}

	models := make([]gormmodel.ResumeTechnology, 0, len(technologies.Entries()))
	seenIDs := make(map[sharedvalueobjects.ID]struct{}, len(technologies.Entries()))
	for _, entry := range technologies.Entries() {
		id := entry.ID()
		if id.IsZero() {
			return fmt.Errorf("%w: technology id is required", domain.ErrResumeInvariant)
		}
		if _, ok := seenIDs[id]; ok {
			return fmt.Errorf("%w: duplicate technology id", domain.ErrInvalidEntity)
		}
		seenIDs[id] = struct{}{}
		model, err := mappers.TechnologyToDB(entry, resume.ID())
		if err != nil {
			return fmt.Errorf("map resume technology: %w", err)
		}
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeTechnology{}).Error; err != nil {
		return fmt.Errorf("replace resume technologies: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume technologies: %w", mapDBError(err))
		}
	}
	return nil
}

func loadTechnologies(tx *gorm.DB, resumeID sharedvalueobjects.ID) (*domain.ResumeTechnologies, error) {
	var rows []gormmodel.ResumeTechnology
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume technologies: %w", mapDBError(err))
	}
	return mappers.TechnologiesToDomain(rows, resumeID)
}
