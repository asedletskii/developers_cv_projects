package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateResumeProfileCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Profile      domain.ProfileParams
}
type GetResumeProfileQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type UpdateResumeProfileCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Profile  domain.ProfileParams
}

type CreateResumeProfileRepository interface {
	CreateProfile(context.Context, CreateResumeProfileRepositoryCommand) error
}
type GetResumeProfileRepository interface {
	GetProfile(context.Context, GetResumeProfileQuery) (*domain.Profile, error)
}
type UpdateResumeProfileRepository interface {
	UpdateProfile(context.Context, UpdateResumeProfileRepositoryCommand) error
}
type CreateResumeProfileRepositoryCommand struct {
	Profile *domain.Profile
	UserID  useridentity.ID
}
type UpdateResumeProfileRepositoryCommand struct {
	Profile *domain.Profile
	UserID  useridentity.ID
}

type CreateResumeProfileHandler struct {
	repository CreateResumeProfileRepository
	uow        UnitOfWork
}

func NewCreateResumeProfileHandler(r CreateResumeProfileRepository, uow UnitOfWork) *CreateResumeProfileHandler {
	return &CreateResumeProfileHandler{repository: r, uow: uow}
}
func (h *CreateResumeProfileHandler) Handle(ctx context.Context, c CreateResumeProfileCommand) (valueobjects.ID, error) {
	if h == nil || h.repository == nil || h.uow == nil {
		return valueobjects.ID{}, fmt.Errorf("create resume profile handler dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	c.Profile.ID = c.ID
	c.Profile.ResumeID = c.ResumeID
	var created *domain.Profile
	err := h.uow.Do(ctx, func(tx context.Context) error {
		p, e := domain.NewProfile(c.Profile)
		if e != nil {
			return e
		}
		created = p
		return h.repository.CreateProfile(tx, CreateResumeProfileRepositoryCommand{Profile: p, UserID: c.UserID})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return created.ID(), nil
}

type GetResumeProfileHandler struct{ repository GetResumeProfileRepository }

func NewGetResumeProfileHandler(r GetResumeProfileRepository) *GetResumeProfileHandler {
	return &GetResumeProfileHandler{repository: r}
}
func (h *GetResumeProfileHandler) Handle(ctx context.Context, q GetResumeProfileQuery) (*domain.Profile, error) {
	if h == nil || h.repository == nil {
		return nil, fmt.Errorf("get resume profile handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("invalid get resume profile query")
	}
	return h.repository.GetProfile(ctx, q)
}

type UpdateResumeProfileHandler struct {
	repository GetResumeProfileRepository
	updater    UpdateResumeProfileRepository
	uow        UnitOfWork
}

func NewUpdateResumeProfileHandler(r GetResumeProfileRepository, u UpdateResumeProfileRepository, uow UnitOfWork) *UpdateResumeProfileHandler {
	return &UpdateResumeProfileHandler{repository: r, updater: u, uow: uow}
}
func (h *UpdateResumeProfileHandler) Handle(ctx context.Context, c UpdateResumeProfileCommand) error {
	if h == nil || h.repository == nil || h.updater == nil || h.uow == nil {
		return fmt.Errorf("update resume profile handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() {
		return fmt.Errorf("invalid update resume profile command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		current, e := h.repository.GetProfile(tx, GetResumeProfileQuery{ResumeID: c.ResumeID, UserID: c.UserID})
		if e != nil {
			return e
		}
		c.Profile.ID = current.ID()
		c.Profile.ResumeID = c.ResumeID
		p, e := domain.NewProfile(c.Profile)
		if e != nil {
			return e
		}
		return h.updater.UpdateProfile(tx, UpdateResumeProfileRepositoryCommand{Profile: p, UserID: c.UserID})
	})
}
