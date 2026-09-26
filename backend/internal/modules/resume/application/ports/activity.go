package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CreateActivityRepositoryCommand struct {
	Activity *domain.Activity
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type GetActivityRepositoryCommand struct {
	ResumeID, ActivityID valueobjects.ID
	UserID               useridentity.ID
}

type UpdateActivityRepositoryCommand struct {
	Activity *domain.Activity
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type DeleteActivityRepositoryCommand struct {
	ResumeID, ActivityID valueobjects.ID
	UserID               useridentity.ID
}

type CreateActivityRepository interface {
	CreateActivity(context.Context, CreateActivityRepositoryCommand) error
}

type GetActivityRepository interface {
	GetActivity(context.Context, GetActivityRepositoryCommand) (*domain.Activity, error)
}

type UpdateActivityRepository interface {
	UpdateActivity(context.Context, UpdateActivityRepositoryCommand) error
}

type DeleteActivityRepository interface {
	DeleteActivity(context.Context, DeleteActivityRepositoryCommand) error
}
