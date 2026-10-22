package gormmodel

import (
	"github.com/google/uuid"
)

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

func (ResumeProject) TableName() string { return schema + ".resume_projects" }

type ResumeProjectScreenshot struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position  int       `gorm:"primaryKey"`
	Path      string    `gorm:"not null"`
}

func (ResumeProjectScreenshot) TableName() string { return schema + ".resume_project_screenshots" }

type ResumeProjectTechnology struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position  int       `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
}

func (ResumeProjectTechnology) TableName() string { return schema + ".resume_project_technologies" }

type ResumeProjectContributor struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position  int       `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Surname   string    `gorm:"not null"`
	Role      string    `gorm:"not null"`
	Phone     *string
	Telegram  *string
	Email     *string
}

func (ResumeProjectContributor) TableName() string { return schema + ".resume_project_contributors" }
