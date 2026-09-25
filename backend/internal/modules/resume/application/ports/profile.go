package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateResumeProfileRepositoryCommand struct {
	Profile *domain.Profile
	UserID  useridentity.ID
}

type GetResumeProfileRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type UpdateResumeProfileRepositoryCommand struct {
	Profile *domain.Profile
	UserID  useridentity.ID
}

type CreateResumeProfileRepository interface {
	CreateProfile(context.Context, CreateResumeProfileRepositoryCommand) error
}

type GetResumeProfileRepository interface {
	GetProfile(context.Context, GetResumeProfileRepositoryCommand) (*domain.Profile, error)
}

type UpdateResumeProfileRepository interface {
	UpdateProfile(context.Context, UpdateResumeProfileRepositoryCommand) error
}
