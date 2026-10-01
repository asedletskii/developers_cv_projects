package gormmodel

import (
	"github.com/google/uuid"
)

type ResumeSkills struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
}

type ResumeSkillHard struct {
	SkillsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}

func (ResumeSkillHard) TableName() string { return "resume_skill_hard" }

type ResumeSkillSoft struct {
	SkillsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value    string    `gorm:"type:text;primaryKey"`
	Position int       `gorm:"not null"`
}

func (ResumeSkillSoft) TableName() string { return "resume_skill_soft" }
