package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateResumeContactsCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Contacts     domain.ContactsParams
}
type GetResumeContactsCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type UpdateResumeContactsCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Contacts domain.ContactsParams
}

type CreateResumeContactsRepositoryCommand struct {
	UserID   useridentity.ID
	Contacts *domain.Contacts
}
type UpdateResumeContactsRepositoryCommand struct {
	UserID   useridentity.ID
	Contacts *domain.Contacts
}

type CreateResumeContactsRepository interface {
	CreateContacts(context.Context, CreateResumeContactsRepositoryCommand) error
}
type GetResumeContactsRepository interface {
	GetContacts(context.Context, GetResumeContactsCommand) (*domain.Contacts, error)
}
type UpdateResumeContactsRepository interface {
	UpdateContacts(context.Context, UpdateResumeContactsRepositoryCommand) error
}

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
	c.Contacts.ID = c.ID
	c.Contacts.ResumeID = c.ResumeID
	var created *domain.Contacts
	err := h.uow.Do(ctx, func(tx context.Context) error {
		a, e := domain.NewContacts(c.Contacts)
		if e != nil {
			return e
		}
		created = a
		return h.repository.CreateContacts(tx, CreateResumeContactsRepositoryCommand{UserID: c.UserID, Contacts: a})
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
		current, e := h.repository.GetContacts(tx, GetResumeContactsCommand{ResumeID: c.ResumeID, UserID: c.UserID})
		if e != nil {
			return e
		}
		if e := current.Replace(c.Contacts); e != nil {
			return e
		}
		return h.updater.UpdateContacts(tx, UpdateResumeContactsRepositoryCommand{UserID: c.UserID, Contacts: current})
	})
}
