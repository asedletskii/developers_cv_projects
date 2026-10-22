package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
)

func persistProfile(tx *gorm.DB, resume *domain.Resume) error {
	p := resume.Profile()
	if p == nil {
		if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeProfile{}).Error; err != nil {
			return fmt.Errorf("delete resume profile: %w", mapDBError(err))
		}
		return nil
	}
	m, formats := mappers.ProfileToDBModels(p)
	result := tx.Model(&gormmodel.ResumeProfile{}).
		Where("id = ? AND resume_id = ?", m.ID, resume.ID().UUID()).
		Updates(map[string]any{
			"name": m.Name, "surname": m.Surname, "patronymic": m.Patronymic,
			"position": m.Position, "photo_url": m.PhotoURL, "date_of_birth": m.DateOfBirth,
			"location": m.Location, "relocation": m.Relocation,
		})
	if result.Error != nil {
		return fmt.Errorf("update resume profile: %w", mapDBError(result.Error))
	}
	if result.RowsAffected == 0 {
		if err := tx.Create(&m).Error; err != nil {
			return fmt.Errorf("create resume profile: %w", mapDBError(err))
		}
	}
	if err := tx.Where("profile_id = ?", m.ID).Delete(&gormmodel.ResumeProfileWorkFormat{}).Error; err != nil {
		return fmt.Errorf("replace resume profile work formats: %w", mapDBError(err))
	}
	for i := range formats {
		formats[i].ProfileID = m.ID
	}
	if len(formats) > 0 {
		if err := tx.Create(&formats).Error; err != nil {
			return fmt.Errorf("create resume profile work formats: %w", mapDBError(err))
		}
	}
	return nil
}

func loadResumeProfile(ctx context.Context, tx *gorm.DB, resumeID uuid.UUID) (*domain.Profile, error) {
	var profileModel gormmodel.ResumeProfile
	err := tx.WithContext(ctx).Where("resume_id = ?", resumeID).First(&profileModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get resume profile: %w", mapDBError(err))
	}

	var formats []gormmodel.ResumeProfileWorkFormat
	if err := tx.WithContext(ctx).Where("profile_id = ?", profileModel.ID).Order("position").Find(&formats).Error; err != nil {
		return nil, fmt.Errorf("get resume profile work formats: %w", mapDBError(err))
	}
	profile, err := mappers.ProfileToDomain(profileModel, formats)
	if err != nil {
		return nil, fmt.Errorf("map resume profile: %w", err)
	}
	return profile, nil
}
