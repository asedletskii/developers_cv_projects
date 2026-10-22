package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type AnalyticsVisitor struct {
	LinkID      uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	VisitorID   uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	FirstSeenAt time.Time `gorm:"not null"`
	FirstViewID uuid.UUID `gorm:"type:uuid;not null"`
	Ordinal     *int64    `gorm:""`
	CreatedAt   time.Time `gorm:"not null"`
}

func (AnalyticsVisitor) TableName() string { return schema + ".analytics_visitors" }
