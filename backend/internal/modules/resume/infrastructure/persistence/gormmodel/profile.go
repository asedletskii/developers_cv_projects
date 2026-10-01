package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

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
