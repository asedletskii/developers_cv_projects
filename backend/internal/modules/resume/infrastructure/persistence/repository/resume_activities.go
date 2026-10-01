package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistActivities(tx *gorm.DB, resume *domain.Resume) error {
	models := make([]gormmodel.ResumeActivity, 0, len(resume.Activities()))
	for position, item := range resume.Activities() {
		if item == nil {
			return fmt.Errorf("%w: nil activity", domain.ErrResumeInvariant)
		}
		if item.ResumeID() != resume.ID() {
			return fmt.Errorf("%w: activity does not belong to resume", domain.ErrResumeInvariant)
		}
		model, err := mappers.ActivityToDB(item)
		if err != nil {
			return err
		}
		model.SortOrder = position
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeActivity{}).Error; err != nil {
		return fmt.Errorf("replace resume activities: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume activities: %w", mapDBError(err))
		}
	}
	return nil
}

func loadActivities(tx *gorm.DB, resumeID sharedvalueobjects.ID) ([]*domain.Activity, error) {
	var rows []gormmodel.ResumeActivity
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume activities: %w", mapDBError(err))
	}
	out := make([]*domain.Activity, 0, len(rows))
	for _, row := range rows {
		x, err := mappers.ActivityToDomain(row)
		if err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	if len(out) == 0 {
		return nil, nil
	}
	return out, nil
}
