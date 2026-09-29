package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type ResumeLink struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"not null"`
	Status    string    `gorm:"type:text;not null"`
	Version   int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (ResumeLink) TableName() string { return "resume_links" }

type ResumeLinkEvent struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	LinkID     uuid.UUID `gorm:"type:uuid;not null;index"`
	EventType  string    `gorm:"type:text;not null"`
	OccurredAt time.Time `gorm:"not null"`
}

func (ResumeLinkEvent) TableName() string { return "resume_link_events" }
