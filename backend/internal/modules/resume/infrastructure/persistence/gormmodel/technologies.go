package gormmodel

import (
	"github.com/google/uuid"
)

type ResumeTechnology struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name     string    `gorm:"type:text;not null"`
	Position int       `gorm:"not null"`
}

func (ResumeTechnology) TableName() string { return schema + ".resume_technologies" }
