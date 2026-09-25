package projects

import "developers-cv/backend/internal/modules/resume/application/queries"

type TechnologyResponse struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}

type ContributorResponse struct {
	Name     string           `json:"name"`
	Surname  string           `json:"surname"`
	Role     string           `json:"role"`
	Contacts ContactsResponse `json:"contacts"`
}

type ContactsResponse struct {
	Phone    *string `json:"phone,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type ProjectResponse struct {
	ID           string                `json:"id"`
	ResumeID     string                `json:"resumeId"`
	Name         string                `json:"name"`
	Description  *string               `json:"description,omitempty"`
	Role         *string               `json:"role,omitempty"`
	Screenshots  []string              `json:"screenshots"`
	Link         *string               `json:"link,omitempty"`
	GitLink      *string               `json:"gitLink,omitempty"`
	Stack        []TechnologyResponse  `json:"stack"`
	Contributors []ContributorResponse `json:"contributors"`
}

func projectQueryView(item queries.ProjectView) ProjectResponse {
	stack := make([]TechnologyResponse, len(item.Stack))
	for i, technology := range item.Stack {
		stack[i] = TechnologyResponse{Name: technology.Name, Grade: technology.Grade}
	}
	contributors := make([]ContributorResponse, len(item.Contributors))
	for i, contributor := range item.Contributors {
		contributors[i] = ContributorResponse{
			Name: contributor.Name, Surname: contributor.Surname, Role: contributor.Role,
			Contacts: ContactsResponse{Phone: contributor.Contacts.Phone, Telegram: contributor.Contacts.Telegram, Email: contributor.Contacts.Email},
		}
	}
	return ProjectResponse{
		ID: item.ID, ResumeID: item.ResumeID, Name: item.Name, Description: item.Description, Role: item.Role,
		Screenshots: item.Screenshots, Link: item.Link, GitLink: item.GitLink, Stack: stack, Contributors: contributors,
	}
}
