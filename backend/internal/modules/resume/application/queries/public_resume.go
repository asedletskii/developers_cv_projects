package queries

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicResumeQuery struct {
	ID sharedvalueobjects.ID
}

type PublicResumeView struct {
	Title        string              `json:"title"`
	Language     string              `json:"language"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
	Overview     ResumeOverviewView  `json:"overview"`
	Profile      *ResumeProfileView  `json:"profile,omitempty"`
	Contacts     *ResumeContactsView `json:"contacts,omitempty"`
	Skills       *SkillsView         `json:"skills,omitempty"`
	Languages    []LanguageView      `json:"languages"`
	Technologies []TechnologyView    `json:"technologies"`
	Experiences  []ExperienceView    `json:"experiences"`
	Educations   []EducationView     `json:"educations"`
	Courses      []CourseView        `json:"courses"`
	Activities   []ActivityView      `json:"activities"`
	Projects     []ProjectView       `json:"projects"`
}

type PublicResumeReader interface {
	GetPublic(context.Context, PublicResumeQuery) (ResumeView, error)
}

type PublicResumeQueryHandler struct {
	reader PublicResumeReader
}

func NewPublicResumeQueryHandler(reader PublicResumeReader) *PublicResumeQueryHandler {
	return &PublicResumeQueryHandler{reader: reader}
}

func (h *PublicResumeQueryHandler) Handle(ctx context.Context, query PublicResumeQuery) (PublicResumeView, error) {
	if h == nil || h.reader == nil {
		return PublicResumeView{}, fmt.Errorf("public resume query handler dependencies are nil")
	}
	if query.ID.IsZero() {
		return PublicResumeView{}, fmt.Errorf("invalid resume id")
	}

	view, err := h.reader.GetPublic(ctx, query)
	if err != nil {
		return PublicResumeView{}, err
	}

	return toPublic(view), nil
}

func toPublic(view ResumeView) PublicResumeView {
	result := PublicResumeView{
		Title:        view.Title,
		Language:     view.Language,
		CreatedAt:    view.CreatedAt,
		UpdatedAt:    view.UpdatedAt,
		Overview:     cloneOverview(view.Overview),
		Languages:    cloneSlice(view.Languages),
		Technologies: cloneSlice(view.Technologies),
		Experiences:  cloneSlice(view.Experiences),
		Educations:   cloneSlice(view.Educations),
		Courses:      cloneSlice(view.Courses),
		Activities:   cloneSlice(view.Activities),
		Projects:     make([]ProjectView, len(view.Projects)),
	}

	if view.Profile != nil {
		profile := *view.Profile
		profile.WorkFormat = cloneSlice(view.Profile.WorkFormat)
		result.Profile = &profile
	}
	if view.Skills != nil {
		skills := *view.Skills
		skills.HardSkills = cloneSlice(view.Skills.HardSkills)
		skills.SoftSkills = cloneSlice(view.Skills.SoftSkills)
		result.Skills = &skills
	}
	if view.Contacts != nil {
		contacts := *view.Contacts
		contacts.HiddenContacts = cloneSlice(view.Contacts.HiddenContacts)
		contacts.ContactsText = clonePublicContactsText(view.Contacts.ContactsText)
		contacts.Phone = cloneStringPointer(view.Contacts.Phone)
		contacts.Telegram = cloneStringPointer(view.Contacts.Telegram)
		contacts.Email = cloneStringPointer(view.Contacts.Email)
		redactHiddenContacts(&contacts)
		result.Contacts = &contacts
	}

	for index, project := range view.Projects {
		result.Projects[index] = project
		result.Projects[index].Role = nil
		result.Projects[index].GitLink = nil
		result.Projects[index].Screenshots = make([]string, 0)
		result.Projects[index].Contributors = make([]ProjectContributorView, 0)
		result.Projects[index].Stack = cloneSlice(project.Stack)
	}

	return result
}

func cloneOverview(value ResumeOverviewView) ResumeOverviewView {
	result := value
	result.Sections = cloneSlice(value.Sections)
	return result
}

func cloneSlice[T any](values []T) []T {
	result := make([]T, len(values))
	copy(result, values)
	return result
}

func redactHiddenContacts(contacts *ResumeContactsView) {
	if contacts == nil {
		return
	}
	for _, key := range contacts.HiddenContacts {
		switch domain.ContactKey(key) {
		case domain.ContactKeyPhone:
			contacts.Phone = nil
		case domain.ContactKeyTelegram:
			contacts.Telegram = nil
		case domain.ContactKeyEmail:
			contacts.Email = nil
		}
	}
}

func clonePublicContactsText(value *ResumeContactsTextView) *ResumeContactsTextView {
	if value == nil {
		return nil
	}
	result := *value
	result.Title = cloneStringPointer(value.Title)
	result.Subtitle = cloneStringPointer(value.Subtitle)
	return &result
}

func cloneStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}
