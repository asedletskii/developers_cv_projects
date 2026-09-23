package commands

import (
	"context"
	"fmt"

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
	repository ResumeRepository
	uow        UnitOfWork
}

func NewUpdateResumeHandler(repository ResumeRepository, uow UnitOfWork) *UpdateResumeHandler {
	return &UpdateResumeHandler{repository: repository, uow: uow}
}

func (h *UpdateResumeHandler) Handle(ctx context.Context, command UpdateResumeCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("update resume handler dependencies are nil")
	}
	if command.ResumeID.IsZero() || command.UserID.IsZero() {
		return fmt.Errorf("invalid update resume command")
	}

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		current, err := h.repository.GetByID(txCtx, command.UserID, command.ResumeID)
		if err != nil {
			return err
		}

		replacement, err := buildResume(command.resumeInput(), command.ResumeID, command.UserID, current.Title())
		if err != nil {
			return err
		}
		return h.repository.Save(txCtx, replacement)
	})
}

type PatchResumeCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Title    string
}

type PatchResumeHandler struct {
	repository ResumeRepository
	uow        UnitOfWork
}

func NewPatchResumeHandler(repository ResumeRepository, uow UnitOfWork) *PatchResumeHandler {
	return &PatchResumeHandler{repository: repository, uow: uow}
}

func (h *PatchResumeHandler) Handle(ctx context.Context, command PatchResumeCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("patch resume handler dependencies are nil")
	}
	if command.ResumeID.IsZero() || command.UserID.IsZero() {
		return fmt.Errorf("invalid patch resume command")
	}

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		resume, err := h.repository.GetByID(txCtx, command.UserID, command.ResumeID)
		if err != nil {
			return err
		}
		if err := resume.Rename(command.Title); err != nil {
			return err
		}
		return h.repository.Save(txCtx, resume)
	})
}
