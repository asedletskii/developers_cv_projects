package queries

import (
	"context"
	"fmt"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	contactkey "developers-cv/backend/internal/shared/domain/contact"
)

type PublicResumeReader interface {
	GetPublic(context.Context, resumecontracts.PublicResumeQuery) (ResumeView, error)
}

type PublicResumeQueryHandler struct {
	reader PublicResumeReader
}

func NewPublicResumeQueryHandler(reader PublicResumeReader) *PublicResumeQueryHandler {
	return &PublicResumeQueryHandler{reader: reader}
}

func (h *PublicResumeQueryHandler) Handle(ctx context.Context, query resumecontracts.PublicResumeQuery) (resumecontracts.PublicResumeView, error) {
	if h == nil || h.reader == nil {
		return resumecontracts.PublicResumeView{}, fmt.Errorf("public resume query handler dependencies are nil")
	}
	if query.ID.IsZero() {
		return resumecontracts.PublicResumeView{}, fmt.Errorf("invalid resume id")
	}
	view, err := h.reader.GetPublic(ctx, query)
	if err != nil {
		return resumecontracts.PublicResumeView{}, err
	}
	return ToPublicResumeView(view), nil
}

func ToPublicResumeView(view ResumeView) resumecontracts.PublicResumeView {
	result := resumecontracts.PublicResumeView{
		Title:        view.Title,
		Language:     view.Language,
		CreatedAt:    view.CreatedAt,
		UpdatedAt:    view.UpdatedAt,
		Overview:     toPublicOverview(view.Overview),
		Languages:    make([]resumecontracts.PublicLanguageView, 0, len(view.Languages)),
		Technologies: make([]resumecontracts.PublicTechnologyView, 0, len(view.Technologies)),
		Experiences:  make([]resumecontracts.PublicExperienceView, 0, len(view.Experiences)),
		Educations:   make([]resumecontracts.PublicEducationView, 0, len(view.Educations)),
		Courses:      make([]resumecontracts.PublicCourseView, 0, len(view.Courses)),
		Activities:   make([]resumecontracts.PublicActivityView, 0, len(view.Activities)),
		Projects:     make([]resumecontracts.PublicProjectView, 0, len(view.Projects)),
	}
	if view.Profile != nil {
		result.Profile = toPublicProfile(*view.Profile)
	}
	if view.Contacts != nil {
		result.Contacts = toPublicContacts(*view.Contacts)
	}
	if view.Skills != nil {
		result.Skills = toPublicSkills(*view.Skills)
	}
	for _, value := range view.Languages {
		result.Languages = append(result.Languages, resumecontracts.PublicLanguageView{ID: value.ID, Language: value.Language, Grade: cloneStringPointer(value.Grade), Position: value.Position})
	}
	for _, value := range view.Technologies {
		result.Technologies = append(result.Technologies, resumecontracts.PublicTechnologyView(value))
	}
	for _, value := range view.Experiences {
		stack := make([]resumecontracts.PublicExperienceStackItemView, 0, len(value.Stack))
		for _, item := range value.Stack {
			stack = append(stack, resumecontracts.PublicExperienceStackItemView(item))
		}
		result.Experiences = append(result.Experiences, resumecontracts.PublicExperienceView{
			ID: value.ID, Company: value.Company, Position: value.Position, Location: cloneStringPointer(value.Location),
			StartDate: value.StartDate, EndDate: cloneStringPointer(value.EndDate), Description: cloneStringPointer(value.Description),
			Achievements: cloneSlice(value.Achievements), Tasks: cloneSlice(value.Tasks), Stack: stack,
		})
	}
	for _, value := range view.Educations {
		result.Educations = append(result.Educations, resumecontracts.PublicEducationView{
			ID: value.ID, Location: value.Location, University: value.University, Grade: value.Grade, ProgramName: value.ProgramName,
			ProfileName: cloneStringPointer(value.ProfileName), Faculty: value.Faculty, Course: cloneIntPointer(value.Course), EndDate: cloneStringPointer(value.EndDate),
		})
	}
	for _, value := range view.Courses {
		result.Courses = append(result.Courses, resumecontracts.PublicCourseView{
			ID: value.ID, Name: value.Name, Organization: value.Organization, EndDate: value.EndDate,
			Hours: cloneIntPointer(value.Hours), CertificateLink: cloneStringPointer(value.CertificateLink),
		})
	}
	for _, value := range view.Activities {
		result.Activities = append(result.Activities, resumecontracts.PublicActivityView{
			ID: value.ID, Type: value.Type, Title: value.Title, Description: cloneStringPointer(value.Description),
			Date: cloneStringPointer(value.Date), Link: cloneStringPointer(value.Link),
		})
	}
	for _, value := range view.Projects {
		stack := make([]resumecontracts.PublicProjectTechnologyView, 0, len(value.Stack))
		for _, item := range value.Stack {
			stack = append(stack, resumecontracts.PublicProjectTechnologyView(item))
		}
		result.Projects = append(result.Projects, resumecontracts.PublicProjectView{
			ID: value.ID, Name: value.Name, Description: cloneStringPointer(value.Description),
			Screenshots: make([]string, 0), Link: cloneStringPointer(value.Link), Stack: stack,
			Contributors: make([]resumecontracts.PublicProjectContributorView, 0),
		})
	}
	return result
}

func toPublicOverview(value ResumeOverviewView) resumecontracts.PublicOverviewView {
	result := resumecontracts.PublicOverviewView{
		Greeting: resumecontracts.PublicGreeting{First: value.Greeting.First, Second: value.Greeting.Second, SecondHighlight: value.Greeting.SecondHighlight, Third: value.Greeting.Third},
		Sections: make([]resumecontracts.PublicSectionSettings, 0, len(value.Sections)),
	}
	if value.Theme != nil {
		result.Theme = &resumecontracts.PublicThemeSettings{Primary: value.Theme.Primary, Secondary: value.Theme.Secondary}
	}
	if value.Availability != nil {
		result.Availability = &resumecontracts.PublicAvailability{Status: value.Availability.Status, From: cloneStringPointer(value.Availability.From), Note: cloneStringPointer(value.Availability.Note)}
	}
	for _, section := range value.Sections {
		result.Sections = append(result.Sections, resumecontracts.PublicSectionSettings(section))
	}
	return result
}

func toPublicProfile(value ResumeProfileView) *resumecontracts.PublicProfileView {
	return &resumecontracts.PublicProfileView{
		ID: value.ID, Name: value.Name, Surname: value.Surname, Patronymic: cloneStringPointer(value.Patronymic),
		Position: value.Position, PhotoURL: value.PhotoURL, DateOfBirth: value.DateOfBirth, Location: value.Location,
		Relocation: value.Relocation, WorkFormat: cloneSlice(value.WorkFormat),
	}
}

func toPublicContacts(value ResumeContactsView) *resumecontracts.PublicContactsView {
	result := &resumecontracts.PublicContactsView{
		ID: value.ID, Phone: cloneStringPointer(value.Phone), Telegram: cloneStringPointer(value.Telegram), Email: cloneStringPointer(value.Email),
		HiddenContacts: cloneSlice(value.HiddenContacts),
	}
	if value.ContactsText != nil {
		result.ContactsText = &resumecontracts.PublicContactsTextView{Title: cloneStringPointer(value.ContactsText.Title), Subtitle: cloneStringPointer(value.ContactsText.Subtitle)}
	}
	redactPublicHiddenContacts(result)
	return result
}

func toPublicSkills(value SkillsView) *resumecontracts.PublicSkillsView {
	return &resumecontracts.PublicSkillsView{ID: value.ID, HardSkills: cloneSlice(value.HardSkills), SoftSkills: cloneSlice(value.SoftSkills)}
}

func redactPublicHiddenContacts(contacts *resumecontracts.PublicContactsView) {
	if contacts == nil {
		return
	}
	for _, key := range contacts.HiddenContacts {
		switch contactkey.Key(key) {
		case contactkey.Phone:
			contacts.Phone = nil
		case contactkey.Telegram:
			contacts.Telegram = nil
		case contactkey.Email:
			contacts.Email = nil
		}
	}
}

func cloneSlice[T any](values []T) []T {
	result := make([]T, len(values))
	copy(result, values)
	return result
}

func cloneStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}

func cloneIntPointer(value *int) *int {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}
