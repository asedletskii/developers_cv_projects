package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type AnalyticsClick struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	ViewID     uuid.UUID `gorm:"type:uuid;not null;index"`
	LinkID     uuid.UUID `gorm:"type:uuid;not null"`
	Section    string    `gorm:"type:text;not null"`
	XPct       int       `gorm:"column:x_pct;not null"`
	YPct       int       `gorm:"column:y_pct;not null"`
	ViewportW  int       `gorm:"column:viewport_w;not null"`
	LayoutHash string    `gorm:"type:text;not null"`
	Target     *string   `gorm:"type:text"`
	OccurredAt time.Time `gorm:"not null"`
}

func (AnalyticsClick) TableName() string { return schema + ".analytics_clicks" }
