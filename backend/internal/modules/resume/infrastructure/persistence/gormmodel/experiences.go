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

func (ResumeExperience) TableName() string { return schema + ".resume_experiences" }

type ResumeExperienceAchievement struct {
	ExperienceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position     int       `gorm:"primaryKey"`
	Value        string    `gorm:"not null"`
}

func (ResumeExperienceAchievement) TableName() string {
	return schema + ".resume_experience_achievements"
}

type ResumeExperienceTask struct {
	ExperienceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position     int       `gorm:"primaryKey"`
	Value        string    `gorm:"not null"`
}

func (ResumeExperienceTask) TableName() string { return schema + ".resume_experience_tasks" }

type ResumeExperienceStack struct {
	ExperienceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position     int       `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
}

func (ResumeExperienceStack) TableName() string { return schema + ".resume_experience_stack" }
