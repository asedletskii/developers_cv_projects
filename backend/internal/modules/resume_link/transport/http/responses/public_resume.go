package responses

import "time"

type PublicResumeResponse struct {
	Title        string                     `json:"title"`
	Language     string                     `json:"language"`
	CreatedAt    time.Time                  `json:"createdAt"`
	UpdatedAt    time.Time                  `json:"updatedAt"`
	Overview     PublicOverviewResponse     `json:"overview"`
	Profile      *PublicProfileResponse     `json:"profile,omitempty"`
	Contacts     *PublicContactsResponse    `json:"contacts,omitempty"`
	Skills       *PublicSkillsResponse      `json:"skills,omitempty"`
	Languages    []PublicLanguageResponse   `json:"languages"`
	Technologies []PublicTechnologyResponse `json:"technologies"`
	Experiences  []PublicExperienceResponse `json:"experiences"`
	Educations   []PublicEducationResponse  `json:"educations"`
	Courses      []PublicCourseResponse     `json:"courses"`
	Activities   []PublicActivityResponse   `json:"activities"`
	Projects     []PublicProjectResponse    `json:"projects"`
}

type PublicOverviewResponse struct {
	Theme        *PublicThemeSettingsResponse    `json:"theme,omitempty"`
	Greeting     PublicGreetingResponse          `json:"greeting"`
	Availability *PublicAvailabilityResponse     `json:"availability,omitempty"`
	Sections     []PublicSectionSettingsResponse `json:"sections"`
}

type PublicGreetingResponse struct {
	First           string `json:"first"`
	Second          string `json:"second"`
	SecondHighlight string `json:"secondHighlight"`
	Third           string `json:"third"`
}

type PublicAvailabilityResponse struct {
	Status string  `json:"status"`
	From   *string `json:"from,omitempty"`
	Note   *string `json:"note,omitempty"`
}

type PublicThemeSettingsResponse struct {
	Primary   string `json:"primary,omitempty"`
	Secondary string `json:"secondary,omitempty"`
}

type PublicSectionSettingsResponse struct {
	ID        string  `json:"id"`
	Enabled   bool    `json:"enabled"`
	Collapsed bool    `json:"collapsed"`
	Style     *string `json:"style,omitempty"`
}

type PublicProfileResponse struct {
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

type PublicContactsResponse struct {
	ID             string                      `json:"id"`
	Phone          *string                     `json:"phone,omitempty"`
	Telegram       *string                     `json:"telegram,omitempty"`
	Email          *string                     `json:"email,omitempty"`
	HiddenContacts []string                    `json:"hiddenContacts"`
	ContactsText   *PublicContactsTextResponse `json:"contactsText,omitempty"`
}

type PublicContactsTextResponse struct {
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
}

type PublicSkillsResponse struct {
	ID         string   `json:"id"`
	HardSkills []string `json:"hardSkills"`
	SoftSkills []string `json:"softSkills"`
}

type PublicLanguageResponse struct {
	ID       string  `json:"id"`
	Language string  `json:"language"`
	Grade    *string `json:"grade,omitempty"`
	Position int     `json:"position"`
}

type PublicTechnologyResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

type PublicExperienceStackItemResponse struct {
	Name string `json:"name"`
}

type PublicExperienceResponse struct {
	ID           string                              `json:"id"`
	Company      string                              `json:"company"`
	Position     string                              `json:"position"`
	Location     *string                             `json:"location,omitempty"`
	StartDate    string                              `json:"startDate"`
	EndDate      *string                             `json:"endDate,omitempty"`
	Description  *string                             `json:"description,omitempty"`
	Achievements []string                            `json:"achievements"`
	Tasks        []string                            `json:"tasks"`
	Stack        []PublicExperienceStackItemResponse `json:"stack"`
}

type PublicEducationResponse struct {
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

type PublicCourseResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Organization    string  `json:"organization"`
	EndDate         string  `json:"endDate"`
	Hours           *int    `json:"hours,omitempty"`
	CertificateLink *string `json:"certificateLink,omitempty"`
}

type PublicActivityResponse struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Date        *string `json:"date,omitempty"`
	Link        *string `json:"link,omitempty"`
}

type PublicProjectResponse struct {
	ID           string                             `json:"id"`
	Name         string                             `json:"name"`
	Description  *string                            `json:"description,omitempty"`
	Screenshots  []string                           `json:"screenshots"`
	Link         *string                            `json:"link,omitempty"`
	Stack        []PublicProjectTechnologyResponse  `json:"stack"`
	Contributors []PublicProjectContributorResponse `json:"contributors"`
}

type PublicProjectTechnologyResponse struct {
	Name string `json:"name"`
}

type PublicProjectContributorResponse struct {
	Name     string                        `json:"name"`
	Surname  string                        `json:"surname"`
	Role     string                        `json:"role"`
	Contacts PublicProjectContactsResponse `json:"contacts"`
}

type PublicProjectContactsResponse struct {
	Phone    *string `json:"phone,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}
