package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
)

func persistContacts(tx *gorm.DB, resume *domain.Resume) error {
	contacts := resume.Contacts()
	if contacts == nil {
		var current gormmodel.ResumeContacts
		err := tx.Where("resume_id = ?", resume.ID().UUID()).First(&current).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("load resume contacts: %w", mapDBError(err))
		}
		if err := tx.Where("contacts_id = ?", current.ID).Delete(&gormmodel.ResumeContactHidden{}).Error; err != nil {
			return fmt.Errorf("delete resume contact hidden: %w", mapDBError(err))
		}
		if err := tx.Where("contacts_id = ?", current.ID).Delete(&gormmodel.ResumeContactsText{}).Error; err != nil {
			return fmt.Errorf("delete resume contacts text: %w", mapDBError(err))
		}
		if err := tx.Delete(&current).Error; err != nil {
			return fmt.Errorf("delete resume contacts: %w", mapDBError(err))
		}
		return nil
	}

	m, hidden, text, err := mappers.ContactsToDBModels(contacts)
	if err != nil {
		return fmt.Errorf("map resume contacts: %w", err)
	}
	var current gormmodel.ResumeContacts
	err = tx.Where("resume_id = ?", resume.ID().UUID()).First(&current).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(&m).Error; err != nil {
			return fmt.Errorf("create resume contacts: %w", mapDBError(err))
		}
	} else if err != nil {
		return fmt.Errorf("load resume contacts: %w", mapDBError(err))
	} else {
		if current.ID != m.ID {
			return fmt.Errorf("%w: contacts identity mismatch", domain.ErrResumeInvariant)
		}
		if err := tx.Model(&gormmodel.ResumeContacts{}).Where("id = ? AND resume_id = ?", current.ID, resume.ID().UUID()).Updates(map[string]any{"phone": m.Phone, "telegram": m.Telegram, "email": m.Email}).Error; err != nil {
			return fmt.Errorf("update resume contacts: %w", mapDBError(err))
		}
	}
	if err := tx.Where("contacts_id = ?", m.ID).Delete(&gormmodel.ResumeContactHidden{}).Error; err != nil {
		return fmt.Errorf("replace resume contact hidden: %w", mapDBError(err))
	}
	if len(hidden) > 0 {
		if err := tx.Create(&hidden).Error; err != nil {
			return fmt.Errorf("create resume contact hidden: %w", mapDBError(err))
		}
	}
	if err := tx.Where("contacts_id = ?", m.ID).Delete(&gormmodel.ResumeContactsText{}).Error; err != nil {
		return fmt.Errorf("replace resume contacts text: %w", mapDBError(err))
	}
	if text != nil {
		if err := tx.Create(text).Error; err != nil {
			return fmt.Errorf("create resume contacts text: %w", mapDBError(err))
		}
	}
	return nil
}
