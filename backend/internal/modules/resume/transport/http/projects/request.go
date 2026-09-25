package projects

import "developers-cv/backend/internal/modules/resume/application/commands"

type TechnologyRequest struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}

type ContributorRequest struct {
	Name     string          `json:"name"`
	Surname  string          `json:"surname"`
	Role     string          `json:"role"`
	Contacts ContactsRequest `json:"contacts"`
}

type ContactsRequest struct {
	Phone    string  `json:"phone"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type ProjectRequest struct {
	Name         string               `json:"name"`
	Description  *string              `json:"description,omitempty"`
	Role         *string              `json:"role,omitempty"`
	Screenshots  []string             `json:"screenshots"`
	Link         *string              `json:"link,omitempty"`
	GitLink      *string              `json:"gitLink,omitempty"`
	Stack        []TechnologyRequest  `json:"stack"`
	Contributors []ContributorRequest `json:"contributors"`
}

func projectInput(body ProjectRequest) commands.ProjectInput {
	stack := make([]commands.ProjectTechnologyInput, len(body.Stack))
	for i, v := range body.Stack {
		stack[i] = commands.ProjectTechnologyInput{Name: v.Name, Grade: v.Grade}
	}
	contributors := make([]commands.ProjectContributorInput, len(body.Contributors))
	for i, v := range body.Contributors {
		contributors[i] = commands.ProjectContributorInput{Name: v.Name, Surname: v.Surname, Role: v.Role, Contacts: commands.ProjectContactsInput{Phone: v.Contacts.Phone, Telegram: v.Contacts.Telegram, Email: v.Contacts.Email}}
	}
	return commands.ProjectInput{Name: body.Name, Description: body.Description, Role: body.Role, Screenshots: body.Screenshots, Link: body.Link, GitLink: body.GitLink, Stack: stack, Contributors: contributors}
}
