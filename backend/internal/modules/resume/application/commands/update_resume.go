package commands

import (
	"context"
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type UpdateResumeCommand struct {
	ResumeID       valueobjects.ID
	UserID         useridentity.ID
	Name           string
	Surname        string
	Patronymic     *string
	Position       string
	Availability   *valueobjects.Availability
	ResumeLink     *valueobjects.AssetPath
	Theme          *valueobjects.ThemeSettings
	Sections       []valueobjects.SectionSettings
	HiddenContacts []string
	ContactsText   *valueobjects.ContactsText
	Greeting       valueobjects.Greeting
	PhotoPath      valueobjects.AssetPath
	DateOfBirth    valueobjects.Date
	Location       string
	Relocation     bool
	WorkFormats    []string
	Contacts       entities.Contacts
	Languages      []entities.ResumeLanguageParams
	Experience     []entities.ExperienceParams
	Education      []entities.EducationParams
	Courses        []entities.CourseParams
	Activities     []entities.ActivityParams
	HardSkills     []string
	SoftSkills     []string
	Technologies   []entities.TechnologyParams
	Projects       []entities.ProjectParams
}

func (c UpdateResumeCommand) resumeInput() resumeInput {
	return resumeInput{
		Name: c.Name, Surname: c.Surname, Patronymic: c.Patronymic, Position: c.Position,
		Availability: c.Availability, ResumeLink: c.ResumeLink, Theme: c.Theme, Sections: c.Sections,
		HiddenContacts: c.HiddenContacts, ContactsText: c.ContactsText, Greeting: c.Greeting, PhotoPath: c.PhotoPath,
		DateOfBirth: c.DateOfBirth, Location: c.Location, Relocation: c.Relocation, WorkFormats: c.WorkFormats,
		Contacts: c.Contacts, Languages: c.Languages, Experience: c.Experience, Education: c.Education,
		Courses: c.Courses, Activities: c.Activities, HardSkills: c.HardSkills, SoftSkills: c.SoftSkills,
		Technologies: c.Technologies, Projects: c.Projects,
	}
}

type UpdateResumeHandler struct {
	repository       GetResumeRepository
	updateRepository UpdateResumeRepository
	uow              UnitOfWork
}

func NewUpdateResumeHandler(repository GetResumeRepository, updateRepository UpdateResumeRepository, uow UnitOfWork) *UpdateResumeHandler {
	return &UpdateResumeHandler{repository: repository, updateRepository: updateRepository, uow: uow}
}

func (h *UpdateResumeHandler) Handle(ctx context.Context, command UpdateResumeCommand) (valueobjects.ID, error) {
	if h == nil || h.repository == nil || h.updateRepository == nil || h.uow == nil {
		return valueobjects.ID{}, fmt.Errorf("update resume handler dependencies are nil")
	}
	if command.ResumeID.IsZero() || command.UserID.IsZero() {
		return valueobjects.ID{}, fmt.Errorf("invalid update resume command")
	}

	err := h.uow.Do(ctx, func(txCtx context.Context) error {
		current, err := h.repository.Get(txCtx, GetResumeRepositoryQuery{UserID: command.UserID, ResumeID: command.ResumeID})
		if err != nil {
			return err
		}

		replacement, err := buildResume(command.resumeInput(), command.ResumeID, command.UserID, current.Title())
		if err != nil {
			return err
		}
		return h.updateRepository.Update(txCtx, UpdateResumeRepositoryCommand{Resume: replacement})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return command.ResumeID, nil
}

type PatchResumeCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Title    string
}

type PatchResumeHandler struct {
	repository UpdateResumeTitleRepository
	uow        UnitOfWork
}

type PatchResumeResult struct {
	ResumeID valueobjects.ID
	Title    string
}

func NewPatchResumeHandler(repository UpdateResumeTitleRepository, uow UnitOfWork) *PatchResumeHandler {
	return &PatchResumeHandler{repository: repository, uow: uow}
}

func (h *PatchResumeHandler) Handle(ctx context.Context, command PatchResumeCommand) (PatchResumeResult, error) {
	if h == nil || h.repository == nil || h.uow == nil {
		return PatchResumeResult{}, fmt.Errorf("patch resume handler dependencies are nil")
	}
	if command.ResumeID.IsZero() || command.UserID.IsZero() {
		return PatchResumeResult{}, fmt.Errorf("invalid patch resume command")
	}

	if err := domain.ValidateResumeTitle(command.Title); err != nil {
		return PatchResumeResult{}, err
	}

	title := strings.TrimSpace(command.Title)
	err := h.uow.Do(ctx, func(txCtx context.Context) error {
		return h.repository.UpdateTitle(txCtx, UpdateResumeTitleRepositoryCommand{UserID: command.UserID, ResumeID: command.ResumeID, Title: title})
	})
	if err != nil {
		return PatchResumeResult{}, err
	}
	return PatchResumeResult{ResumeID: command.ResumeID, Title: title}, nil
}
