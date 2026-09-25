package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type ResumeRoot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Title     string    `gorm:"not null"`
	Language  string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (ResumeRoot) TableName() string { return "resume_roots" }

type ResumePresentation struct {
	ResumeID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	GreetingFirst      string    `gorm:"not null"`
	GreetingSecond     string    `gorm:"not null"`
	GreetingHighlight  string    `gorm:"not null"`
	GreetingThird      string    `gorm:"not null"`
	AvailabilityStatus *string
	AvailabilityFrom   *time.Time `gorm:"type:date"`
	AvailabilityNote   *string
	ResumeLink         *string
	ThemePrimary       *string
	ThemeSecondary     *string
}

func (ResumePresentation) TableName() string { return "resume_presentations" }

type ResumePresentationSection struct {
	ResumeID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	SectionID string    `gorm:"type:text;primaryKey"`
	Position  int       `gorm:"not null"`
	Enabled   bool      `gorm:"not null"`
	Collapsed bool      `gorm:"not null"`
}

func (ResumePresentationSection) TableName() string { return "resume_presentation_sections" }

type ResumeProfile struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Name        string    `gorm:"not null"`
	Surname     string    `gorm:"not null"`
	Patronymic  *string
	Position    string    `gorm:"not null"`
	PhotoPath   string    `gorm:"not null"`
	DateOfBirth time.Time `gorm:"type:date;not null"`
	Location    string    `gorm:"not null"`
	Relocation  bool      `gorm:"not null"`
}

func (ResumeProfile) TableName() string { return "resume_profiles" }

type ResumeProfileWorkFormat struct {
	ProfileID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value     string    `gorm:"type:text;primaryKey"`
	Position  int       `gorm:"not null"`
}

func (ResumeProfileWorkFormat) TableName() string { return "resume_profile_work_formats" }

type ResumeContacts struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Phone    *string
	Telegram *string
	Email    *string
}

func (ResumeContacts) TableName() string { return "resume_contacts" }

type ResumeContactHidden struct {
	ContactsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value      string    `gorm:"type:text;primaryKey"`
	Position   int       `gorm:"not null"`
}

func (ResumeContactHidden) TableName() string { return "resume_contact_hidden" }

type ResumeContactsText struct {
	ContactsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title      *string
	Subtitle   *string
}

func (ResumeContactsText) TableName() string { return "resume_contacts_text" }

type ResumeSkills struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
}
type ResumeSkillHard struct {
	SkillsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type ResumeSkillSoft struct {
	SkillsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type ResumeLanguage struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID uuid.UUID `gorm:"type:uuid;not null;index"`
	Language string    `gorm:"type:text;not null"`
	Position int       `gorm:"not null"`
	Grade    *string
}
type ResumeTechnology struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name     string    `gorm:"type:text;not null"`
	Position int       `gorm:"not null"`
	Grade    *string
}

func (ResumeSkills) TableName() string     { return "resume_skills" }
func (ResumeSkillHard) TableName() string  { return "resume_skill_hard" }
func (ResumeSkillSoft) TableName() string  { return "resume_skill_soft" }
func (ResumeLanguage) TableName() string   { return "resume_languages" }
func (ResumeTechnology) TableName() string { return "resume_technologies" }

type ResumeExperience struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
	SortOrder   int       `gorm:"not null"`
	Company     string    `gorm:"not null"`
	Position    string    `gorm:"not null"`
	Location    *string
	StartDate   string  `gorm:"type:date;not null"`
	EndDate     *string `gorm:"type:date"`
	Description *string
}

type ResumeExperienceAchievement struct {
	ExperienceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position     int       `gorm:"primaryKey"`
	Value        string    `gorm:"not null"`
}

type ResumeExperienceTechnology struct {
	ExperienceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position     int       `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	Grade        *string
}

func (ResumeExperience) TableName() string            { return "resume_experiences" }
func (ResumeExperienceAchievement) TableName() string { return "resume_experience_achievements" }
func (ResumeExperienceTechnology) TableName() string  { return "resume_experience_technologies" }

type ResumeEducation struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
	SortOrder   int       `gorm:"not null"`
	Location    string    `gorm:"not null"`
	University  string    `gorm:"not null"`
	Grade       int       `gorm:"not null"`
	ProgramName string    `gorm:"not null"`
	ProfileName *string
	Faculty     string `gorm:"not null"`
	Course      *int
	EndDate     *string `gorm:"type:date"`
}

func (ResumeEducation) TableName() string { return "resume_educations" }

type ResumeCourse struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID        uuid.UUID `gorm:"type:uuid;not null;index"`
	SortOrder       int       `gorm:"not null"`
	Name            string    `gorm:"not null"`
	Organization    string    `gorm:"not null"`
	EndDate         string    `gorm:"type:date;not null"`
	Hours           *int
	CertificateLink *string
}

func (ResumeCourse) TableName() string { return "resume_courses" }

type ResumeActivity struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
	SortOrder   int       `gorm:"not null"`
	Type        string    `gorm:"not null"`
	Title       string    `gorm:"not null"`
	Description *string
	Date        *string `gorm:"type:date"`
	Link        *string
}

func (ResumeActivity) TableName() string { return "resume_activities" }

type ResumeProject struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
	SortOrder   int       `gorm:"not null"`
	Name        string    `gorm:"not null"`
	Description *string
	Role        *string
	Link        *string
	GitLink     *string
}

type ResumeProjectScreenshot struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position  int       `gorm:"primaryKey"`
	Path      string    `gorm:"not null"`
}

type ResumeProjectTechnology struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position  int       `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Grade     *string
}

type ResumeProjectContributor struct {
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

func (ResumeProject) TableName() string            { return "resume_projects" }
func (ResumeProjectScreenshot) TableName() string  { return "resume_project_screenshots" }
func (ResumeProjectTechnology) TableName() string  { return "resume_project_technologies" }
func (ResumeProjectContributor) TableName() string { return "resume_project_contributors" }
