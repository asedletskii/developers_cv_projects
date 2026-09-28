package queries

type ExperienceStackItemView struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}
type ExperienceView struct {
	ID           string                    `json:"id"`
	ResumeID     string                    `json:"resumeId"`
	Company      string                    `json:"company"`
	Position     string                    `json:"position"`
	Location     *string                   `json:"location,omitempty"`
	StartDate    string                    `json:"startDate"`
	EndDate      *string                   `json:"endDate,omitempty"`
	Description  *string                   `json:"description,omitempty"`
	Achievements []string                  `json:"achievements"`
	Stack        []ExperienceStackItemView `json:"stack"`
}
