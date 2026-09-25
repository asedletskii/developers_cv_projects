package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateExperienceRepositoryCommand struct {
	Experience *domain.Experience
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}

type GetExperienceRepositoryCommand struct {
	ResumeID, ExperienceID valueobjects.ID
	UserID                 useridentity.ID
}

type UpdateExperienceRepositoryCommand struct {
	Experience *domain.Experience
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}

type DeleteExperienceRepositoryCommand struct {
	ResumeID, ExperienceID valueobjects.ID
	UserID                 useridentity.ID
}

type CreateExperienceRepository interface {
	CreateExperience(context.Context, CreateExperienceRepositoryCommand) error
}

type GetExperienceRepository interface {
	GetExperience(context.Context, GetExperienceRepositoryCommand) (*domain.Experience, error)
}

type UpdateExperienceRepository interface {
	UpdateExperience(context.Context, UpdateExperienceRepositoryCommand) error
}

type DeleteExperienceRepository interface {
	DeleteExperience(context.Context, DeleteExperienceRepositoryCommand) error
}
