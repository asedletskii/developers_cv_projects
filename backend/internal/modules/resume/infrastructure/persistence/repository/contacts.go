package repository

import (
	"context"
	"errors"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"

	"gorm.io/gorm"
)

var _ ports.CreateResumeContactsRepository = (*ContactsRepository)(nil)
var _ ports.GetResumeContactsRepository = (*ContactsRepository)(nil)
var _ ports.UpdateResumeContactsRepository = (*ContactsRepository)(nil)

func (r *ContactsRepository) CreateContacts(ctx context.Context, c ports.CreateResumeContactsRepositoryCommand) error {
	tx, e := database.DBFromTransactionContext(ctx)
	if e != nil {
		return fmt.Errorf("create contacts: %w", e)
	}
	m, h, t, e := contactsModels(c.Contacts)
	if e != nil {
		return e
	}
	var n int64
	if e = tx.WithContext(ctx).Model(&gormmodel.ResumeRoot{}).Where("id = ? AND user_id = ?", m.ResumeID, c.UserID.UUID()).Count(&n).Error; e != nil {
		return fmt.Errorf("create contacts ownership: %w", e)
	}
	if n == 0 {
		return domain.ErrEntityNotFound
	}
	if e = tx.WithContext(ctx).Model(&gormmodel.ResumeContacts{}).Where("resume_id = ?", m.ResumeID).Count(&n).Error; e != nil {
		return fmt.Errorf("check contacts uniqueness: %w", e)
	}
	if n > 0 {
		return domain.ErrDuplicateEntity
	}
	if e = tx.WithContext(ctx).Create(&m).Error; e != nil {
		return fmt.Errorf("create contacts: %w", e)
	}
	if len(h) > 0 {
		if e = tx.WithContext(ctx).Create(&h).Error; e != nil {
			return fmt.Errorf("create hidden contacts: %w", e)
		}
	}
	if t != nil {
		if e = tx.WithContext(ctx).Create(t).Error; e != nil {
			return fmt.Errorf("create contacts text: %w", e)
		}
	}
	return nil
}
func (r *ContactsRepository) GetContacts(ctx context.Context, q ports.GetResumeContactsRepositoryCommand) (*domain.Contacts, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get contacts: repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db)
	var m gormmodel.ResumeContacts
	if e := tx.WithContext(ctx).Where("resume_id = ? AND resume_id IN (SELECT id FROM resume_roots WHERE id = ? AND user_id = ?)", q.ResumeID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).First(&m).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get contacts: %w", e)
	}
	var hidden []gormmodel.ResumeContactHidden
	if e := tx.WithContext(ctx).Where("contacts_id = ?", m.ID).Order("position").Find(&hidden).Error; e != nil {
		return nil, fmt.Errorf("get hidden contacts: %w", e)
	}
	var t gormmodel.ResumeContactsText
	te := tx.WithContext(ctx).Where("contacts_id = ?", m.ID).First(&t).Error
	var tp *gormmodel.ResumeContactsText
	if te == nil {
		tp = &t
	} else if !errors.Is(te, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get contacts text: %w", te)
	}
	return mapContacts(m, hidden, tp)
}
func (r *ContactsRepository) UpdateContacts(ctx context.Context, c ports.UpdateResumeContactsRepositoryCommand) error {
	tx, e := database.DBFromTransactionContext(ctx)
	if e != nil {
		return fmt.Errorf("update contacts: %w", e)
	}
	m, h, t, e := contactsModels(c.Contacts)
	if e != nil {
		return e
	}
	result := tx.WithContext(ctx).Model(&gormmodel.ResumeContacts{}).Where("id = ? AND resume_id IN (SELECT id FROM resume_roots WHERE id = ? AND user_id = ?)", m.ID, m.ResumeID, c.UserID.UUID()).Updates(map[string]any{"phone": m.Phone, "telegram": m.Telegram, "email": m.Email})
	if result.Error != nil {
		return fmt.Errorf("update contacts: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	if e = tx.WithContext(ctx).Where("contacts_id = ?", m.ID).Delete(&gormmodel.ResumeContactHidden{}).Error; e != nil {
		return fmt.Errorf("replace hidden contacts: %w", e)
	}
	if len(h) > 0 {
		if e = tx.WithContext(ctx).Create(&h).Error; e != nil {
			return fmt.Errorf("create hidden contacts: %w", e)
		}
	}
	if e = tx.WithContext(ctx).Where("contacts_id = ?", m.ID).Delete(&gormmodel.ResumeContactsText{}).Error; e != nil {
		return fmt.Errorf("replace contacts text: %w", e)
	}
	if t != nil {
		if e = tx.WithContext(ctx).Create(t).Error; e != nil {
			return fmt.Errorf("create contacts text: %w", e)
		}
	}
	return nil
}
