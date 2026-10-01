package queries

type TechnologyView struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Grade    *string `json:"grade,omitempty"`
	Position int     `json:"position"`
}
