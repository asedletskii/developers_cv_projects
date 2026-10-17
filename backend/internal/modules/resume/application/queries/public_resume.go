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
	Title        string                 `json:"title"`
	Language     string                 `json:"language"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
	Overview     PublicOverviewView     `json:"overview"`
	Profile      *PublicProfileView     `json:"profile,omitempty"`
	Contacts     *PublicContactsView    `json:"contacts,omitempty"`
	Skills       *PublicSkillsView      `json:"skills,omitempty"`
	Languages    []PublicLanguageView   `json:"languages"`
	Technologies []PublicTechnologyView `json:"technologies"`
	Experiences  []PublicExperienceView `json:"experiences"`
	Educations   []PublicEducationView  `json:"educations"`
	Courses      []PublicCourseView     `json:"courses"`
	Activities   []PublicActivityView   `json:"activities"`
	Projects     []PublicProjectView    `json:"projects"`
}

type PublicOverviewView struct {
	Theme        *PublicThemeSettings    `json:"theme,omitempty"`
	Greeting     PublicGreeting          `json:"greeting"`
	Availability *PublicAvailability     `json:"availability,omitempty"`
	ResumeLink   *string                 `json:"resumeLink,omitempty"`
	Sections     []PublicSectionSettings `json:"sections"`
}

type PublicGreeting struct {
	First           string `json:"first"`
	Second          string `json:"second"`
	SecondHighlight string `json:"secondHighlight"`
	Third           string `json:"third"`
}

type PublicAvailability struct {
	Status string  `json:"status"`
	From   *string `json:"from,omitempty"`
	Note   *string `json:"note,omitempty"`
}

type PublicThemeSettings struct {
	Primary   string `json:"primary,omitempty"`
	Secondary string `json:"secondary,omitempty"`
}

type PublicSectionSettings struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}

type PublicProfileView struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Patronymic  *string  `json:"patronymic,omitempty"`
	Position    string   `json:"position"`
	PhotoURL    string   `json:"photoPath"`
	DateOfBirth string   `json:"dateOfBirth"`
	Location    string   `json:"location"`
	Relocation  bool     `json:"relocation"`
	WorkFormat  []string `json:"workFormat"`
}

type PublicContactsView struct {
	ID             string                  `json:"id"`
	Phone          *string                 `json:"phone,omitempty"`
	Telegram       *string                 `json:"telegram,omitempty"`
	Email          *string                 `json:"email,omitempty"`
	HiddenContacts []string                `json:"hiddenContacts"`
	ContactsText   *PublicContactsTextView `json:"contactsText,omitempty"`
}

type PublicContactsTextView struct {
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
}

type PublicSkillsView struct {
	ID         string   `json:"id"`
	HardSkills []string `json:"hardSkills"`
	SoftSkills []string `json:"softSkills"`
}

type PublicLanguageView struct {
	ID       string  `json:"id"`
	Language string  `json:"language"`
	Grade    *string `json:"grade,omitempty"`
	Position int     `json:"position"`
}

type PublicTechnologyView struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

type PublicExperienceStackItemView struct {
	Name string `json:"name"`
}

type PublicExperienceView struct {
	ID           string                          `json:"id"`
	Company      string                          `json:"company"`
	Position     string                          `json:"position"`
	Location     *string                         `json:"location,omitempty"`
	StartDate    string                          `json:"startDate"`
	EndDate      *string                         `json:"endDate,omitempty"`
	Description  *string                         `json:"description,omitempty"`
	Achievements []string                        `json:"achievements"`
	Tasks        []string                        `json:"tasks"`
	Stack        []PublicExperienceStackItemView `json:"stack"`
}

type PublicEducationView struct {
	ID          string  `json:"id"`
	Location    string  `json:"location"`
	University  string  `json:"university"`
	Grade       int     `json:"grade"`
	ProgramName string  `json:"programName"`
	ProfileName *string `json:"profileName,omitempty"`
	Faculty     string  `json:"faculty"`
	Course      *int    `json:"course,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
}

type PublicCourseView struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Organization    string  `json:"organization"`
	EndDate         string  `json:"endDate"`
	Hours           *int    `json:"hours,omitempty"`
	CertificateLink *string `json:"certificateLink,omitempty"`
}

type PublicActivityView struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Date        *string `json:"date,omitempty"`
	Link        *string `json:"link,omitempty"`
}

type PublicProjectView struct {
	ID           string                         `json:"id"`
	Name         string                         `json:"name"`
	Description  *string                        `json:"description,omitempty"`
	Screenshots  []string                       `json:"screenshots"`
	Link         *string                        `json:"link,omitempty"`
	Stack        []PublicProjectTechnologyView  `json:"stack"`
	Contributors []PublicProjectContributorView `json:"contributors"`
}

type PublicProjectTechnologyView struct {
	Name string `json:"name"`
}

type PublicProjectContributorView struct {
	Name     string                    `json:"name"`
	Surname  string                    `json:"surname"`
	Role     string                    `json:"role"`
	Contacts PublicProjectContactsView `json:"contacts"`
}

type PublicProjectContactsView struct {
	Phone    *string `json:"phone,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type PublicResumeReader interface {
	GetPublic(context.Context, PublicResumeQuery) (PublicResumeView, error)
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
	return h.reader.GetPublic(ctx, query)
}

func ToPublicResumeView(view ResumeView) PublicResumeView {
	result := PublicResumeView{
		Title:        view.Title,
		Language:     view.Language,
		CreatedAt:    view.CreatedAt,
		UpdatedAt:    view.UpdatedAt,
		Overview:     toPublicOverview(view.Overview),
		Languages:    make([]PublicLanguageView, 0, len(view.Languages)),
		Technologies: make([]PublicTechnologyView, 0, len(view.Technologies)),
		Experiences:  make([]PublicExperienceView, 0, len(view.Experiences)),
		Educations:   make([]PublicEducationView, 0, len(view.Educations)),
		Courses:      make([]PublicCourseView, 0, len(view.Courses)),
		Activities:   make([]PublicActivityView, 0, len(view.Activities)),
		Projects:     make([]PublicProjectView, 0, len(view.Projects)),
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
		result.Languages = append(result.Languages, PublicLanguageView{ID: value.ID, Language: value.Language, Grade: cloneStringPointer(value.Grade), Position: value.Position})
	}
	for _, value := range view.Technologies {
		result.Technologies = append(result.Technologies, PublicTechnologyView(value))
	}
	for _, value := range view.Experiences {
		stack := make([]PublicExperienceStackItemView, 0, len(value.Stack))
		for _, item := range value.Stack {
			stack = append(stack, PublicExperienceStackItemView(item))
		}
		result.Experiences = append(result.Experiences, PublicExperienceView{
			ID: value.ID, Company: value.Company, Position: value.Position, Location: cloneStringPointer(value.Location),
			StartDate: value.StartDate, EndDate: cloneStringPointer(value.EndDate), Description: cloneStringPointer(value.Description),
			Achievements: cloneSlice(value.Achievements), Tasks: cloneSlice(value.Tasks), Stack: stack,
		})
	}
	for _, value := range view.Educations {
		result.Educations = append(result.Educations, PublicEducationView{
			ID: value.ID, Location: value.Location, University: value.University, Grade: value.Grade, ProgramName: value.ProgramName,
			ProfileName: cloneStringPointer(value.ProfileName), Faculty: value.Faculty, Course: cloneIntPointer(value.Course), EndDate: cloneStringPointer(value.EndDate),
		})
	}
	for _, value := range view.Courses {
		result.Courses = append(result.Courses, PublicCourseView{
			ID: value.ID, Name: value.Name, Organization: value.Organization, EndDate: value.EndDate,
			Hours: cloneIntPointer(value.Hours), CertificateLink: cloneStringPointer(value.CertificateLink),
		})
	}
	for _, value := range view.Activities {
		result.Activities = append(result.Activities, PublicActivityView{
			ID: value.ID, Type: value.Type, Title: value.Title, Description: cloneStringPointer(value.Description),
			Date: cloneStringPointer(value.Date), Link: cloneStringPointer(value.Link),
		})
	}
	for _, value := range view.Projects {
		stack := make([]PublicProjectTechnologyView, 0, len(value.Stack))
		for _, item := range value.Stack {
			stack = append(stack, PublicProjectTechnologyView(item))
		}
		result.Projects = append(result.Projects, PublicProjectView{
			ID: value.ID, Name: value.Name, Description: cloneStringPointer(value.Description),
			Screenshots: make([]string, 0), Link: cloneStringPointer(value.Link), Stack: stack,
			Contributors: make([]PublicProjectContributorView, 0),
		})
	}
	return result
}

func toPublicOverview(value ResumeOverviewView) PublicOverviewView {
	result := PublicOverviewView{
		Greeting:   PublicGreeting{First: value.Greeting.First, Second: value.Greeting.Second, SecondHighlight: value.Greeting.SecondHighlight, Third: value.Greeting.Third},
		ResumeLink: cloneStringPointer(value.ResumeLink),
		Sections:   make([]PublicSectionSettings, 0, len(value.Sections)),
	}
	if value.Theme != nil {
		result.Theme = &PublicThemeSettings{Primary: value.Theme.Primary, Secondary: value.Theme.Secondary}
	}
	if value.Availability != nil {
		result.Availability = &PublicAvailability{Status: value.Availability.Status, From: cloneStringPointer(value.Availability.From), Note: cloneStringPointer(value.Availability.Note)}
	}
	for _, section := range value.Sections {
		result.Sections = append(result.Sections, PublicSectionSettings(section))
	}
	return result
}

func toPublicProfile(value ResumeProfileView) *PublicProfileView {
	return &PublicProfileView{
		ID: value.ID, Name: value.Name, Surname: value.Surname, Patronymic: cloneStringPointer(value.Patronymic),
		Position: value.Position, PhotoURL: value.PhotoURL, DateOfBirth: value.DateOfBirth, Location: value.Location,
		Relocation: value.Relocation, WorkFormat: cloneSlice(value.WorkFormat),
	}
}

func toPublicContacts(value ResumeContactsView) *PublicContactsView {
	result := &PublicContactsView{
		ID: value.ID, Phone: cloneStringPointer(value.Phone), Telegram: cloneStringPointer(value.Telegram), Email: cloneStringPointer(value.Email),
		HiddenContacts: cloneSlice(value.HiddenContacts),
	}
	if value.ContactsText != nil {
		result.ContactsText = &PublicContactsTextView{Title: cloneStringPointer(value.ContactsText.Title), Subtitle: cloneStringPointer(value.ContactsText.Subtitle)}
	}
	redactPublicHiddenContacts(result)
	return result
}

func toPublicSkills(value SkillsView) *PublicSkillsView {
	return &PublicSkillsView{ID: value.ID, HardSkills: cloneSlice(value.HardSkills), SoftSkills: cloneSlice(value.SoftSkills)}
}

func redactPublicHiddenContacts(contacts *PublicContactsView) {
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
