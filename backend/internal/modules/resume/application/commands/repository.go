package commands

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateResumeRepositoryCommand struct {
	Resume *domain.Resume
}

type UpdateResumeRepositoryCommand struct {
	Resume *domain.Resume
}

type GetResumeRepositoryQuery struct {
	UserID   useridentity.ID
	ResumeID valueobjects.ID
}

type UpdateResumeTitleRepositoryCommand struct {
	UserID   useridentity.ID
	ResumeID valueobjects.ID
	Title    string
}

type CreateResumeRepository interface {
	Create(ctx context.Context, command CreateResumeRepositoryCommand) error
}

type UpdateResumeRepository interface {
	Update(ctx context.Context, command UpdateResumeRepositoryCommand) error
}

type GetResumeRepository interface {
	Get(ctx context.Context, query GetResumeRepositoryQuery) (*domain.Resume, error)
}

type UpdateResumeTitleRepository interface {
	UpdateTitle(ctx context.Context, command UpdateResumeTitleRepositoryCommand) error
}
