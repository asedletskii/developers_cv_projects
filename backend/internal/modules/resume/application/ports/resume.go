package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateResumeRepositoryCommand struct{ Resume *domain.Resume }

type GetResumeRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type UpdateResumeRepositoryCommand struct{ Resume *domain.Resume }

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

type UpdateResumeRepository interface {
	UpdateResume(context.Context, UpdateResumeRepositoryCommand) error
}

type DeleteResumeRepository interface {
	DeleteResume(context.Context, DeleteResumeRepositoryCommand) error
}
