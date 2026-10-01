package projects

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type technologyRequest struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}

type contributorRequest struct {
	Name     string          `json:"name"`
	Surname  string          `json:"surname"`
	Role     string          `json:"role"`
	Contacts contactsRequest `json:"contacts"`
}

type contactsRequest struct {
	Phone    string  `json:"phone"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type projectRequest struct {
	ID           string               `json:"id,omitempty"`
	Name         string               `json:"name"`
	Description  *string              `json:"description,omitempty"`
	Role         *string              `json:"role,omitempty"`
	Screenshots  []string             `json:"screenshots"`
	Link         *string              `json:"link,omitempty"`
	GitLink      *string              `json:"gitLink,omitempty"`
	Stack        []technologyRequest  `json:"stack"`
	Contributors []contributorRequest `json:"contributors"`
}

type replaceProjectsRequest struct {
	Version int64            `json:"version"`
	Items   []projectRequest `json:"items"`
}

func projectInput(body projectRequest) (commands.ProjectInput, error) {
	var id identity.ID
	if body.ID != "" {
		parsed, err := identity.ParseID(body.ID)
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
		var grade *domain.TechnologyGrade
		if value.Grade != nil {
			parsed := domain.TechnologyGrade(*value.Grade)
			grade = &parsed
		}
		stack[i] = commands.ProjectTechnologyInput{Name: value.Name, Grade: grade}
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
