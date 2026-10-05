package gormmodel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ResumeLinkView struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	LinkID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	VisitorID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	IPHash        []byte         `gorm:"type:bytea;not null"`
	Device        string         `gorm:"type:text;not null"`
	Browser       string         `gorm:"type:text;not null"`
	OS            string         `gorm:"type:text;not null"`
	ReferrerHost  string         `gorm:"type:text"`
	Traffic       string         `gorm:"type:text;not null"`
	TrafficReason string         `gorm:"type:text"`
	Layout        datatypes.JSON `gorm:"type:jsonb;not null"`
	StartedAt     time.Time      `gorm:"not null"`
	ConfirmedAt   *time.Time     `gorm:""`
	LastSeenAt    *time.Time     `gorm:""`
	ActiveMS      int64          `gorm:"not null"`
	LastStateSeq  int64          `gorm:"not null"`
	MaxSection    string         `gorm:"type:text"`
	ExitSection   string         `gorm:"type:text"`
	EndedAt       *time.Time     `gorm:""`
	CreatedAt     time.Time      `gorm:"not null"`
}

func (ResumeLinkView) TableName() string { return "resume_link_views" }

type ResumeLinkAnalyticsEvent struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	ViewID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	LinkID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	ClientEventID *uuid.UUID     `gorm:"type:uuid"`
	EventType     string         `gorm:"type:text;not null"`
	Payload       datatypes.JSON `gorm:"type:jsonb;not null"`
	OccurredAt    time.Time      `gorm:"not null"`
}

func (ResumeLinkAnalyticsEvent) TableName() string { return "resume_link_events" }
