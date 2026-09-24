package request

type ResumeRequest struct {
	Name           string               `json:"name"`
	Surname        string               `json:"surname"`
	Patronymic     *string              `json:"patronymic,omitempty"`
	Position       string               `json:"position"`
	Availability   *AvailabilityRequest `json:"availability,omitempty"`
	ResumeLink     *string              `json:"resumeLink,omitempty"`
	Theme          *ThemeRequest        `json:"theme,omitempty"`
	Sections       []SectionRequest     `json:"sections"`
	HiddenContacts []string             `json:"hiddenContacts"`
	ContactsText   *ContactsTextRequest `json:"contactsText,omitempty"`
	Greeting       GreetingRequest      `json:"greeting"`
	PhotoPath      string               `json:"photoPath"`
	DateOfBirth    string               `json:"dateOfBirth"`
	Location       string               `json:"location"`
	Relocation     bool                 `json:"relocation"`
	WorkFormat     []string             `json:"workFormat"`
	Contacts       ContactsRequest      `json:"contacts"`
	Languages      []LanguageRequest    `json:"languages"`
	Experience     []ExperienceRequest  `json:"experience"`
	Education      []EducationRequest   `json:"education"`
	Courses        []CourseRequest      `json:"courses"`
	Activities     []ActivityRequest    `json:"activities"`
	HardSkills     []string             `json:"hardSkills"`
	SoftSkills     []string             `json:"softSkills"`
	Technologies   []TechnologyRequest  `json:"technologies"`
	Projects       []ProjectRequest     `json:"projects"`
}

type LocalizedTextRequest struct {
	RU string `json:"ru"`
	EN string `json:"en"`
}
type GreetingRequest struct {
	First           LocalizedTextRequest `json:"first"`
	Second          LocalizedTextRequest `json:"second"`
	SecondHighlight LocalizedTextRequest `json:"secondHighlight"`
	Third           LocalizedTextRequest `json:"third"`
}
type AvailabilityRequest struct {
	Status string                `json:"status"`
	From   *string               `json:"from,omitempty"`
	Note   *LocalizedTextRequest `json:"note,omitempty"`
}
type ThemeRequest struct {
	Primary   *string `json:"primary,omitempty"`
	Secondary *string `json:"secondary,omitempty"`
}
type SectionRequest struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}
type ContactsTextRequest struct {
	Title    *LocalizedTextRequest `json:"title,omitempty"`
	Subtitle *LocalizedTextRequest `json:"subtitle,omitempty"`
}
type ContactsRequest struct {
	Phone    *string `json:"phone,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}
type TechnologyRequest struct {
	Name string `json:"name"`
}
type LanguageRequest struct {
	Language string  `json:"language"`
	Grade    *string `json:"grade,omitempty"`
}
type ExperienceRequest struct {
	ID           string              `json:"id,omitempty"`
	Company      string              `json:"company"`
	Position     string              `json:"position"`
	Location     *string             `json:"location,omitempty"`
	StartDate    string              `json:"startDate"`
	EndDate      *string             `json:"endDate,omitempty"`
	Description  *string             `json:"description,omitempty"`
	Achievements []string            `json:"achievements"`
	Stack        []TechnologyRequest `json:"stack"`
}
type EducationRequest struct {
	ID          string  `json:"id,omitempty"`
	Location    string  `json:"location"`
	University  string  `json:"university"`
	Grade       int     `json:"grade"`
	ProgramName string  `json:"programName"`
	ProfileName *string `json:"profileName,omitempty"`
	Faculty     string  `json:"faculty"`
	Course      *int    `json:"course,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
}
type CourseRequest struct {
	ID              string  `json:"id,omitempty"`
	Name            string  `json:"name"`
	Organization    string  `json:"organization"`
	EndDate         string  `json:"endDate"`
	Hours           *int    `json:"hours,omitempty"`
	CertificateLink *string `json:"certificateLink,omitempty"`
}
type ActivityRequest struct {
	ID          string  `json:"id,omitempty"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Date        *string `json:"date,omitempty"`
	Link        *string `json:"link,omitempty"`
}
type ContributorRequest struct {
	Name     string          `json:"name"`
	Surname  string          `json:"surname"`
	Role     string          `json:"role"`
	Contacts ContactsRequest `json:"contacts"`
}
type ProjectRequest struct {
	ID           string               `json:"id,omitempty"`
	Name         string               `json:"name"`
	Description  *string              `json:"description,omitempty"`
	Role         *string              `json:"role,omitempty"`
	Screenshots  []string             `json:"screenshots"`
	Link         *string              `json:"link,omitempty"`
	GitLink      *string              `json:"gitLink,omitempty"`
	Stack        []TechnologyRequest  `json:"stack"`
	Contributors []ContributorRequest `json:"contributors"`
}
