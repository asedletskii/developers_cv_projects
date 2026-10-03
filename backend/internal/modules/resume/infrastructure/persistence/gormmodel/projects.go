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

type ResumeProjectScreenshot struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position  int       `gorm:"primaryKey"`
	Path      string    `gorm:"not null"`
}

type ResumeProjectTechnology struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position  int       `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
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

func (ResumeProjectContributor) TableName() string { return "resume_project_contributors" }
