package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CreateLanguageRepositoryCommand struct {
	Language domain.ResumeLanguageEntry
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type GetLanguageRepositoryCommand struct {
	ResumeID, LanguageID valueobjects.ID
	UserID               useridentity.ID
}

type UpdateLanguageRepositoryCommand struct {
	Language domain.ResumeLanguageEntry
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type ReorderLanguagesRepositoryCommand struct {
	ResumeID    valueobjects.ID
	UserID      useridentity.ID
	LanguageIDs []valueobjects.ID
}

type DeleteLanguageRepositoryCommand struct {
	ResumeID, LanguageID valueobjects.ID
	UserID               useridentity.ID
}

type CreateLanguageRepository interface {
	CreateLanguage(context.Context, CreateLanguageRepositoryCommand) error
}

type GetLanguageRepository interface {
	GetLanguage(context.Context, GetLanguageRepositoryCommand) (domain.ResumeLanguageEntry, error)
}

type UpdateLanguageRepository interface {
	UpdateLanguage(context.Context, UpdateLanguageRepositoryCommand) error
}

type ReorderLanguagesRepository interface {
	ReorderLanguages(context.Context, ReorderLanguagesRepositoryCommand) error
}

type DeleteLanguageRepository interface {
	DeleteLanguage(context.Context, DeleteLanguageRepositoryCommand) error
}
