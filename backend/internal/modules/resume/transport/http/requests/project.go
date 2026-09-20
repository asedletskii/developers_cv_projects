package requests

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ProjectTechnologyRequest struct {
	Name string `json:"name"`
}

type ProjectContributorRequest struct {
	Name     string                 `json:"name"`
	Surname  string                 `json:"surname"`
	Role     string                 `json:"role"`
	Contacts ProjectContactsRequest `json:"contacts"`
}

type ProjectContactsRequest struct {
	Phone    string  `json:"phone"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type ProjectRequest struct {
	ID           string                      `json:"id,omitempty"`
	Name         string                      `json:"name"`
	Description  *string                     `json:"description,omitempty"`
	Role         *string                     `json:"role,omitempty"`
	Screenshots  []string                    `json:"screenshots"`
	Link         *string                     `json:"link,omitempty"`
	GitLink      *string                     `json:"gitLink,omitempty"`
	Stack        []ProjectTechnologyRequest  `json:"stack"`
	Contributors []ProjectContributorRequest `json:"contributors"`
}

type ReplaceProjectsRequest struct {
	Version int64            `json:"version"`
	Items   []ProjectRequest `json:"items"`
}

func (body ProjectRequest) ToInput() (commands.ProjectInput, error) {
	var id sharedvalueobjects.ID
	if body.ID != "" {
		parsed, err := sharedvalueobjects.ParseID(body.ID)
		if err != nil {
			return commands.ProjectInput{}, err
		}
		id = parsed
	}

	screenshots := make([]domain.ExternalURL, len(body.Screenshots))
	for i, value := range body.Screenshots {
		url, err := domain.NewExternalURL(value)
		if err != nil {
			return commands.ProjectInput{}, err
		}
		screenshots[i] = url
	}

	var link *domain.ExternalURL
	if body.Link != nil {
		value, err := domain.NewExternalURL(*body.Link)
		if err != nil {
			return commands.ProjectInput{}, err
		}
		link = &value
	}

	var gitLink *domain.ExternalURL
	if body.GitLink != nil {
		value, err := domain.NewExternalURL(*body.GitLink)
		if err != nil {
			return commands.ProjectInput{}, err
		}
		gitLink = &value
	}

	stack := make([]commands.ProjectTechnologyInput, len(body.Stack))
	for i, value := range body.Stack {
		stack[i] = commands.ProjectTechnologyInput{Name: value.Name}
	}

	contributors := make([]commands.ProjectContributorInput, len(body.Contributors))
	for i, value := range body.Contributors {
		var telegram *domain.Telegram
		if value.Contacts.Telegram != nil {
			parsed, err := domain.NewTelegram(*value.Contacts.Telegram)
			if err != nil {
				return commands.ProjectInput{}, err
			}
			telegram = &parsed
		}
		var email *domain.Email
		if value.Contacts.Email != nil {
			parsed, err := domain.NewEmail(*value.Contacts.Email)
			if err != nil {
				return commands.ProjectInput{}, err
			}
			email = &parsed
		}
		contributors[i] = commands.ProjectContributorInput{
			Name: value.Name, Surname: value.Surname, Role: value.Role,
			Contacts: commands.ProjectContactsInput{Phone: value.Contacts.Phone, Telegram: telegram, Email: email},
		}
	}

	return commands.ProjectInput{
		ID: id, Name: body.Name, Description: body.Description, Role: body.Role,
		Screenshots: screenshots, Link: link, GitLink: gitLink,
		Stack: stack, Contributors: contributors,
	}, nil
}
