package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ProfileInput struct {
	ID          valueobjects.ID
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

func (i ProfileInput) toDomain(resumeID valueobjects.ID) domain.ProfileParams {
	return domain.ProfileParams{ID: i.ID, ResumeID: resumeID, Name: i.Name, Surname: i.Surname, Patronymic: i.Patronymic, Position: i.Position, PhotoPath: i.PhotoPath, DateOfBirth: i.DateOfBirth, Location: i.Location, Relocation: i.Relocation, WorkFormats: append([]enums.WorkFormat(nil), i.WorkFormats...)}
}

type ReplaceProfileCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Input    ProfileInput
}

type ReplaceProfileHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceProfileHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceProfileHandler {
	return &ReplaceProfileHandler{repository: repository, uow: uow}
}

func (h *ReplaceProfileHandler) Handle(ctx context.Context, c ReplaceProfileCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace profile handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace profile command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if resume.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := resume.ReplaceProfile(c.Input.toDomain(c.ResumeID)); err != nil {
			return err
		}
		return h.repository.Save(tx, resume)
	})
}
