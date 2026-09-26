package readers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type ContactsReader struct{ db *gorm.DB }

func NewContactsReader(db *gorm.DB) *ContactsReader { return &ContactsReader{db: db} }

func (r *ContactsReader) GetResumeContacts(ctx context.Context, q queries.GetResumeContactsQuery) (queries.ResumeContactsView, error) {
	db, err := resolveDB(ctx, r.db, "get resume contacts")
	if err != nil {
		return queries.ResumeContactsView{}, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return queries.ResumeContactsView{}, fmt.Errorf("get resume contacts: invalid query")
	}

	var model gormmodel.ResumeContacts
	err = db.Joins("JOIN resumes ON resumes.id = resume_contacts.resume_id").Where("resume_contacts.resume_id = ? AND resumes.user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.ResumeContactsView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.ResumeContactsView{}, fmt.Errorf("get resume contacts: %w", err)
	}

	var hidden []gormmodel.ResumeContactHidden
	if err := db.Where("contacts_id = ?", model.ID).Order("position").Find(&hidden).Error; err != nil {
		return queries.ResumeContactsView{}, fmt.Errorf("get hidden contacts: %w", err)
	}
	hiddenContacts := make([]string, 0, len(hidden))
	for _, item := range hidden {
		hiddenContacts = append(hiddenContacts, item.Value)
	}

	view := queries.ResumeContactsView{
		ID: model.ID.String(), ResumeID: model.ResumeID.String(), Phone: model.Phone,
		Telegram: model.Telegram, Email: model.Email, HiddenContacts: hiddenContacts,
	}
	var text gormmodel.ResumeContactsText
	if err := db.Where("contacts_id = ?", model.ID).First(&text).Error; err == nil {
		view.ContactsText = &queries.ResumeContactsTextView{Title: text.Title, Subtitle: text.Subtitle}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.ResumeContactsView{}, fmt.Errorf("get contacts text: %w", err)
	}
	return view, nil
}

var _ = uuid.Nil
