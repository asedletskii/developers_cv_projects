package queries

type LanguageView struct {
	ID       string  `json:"id"`
	ResumeID string  `json:"resumeId"`
	Language string  `json:"language"`
	Grade    *string `json:"grade,omitempty"`
	Position int     `json:"position"`
}
