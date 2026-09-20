package contracts

import (
	"context"
	"time"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicResumeQuery struct {
	ID sharedvalueobjects.ID
}

type PublicResumeView struct {
	Title        string
	Language     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Overview     PublicOverviewView
	Profile      *PublicProfileView
	Contacts     *PublicContactsView
	Skills       *PublicSkillsView
	Languages    []PublicLanguageView
	Technologies []PublicTechnologyView
	Experiences  []PublicExperienceView
	Educations   []PublicEducationView
	Courses      []PublicCourseView
	Activities   []PublicActivityView
	Projects     []PublicProjectView
}

type PublicOverviewView struct {
	Theme        *PublicThemeSettings
	Greeting     PublicGreeting
	Availability *PublicAvailability
	Sections     []PublicSectionSettings
}

type PublicGreeting struct {
	First           string
	Second          string
	SecondHighlight string
	Third           string
}

type PublicAvailability struct {
	Status string
	From   *string
	Note   *string
}

type PublicThemeSettings struct {
	Primary   string
	Secondary string
}

type PublicSectionSettings struct {
	ID        string
	Enabled   bool
	Collapsed bool
	Style     *string
}

type PublicProfileView struct {
	ID          string
	Name        string
	Surname     string
	Patronymic  *string
	Position    string
	PhotoURL    string
	DateOfBirth string
	Location    string
	Relocation  bool
	WorkFormat  []string
}

type PublicContactsView struct {
	ID             string
	Phone          *string
	Telegram       *string
	Email          *string
	HiddenContacts []string
	ContactsText   *PublicContactsTextView
}

type PublicContactsTextView struct {
	Title    *string
	Subtitle *string
}

type PublicSkillsView struct {
	ID         string
	HardSkills []string
	SoftSkills []string
}

type PublicLanguageView struct {
	ID       string
	Language string
	Grade    *string
	Position int
}

type PublicTechnologyView struct {
	ID       string
	Name     string
	Position int
}

type PublicExperienceStackItemView struct {
	Name string
}

type PublicExperienceView struct {
	ID           string
	Company      string
	Position     string
	Location     *string
	StartDate    string
	EndDate      *string
	Description  *string
	Achievements []string
	Tasks        []string
	Stack        []PublicExperienceStackItemView
}

type PublicEducationView struct {
	ID          string
	Location    string
	University  string
	Grade       int
	ProgramName string
	ProfileName *string
	Faculty     string
	Course      *int
	EndDate     *string
}

type PublicCourseView struct {
	ID              string
	Name            string
	Organization    string
	EndDate         string
	Hours           *int
	CertificateLink *string
}

type PublicActivityView struct {
	ID          string
	Type        string
	Title       string
	Description *string
	Date        *string
	Link        *string
}

type PublicProjectView struct {
	ID           string
	Name         string
	Description  *string
	Screenshots  []string
	Link         *string
	Stack        []PublicProjectTechnologyView
	Contributors []PublicProjectContributorView
}

type PublicProjectTechnologyView struct {
	Name string
}

type PublicProjectContributorView struct {
	Name     string
	Surname  string
	Role     string
	Contacts PublicProjectContactsView
}

type PublicProjectContactsView struct {
	Phone    *string
	Telegram *string
	Email    *string
}

type PublicResumeQueryHandler interface {
	Handle(context.Context, PublicResumeQuery) (PublicResumeView, error)
}
