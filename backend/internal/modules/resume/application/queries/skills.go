package queries

type SkillsView struct {
	ID         string   `json:"id"`
	ResumeID   string   `json:"resumeId"`
	HardSkills []string `json:"hardSkills"`
	SoftSkills []string `json:"softSkills"`
}
