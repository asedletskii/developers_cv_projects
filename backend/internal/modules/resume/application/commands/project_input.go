package commands

import (
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ProjectInput struct {
	ID           sharedvalueobjects.ID
	Name         string
	Description  *string
	Role         *string
	Screenshots  []domain.ExternalURL
	Link         *domain.ExternalURL
	GitLink      *domain.ExternalURL
	Stack        []ProjectTechnologyInput
	Contributors []ProjectContributorInput
}

type ProjectsInput struct {
	Items []ProjectInput
}

type ProjectTechnologyInput struct {
	Name string
}

type ProjectContributorInput struct {
	Name, Surname, Role string
	Contacts            ProjectContactsInput
}

type ProjectContactsInput struct {
	Phone    string
	Telegram *domain.Telegram
	Email    *domain.Email
}

func projectInputToDomain(input ProjectInput) (domain.ProjectParams, error) {
	screenshots := append([]domain.ExternalURL(nil), input.Screenshots...)
	var link, gitLink *domain.ExternalURL
	if input.Link != nil {
		copy := *input.Link
		link = &copy
	}
	if input.GitLink != nil {
		copy := *input.GitLink
		gitLink = &copy
	}

	stack := make([]domain.Technology, 0, len(input.Stack))
	for _, item := range input.Stack {
		technology, err := domain.NewTechnology(domain.TechnologyParams{Name: item.Name})
		if err != nil {
			return domain.ProjectParams{}, err
		}
		stack = append(stack, technology)
	}

	contributors := make([]domain.Contributor, 0, len(input.Contributors))
	for _, item := range input.Contributors {
		var telegram *domain.Telegram
		if item.Contacts.Telegram != nil {
			copy := *item.Contacts.Telegram
			telegram = &copy
		}
		var email *domain.Email
		if item.Contacts.Email != nil {
			copy := *item.Contacts.Email
			email = &copy
		}
		contacts, err := domain.NewContactMethods(domain.ContactMethodsParams{
			Phone: item.Contacts.Phone, Telegram: telegram, Email: email,
		})
		if err != nil {
			return domain.ProjectParams{}, err
		}
		contributor, err := domain.NewContributor(domain.ContributorParams{
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
