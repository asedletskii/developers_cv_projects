package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ContactsInput struct {
	ID             identity.ID
	Phone          *string
	Telegram       *domain.Telegram
	Email          *domain.Email
	HiddenContacts []domain.ContactKey
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
	contacts, err := domain.NewContactMethods(domain.ContactMethodsParams{Phone: phone, Telegram: i.Telegram, Email: i.Email})
	if err != nil {
		return domain.ContactsParams{}, err
	}

	var text *domain.ContactsText
	if i.ContactsText != nil {
		value, err := domain.NewContactsText(domain.ContactsTextParams{Title: i.ContactsText.Title, Subtitle: i.ContactsText.Subtitle})
		if err != nil {
			return domain.ContactsParams{}, err
		}
		text = &value
	}

	return domain.ContactsParams{
		ID:             i.ID,
		Contacts:       contacts,
		HiddenContacts: append([]domain.ContactKey(nil), i.HiddenContacts...),
		Text:           text,
	}, nil
}

type ReplaceContactsCommand struct {
	ResumeID identity.ID
	UserID   identity.ID
	Version  int64
	Input    ContactsInput
}

type ReplaceContactsHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewReplaceContactsHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceContactsHandler {
	return &ReplaceContactsHandler{repository: repository, uow: uow}
}

func (h *ReplaceContactsHandler) Handle(ctx context.Context, c ReplaceContactsCommand) error {
	if h == nil {
		return fmt.Errorf("replace contacts handler dependencies are nil")
	}
	input, err := c.Input.toDomain()
	if err != nil {
		return err
	}
	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace contacts", func(resume *domain.Resume) error {
		return resume.ReplaceContacts(input)
	})
}
