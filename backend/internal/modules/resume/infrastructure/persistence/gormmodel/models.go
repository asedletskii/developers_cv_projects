package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	Surname     string    `gorm:"not null"`
	Patronymic  *string
	Position    string `gorm:"not null"`
	ResumeLink  *string
	PhotoPath   string    `gorm:"not null"`
	DateOfBirth string    `gorm:"type:date;not null"`
	Location    string    `gorm:"not null"`
	Relocation  bool      `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

type Greeting struct {
	PersonID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstRU           string    `gorm:"not null"`
	FirstEN           string    `gorm:"not null"`
	SecondRU          string    `gorm:"not null"`
	SecondEN          string    `gorm:"not null"`
	SecondHighlightRU string    `gorm:"not null"`
	SecondHighlightEN string    `gorm:"not null"`
	ThirdRU           string    `gorm:"not null"`
	ThirdEN           string    `gorm:"not null"`
}

type Availability struct {
	PersonID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Status   string    `gorm:"not null"`
	From     *string   `gorm:"type:date"`
	NoteRU   *string
	NoteEN   *string
}

type Theme struct {
	PersonID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	Primary   *string
	Secondary *string
}

type ContactsText struct {
	PersonID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	TitleRU    *string
	TitleEN    *string
	SubtitleRU *string
	SubtitleEN *string
}

type Section struct {
	PersonID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	SectionID string    `gorm:"type:text;primaryKey"`
	Position  int       `gorm:"not null"`
	Enabled   bool      `gorm:"not null"`
	Collapsed bool      `gorm:"not null"`
}

type Contact struct {
	PersonID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Phone    *string
	Telegram *string
	Email    *string
}

type PersonWorkFormat struct {
	PersonID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type PersonHiddenContact struct {
	PersonID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type PersonHardSkill struct {
	PersonID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type PersonSoftSkill struct {
	PersonID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type PersonLanguage struct {
	PersonID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Language string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
	Grade    *string
}
type PersonTechnology struct {
	PersonID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
	Grade    *string
}

type Experience struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	PersonID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Position    int       `gorm:"not null"`
	Company     string    `gorm:"not null"`
	JobPosition string    `gorm:"not null"`
	Location    *string
	StartDate   string  `gorm:"type:date;not null"`
	EndDate     *string `gorm:"type:date"`
	Description *string
}
type ExperienceAchievement struct {
	ExperienceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position     int       `gorm:"primaryKey"`
	Value        string    `gorm:"not null"`
}
type ExperienceTechnology struct {
	ExperienceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position     int       `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	Grade        *string
}

type Education struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	PersonID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Position    int       `gorm:"not null"`
	Location    string    `gorm:"not null"`
	University  string    `gorm:"not null"`
	Grade       int       `gorm:"not null"`
	ProgramName string    `gorm:"not null"`
	ProfileName *string
	Faculty     string `gorm:"not null"`
	Course      *int
	EndDate     *string `gorm:"type:date"`
}

type Course struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	PersonID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Position        int       `gorm:"not null"`
	Name            string    `gorm:"not null"`
	Organization    string    `gorm:"not null"`
	EndDate         string    `gorm:"type:date;not null"`
	Hours           *int
	CertificateLink *string
}

type Activity struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	PersonID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Position    int       `gorm:"not null"`
	Type        string    `gorm:"not null"`
	Title       string    `gorm:"not null"`
	Description *string
	Date        *string `gorm:"type:date"`
	Link        *string
}

type Project struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	PersonID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Position    int       `gorm:"not null"`
	Name        string    `gorm:"not null"`
	Description *string
	Role        *string
	Link        *string
	GitLink     *string
}
type ProjectScreenshot struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position  int       `gorm:"primaryKey"`
	Path      string    `gorm:"not null"`
}
type ProjectTechnology struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position  int       `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Grade     *string
}
type ProjectContributor struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	Position  int       `gorm:"not null"`
	Name      string    `gorm:"not null"`
	Surname   string    `gorm:"not null"`
	Role      string    `gorm:"not null"`
	Phone     *string
	Telegram  *string
	Email     *string
}
