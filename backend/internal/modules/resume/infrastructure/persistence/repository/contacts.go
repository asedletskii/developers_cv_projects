package repository

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	"errors"
	"fmt"

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
	if e = tx.WithContext(ctx).Model(&gormmodel.Resume{}).Where("id = ? AND user_id = ?", m.ResumeID, c.UserID.UUID()).Count(&n).Error; e != nil {
		return fmt.Errorf("create contacts ownership: %w", e)
	}
	if n == 0 {
		return domain.ErrEntityNotFound
	}
	if e = tx.WithContext(ctx).Create(&m).Error; e != nil {
		if errors.Is(e, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntity
		}
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
	var root gormmodel.Resume
	if e := tx.WithContext(ctx).Where("id = ? AND user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).First(&root).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get contacts root: %w", e)
	}
	if e := tx.WithContext(ctx).Where("resume_id = ?", q.ResumeID.UUID()).First(&m).Error; e != nil {
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
	return mapContacts(m, hidden, tp, root.Language)
}
func (r *ContactsRepository) UpdateContacts(ctx context.Context, c ports.UpdateResumeContactsRepositoryCommand) error {
	tx, e := database.DBFromTransactionContext(ctx)
	if e != nil {
		return fmt.Errorf("update contacts: %w", e)
	}
	var current gormmodel.ResumeContacts
	if e = tx.WithContext(ctx).Where("resume_id = ? AND resume_id IN (SELECT id FROM resumes WHERE id = ? AND user_id = ?)", c.Contacts.ResumeID().UUID(), c.Contacts.ResumeID().UUID(), c.UserID.UUID()).First(&current).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return domain.ErrEntityNotFound
		}
		return fmt.Errorf("update contacts lookup: %w", e)
	}
	m, h, t, e := contactsModels(c.Contacts)
	if e != nil {
		return e
	}
	m.ID = current.ID
	m.ResumeID = current.ResumeID
	if t != nil {
		t.ContactsID = current.ID
	}
	for i := range h {
		h[i].ContactsID = current.ID
	}
	result := tx.WithContext(ctx).Model(&gormmodel.ResumeContacts{}).Where("id = ? AND resume_id IN (SELECT id FROM resumes WHERE id = ? AND user_id = ?)", m.ID, m.ResumeID, c.UserID.UUID()).Updates(map[string]any{"phone": m.Phone, "telegram": m.Telegram, "email": m.Email})
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

func contactsModels(c *domain.Contacts) (gormmodel.ResumeContacts, []gormmodel.ResumeContactHidden, *gormmodel.ResumeContactsText, error) {
	if c == nil {
		return gormmodel.ResumeContacts{}, nil, nil, fmt.Errorf("contacts is nil")
	}
	d := c.Contacts()
	m := gormmodel.ResumeContacts{ID: c.ID().UUID(), ResumeID: c.ResumeID().UUID()}
	if v := d.Phone(); v != nil {
		m.Phone = v
	}
	if v := d.Telegram(); v != nil {
		x := v.String()
		m.Telegram = &x
	}
	if v := d.Email(); v != nil {
		x := v.String()
		m.Email = &x
	}
	h := make([]gormmodel.ResumeContactHidden, 0, len(c.HiddenContacts()))
	for i, v := range c.HiddenContacts() {
		h = append(h, gormmodel.ResumeContactHidden{ContactsID: m.ID, Value: string(v), Position: i})
	}
	var t *gormmodel.ResumeContactsText
	if x := c.Text(); x != nil {
		t = &gormmodel.ResumeContactsText{ContactsID: m.ID}
		if v := x.Title(); v != nil {
			value := *v
			t.Title = &value
		}
		if v := x.Subtitle(); v != nil {
			value := *v
			t.Subtitle = &value
		}
	}
	return m, h, t, nil
}
func mapContacts(m gormmodel.ResumeContacts, h []gormmodel.ResumeContactHidden, t *gormmodel.ResumeContactsText, _ string) (*domain.Contacts, error) {
	var tg *valueobjects.Telegram
	if m.Telegram != nil {
		x, e := valueobjects.NewTelegram(*m.Telegram)
		if e != nil {
			return nil, e
		}
		tg = &x
	}
	var em *valueobjects.Email
	if m.Email != nil {
		x, e := valueobjects.NewEmail(*m.Email)
		if e != nil {
			return nil, e
		}
		em = &x
	}
	contacts, e := entities.NewContacts(entities.ContactsParams{Phone: pointerStringValue(m.Phone), Telegram: tg, Email: em})
	if e != nil {
		return nil, e
	}
	hidden := make([]enums.ContactKey, len(h))
	for i, v := range h {
		hidden[i] = enums.ContactKey(v.Value)
	}
	var text *valueobjects.ContactsText
	if t != nil {
		x, e := valueobjects.NewContactsText(valueobjects.ContactsTextParams{Title: t.Title, Subtitle: t.Subtitle})
		if e != nil {
			return nil, e
		}
		text = &x
	}

	id, e := valueobjects.IDFromUUID(m.ID)
	if e != nil {
		return nil, e
	}
	rid, e := valueobjects.IDFromUUID(m.ResumeID)
	if e != nil {
		return nil, e
	}
	return domain.NewContacts(domain.ContactsParams{ID: id, ResumeID: rid, Contacts: contacts, HiddenContacts: hidden, Text: text})
}

func pointerStringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
