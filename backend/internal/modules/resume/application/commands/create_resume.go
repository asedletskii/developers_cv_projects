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

func (h *CreateResumeHandler) Handle(ctx context.Context, command CreateResumeCommand) (valueobjects.ID, error) {
	if h == nil || h.repository == nil || h.uow == nil {
		return valueobjects.ID{}, fmt.Errorf("create resume handler dependencies are nil")
	}

	if command.ID.IsZero() {
		command.ID = valueobjects.NewID()
	}

	err := h.uow.Do(ctx, func(txCtx context.Context) error {
		resume, err := buildResume(command.resumeInput(), command.ID, command.UserID, "")
		if err != nil {
			return err
		}
		return h.repository.Save(txCtx, resume)
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return command.ID, nil
}

type resumeInput struct {
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

func (c CreateResumeCommand) resumeInput() resumeInput {
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

func buildResume(
	input resumeInput,
	id valueobjects.ID,
	userID useridentity.ID,
	title string,
) (*domain.Resume, error) {
	languages := make([]entities.ResumeLanguage, 0, len(input.Languages))
	for _, params := range input.Languages {
		language, err := entities.NewResumeLanguage(params)
		if err != nil {
			return nil, err
		}
		languages = append(languages, language)
	}

	experience := make([]entities.Experience, 0, len(input.Experience))
	for _, params := range input.Experience {
		if params.ID.IsZero() {
			params.ID = valueobjects.NewID()
		}
		item, err := entities.NewExperience(params)
		if err != nil {
			return nil, err
		}
		experience = append(experience, item)
	}

	education := make([]entities.Education, 0, len(input.Education))
	for _, params := range input.Education {
		if params.ID.IsZero() {
			params.ID = valueobjects.NewID()
		}
		item, err := entities.NewEducation(params)
		if err != nil {
			return nil, err
		}
		education = append(education, item)
	}

	courses := make([]entities.Course, 0, len(input.Courses))
	for _, params := range input.Courses {
		if params.ID.IsZero() {
			params.ID = valueobjects.NewID()
		}
		item, err := entities.NewCourse(params)
		if err != nil {
			return nil, err
		}
		courses = append(courses, item)
	}

	activities := make([]entities.Activity, 0, len(input.Activities))
	for _, params := range input.Activities {
		if params.ID.IsZero() {
			params.ID = valueobjects.NewID()
		}
		item, err := entities.NewActivity(params)
		if err != nil {
			return nil, err
		}
		activities = append(activities, item)
	}

	technologies := make([]entities.Technology, 0, len(input.Technologies))
	for _, params := range input.Technologies {
		item, err := entities.NewTechnology(params)
		if err != nil {
			return nil, err
		}
		technologies = append(technologies, item)
	}

	projects := make([]entities.Project, 0, len(input.Projects))
	for _, params := range input.Projects {
		if params.ID.IsZero() {
			params.ID = valueobjects.NewID()
		}
		item, err := entities.NewProject(params)
		if err != nil {
			return nil, err
		}
		projects = append(projects, item)
	}

	return domain.NewResume(domain.ResumeParams{
		ID: id, UserID: userID, Name: input.Name, Surname: input.Surname, Patronymic: input.Patronymic,
		Position: input.Position, Title: title, Availability: input.Availability, ResumeLink: input.ResumeLink,
		Theme: input.Theme, Sections: input.Sections, HiddenContacts: toContactKeys(input.HiddenContacts),
		ContactsText: input.ContactsText, Greeting: input.Greeting, PhotoPath: input.PhotoPath, DateOfBirth: input.DateOfBirth,
		Location: input.Location, Relocation: input.Relocation, WorkFormats: toWorkFormats(input.WorkFormats),
		Contacts: input.Contacts, Languages: languages, Experience: experience, Education: education, Courses: courses,
		Activities: activities, HardSkills: toHardSkills(input.HardSkills), SoftSkills: toSoftSkills(input.SoftSkills),
		Technologies: technologies, Projects: projects,
	})
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
