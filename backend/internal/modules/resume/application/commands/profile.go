package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ProfileInput struct {
	ID          identity.ID
	Name        string
	Surname     string
	Patronymic  *string
	Position    string
	PhotoPath   domain.ExternalURL
	DateOfBirth domain.Date
	Location    string
	Relocation  bool
	WorkFormats []domain.WorkFormat
}

func (i ProfileInput) toDomain(resumeID identity.ID) domain.ProfileParams {
	return domain.ProfileParams{ID: i.ID, ResumeID: resumeID, Name: i.Name, Surname: i.Surname, Patronymic: i.Patronymic, Position: i.Position, PhotoPath: i.PhotoPath, DateOfBirth: i.DateOfBirth, Location: i.Location, Relocation: i.Relocation, WorkFormats: append([]domain.WorkFormat(nil), i.WorkFormats...)}
}

type ReplaceProfileCommand struct {
	ResumeID identity.ID
	UserID   identity.ID
	Version  int64
	Input    ProfileInput
}

type ReplaceProfileHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewReplaceProfileHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceProfileHandler {
	return &ReplaceProfileHandler{repository: repository, uow: uow}
}

func (h *ReplaceProfileHandler) Handle(ctx context.Context, c ReplaceProfileCommand) error {
	if h == nil {
		return fmt.Errorf("replace profile handler dependencies are nil")
	}

	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace profile", func(resume *domain.Resume) error {
		return resume.ReplaceProfile(c.Input.toDomain(c.ResumeID))
	})
}
