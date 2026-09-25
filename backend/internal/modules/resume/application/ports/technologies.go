package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateTechnologyRepositoryCommand struct {
	Technology domain.ResumeTechnologyEntry
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}

type GetTechnologyRepositoryCommand struct {
	ResumeID, TechnologyID valueobjects.ID
	UserID                 useridentity.ID
}

type UpdateTechnologyRepositoryCommand struct {
	Technology domain.ResumeTechnologyEntry
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}

type DeleteTechnologyRepositoryCommand struct {
	ResumeID, TechnologyID valueobjects.ID
	UserID                 useridentity.ID
}

type CreateTechnologyRepository interface {
	CreateTechnology(context.Context, CreateTechnologyRepositoryCommand) error
}

type GetTechnologyRepository interface {
	GetTechnology(context.Context, GetTechnologyRepositoryCommand) (domain.ResumeTechnologyEntry, error)
}

type UpdateTechnologyRepository interface {
	UpdateTechnology(context.Context, UpdateTechnologyRepositoryCommand) error
}

type DeleteTechnologyRepository interface {
	DeleteTechnology(context.Context, DeleteTechnologyRepositoryCommand) error
}
