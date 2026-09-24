package commands

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type ProjectInput struct {
	ID           valueobjects.ID
	ResumeID     valueobjects.ID
	Name         string
	Description  *string
	Role         *string
	Screenshots  []valueobjects.AssetPath
	Link         *valueobjects.ExternalURL
	GitLink      *valueobjects.ExternalURL
	Stack        []entities.Technology
	Contributors []entities.Contributor
}

type CreateProjectCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Input        ProjectInput
}
type ListProjectsQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetProjectQuery struct {
	ResumeID, ProjectID valueobjects.ID
	UserID              useridentity.ID
}
type UpdateProjectCommand struct {
	ResumeID, ProjectID valueobjects.ID
	UserID              useridentity.ID
	Input               ProjectInput
}
type DeleteProjectCommand struct {
	ResumeID, ProjectID valueobjects.ID
	UserID              useridentity.ID
}

type ProjectsRepository interface {
	CreateProject(context.Context, CreateProjectCommand) error
	ListProjects(context.Context, ListProjectsQuery) ([]*domain.Project, error)
	GetProject(context.Context, GetProjectQuery) (*domain.Project, error)
	UpdateProject(context.Context, UpdateProjectCommand) error
	DeleteProject(context.Context, DeleteProjectCommand) error
}

type CreateProjectHandler struct {
	repository ProjectsRepository
	uow        UnitOfWork
}

func NewCreateProjectHandler(repository ProjectsRepository, uow UnitOfWork) *CreateProjectHandler {
	return &CreateProjectHandler{repository: repository, uow: uow}
}
func (h *CreateProjectHandler) Handle(ctx context.Context, command CreateProjectCommand) (valueobjects.ID, error) {
	if command.ID.IsZero() {
		command.ID = valueobjects.NewID()
	}
	command.Input.ID, command.Input.ResumeID = command.ID, command.ResumeID
	if _, err := projectFromInput(command.Input); err != nil {
		return valueobjects.ID{}, err
	}
	return command.ID, h.uow.Do(ctx, func(tx context.Context) error { return h.repository.CreateProject(tx, command) })
}

type ListProjectsHandler struct{ repository ProjectsRepository }

func NewListProjectsHandler(repository ProjectsRepository) *ListProjectsHandler {
	return &ListProjectsHandler{repository: repository}
}
func (h *ListProjectsHandler) Handle(ctx context.Context, query ListProjectsQuery) ([]*domain.Project, error) {
	return h.repository.ListProjects(ctx, query)
}

type GetProjectHandler struct{ repository ProjectsRepository }

func NewGetProjectHandler(repository ProjectsRepository) *GetProjectHandler {
	return &GetProjectHandler{repository: repository}
}
func (h *GetProjectHandler) Handle(ctx context.Context, query GetProjectQuery) (*domain.Project, error) {
	return h.repository.GetProject(ctx, query)
}

type UpdateProjectHandler struct {
	repository ProjectsRepository
	uow        UnitOfWork
}

func NewUpdateProjectHandler(repository ProjectsRepository, uow UnitOfWork) *UpdateProjectHandler {
	return &UpdateProjectHandler{repository: repository, uow: uow}
}
func (h *UpdateProjectHandler) Handle(ctx context.Context, command UpdateProjectCommand) error {
	command.Input.ID, command.Input.ResumeID = command.ProjectID, command.ResumeID
	if _, err := projectFromInput(command.Input); err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error { return h.repository.UpdateProject(tx, command) })
}

type DeleteProjectHandler struct {
	repository ProjectsRepository
	uow        UnitOfWork
}

func NewDeleteProjectHandler(repository ProjectsRepository, uow UnitOfWork) *DeleteProjectHandler {
	return &DeleteProjectHandler{repository: repository, uow: uow}
}
func (h *DeleteProjectHandler) Handle(ctx context.Context, command DeleteProjectCommand) error {
	return h.uow.Do(ctx, func(tx context.Context) error { return h.repository.DeleteProject(tx, command) })
}

func projectFromInput(input ProjectInput) (*domain.Project, error) {
	return domain.NewProject(domain.ProjectParams{ID: input.ID, ResumeID: input.ResumeID, Name: input.Name, Description: input.Description, Role: input.Role, Screenshots: input.Screenshots, Link: input.Link, GitLink: input.GitLink, Stack: input.Stack, Contributors: input.Contributors})
}
