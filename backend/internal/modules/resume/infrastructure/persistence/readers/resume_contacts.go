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
	"developers-cv/backend/internal/platform/database"
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
		ID: model.ID.String(),

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
	var result queries.ContactView
	err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		db, err := resolveDB(snapshotCtx, r.db, "get contact for view")
		if err != nil {
			return err
		}
		var resume gormmodel.Resume
		if err := db.WithContext(snapshotCtx).Where("id = ?", query.ResumeID.UUID()).First(&resume).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return queries.ErrResumeNotFound
			}
			return fmt.Errorf("get resume for contact view: %w", err)
		}
		if resume.Status != string(domain.ResumeStatusOpened) {
			return queries.ErrResumeNotFound
		}

		var model gormmodel.ResumeContacts
		if err := db.WithContext(snapshotCtx).
			Where("resume_id = ?", query.ResumeID.UUID()).
			First(&model).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return queries.ErrContactNotFound
			}
			return fmt.Errorf("get contact for view: %w", err)
		}

		var hidden gormmodel.ResumeContactHidden
		if err := db.WithContext(snapshotCtx).
			Where("contacts_id = ? AND value = ?", model.ID, string(query.Contact)).
			First(&hidden).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return queries.ErrContactNotHidden
			}
			return fmt.Errorf("get hidden contact for view: %w", err)
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
			return queries.ErrContactNotFound
		}
		result = queries.ContactView{Contact: query.Contact, Value: *value}
		return nil
	})
	if err != nil {
		return queries.ContactView{}, err
	}
	return result, nil
}

func (r *ResumeReader) GetContactsForDownload(
	ctx context.Context,
	query queries.DownloadContactsQuery,
) (queries.DownloadContactsView, error) {
	if r == nil {
		return queries.DownloadContactsView{}, fmt.Errorf("get contacts for download: reader is nil")
	}
	var result queries.DownloadContactsView
	err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		db, err := resolveDB(snapshotCtx, r.db, "get contacts for download")
		if err != nil {
			return err
		}
		var resume gormmodel.Resume
		if err := db.WithContext(snapshotCtx).Where("id = ?", query.ResumeID.UUID()).First(&resume).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return queries.ErrResumeNotFound
			}
			return fmt.Errorf("get resume for contact download: %w", err)
		}
		if resume.Status != string(domain.ResumeStatusOpened) {
			return queries.ErrResumeNotFound
		}

		var contacts gormmodel.ResumeContacts
		if err := db.WithContext(snapshotCtx).
			Select("phone", "telegram", "email").
			Where("resume_id = ?", query.ResumeID.UUID()).
			First(&contacts).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("get resume contacts for download: %w", err)
		}
		result = queries.DownloadContactsView{
			Phone:    contacts.Phone,
			Telegram: contacts.Telegram,
			Email:    contacts.Email,
		}
		return nil
	})
	if err != nil {
		return queries.DownloadContactsView{}, err
	}
	return result, nil
}

func (r *ResumeReader) loadPublicContacts(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var metadata gormmodel.ResumeContacts
	if err := db.WithContext(ctx).
		Select("id").
		Where("resume_id = ?", resumeID).
		First(&metadata).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("get public resume contacts metadata: %w", err)
	}

	var hidden []gormmodel.ResumeContactHidden
	if err := db.WithContext(ctx).
		Select("value").
		Where("contacts_id = ?", metadata.ID).
		Order("position").
		Find(&hidden).Error; err != nil {
		return fmt.Errorf("get public hidden contacts: %w", err)
	}

	hiddenContacts := make([]string, 0, len(hidden))
	hiddenSet := make(map[string]struct{}, len(hidden))
	for _, item := range hidden {
		hiddenContacts = append(hiddenContacts, item.Value)
		hiddenSet[item.Value] = struct{}{}
	}

	columns := []string{"id", "phone", "telegram", "email"}
	selectedColumns := []string{"id"}
	for _, field := range columns[1:] {
		if _, isHidden := hiddenSet[field]; !isHidden {
			selectedColumns = append(selectedColumns, field)
		}
	}

	var model gormmodel.ResumeContacts
	if err := db.WithContext(ctx).
		Select(selectedColumns).
		Where("id = ?", metadata.ID).
		First(&model).Error; err != nil {
		return fmt.Errorf("get visible public contacts: %w", err)
	}

	result := queries.ResumeContactsView{
		ID:             model.ID.String(),
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
