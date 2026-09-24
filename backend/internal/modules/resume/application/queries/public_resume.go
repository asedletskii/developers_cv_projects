package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PublicResumeView struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Surname        string            `json:"surname"`
	Patronymic     *string           `json:"patronymic,omitempty"`
	Position       string            `json:"position"`
	Title          string            `json:"title"`
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
	ID uuid.UUID
}

type PublicResumeQueryHandler struct {
	reader PublicResumeReader
}

type PublicResumeReader interface {
	GetPublicResume(ctx context.Context, query PublicResumeQuery) (PublicResumeView, error)
}
