package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type ContactsInput struct {
	Phone          *string
	Telegram       *string
	Email          *string
	HiddenContacts []string
	ContactsText   *ContactsTextInput
}
type ContactsTextInput struct{ Title, Subtitle *LocalizedTextInput }
type LocalizedTextInput struct{ RU, EN string }

func (i ContactsInput) params(id, resumeID valueobjects.ID) (domain.ContactsParams, error) {
	var tg *valueobjects.Telegram
	if i.Telegram != nil {
		v, e := valueobjects.NewTelegram(*i.Telegram)
		if e != nil {
			return domain.ContactsParams{}, e
		}
		tg = &v
	}
	var em *valueobjects.Email
	if i.Email != nil {
		v, e := valueobjects.NewEmail(*i.Email)
		if e != nil {
			return domain.ContactsParams{}, e
		}
		em = &v
	}
	cs, e := entities.NewContacts(entities.ContactsParams{Phone: stringValue(i.Phone), Telegram: tg, Email: em})
	if e != nil {
		return domain.ContactsParams{}, e
	}
	hidden := make([]enums.ContactKey, len(i.HiddenContacts))
	for n, v := range i.HiddenContacts {
		hidden[n] = enums.ContactKey(v)
	}
	var text *valueobjects.ContactsText
	if i.ContactsText != nil {
		var title, subtitle *valueobjects.LocalizedText
		if i.ContactsText.Title != nil {
			v, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: i.ContactsText.Title.RU, EN: i.ContactsText.Title.EN})
			if e != nil {
				return domain.ContactsParams{}, e
			}
			title = &v
		}
		if i.ContactsText.Subtitle != nil {
			v, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: i.ContactsText.Subtitle.RU, EN: i.ContactsText.Subtitle.EN})
			if e != nil {
				return domain.ContactsParams{}, e
			}
			subtitle = &v
		}
		v, e := valueobjects.NewContactsText(valueobjects.ContactsTextParams{Title: title, Subtitle: subtitle})
		if e != nil {
			return domain.ContactsParams{}, e
		}
		text = &v
	}
	return domain.ContactsParams{ID: id, ResumeID: resumeID, Contacts: cs, HiddenContacts: hidden, Text: text}, nil
}
func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
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

type CreateResumeContactsRepository = ports.CreateResumeContactsRepository
type CreateResumeContactsRepositoryCommand = ports.CreateResumeContactsRepositoryCommand
type GetResumeContactsRepository = ports.GetResumeContactsRepository
type UpdateResumeContactsRepository = ports.UpdateResumeContactsRepository
type UpdateResumeContactsRepositoryCommand = ports.UpdateResumeContactsRepositoryCommand

type CreateResumeContactsHandler struct {
	repository CreateResumeContactsRepository
	uow        UnitOfWork
}

func NewCreateResumeContactsHandler(r CreateResumeContactsRepository, uow UnitOfWork) *CreateResumeContactsHandler {
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
	repository GetResumeContactsRepository
	updater    UpdateResumeContactsRepository
	uow        UnitOfWork
}

func NewUpdateResumeContactsHandler(r GetResumeContactsRepository, u UpdateResumeContactsRepository, uow UnitOfWork) *UpdateResumeContactsHandler {
	return &UpdateResumeContactsHandler{repository: r, updater: u, uow: uow}
}
func (h *UpdateResumeContactsHandler) Handle(ctx context.Context, c UpdateResumeContactsCommand) error {
	if h == nil || h.repository == nil || h.updater == nil || h.uow == nil {
		return fmt.Errorf("update resume contacts handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() {
		return fmt.Errorf("invalid update resume contacts command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		current, e := h.repository.GetContacts(tx, ports.GetResumeContactsRepositoryCommand{ResumeID: c.ResumeID, UserID: c.UserID})
		if e != nil {
			return e
		}
		params, e := c.Contacts.params(current.ID(), current.ResumeID())
		if e != nil {
			return e
		}
		if e := current.Replace(params); e != nil {
			return e
		}
		return h.updater.UpdateContacts(tx, ports.UpdateResumeContactsRepositoryCommand{UserID: c.UserID, Contacts: current})
	})
}
