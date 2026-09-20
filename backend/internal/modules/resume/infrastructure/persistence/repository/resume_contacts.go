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

func persistContacts(tx *gorm.DB, resume *domain.Resume) error {
	contacts := resume.Contacts()
	if contacts == nil {
		if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeContacts{}).Error; err != nil {
			return fmt.Errorf("delete resume contacts: %w", mapDBError(err))
		}
		return nil
	}

	m, hidden, text := mappers.ContactsToDBModels(contacts)
	result := tx.Model(&gormmodel.ResumeContacts{}).
		Where("id = ? AND resume_id = ?", m.ID, resume.ID().UUID()).
		Updates(map[string]any{"phone": m.Phone, "telegram": m.Telegram, "email": m.Email})
	if result.Error != nil {
		return fmt.Errorf("update resume contacts: %w", mapDBError(result.Error))
	}
	if result.RowsAffected == 0 {
		if err := tx.Create(&m).Error; err != nil {
			return fmt.Errorf("create resume contacts: %w", mapDBError(err))
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

func loadResumeContacts(ctx context.Context, tx *gorm.DB, resumeID uuid.UUID) (*domain.Contacts, error) {
	var contactsModel gormmodel.ResumeContacts
	err := tx.WithContext(ctx).Where("resume_id = ?", resumeID).First(&contactsModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get resume contacts: %w", mapDBError(err))
	}

	var hidden []gormmodel.ResumeContactHidden
	if err := tx.WithContext(ctx).Where("contacts_id = ?", contactsModel.ID).Order("position").Find(&hidden).Error; err != nil {
		return nil, fmt.Errorf("get resume contact hidden: %w", mapDBError(err))
	}
	var textModel gormmodel.ResumeContactsText
	textErr := tx.WithContext(ctx).Where("contacts_id = ?", contactsModel.ID).First(&textModel).Error
	var text *gormmodel.ResumeContactsText
	if textErr == nil {
		text = &textModel
	} else if !errors.Is(textErr, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get resume contacts text: %w", mapDBError(textErr))
	}

	contacts, err := mappers.ContactsToDomain(contactsModel, hidden, text)
	if err != nil {
		return nil, fmt.Errorf("map resume contacts: %w", err)
	}
	return contacts, nil
}
