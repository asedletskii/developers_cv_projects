package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type ResumeExperience struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
	SortOrder   int       `gorm:"not null"`
	Company     string    `gorm:"not null"`
	Position    string    `gorm:"not null"`
	Location    *string
	StartDate   time.Time  `gorm:"type:date;not null"`
	EndDate     *time.Time `gorm:"type:date"`
	Description *string
}

type ResumeExperienceAchievement struct {
	ExperienceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position     int       `gorm:"primaryKey"`
	Value        string    `gorm:"not null"`
}

func (ResumeExperienceAchievement) TableName() string { return "resume_experience_achievements" }

type ResumeExperienceTask struct {
	ExperienceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position     int       `gorm:"primaryKey"`
	Value        string    `gorm:"not null"`
}

type ResumeExperienceStack struct {
	ExperienceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position     int       `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	Grade        *string
}

func (ResumeExperienceStack) TableName() string { return "resume_experience_stack" }
