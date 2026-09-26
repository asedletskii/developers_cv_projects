package queries

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type PublicResumeView struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Surname        string            `json:"surname"`
	Patronymic     *string           `json:"patronymic,omitempty"`
	Position       string            `json:"position"`
	Title          string            `json:"title"`
	Language       string            `json:"language"`
	Availability   *Availability     `json:"availability,omitempty"`
	UpdatedAt      *time.Time        `json:"updatedAt,omitempty"`
	ResumeLink     *string           `json:"resumeLink,omitempty"`
	Theme          *ThemeSettings    `json:"theme,omitempty"`
	Sections       []SectionSettings `json:"sections"`
	HiddenContacts []string          `json:"hiddenContacts"`
	ContactsText   *ContactsText     `json:"contactsText,omitempty"`
	Greeting       Greeting          `json:"greeting"`
	PhotoPath      string            `json:"photoPath"`
	DateOfBirth    string            `json:"dateOfBirth"`
	Location       string            `json:"location"`
	Relocation     bool              `json:"relocation"`
	WorkFormat     []string          `json:"workFormat"`
	Contacts       Contacts          `json:"contacts"`
	Languages      []ResumeLanguage  `json:"languages"`
	Experience     []Experience      `json:"experience"`
	Education      []Education       `json:"education"`
	Courses        []Course          `json:"courses"`
	Activities     []Activity        `json:"activities"`
	HardSkills     []string          `json:"hardSkills"`
	SoftSkills     []string          `json:"softSkills"`
	Technologies   []Technology      `json:"technologies"`
	Projects       []ProjectSummary  `json:"projects"`
}

type PublicResumeQuery struct {
	ID valueobjects.ID
}

type PublicResumeQueryHandler struct {
	reader PublicResumeReader
}

type Contacts struct {
	Phone    *string `json:"phone,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type ContactsText struct {
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
}

type Greeting struct {
	First           string `json:"first"`
	Second          string `json:"second"`
	SecondHighlight string `json:"secondHighlight"`
	Third           string `json:"third"`
}

type Availability struct {
	Status string  `json:"status"`
	From   *string `json:"from,omitempty"`
	Note   *string `json:"note,omitempty"`
}

type ThemeSettings struct {
	Primary   string `json:"primary,omitempty"`
	Secondary string `json:"secondary,omitempty"`
}

type SectionSettings struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed,omitempty"`
}

type ResumeLanguage struct {
	Language string  `json:"language"`
	Grade    *string `json:"grade,omitempty"`
}

type Technology struct {
	Name string `json:"name"`
}

type Experience struct {
	ID           string       `json:"id"`
	Company      string       `json:"company"`
	Position     string       `json:"position"`
	Location     *string      `json:"location,omitempty"`
	StartDate    string       `json:"startDate"`
	EndDate      *string      `json:"endDate,omitempty"`
	Description  *string      `json:"description,omitempty"`
	Achievements []string     `json:"achievements"`
	Stack        []Technology `json:"stack"`
}

type Education struct {
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

type Course struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Organization    string  `json:"organization"`
	EndDate         string  `json:"endDate"`
	Hours           *int    `json:"hours,omitempty"`
	CertificateLink *string `json:"certificateLink,omitempty"`
}

type Activity struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Date        *string `json:"date,omitempty"`
	Link        *string `json:"link,omitempty"`
}

type ProjectSummary struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	Stack       []Technology `json:"stack"`
	Link        *string      `json:"link,omitempty"`
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

	view, err := h.reader.GetPublic(ctx, query.ID)
	if err != nil {
		return PublicResumeView{}, err
	}
	return projectPublicResume(view), nil
}

func projectPublicResume(view ResumeView) PublicResumeView {
	result := PublicResumeView{
		ID: view.ID, Title: view.Title, Language: view.Language, UpdatedAt: &view.UpdatedAt,
		Sections:       make([]SectionSettings, 0, len(view.Overview.Sections)),
		HiddenContacts: make([]string, 0), Languages: make([]ResumeLanguage, 0, len(view.Languages)),
		Experience: make([]Experience, 0, len(view.Experiences)), Education: make([]Education, 0, len(view.Educations)),
		Courses: make([]Course, 0, len(view.Courses)), Activities: make([]Activity, 0, len(view.Activities)),
		HardSkills: make([]string, 0), SoftSkills: make([]string, 0), Technologies: make([]Technology, 0, len(view.Technologies)),
		Projects: make([]ProjectSummary, 0, len(view.Projects)),
	}
	for _, section := range view.Overview.Sections {
		result.Sections = append(result.Sections, SectionSettings(section))
	}
	result.Greeting = Greeting{First: view.Overview.Greeting.First, Second: view.Overview.Greeting.Second, SecondHighlight: view.Overview.Greeting.SecondHighlight, Third: view.Overview.Greeting.Third}
	result.ResumeLink = view.Overview.ResumeLink
	if availability := view.Overview.Availability; availability != nil {
		result.Availability = &Availability{Status: availability.Status, From: availability.From, Note: availability.Note}
	}
	if theme := view.Overview.Theme; theme != nil {
		result.Theme = &ThemeSettings{Primary: theme.Primary, Secondary: theme.Secondary}
	}
	if profile := view.Profile; profile != nil {
		result.Name, result.Surname, result.Patronymic = profile.Name, profile.Surname, profile.Patronymic
		result.Position, result.PhotoPath, result.DateOfBirth = profile.Position, profile.PhotoPath, profile.DateOfBirth
		result.Location, result.Relocation = profile.Location, profile.Relocation
		result.WorkFormat = append(result.WorkFormat[:0], profile.WorkFormat...)
	} else {
		result.WorkFormat = make([]string, 0)
	}
	if contacts := view.Contacts; contacts != nil {
		result.Contacts = Contacts{Phone: contacts.Phone, Telegram: contacts.Telegram, Email: contacts.Email}
		result.HiddenContacts = append(result.HiddenContacts, contacts.HiddenContacts...)
		if text := contacts.ContactsText; text != nil {
			result.ContactsText = &ContactsText{Title: text.Title, Subtitle: text.Subtitle}
		}
	}
	if skills := view.Skills; skills != nil {
		result.HardSkills = append(result.HardSkills, skills.HardSkills...)
		result.SoftSkills = append(result.SoftSkills, skills.SoftSkills...)
	}
	for _, language := range view.Languages {
		result.Languages = append(result.Languages, ResumeLanguage{Language: language.Language, Grade: language.Grade})
	}
	for _, item := range view.Experiences {
		stack := make([]Technology, 0, len(item.Technologies))
		for _, tech := range item.Technologies {
			stack = append(stack, Technology{Name: tech.Name})
		}
		result.Experience = append(result.Experience, Experience{ID: item.ID, Company: item.Company, Position: item.Position, Location: item.Location, StartDate: item.StartDate, EndDate: item.EndDate, Description: item.Description, Achievements: append([]string(nil), item.Achievements...), Stack: stack})
	}
	for _, item := range view.Educations {
		result.Education = append(result.Education, Education{ID: item.ID, Location: item.Location, University: item.University, Grade: item.Grade, ProgramName: item.ProgramName, ProfileName: item.ProfileName, Faculty: item.Faculty, Course: item.Course, EndDate: item.EndDate})
	}
	for _, item := range view.Courses {
		result.Courses = append(result.Courses, Course{ID: item.ID, Name: item.Name, Organization: item.Organization, EndDate: item.EndDate, Hours: item.Hours, CertificateLink: item.CertificateLink})
	}
	for _, item := range view.Activities {
		result.Activities = append(result.Activities, Activity{ID: item.ID, Type: item.Type, Title: item.Title, Description: item.Description, Date: item.Date, Link: item.Link})
	}
	for _, item := range view.Technologies {
		result.Technologies = append(result.Technologies, Technology{Name: item.Name})
	}
	for _, item := range view.Projects {
		stack := make([]Technology, 0, len(item.Stack))
		for _, tech := range item.Stack {
			stack = append(stack, Technology{Name: tech.Name})
		}
		result.Projects = append(result.Projects, ProjectSummary{ID: item.ID, Name: item.Name, Description: item.Description, Stack: stack, Link: item.Link})
	}
	return result
}
