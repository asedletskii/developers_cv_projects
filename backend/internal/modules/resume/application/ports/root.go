package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateResumeRootRepositoryCommand struct{ Resume *domain.ResumeRoot }

type GetResumeRootRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type UpdateResumeRootRepositoryCommand struct{ Resume *domain.ResumeRoot }

type DeleteResumeRootRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type CreateResumeRootRepository interface {
	CreateRoot(context.Context, CreateResumeRootRepositoryCommand) error
}

type GetResumeRootRepository interface {
	GetRoot(context.Context, GetResumeRootRepositoryCommand) (*domain.ResumeRoot, error)
}

type UpdateResumeRootRepository interface {
	UpdateRoot(context.Context, UpdateResumeRootRepositoryCommand) error
}

type DeleteResumeRootRepository interface {
	DeleteRoot(context.Context, DeleteResumeRootRepositoryCommand) error
}
