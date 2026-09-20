package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type ResumeOverview struct {
	ResumeID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	GreetingFirst      string    `gorm:"not null"`
	GreetingSecond     string    `gorm:"not null"`
	GreetingHighlight  string    `gorm:"not null"`
	GreetingThird      string    `gorm:"not null"`
	AvailabilityStatus *string
	AvailabilityFrom   *time.Time `gorm:"type:date"`
	AvailabilityNote   *string
	ThemePrimary       *string
	ThemeSecondary     *string
}

func (ResumeOverview) TableName() string { return schema + ".resume_overviews" }

type ResumeOverviewSection struct {
	ResumeID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	SectionID string    `gorm:"type:text;primaryKey"`
	Position  int       `gorm:"not null"`
	Enabled   bool      `gorm:"not null"`
	Collapsed bool      `gorm:"not null"`
	Style     *string
}

func (ResumeOverviewSection) TableName() string { return schema + ".resume_overview_sections" }
