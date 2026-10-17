package gormmodel

import (
	"github.com/google/uuid"
)

type ResumeSpokenLanguage struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID uuid.UUID `gorm:"type:uuid;not null;index"`
	Language string    `gorm:"type:text;not null"`
	Position int       `gorm:"not null"`
	Grade    *string
}

func (ResumeSpokenLanguage) TableName() string { return "resume_spoken_languages" }
