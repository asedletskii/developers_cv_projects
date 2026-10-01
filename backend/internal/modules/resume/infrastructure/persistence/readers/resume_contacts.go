package readers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func (r *ResumeReader) loadContacts(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var model gormmodel.ResumeContacts
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("get resume contacts: %w", err)
	}

	var hidden []gormmodel.ResumeContactHidden
	if err := db.WithContext(ctx).Where("contacts_id = ?", model.ID).Order("position").Find(&hidden).Error; err != nil {
		return fmt.Errorf("get hidden contacts: %w", err)
	}

	hiddenContacts := make([]string, 0, len(hidden))
	for _, item := range hidden {
		hiddenContacts = append(hiddenContacts, item.Value)
	}

	result := queries.ResumeContactsView{
		ID:             model.ID.String(),
		ResumeID:       model.ResumeID.String(),
		Phone:          model.Phone,
		Telegram:       model.Telegram,
		Email:          model.Email,
		HiddenContacts: hiddenContacts,
	}
	if err := r.loadContactsText(ctx, db, model.ID, &result); err != nil {
		return err
	}
	view.Contacts = &result
	return nil
}

func (r *ResumeReader) loadContactsText(ctx context.Context, db *gorm.DB, contactsID uuid.UUID, result *queries.ResumeContactsView) error {
	var text gormmodel.ResumeContactsText
	if err := db.WithContext(ctx).Where("contacts_id = ?", contactsID).First(&text).Error; err == nil {
		result.ContactsText = &queries.ResumeContactsTextView{Title: text.Title, Subtitle: text.Subtitle}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get contacts text: %w", err)
	}
	return nil
}

func (r *ResumeReader) GetContactForView(ctx context.Context, query queries.ContactViewQuery) (queries.ContactView, error) {
	if r == nil {
		return queries.ContactView{}, fmt.Errorf("get contact for view: reader is nil")
	}
	if query.ResumeID.IsZero() || !query.Contact.IsValid() {
		return queries.ContactView{}, fmt.Errorf("get contact for view: invalid query")
	}
	db, err := resolveDB(ctx, r.db, "get contact for view")
	if err != nil {
		return queries.ContactView{}, err
	}

	var resume gormmodel.Resume
	if err := db.WithContext(ctx).Where("id = ?", query.ResumeID.UUID()).First(&resume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ContactView{}, queries.ErrResumeNotFound
		}
		return queries.ContactView{}, fmt.Errorf("get resume for contact view: %w", err)
	}
	if resume.Status != string(domain.ResumeStatusOpened) {
		return queries.ContactView{}, queries.ErrResumeNotFound
	}

	var model gormmodel.ResumeContacts
	if err := db.WithContext(ctx).
		Where("resume_id = ?", query.ResumeID.UUID()).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ContactView{}, queries.ErrContactNotFound
		}
		return queries.ContactView{}, fmt.Errorf("get contact for view: %w", err)
	}

	var hidden gormmodel.ResumeContactHidden
	if err := db.WithContext(ctx).
		Where("contacts_id = ? AND value = ?", model.ID, string(query.Contact)).
		First(&hidden).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ContactView{}, queries.ErrContactNotHidden
		}
		return queries.ContactView{}, fmt.Errorf("get hidden contact for view: %w", err)
	}

	var value *string
	switch query.Contact {
	case domain.ContactKeyPhone:
		value = model.Phone
	case domain.ContactKeyTelegram:
		value = model.Telegram
	case domain.ContactKeyEmail:
		value = model.Email
	}
	if value == nil || *value == "" {
		return queries.ContactView{}, queries.ErrContactNotFound
	}

	return queries.ContactView{Contact: query.Contact, Value: *value}, nil
}
