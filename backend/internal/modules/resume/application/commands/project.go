package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ProjectInput struct {
	Name         string
	Description  *string
	Role         *string
	Screenshots  []valueobjects.AssetPath
	Link         *valueobjects.ExternalURL
	GitLink      *valueobjects.ExternalURL
	Stack        []ProjectTechnologyInput
	Contributors []ProjectContributorInput
}
type ProjectTechnologyInput struct {
	Name  string
	Grade *enums.TechnologyGrade
}
type ProjectContributorInput struct {
	Name, Surname, Role string
	Contacts            ProjectContactsInput
}
type ProjectContactsInput struct {
	Phone    string
	Telegram *valueobjects.Telegram
	Email    *valueobjects.Email
}

func projectFromInput(input ProjectInput, id, resumeID valueobjects.ID) (*domain.Project, error) {
	screenshots := append([]valueobjects.AssetPath(nil), input.Screenshots...)
	var link, gitLink *valueobjects.ExternalURL
	if input.Link != nil {
		copy := *input.Link
		link = &copy
	}
	if input.GitLink != nil {
		copy := *input.GitLink
		gitLink = &copy
	}
	stack := make([]entities.Technology, 0, len(input.Stack))
	for _, v := range input.Stack {
		x, e := entities.NewTechnology(entities.TechnologyParams{Name: v.Name, Grade: v.Grade})
		if e != nil {
			return nil, e
		}
		stack = append(stack, x)
	}
	contributors := make([]entities.Contributor, 0, len(input.Contributors))
	for _, v := range input.Contributors {
		var tg *valueobjects.Telegram
		if v.Contacts.Telegram != nil {
			copy := *v.Contacts.Telegram
			tg = &copy
		}
		var em *valueobjects.Email
		if v.Contacts.Email != nil {
			copy := *v.Contacts.Email
			em = &copy
		}
		var cs entities.Contacts
		if v.Contacts.Phone == "" && tg == nil && em == nil {
			cs = entities.EmptyContacts()
		} else {
			x, e := entities.NewContacts(entities.ContactsParams{Phone: v.Contacts.Phone, Telegram: tg, Email: em})
			if e != nil {
				return nil, e
			}
			cs = x
		}
		x, e := entities.NewContributor(entities.ContributorParams{Name: v.Name, Surname: v.Surname, Role: v.Role, Contacts: cs})
		if e != nil {
			return nil, e
		}
		contributors = append(contributors, x)
	}
	return domain.NewProject(domain.ProjectParams{ID: id, ResumeID: resumeID, Name: input.Name, Description: input.Description, Role: input.Role, Screenshots: screenshots, Link: link, GitLink: gitLink, Stack: stack, Contributors: contributors})
}

type CreateProjectCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Input        ProjectInput
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

type CreateProjectHandler struct {
	repository ports.CreateProjectRepository
	uow        UnitOfWork
}

func NewCreateProjectHandler(repository ports.CreateProjectRepository, uow UnitOfWork) *CreateProjectHandler {
	return &CreateProjectHandler{repository: repository, uow: uow}
}

func (h *CreateProjectHandler) Handle(ctx context.Context, command CreateProjectCommand) (valueobjects.ID, error) {
	if command.ID.IsZero() {
		command.ID = valueobjects.NewID()
	}
	project, err := projectFromInput(command.Input, command.ID, command.ResumeID)
	if err != nil {
		return valueobjects.ID{}, err
	}
	return command.ID, h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.CreateProject(tx, ports.CreateProjectRepositoryCommand{Project: project, ResumeID: command.ResumeID, UserID: command.UserID})
	})
}

type UpdateProjectHandler struct {
	repository ports.UpdateProjectRepository
	uow        UnitOfWork
}

func NewUpdateProjectHandler(repository ports.UpdateProjectRepository, uow UnitOfWork) *UpdateProjectHandler {
	return &UpdateProjectHandler{repository: repository, uow: uow}
}

func (h *UpdateProjectHandler) Handle(ctx context.Context, command UpdateProjectCommand) error {
	project, err := projectFromInput(command.Input, command.ProjectID, command.ResumeID)
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.UpdateProject(tx, ports.UpdateProjectRepositoryCommand{Project: project, ResumeID: command.ResumeID, UserID: command.UserID})
	})
}

type DeleteProjectHandler struct {
	repository ports.DeleteProjectRepository
	uow        UnitOfWork
}

func NewDeleteProjectHandler(repository ports.DeleteProjectRepository, uow UnitOfWork) *DeleteProjectHandler {
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
