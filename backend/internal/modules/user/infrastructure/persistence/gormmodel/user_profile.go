package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	PhotoURL  *string   `gorm:"column:photo_path"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (UserProfile) TableName() string { return "user_profiles" }

type UserProfileLocalization struct {
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Language   string    `gorm:"primaryKey"`
	Name       *string
	Surname    *string
	Patronymic *string
	Position   *string
	Location   *string
}

func (UserProfileLocalization) TableName() string { return "user_profile_localizations" }
