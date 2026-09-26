package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ProfileInput struct {
	Name        string
	Surname     string
	Patronymic  *string
	Position    string
	PhotoPath   valueobjects.ExternalURL
	DateOfBirth valueobjects.Date
	Location    string
	Relocation  bool
	WorkFormats []enums.WorkFormat
}

func (i ProfileInput) params(id, resumeID valueobjects.ID) domain.ProfileParams {
	return domain.ProfileParams{
		ID:          id,
		ResumeID:    resumeID,
		Name:        i.Name,
		Surname:     i.Surname,
		Patronymic:  i.Patronymic,
		Position:    i.Position,
		PhotoPath:   i.PhotoPath,
		DateOfBirth: i.DateOfBirth,
		Location:    i.Location,
		Relocation:  i.Relocation,
		WorkFormats: append([]enums.WorkFormat(nil), i.WorkFormats...),
	}
}

type CreateResumeProfileCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Profile      ProfileInput
}
type UpdateResumeProfileCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Profile  ProfileInput
}

type CreateResumeProfileHandler struct {
	repository ports.CreateResumeProfileRepository
	uow        UnitOfWork
}

func NewCreateResumeProfileHandler(r ports.CreateResumeProfileRepository, uow UnitOfWork) *CreateResumeProfileHandler {
	return &CreateResumeProfileHandler{repository: r, uow: uow}
}
func (h *CreateResumeProfileHandler) Handle(ctx context.Context, c CreateResumeProfileCommand) (valueobjects.ID, error) {
	if h == nil || h.repository == nil || h.uow == nil {
		return valueobjects.ID{}, fmt.Errorf("create resume profile handler dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	var created *domain.Profile
	params := c.Profile.params(c.ID, c.ResumeID)
	err := h.uow.Do(ctx, func(tx context.Context) error {
		p, e := domain.NewProfile(params)
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
	updater ports.UpdateResumeProfileRepository
	uow     UnitOfWork
}

func NewUpdateResumeProfileHandler(u ports.UpdateResumeProfileRepository, uow UnitOfWork) *UpdateResumeProfileHandler {
	return &UpdateResumeProfileHandler{updater: u, uow: uow}
}
func (h *UpdateResumeProfileHandler) Handle(ctx context.Context, c UpdateResumeProfileCommand) error {
	if h == nil || h.updater == nil || h.uow == nil {
		return fmt.Errorf("update resume profile handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() {
		return fmt.Errorf("invalid update resume profile command")
	}
	params := c.Profile.params(valueobjects.NewID(), c.ResumeID)
	current, e := domain.NewProfile(params)
	if e != nil {
		return e
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.updater.UpdateProfile(tx, ports.UpdateResumeProfileRepositoryCommand{Profile: current, UserID: c.UserID})
	})
}
