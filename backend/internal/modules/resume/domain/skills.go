package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type Skills struct {
	id, resumeID valueobjects.ID
	hard         []enums.HardSkill
	soft         []enums.SoftSkill
}
type SkillsParams struct {
	ID, ResumeID valueobjects.ID
	Hard         []enums.HardSkill
	Soft         []enums.SoftSkill
}

func NewSkills(p SkillsParams) (*Skills, error) {
	if p.ID.IsZero() || p.ResumeID.IsZero() {
		return nil, fmt.Errorf("%w: skills ids are required", ErrInvalidEntity)
	}
	h := make([]enums.HardSkill, 0, len(p.Hard))
	seen := map[string]bool{}
	for _, v := range p.Hard {
		s := strings.TrimSpace(string(v))
		if s == "" || !enums.HardSkill(s).IsValid() || seen[s] {
			return nil, fmt.Errorf("%w: invalid or duplicate hard skill", ErrInvalidEntity)
		}
		seen[s] = true
		h = append(h, enums.HardSkill(s))
	}
	sft := make([]enums.SoftSkill, 0, len(p.Soft))
	seen = map[string]bool{}
	for _, v := range p.Soft {
		s := strings.TrimSpace(string(v))
		if s == "" || !enums.SoftSkill(s).IsValid() || seen[s] {
			return nil, fmt.Errorf("%w: invalid or duplicate soft skill", ErrInvalidEntity)
		}
		seen[s] = true
		sft = append(sft, enums.SoftSkill(s))
	}
	return &Skills{id: p.ID, resumeID: p.ResumeID, hard: h, soft: sft}, nil
}
func (s *Skills) ID() valueobjects.ID           { return s.id }
func (s *Skills) ResumeID() valueobjects.ID     { return s.resumeID }
func (s *Skills) HardSkills() []enums.HardSkill { return append([]enums.HardSkill(nil), s.hard...) }
func (s *Skills) SoftSkills() []enums.SoftSkill { return append([]enums.SoftSkill(nil), s.soft...) }
