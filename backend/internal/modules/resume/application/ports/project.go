package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateProjectRepositoryCommand struct {
	Project  *domain.Project
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type GetProjectRepositoryCommand struct {
	ResumeID, ProjectID valueobjects.ID
	UserID              useridentity.ID
}

type UpdateProjectRepositoryCommand struct {
	Project  *domain.Project
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type DeleteProjectRepositoryCommand struct {
	ResumeID, ProjectID valueobjects.ID
	UserID              useridentity.ID
}

type CreateProjectRepository interface {
	CreateProject(context.Context, CreateProjectRepositoryCommand) error
}

type GetProjectRepository interface {
	GetProject(context.Context, GetProjectRepositoryCommand) (*domain.Project, error)
}

type UpdateProjectRepository interface {
	UpdateProject(context.Context, UpdateProjectRepositoryCommand) error
}

type DeleteProjectRepository interface {
	DeleteProject(context.Context, DeleteProjectRepositoryCommand) error
}
