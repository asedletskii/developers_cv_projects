package gormmodel

import (
	"time"

	"github.com/google/uuid"
)

type UserIdentity struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	Provider        string    `gorm:"type:text;not null"`
	ProviderSubject string    `gorm:"type:text;not null"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
	LastLoginAt     time.Time `gorm:"not null"`
}

func (UserIdentity) TableName() string { return "user_identities" }

type AuthSession struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash  []byte    `gorm:"type:bytea;not null;uniqueIndex"`
	CreatedAt  time.Time `gorm:"not null"`
	LastSeenAt time.Time `gorm:"not null"`
	ExpiresAt  time.Time `gorm:"not null;index"`
	RevokedAt  *time.Time
	UserAgent  *string
	IPAddress  *string `gorm:"type:inet"`
}

func (AuthSession) TableName() string { return "auth_sessions" }

type AuthTransaction struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Provider     string    `gorm:"type:text;not null"`
	StateHash    []byte    `gorm:"type:bytea;not null;uniqueIndex"`
	NonceHash    []byte    `gorm:"type:bytea;not null"`
	CodeVerifier string    `gorm:"type:text;not null"`
	RedirectURI  string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"not null"`
	ExpiresAt    time.Time `gorm:"not null;index"`
	ConsumedAt   *time.Time
}

func (AuthTransaction) TableName() string { return "auth_transactions" }
