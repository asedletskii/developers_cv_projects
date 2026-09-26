package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ContactsInput struct {
	Phone          *string
	Telegram       *valueobjects.Telegram
	Email          *valueobjects.Email
	HiddenContacts []enums.ContactKey
	ContactsText   *ContactsTextInput
}
type ContactsTextInput struct{ Title, Subtitle *string }

func (i ContactsInput) params(id, resumeID valueobjects.ID) (domain.ContactsParams, error) {
	phone := ""
	if i.Phone != nil {
		phone = *i.Phone
	}
	cs, e := entities.NewContacts(entities.ContactsParams{Phone: phone, Telegram: i.Telegram, Email: i.Email})
	if e != nil {
		return domain.ContactsParams{}, e
	}
	var text *valueobjects.ContactsText
	if i.ContactsText != nil {
		v, e := valueobjects.NewContactsText(valueobjects.ContactsTextParams{Title: i.ContactsText.Title, Subtitle: i.ContactsText.Subtitle})
		if e != nil {
			return domain.ContactsParams{}, e
		}
		text = &v
	}
	return domain.ContactsParams{ID: id, ResumeID: resumeID, Contacts: cs, HiddenContacts: append([]enums.ContactKey(nil), i.HiddenContacts...), Text: text}, nil
}

type CreateResumeContactsCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Contacts     ContactsInput
}
type UpdateResumeContactsCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Contacts ContactsInput
}

type CreateResumeContactsHandler struct {
	repository ports.CreateResumeContactsRepository
	uow        UnitOfWork
}

func NewCreateResumeContactsHandler(r ports.CreateResumeContactsRepository, uow UnitOfWork) *CreateResumeContactsHandler {
	return &CreateResumeContactsHandler{repository: r, uow: uow}
}
func (h *CreateResumeContactsHandler) Handle(ctx context.Context, c CreateResumeContactsCommand) (valueobjects.ID, error) {
	if h == nil || h.repository == nil || h.uow == nil {
		return valueobjects.ID{}, fmt.Errorf("create resume contacts handler dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	var created *domain.Contacts
	params, e := c.Contacts.params(c.ID, c.ResumeID)
	if e != nil {
		return valueobjects.ID{}, e
	}
	err := h.uow.Do(ctx, func(tx context.Context) error {
		a, e := domain.NewContacts(params)
		if e != nil {
			return e
		}
		created = a
		return h.repository.CreateContacts(tx, ports.CreateResumeContactsRepositoryCommand{UserID: c.UserID, Contacts: a})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return created.ID(), nil
}

type UpdateResumeContactsHandler struct {
	updater ports.UpdateResumeContactsRepository
	uow     UnitOfWork
}

func NewUpdateResumeContactsHandler(u ports.UpdateResumeContactsRepository, uow UnitOfWork) *UpdateResumeContactsHandler {
	return &UpdateResumeContactsHandler{updater: u, uow: uow}
}
func (h *UpdateResumeContactsHandler) Handle(ctx context.Context, c UpdateResumeContactsCommand) error {
	if h == nil || h.updater == nil || h.uow == nil {
		return fmt.Errorf("update resume contacts handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() {
		return fmt.Errorf("invalid update resume contacts command")
	}
	params, e := c.Contacts.params(valueobjects.NewID(), c.ResumeID)
	if e != nil {
		return e
	}
	current, e := domain.NewContacts(params)
	if e != nil {
		return e
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.updater.UpdateContacts(tx, ports.UpdateResumeContactsRepositoryCommand{UserID: c.UserID, Contacts: current})
	})
}
