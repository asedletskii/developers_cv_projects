package projects

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
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
	Name         string               `json:"name"`
	Description  *string              `json:"description,omitempty"`
	Role         *string              `json:"role,omitempty"`
	Screenshots  []string             `json:"screenshots"`
	Link         *string              `json:"link,omitempty"`
	GitLink      *string              `json:"gitLink,omitempty"`
	Stack        []technologyRequest  `json:"stack"`
	Contributors []contributorRequest `json:"contributors"`
}

func projectInput(body projectRequest) (commands.ProjectInput, error) {
	screenshots := make([]valueobjects.AssetPath, len(body.Screenshots))
	for i, value := range body.Screenshots {
		path, err := valueobjects.NewAssetPath(value)
		if err != nil {
			return commands.ProjectInput{}, err
		}
		screenshots[i] = path
	}

	var link *valueobjects.ExternalURL
	if body.Link != nil {
		value, err := valueobjects.NewExternalURL(*body.Link)
		if err != nil {
			return commands.ProjectInput{}, err
		}
		link = &value
	}

	var gitLink *valueobjects.ExternalURL
	if body.GitLink != nil {
		value, err := valueobjects.NewExternalURL(*body.GitLink)
		if err != nil {
			return commands.ProjectInput{}, err
		}
		gitLink = &value
	}

	stack := make([]commands.ProjectTechnologyInput, len(body.Stack))
	for i, value := range body.Stack {
		var grade *enums.TechnologyGrade
		if value.Grade != nil {
			parsed := enums.TechnologyGrade(*value.Grade)
			grade = &parsed
		}
		stack[i] = commands.ProjectTechnologyInput{Name: value.Name, Grade: grade}
	}

	contributors := make([]commands.ProjectContributorInput, len(body.Contributors))
	for i, value := range body.Contributors {
		var telegram *valueobjects.Telegram
		if value.Contacts.Telegram != nil {
			parsed, err := valueobjects.NewTelegram(*value.Contacts.Telegram)
			if err != nil {
				return commands.ProjectInput{}, err
			}
			telegram = &parsed
		}
		var email *valueobjects.Email
		if value.Contacts.Email != nil {
			parsed, err := valueobjects.NewEmail(*value.Contacts.Email)
			if err != nil {
				return commands.ProjectInput{}, err
			}
			email = &parsed
		}
		contributors[i] = commands.ProjectContributorInput{
			Name:    value.Name,
			Surname: value.Surname,
			Role:    value.Role,
			Contacts: commands.ProjectContactsInput{
				Phone:    value.Contacts.Phone,
				Telegram: telegram,
				Email:    email,
			},
		}
	}

	return commands.ProjectInput{
		Name:         body.Name,
		Description:  body.Description,
		Role:         body.Role,
		Screenshots:  screenshots,
		Link:         link,
		GitLink:      gitLink,
		Stack:        stack,
		Contributors: contributors,
	}, nil
}
