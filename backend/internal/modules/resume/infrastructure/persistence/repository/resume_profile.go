package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
)

func persistProfile(tx *gorm.DB, resume *domain.Resume) error {
	p := resume.Profile()
	if p == nil {
		var ids []string
		if err := tx.Model(&gormmodel.ResumeProfile{}).Where("resume_id = ?", resume.ID().UUID()).Pluck("id", &ids).Error; err != nil {
			return fmt.Errorf("load resume profile ids: %w", mapDBError(err))
		}
		if len(ids) > 0 {
			if err := tx.Where("profile_id IN ?", ids).Delete(&gormmodel.ResumeProfileWorkFormat{}).Error; err != nil {
				return fmt.Errorf("delete resume profile work formats: %w", mapDBError(err))
			}
		}
		if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeProfile{}).Error; err != nil {
			return fmt.Errorf("delete resume profile: %w", mapDBError(err))
		}
		return nil
	}
	m, formats, err := mappers.ProfileToDBModels(p)
	if err != nil {
		return fmt.Errorf("map resume profile: %w", err)
	}
	var current gormmodel.ResumeProfile
	err = tx.Where("resume_id = ?", resume.ID().UUID()).First(&current).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(&m).Error; err != nil {
			return fmt.Errorf("create resume profile: %w", mapDBError(err))
		}
	} else if err != nil {
		return fmt.Errorf("load resume profile: %w", mapDBError(err))
	} else {
		if err := tx.Model(&gormmodel.ResumeProfile{}).Where("id = ?", current.ID).Updates(map[string]any{"name": m.Name, "surname": m.Surname, "patronymic": m.Patronymic, "position": m.Position, "photo_path": m.PhotoPath, "date_of_birth": m.DateOfBirth, "location": m.Location, "relocation": m.Relocation}).Error; err != nil {
			return fmt.Errorf("update resume profile: %w", mapDBError(err))
		}
		m.ID = current.ID
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
