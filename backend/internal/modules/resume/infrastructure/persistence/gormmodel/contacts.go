package gormmodel

import (
	"github.com/google/uuid"
)

type ResumeContacts struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Phone    *string
	Telegram *string
	Email    *string
}

func (ResumeContacts) TableName() string { return schema + ".resume_contacts" }

type ResumeContactHidden struct {
	ContactsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value      string    `gorm:"type:text;primaryKey"`
	Position   int       `gorm:"not null"`
}

func (ResumeContactHidden) TableName() string { return schema + ".resume_contact_hidden" }

type ResumeContactsText struct {
	ContactsID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title      *string
	Subtitle   *string
}

func (ResumeContactsText) TableName() string { return schema + ".resume_contacts_text" }
