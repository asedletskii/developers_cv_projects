package queries

import "time"

type PublicResumeResult struct {
	Title        string
	Language     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Overview     PublicOverviewResult
	Profile      *PublicProfileResult
	Contacts     *PublicContactsResult
	Skills       *PublicSkillsResult
	Languages    []PublicLanguageResult
	Technologies []PublicTechnologyResult
	Experiences  []PublicExperienceResult
	Educations   []PublicEducationResult
	Courses      []PublicCourseResult
	Activities   []PublicActivityResult
	Projects     []PublicProjectResult
}

type PublicOverviewResult struct {
	Theme        *PublicThemeSettingsResult
	Greeting     PublicGreetingResult
	Availability *PublicAvailabilityResult
	Sections     []PublicSectionSettingsResult
}

type PublicGreetingResult struct {
	First           string
	Second          string
	SecondHighlight string
	Third           string
}

type PublicAvailabilityResult struct {
	Status string
	From   *string
	Note   *string
}

type PublicThemeSettingsResult struct {
	Primary   string
	Secondary string
}

type PublicSectionSettingsResult struct {
	ID        string
	Enabled   bool
	Collapsed bool
	Style     *string
}

type PublicProfileResult struct {
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

type PublicContactsResult struct {
	ID             string
	Phone          *string
	Telegram       *string
	Email          *string
	HiddenContacts []string
	ContactsText   *PublicContactsTextResult
}

type PublicContactsTextResult struct {
	Title    *string
	Subtitle *string
}

type PublicSkillsResult struct {
	ID         string
	HardSkills []string
	SoftSkills []string
}

type PublicLanguageResult struct {
	ID       string
	Language string
	Grade    *string
	Position int
}

type PublicTechnologyResult struct {
	ID       string
	Name     string
	Position int
}

type PublicExperienceStackItemResult struct {
	Name string
}

type PublicExperienceResult struct {
	ID           string
	Company      string
	Position     string
	Location     *string
	StartDate    string
	EndDate      *string
	Description  *string
	Achievements []string
	Tasks        []string
	Stack        []PublicExperienceStackItemResult
}

type PublicEducationResult struct {
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

type PublicCourseResult struct {
	ID              string
	Name            string
	Organization    string
	EndDate         string
	Hours           *int
	CertificateLink *string
}

type PublicActivityResult struct {
	ID          string
	Type        string
	Title       string
	Description *string
	Date        *string
	Link        *string
}

type PublicProjectResult struct {
	ID           string
	Name         string
	Description  *string
	Screenshots  []string
	Link         *string
	Stack        []PublicProjectTechnologyResult
	Contributors []PublicProjectContributorResult
}

type PublicProjectTechnologyResult struct {
	Name string
}

type PublicProjectContributorResult struct {
	Name     string
	Surname  string
	Role     string
	Contacts PublicProjectContactsResult
}

type PublicProjectContactsResult struct {
	Phone    *string
	Telegram *string
	Email    *string
}
