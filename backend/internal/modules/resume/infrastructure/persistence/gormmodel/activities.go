package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type ResumeActivity struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
	SortOrder   int       `gorm:"not null"`
	Type        string    `gorm:"not null"`
	Title       string    `gorm:"not null"`
	Description *string
	Date        *time.Time `gorm:"type:date"`
	Link        *string
}

func (ResumeActivity) TableName() string { return schema + ".resume_activities" }
