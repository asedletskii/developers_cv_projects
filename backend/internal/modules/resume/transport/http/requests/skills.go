package requests

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type SkillsRequest struct {
	Version    int64    `json:"version"`
	ID         string   `json:"id,omitempty"`
	HardSkills []string `json:"hardSkills"`
	SoftSkills []string `json:"softSkills"`
}

func (r SkillsRequest) ToInput() (commands.SkillsInput, error) {
	var id sharedvalueobjects.ID
	if r.ID != "" {
		parsed, err := sharedvalueobjects.ParseID(r.ID)
		if err != nil {
			return commands.SkillsInput{}, err
		}
		id = parsed
	}
	return commands.SkillsInput{
		ID: id, Hard: append([]string(nil), r.HardSkills...), Soft: append([]string(nil), r.SoftSkills...),
	}, nil
}
