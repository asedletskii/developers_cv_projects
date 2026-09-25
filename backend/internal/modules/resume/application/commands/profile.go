package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type ProfileInput struct {
	ResumeID    valueobjects.ID
	Name        string
	Surname     string
	Patronymic  *string
	Position    string
	PhotoPath   string
	DateOfBirth string
	Location    string
	Relocation  bool
	WorkFormats []string
}

func (i ProfileInput) params(id, resumeID valueobjects.ID) (domain.ProfileParams, error) {
	dob, err := valueobjects.ParseDate(i.DateOfBirth)
	if err != nil {
		return domain.ProfileParams{}, err
	}
	photo, err := valueobjects.NewAssetPath(i.PhotoPath)
	if err != nil {
		return domain.ProfileParams{}, err
	}
	formats := make([]enums.WorkFormat, len(i.WorkFormats))
	for n, v := range i.WorkFormats {
		formats[n] = enums.WorkFormat(v)
	}
	return domain.ProfileParams{ID: id, ResumeID: resumeID, Name: i.Name, Surname: i.Surname, Patronymic: i.Patronymic, Position: i.Position, PhotoPath: photo, DateOfBirth: dob, Location: i.Location, Relocation: i.Relocation, WorkFormats: formats}, nil
}

type CreateResumeProfileCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Profile      ProfileInput
}
type GetResumeProfileCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type UpdateResumeProfileCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Profile  ProfileInput
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
	var created *domain.Profile
	params, err := c.Profile.params(c.ID, c.ResumeID)
	if err != nil {
		return valueobjects.ID{}, err
	}
	err = h.uow.Do(ctx, func(tx context.Context) error {
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
		params, e := c.Profile.params(current.ID(), current.ResumeID())
		if e != nil {
			return e
		}
		if e := current.Replace(params); e != nil {
			return e
		}
		return h.updater.UpdateProfile(tx, ports.UpdateResumeProfileRepositoryCommand{Profile: current, UserID: c.UserID})
	})
}
