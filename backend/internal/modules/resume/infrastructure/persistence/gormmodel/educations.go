package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

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
	EndDate     *time.Time `gorm:"type:date"`
}

func (ResumeEducation) TableName() string { return "resume_educations" }
