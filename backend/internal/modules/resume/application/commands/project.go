package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"fmt"

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
type GetProjectCommand struct {
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

type CreateProjectRepository = ports.CreateProjectRepository
type GetProjectRepository = ports.GetProjectRepository
type UpdateProjectRepository = ports.UpdateProjectRepository
type DeleteProjectRepository = ports.DeleteProjectRepository

type CreateProjectHandler struct {
	repository CreateProjectRepository
	uow        UnitOfWork
}

func NewCreateProjectHandler(repository CreateProjectRepository, uow UnitOfWork) *CreateProjectHandler {
	return &CreateProjectHandler{repository: repository, uow: uow}
}

func (h *CreateProjectHandler) Handle(ctx context.Context, command CreateProjectCommand) (valueobjects.ID, error) {
	if command.ID.IsZero() {
		command.ID = valueobjects.NewID()
	}
	command.Input.ID, command.Input.ResumeID = command.ID, command.ResumeID
	project, err := projectFromInput(command.Input)
	if err != nil {
		return valueobjects.ID{}, err
	}
	return command.ID, h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.CreateProject(tx, ports.CreateProjectRepositoryCommand{Project: project, ResumeID: command.ResumeID, UserID: command.UserID})
	})
}

type UpdateProjectHandler struct {
	get        GetProjectRepository
	repository UpdateProjectRepository
	uow        UnitOfWork
}

func NewUpdateProjectHandler(get GetProjectRepository, repository UpdateProjectRepository, uow UnitOfWork) *UpdateProjectHandler {
	return &UpdateProjectHandler{get: get, repository: repository, uow: uow}
}

func (h *UpdateProjectHandler) Handle(ctx context.Context, command UpdateProjectCommand) error {
	return h.uow.Do(ctx, func(tx context.Context) error {
		current, err := h.get.GetProject(tx, ports.GetProjectRepositoryCommand{ResumeID: command.ResumeID, ProjectID: command.ProjectID, UserID: command.UserID})
		if err != nil {
			return err
		}
		if err := current.Replace(domain.ProjectParams{ID: current.ID(), ResumeID: command.ResumeID, Name: command.Input.Name, Description: command.Input.Description, Role: command.Input.Role, Screenshots: command.Input.Screenshots, Link: command.Input.Link, GitLink: command.Input.GitLink, Stack: command.Input.Stack, Contributors: command.Input.Contributors}); err != nil {
			return err
		}
		return h.repository.UpdateProject(tx, ports.UpdateProjectRepositoryCommand{Project: current, ResumeID: command.ResumeID, ProjectID: command.ProjectID, UserID: command.UserID})
	})
}

type DeleteProjectHandler struct {
	repository DeleteProjectRepository
	uow        UnitOfWork
}

func NewDeleteProjectHandler(repository DeleteProjectRepository, uow UnitOfWork) *DeleteProjectHandler {
	return &DeleteProjectHandler{repository: repository, uow: uow}
}

func (h *DeleteProjectHandler) Handle(ctx context.Context, command DeleteProjectCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("delete project handler dependencies are nil")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.DeleteProject(tx, ports.DeleteProjectRepositoryCommand{ResumeID: command.ResumeID, ProjectID: command.ProjectID, UserID: command.UserID})
	})
}

func projectFromInput(input ProjectInput) (*domain.Project, error) {
	return domain.NewProject(domain.ProjectParams{ID: input.ID, ResumeID: input.ResumeID, Name: input.Name, Description: input.Description, Role: input.Role, Screenshots: input.Screenshots, Link: input.Link, GitLink: input.GitLink, Stack: input.Stack, Contributors: input.Contributors})
}
