package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type Resume struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"not null"`
	Surname     string    `gorm:"not null"`
	Patronymic  *string
	Position    string `gorm:"not null"`
	Title       string `gorm:"not null"`
	ResumeLink  *string
	PhotoPath   string    `gorm:"not null"`
	DateOfBirth time.Time `gorm:"type:date;not null"`
	Location    string    `gorm:"not null"`
	Relocation  bool      `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

type Greeting struct {
	ResumeID          uuid.UUID `gorm:"type:uuid;primaryKey"`
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
	ResumeID uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Status   string     `gorm:"not null"`
	From     *time.Time `gorm:"type:date"`
	NoteRU   *string
	NoteEN   *string
}

type Theme struct {
	ResumeID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	Primary   *string
	Secondary *string
}

type ContactsText struct {
	ResumeID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	TitleRU    *string
	TitleEN    *string
	SubtitleRU *string
	SubtitleEN *string
}

type Section struct {
	ResumeID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	SectionID string    `gorm:"type:text;primaryKey"`
	Position  int       `gorm:"not null"`
	Enabled   bool      `gorm:"not null"`
	Collapsed bool      `gorm:"not null"`
}

type Contact struct {
	ResumeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Phone    *string
	Telegram *string
	Email    *string
}

func (Greeting) TableName() string {
	return "resume_greetings"
}

func (Availability) TableName() string {
	return "resume_availability"
}

func (Theme) TableName() string {
	return "resume_themes"
}

func (ContactsText) TableName() string {
	return "resume_contacts_text"
}

func (Section) TableName() string {
	return "resume_sections"
}

func (Contact) TableName() string {
	return "resume_contacts"
}

type ResumeWorkFormat struct {
	ResumeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type ResumeHiddenContact struct {
	ResumeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type ResumeHardSkill struct {
	ResumeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type ResumeSoftSkill struct {
	ResumeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type ResumeLanguage struct {
	ResumeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Language string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
	Grade    *string
}
type ResumeTechnology struct {
	ResumeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
	Grade    *string
}

type Experience struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
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
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
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
	ResumeID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Position        int       `gorm:"not null"`
	Name            string    `gorm:"not null"`
	Organization    string    `gorm:"not null"`
	EndDate         string    `gorm:"type:date;not null"`
	Hours           *int
	CertificateLink *string
}

type Activity struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Position    int       `gorm:"not null"`
	Type        string    `gorm:"not null"`
	Title       string    `gorm:"not null"`
	Description *string
	Date        *string `gorm:"type:date"`
	Link        *string
}

type Project struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
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

type ResumeRoot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Title     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (ResumeRoot) TableName() string { return "resume_roots" }

type ResumePresentation struct {
	ResumeID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	GreetingFirstRU     string    `gorm:"not null"`
	GreetingFirstEN     string    `gorm:"not null"`
	GreetingSecondRU    string    `gorm:"not null"`
	GreetingSecondEN    string    `gorm:"not null"`
	GreetingHighlightRU string    `gorm:"not null"`
	GreetingHighlightEN string    `gorm:"not null"`
	GreetingThirdRU     string    `gorm:"not null"`
	GreetingThirdEN     string    `gorm:"not null"`
	AvailabilityStatus  *string
	AvailabilityFrom    *time.Time `gorm:"type:date"`
	AvailabilityNoteRU  *string
	AvailabilityNoteEN  *string
	ResumeLink          *string
	ThemePrimary        *string
	ThemeSecondary      *string
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

func (ResumeContacts) TableName() string { return "resume_contacts_v2" }

type ResumeContactHidden struct {
	ContactsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value      string    `gorm:"type:text;primaryKey"`
	Position   int       `gorm:"not null"`
}

func (ResumeContactHidden) TableName() string { return "resume_contact_hidden" }

type ResumeContactsText struct {
	ContactsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	TitleRU    *string
	TitleEN    *string
	SubtitleRU *string
	SubtitleEN *string
}

func (ResumeContactsText) TableName() string { return "resume_contacts_text_v2" }

type ResumeSkillsV2 struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
}
type ResumeSkillHardV2 struct {
	SkillsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type ResumeSkillSoftV2 struct {
	SkillsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}
type ResumeLanguageV2 struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID uuid.UUID `gorm:"type:uuid;not null;index"`
	Language string    `gorm:"type:text;not null"`
	Position int       `gorm:"not null"`
	Grade    *string
}
type ResumeTechnologyV2 struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name     string    `gorm:"type:text;not null"`
	Position int       `gorm:"not null"`
	Grade    *string
}

func (ResumeSkillsV2) TableName() string     { return "resume_skills_v2" }
func (ResumeSkillHardV2) TableName() string  { return "resume_skill_hard_v2" }
func (ResumeSkillSoftV2) TableName() string  { return "resume_skill_soft_v2" }
func (ResumeLanguageV2) TableName() string   { return "resume_languages_v2" }
func (ResumeTechnologyV2) TableName() string { return "resume_technologies_v2" }

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

func (ResumeExperience) TableName() string            { return "resume_experiences_v2" }
func (ResumeExperienceAchievement) TableName() string { return "resume_experience_achievements_v2" }
func (ResumeExperienceTechnology) TableName() string  { return "resume_experience_technologies_v2" }

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

func (ResumeEducation) TableName() string { return "resume_educations_v2" }

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

func (ResumeCourse) TableName() string { return "resume_courses_v2" }

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

func (ResumeActivity) TableName() string { return "resume_activities_v2" }

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

func (ResumeProject) TableName() string            { return "resume_projects_v2" }
func (ResumeProjectScreenshot) TableName() string  { return "resume_project_screenshots_v2" }
func (ResumeProjectTechnology) TableName() string  { return "resume_project_technologies_v2" }
func (ResumeProjectContributor) TableName() string { return "resume_project_contributors_v2" }
