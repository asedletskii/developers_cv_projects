package queries

type ProjectView struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
	Description  *string                  `json:"description,omitempty"`
	Role         *string                  `json:"role,omitempty"`
	Screenshots  []string                 `json:"screenshots"`
	Link         *string                  `json:"link,omitempty"`
	GitLink      *string                  `json:"gitLink,omitempty"`
	Stack        []ProjectTechnologyView  `json:"stack"`
	Contributors []ProjectContributorView `json:"contributors"`
}

type ProjectTechnologyView struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}

type ProjectContributorView struct {
	Name     string              `json:"name"`
	Surname  string              `json:"surname"`
	Role     string              `json:"role"`
	Contacts ProjectContactsView `json:"contacts"`
}

type ProjectContactsView struct {
	Phone    *string `json:"phone,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}
