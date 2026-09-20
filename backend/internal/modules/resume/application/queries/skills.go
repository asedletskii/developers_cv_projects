package queries

type SkillsView struct {
	ID         string   `json:"id"`
	HardSkills []string `json:"hardSkills"`
	SoftSkills []string `json:"softSkills"`
}
