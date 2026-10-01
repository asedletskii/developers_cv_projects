package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type ResumeCourse struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID        uuid.UUID `gorm:"type:uuid;not null;index"`
	SortOrder       int       `gorm:"not null"`
	Name            string    `gorm:"not null"`
	Organization    string    `gorm:"not null"`
	EndDate         time.Time `gorm:"type:date;not null"`
	Hours           *int
	CertificateLink *string
}

func (ResumeCourse) TableName() string { return "resume_courses" }
