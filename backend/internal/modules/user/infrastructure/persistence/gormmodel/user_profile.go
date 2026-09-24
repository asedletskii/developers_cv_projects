package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name       string    `gorm:"not null"`
	Surname    string    `gorm:"not null"`
	Patronymic *string
	Position   *string
	Location   *string
	PhotoPath  *string
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
