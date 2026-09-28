package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ProjectInput struct {
	ID           valueobjects.ID
	Name         string
	Description  *string
	Role         *string
	Screenshots  []valueobjects.ExternalURL
	Link         *valueobjects.ExternalURL
	GitLink      *valueobjects.ExternalURL
	Stack        []ProjectTechnologyInput
	Contributors []ProjectContributorInput
}

type ProjectsInput struct {
	Items []ProjectInput
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

func projectInputToDomain(input ProjectInput) (domain.ProjectParams, error) {
	screenshots := append([]valueobjects.ExternalURL(nil), input.Screenshots...)
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
	for _, item := range input.Stack {
		technology, err := entities.NewTechnology(entities.TechnologyParams{Name: item.Name, Grade: item.Grade})
		if err != nil {
			return domain.ProjectParams{}, err
		}
		stack = append(stack, technology)
	}

	contributors := make([]entities.Contributor, 0, len(input.Contributors))
	for _, item := range input.Contributors {
		var telegram *valueobjects.Telegram
		if item.Contacts.Telegram != nil {
			copy := *item.Contacts.Telegram
			telegram = &copy
		}
		var email *valueobjects.Email
		if item.Contacts.Email != nil {
			copy := *item.Contacts.Email
			email = &copy
		}
		contacts, err := entities.NewContacts(entities.ContactsParams{
			Phone: item.Contacts.Phone, Telegram: telegram, Email: email,
		})
		if err != nil {
			return domain.ProjectParams{}, err
		}
		contributor, err := entities.NewContributor(entities.ContributorParams{
			Name: item.Name, Surname: item.Surname, Role: item.Role, Contacts: contacts,
		})
		if err != nil {
			return domain.ProjectParams{}, err
		}
		contributors = append(contributors, contributor)
	}

	return domain.ProjectParams{
		ID: input.ID, Name: input.Name, Description: input.Description, Role: input.Role,
		Screenshots: screenshots, Link: link, GitLink: gitLink,
		Stack: stack, Contributors: contributors,
	}, nil
}

func (i ProjectsInput) toDomain() ([]domain.ProjectParams, error) {
	out := make([]domain.ProjectParams, 0, len(i.Items))
	for _, item := range i.Items {
		params, err := projectInputToDomain(item)
		if err != nil {
			return nil, err
		}
		out = append(out, params)
	}
	return out, nil
}

type ReplaceProjectsCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Input    ProjectsInput
}

type ReplaceProjectsHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceProjectsHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceProjectsHandler {
	return &ReplaceProjectsHandler{repository: repository, uow: uow}
}

func (h *ReplaceProjectsHandler) Handle(ctx context.Context, c ReplaceProjectsCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace projects handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace projects command")
	}
	params, err := c.Input.toDomain()
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if resume.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := resume.ReplaceProjects(params); err != nil {
			return err
		}
		return h.repository.Save(tx, resume)
	})
}
