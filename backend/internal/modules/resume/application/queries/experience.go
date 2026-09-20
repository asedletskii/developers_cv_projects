package queries

type ExperienceStackItemView struct {
	Name string `json:"name"`
}
type ExperienceView struct {
	ID           string                    `json:"id"`
	Company      string                    `json:"company"`
	Position     string                    `json:"position"`
	Location     *string                   `json:"location,omitempty"`
	StartDate    string                    `json:"startDate"`
	EndDate      *string                   `json:"endDate,omitempty"`
	Description  *string                   `json:"description,omitempty"`
	Achievements []string                  `json:"achievements"`
	Tasks        []string                  `json:"tasks"`
	Stack        []ExperienceStackItemView `json:"stack"`
}
