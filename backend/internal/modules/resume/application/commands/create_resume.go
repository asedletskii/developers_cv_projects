package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateResumeCommand struct {
	ID             valueobjects.ID
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

type CreateResumeHandler struct {
	repository ResumeRepository
	uow        UnitOfWork
}

func NewCreateResumeHandler(repository ResumeRepository, uow UnitOfWork) *CreateResumeHandler {
	return &CreateResumeHandler{repository: repository, uow: uow}
}

func (h *CreateResumeHandler) Handle(ctx context.Context, command CreateResumeCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("create resume handler dependencies are nil")
	}

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		resume, err := h.buildResume(command)
		if err != nil {
			return err
		}
		return h.repository.Save(txCtx, resume)
	})
}

func (h *CreateResumeHandler) buildResume(command CreateResumeCommand) (*domain.Resume, error) {
	if command.ID.IsZero() {
		command.ID = valueobjects.NewID()
	}

	languages := make([]entities.ResumeLanguage, 0, len(command.Languages))
	for _, params := range command.Languages {
		language, err := entities.NewResumeLanguage(params)
		if err != nil {
			return nil, err
		}
		languages = append(languages, language)
	}

	experience := make([]entities.Experience, 0, len(command.Experience))
	for _, params := range command.Experience {
		if params.ID.IsZero() {
			params.ID = valueobjects.NewID()
		}
		item, err := entities.NewExperience(params)
		if err != nil {
			return nil, err
		}
		experience = append(experience, item)
	}

	education := make([]entities.Education, 0, len(command.Education))
	for _, params := range command.Education {
		if params.ID.IsZero() {
			params.ID = valueobjects.NewID()
		}
		item, err := entities.NewEducation(params)
		if err != nil {
			return nil, err
		}
		education = append(education, item)
	}

	courses := make([]entities.Course, 0, len(command.Courses))
	for _, params := range command.Courses {
		if params.ID.IsZero() {
			params.ID = valueobjects.NewID()
		}
		item, err := entities.NewCourse(params)
		if err != nil {
			return nil, err
		}
		courses = append(courses, item)
	}

	activities := make([]entities.Activity, 0, len(command.Activities))
	for _, params := range command.Activities {
		if params.ID.IsZero() {
			params.ID = valueobjects.NewID()
		}
		item, err := entities.NewActivity(params)
		if err != nil {
			return nil, err
		}
		activities = append(activities, item)
	}

	technologies := make([]entities.Technology, 0, len(command.Technologies))
	for _, params := range command.Technologies {
		item, err := entities.NewTechnology(params)
		if err != nil {
			return nil, err
		}
		technologies = append(technologies, item)
	}

	projects := make([]entities.Project, 0, len(command.Projects))
	for _, params := range command.Projects {
		if params.ID.IsZero() {
			params.ID = valueobjects.NewID()
		}
		item, err := entities.NewProject(params)
		if err != nil {
			return nil, err
		}
		projects = append(projects, item)
	}

	resume, err := domain.NewResume(domain.ResumeParams{
		ID:             command.ID,
		UserID:         command.UserID,
		Name:           command.Name,
		Surname:        command.Surname,
		Patronymic:     command.Patronymic,
		Position:       command.Position,
		Availability:   command.Availability,
		ResumeLink:     command.ResumeLink,
		Theme:          command.Theme,
		Sections:       command.Sections,
		HiddenContacts: toContactKeys(command.HiddenContacts),
		ContactsText:   command.ContactsText,
		Greeting:       command.Greeting,
		PhotoPath:      command.PhotoPath,
		DateOfBirth:    command.DateOfBirth,
		Location:       command.Location,
		Relocation:     command.Relocation,
		WorkFormats:    toWorkFormats(command.WorkFormats),
		Contacts:       command.Contacts,
		Languages:      languages,
		Experience:     experience,
		Education:      education,
		Courses:        courses,
		Activities:     activities,
		HardSkills:     toHardSkills(command.HardSkills),
		SoftSkills:     toSoftSkills(command.SoftSkills),
		Technologies:   technologies,
		Projects:       projects,
	})
	if err != nil {
		return nil, err
	}
	return resume, nil
}

func toContactKeys(values []string) []enums.ContactKey {
	result := make([]enums.ContactKey, len(values))
	for i, value := range values {
		result[i] = enums.ContactKey(value)
	}
	return result
}

func toWorkFormats(values []string) []enums.WorkFormat {
	result := make([]enums.WorkFormat, len(values))
	for i, value := range values {
		result[i] = enums.WorkFormat(value)
	}
	return result
}

func toHardSkills(values []string) []enums.HardSkill {
	result := make([]enums.HardSkill, len(values))
	for i, value := range values {
		result[i] = enums.HardSkill(value)
	}
	return result
}

func toSoftSkills(values []string) []enums.SoftSkill {
	result := make([]enums.SoftSkill, len(values))
	for i, value := range values {
		result[i] = enums.SoftSkill(value)
	}
	return result
}
