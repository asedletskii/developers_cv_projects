package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
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
type GetResumeProfileCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type UpdateResumeProfileCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Profile  domain.ProfileParams
}

type CreateResumeProfileRepository = ports.CreateResumeProfileRepository
type CreateResumeProfileRepositoryCommand = ports.CreateResumeProfileRepositoryCommand
type GetResumeProfileRepository = ports.GetResumeProfileRepository
type UpdateResumeProfileRepository = ports.UpdateResumeProfileRepository
type UpdateResumeProfileRepositoryCommand = ports.UpdateResumeProfileRepositoryCommand

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
		return h.repository.CreateProfile(tx, ports.CreateResumeProfileRepositoryCommand{Profile: p, UserID: c.UserID})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return created.ID(), nil
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
		current, e := h.repository.GetProfile(tx, ports.GetResumeProfileRepositoryCommand{ResumeID: c.ResumeID, UserID: c.UserID})
		if e != nil {
			return e
		}
		if e := current.Replace(c.Profile); e != nil {
			return e
		}
		return h.updater.UpdateProfile(tx, ports.UpdateResumeProfileRepositoryCommand{Profile: current, UserID: c.UserID})
	})
}
