package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ContactsInput struct {
	ID             valueobjects.ID
	Phone          *string
	Telegram       *valueobjects.Telegram
	Email          *valueobjects.Email
	HiddenContacts []enums.ContactKey
	ContactsText   *ContactsTextInput
}

type ContactsTextInput struct {
	Title    *string
	Subtitle *string
}

func (i ContactsInput) toDomain() (domain.ContactsParams, error) {
	phone := ""
	if i.Phone != nil {
		phone = *i.Phone
	}
	contacts, err := entities.NewContacts(entities.ContactsParams{Phone: phone, Telegram: i.Telegram, Email: i.Email})
	if err != nil {
		return domain.ContactsParams{}, err
	}

	var text *valueobjects.ContactsText
	if i.ContactsText != nil {
		value, err := valueobjects.NewContactsText(valueobjects.ContactsTextParams{Title: i.ContactsText.Title, Subtitle: i.ContactsText.Subtitle})
		if err != nil {
			return domain.ContactsParams{}, err
		}
		text = &value
	}

	return domain.ContactsParams{
		ID:             i.ID,
		Contacts:       contacts,
		HiddenContacts: append([]enums.ContactKey(nil), i.HiddenContacts...),
		Text:           text,
	}, nil
}

type ReplaceContactsCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Input    ContactsInput
}

type ReplaceContactsHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceContactsHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceContactsHandler {
	return &ReplaceContactsHandler{repository: repository, uow: uow}
}

func (h *ReplaceContactsHandler) Handle(ctx context.Context, c ReplaceContactsCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace contacts handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace contacts command")
	}
	input, err := c.Input.toDomain()
	if err != nil {
		return err
	}

	return h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if resume.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := resume.ReplaceContacts(input); err != nil {
			return err
		}
		return h.repository.Save(tx, resume)
	})
}
