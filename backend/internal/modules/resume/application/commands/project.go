package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type ProjectInput struct {
	Name         string
	Description  *string
	Role         *string
	Screenshots  []string
	Link         *string
	GitLink      *string
	Stack        []ProjectTechnologyInput
	Contributors []ProjectContributorInput
}
type ProjectTechnologyInput struct {
	Name  string
	Grade *string
}
type ProjectContributorInput struct {
	Name, Surname, Role string
	Contacts            ProjectContactsInput
}
type ProjectContactsInput struct {
	Phone           string
	Telegram, Email *string
}

func projectFromInput(input ProjectInput, id, resumeID valueobjects.ID) (*domain.Project, error) {
	screenshots := make([]valueobjects.AssetPath, len(input.Screenshots))
	for i, v := range input.Screenshots {
		x, e := valueobjects.NewAssetPath(v)
		if e != nil {
			return nil, e
		}
		screenshots[i] = x
	}
	var link, gitLink *valueobjects.ExternalURL
	if input.Link != nil {
		x, e := valueobjects.NewExternalURL(*input.Link)
		if e != nil {
			return nil, e
		}
		link = &x
	}
	if input.GitLink != nil {
		x, e := valueobjects.NewExternalURL(*input.GitLink)
		if e != nil {
			return nil, e
		}
		gitLink = &x
	}
	stack := make([]entities.Technology, 0, len(input.Stack))
	for _, v := range input.Stack {
		var g *enums.TechnologyGrade
		if v.Grade != nil {
			x := enums.TechnologyGrade(*v.Grade)
			g = &x
		}
		x, e := entities.NewTechnology(entities.TechnologyParams{Name: v.Name, Grade: g})
		if e != nil {
			return nil, e
		}
		stack = append(stack, x)
	}
	contributors := make([]entities.Contributor, 0, len(input.Contributors))
	for _, v := range input.Contributors {
		var tg *valueobjects.Telegram
		if v.Contacts.Telegram != nil {
			x, e := valueobjects.NewTelegram(*v.Contacts.Telegram)
			if e != nil {
				return nil, e
			}
			tg = &x
		}
		var em *valueobjects.Email
		if v.Contacts.Email != nil {
			x, e := valueobjects.NewEmail(*v.Contacts.Email)
			if e != nil {
				return nil, e
			}
			em = &x
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
