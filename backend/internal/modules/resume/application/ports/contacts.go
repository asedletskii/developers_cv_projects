package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateResumeContactsRepositoryCommand struct {
	Contacts *domain.Contacts
	UserID   useridentity.ID
}

type GetResumeContactsRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type UpdateResumeContactsRepositoryCommand struct {
	Contacts *domain.Contacts
	UserID   useridentity.ID
}

type CreateResumeContactsRepository interface {
	CreateContacts(context.Context, CreateResumeContactsRepositoryCommand) error
}

type GetResumeContactsRepository interface {
	GetContacts(context.Context, GetResumeContactsRepositoryCommand) (*domain.Contacts, error)
}

type UpdateResumeContactsRepository interface {
	UpdateContacts(context.Context, UpdateResumeContactsRepositoryCommand) error
}
