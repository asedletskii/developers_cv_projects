package queries

type ActivityView struct {
	ID          string  `json:"id"`
	ResumeID    string  `json:"resumeId"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Date        *string `json:"date,omitempty"`
	Link        *string `json:"link,omitempty"`
}
