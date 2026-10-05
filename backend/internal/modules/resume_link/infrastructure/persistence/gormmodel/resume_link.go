package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type ResumeLink struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"type:text;not null"`
	Status    string    `gorm:"type:text;not null"`
	Version   int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (ResumeLink) TableName() string { return "resume_links" }
