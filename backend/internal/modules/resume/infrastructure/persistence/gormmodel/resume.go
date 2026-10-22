package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type Resume struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Title     string    `gorm:"not null"`
	Language  string    `gorm:"type:text;not null"`
	Status    string    `gorm:"type:text;not null"`
	Version   int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (Resume) TableName() string { return schema + ".resumes" }
