package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CreateResumeRepositoryCommand struct{ Resume *domain.Resume }

type GetResumeRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type ReplaceResumeTitleRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Title    string
}

type DeleteResumeRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type CreateResumeRepository interface {
	CreateResume(context.Context, CreateResumeRepositoryCommand) error
}

type GetResumeRepository interface {
	GetResume(context.Context, GetResumeRepositoryCommand) (*domain.Resume, error)
}

type ReplaceResumeTitleRepository interface {
	ReplaceResumeTitle(context.Context, ReplaceResumeTitleRepositoryCommand) error
}

type DeleteResumeRepository interface {
	DeleteResume(context.Context, DeleteResumeRepositoryCommand) error
}
