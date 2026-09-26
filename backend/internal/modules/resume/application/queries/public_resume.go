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

type PublicResumeReader interface {
	GetPublicResume(ctx context.Context, query PublicResumeQuery) (PublicResumeView, error)
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
	return h.reader.GetPublicResume(ctx, query)
}
