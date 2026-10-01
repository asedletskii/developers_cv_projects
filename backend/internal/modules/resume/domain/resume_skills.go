package domain

import (
	"fmt"
	"strings"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type SoftSkill string

const (
	SoftSkillCommunication  SoftSkill = "communication"
	SoftSkillTeamwork       SoftSkill = "teamwork"
	SoftSkillProblemSolving SoftSkill = "problemSolving"
	SoftSkillTimeManagement SoftSkill = "timeManagement"
	SoftSkillMentoring      SoftSkill = "mentoring"
	SoftSkillAdaptability   SoftSkill = "adaptability"
	SoftSkillOwnership      SoftSkill = "ownership"
)

func (v SoftSkill) IsValid() bool {
	switch v {
	case SoftSkillCommunication, SoftSkillTeamwork, SoftSkillProblemSolving, SoftSkillTimeManagement,
		SoftSkillMentoring, SoftSkillAdaptability, SoftSkillOwnership:
		return true
	default:
		return false
	}
}

type HardSkill string

const (
	HardSkillArchitecture  HardSkill = "architecture"
	HardSkillTesting       HardSkill = "testing"
	HardSkillCICD          HardSkill = "cicd"
	HardSkillPerformance   HardSkill = "performance"
	HardSkillAccessibility HardSkill = "accessibility"
	HardSkillCodeReview    HardSkill = "codeReview"
	HardSkillUIDesign      HardSkill = "uiDesign"
	HardSkillAPI           HardSkill = "api"
)

func (v HardSkill) IsValid() bool {
	switch v {
	case HardSkillArchitecture, HardSkillTesting, HardSkillCICD, HardSkillPerformance,
		HardSkillAccessibility, HardSkillCodeReview, HardSkillUIDesign, HardSkillAPI:
		return true
	default:
		return false
	}
}

type Skills struct {
	id, resumeID identity.ID
	hard         []HardSkill
	soft         []SoftSkill
}
type SkillsParams struct {
	ID, ResumeID identity.ID
	Hard         []HardSkill
	Soft         []SoftSkill
}

func NewSkills(p SkillsParams) (*Skills, error) {
	if p.ID.IsZero() || p.ResumeID.IsZero() {
		return nil, fmt.Errorf("%w: skills ids are required", ErrInvalidEntity)
	}
	h := make([]HardSkill, 0, len(p.Hard))
	seen := map[string]bool{}
	for _, v := range p.Hard {
		s := strings.TrimSpace(string(v))
		if s == "" || !HardSkill(s).IsValid() || seen[s] {
			return nil, fmt.Errorf("%w: invalid or duplicate hard skill", ErrInvalidEntity)
		}
		seen[s] = true
		h = append(h, HardSkill(s))
	}
	sft := make([]SoftSkill, 0, len(p.Soft))
	seen = map[string]bool{}
	for _, v := range p.Soft {
		s := strings.TrimSpace(string(v))
		if s == "" || !SoftSkill(s).IsValid() || seen[s] {
			return nil, fmt.Errorf("%w: invalid or duplicate soft skill", ErrInvalidEntity)
		}
		seen[s] = true
		sft = append(sft, SoftSkill(s))
	}
	return &Skills{id: p.ID, resumeID: p.ResumeID, hard: h, soft: sft}, nil
}
func (s *Skills) Replace(params SkillsParams) error {
	params.ID, params.ResumeID = s.id, s.resumeID
	updated, err := NewSkills(params)
	if err != nil {
		return err
	}
	*s = *updated
	return nil
}

func (s *Skills) ID() identity.ID         { return s.id }
func (s *Skills) ResumeID() identity.ID   { return s.resumeID }
func (s *Skills) HardSkills() []HardSkill { return append([]HardSkill(nil), s.hard...) }
func (s *Skills) SoftSkills() []SoftSkill { return append([]SoftSkill(nil), s.soft...) }

func (r *Resume) ReplaceSkills(params SkillsParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if params.ID.IsZero() {
		if r.skills != nil {
			params.ID = r.skills.ID()
		} else {
			params.ID = identity.NewID()
		}
	} else if r.skills == nil || params.ID != r.skills.ID() {
		return fmt.Errorf("%w: skills do not belong to resume", ErrResumeInvariant)
	}
	params.ResumeID = r.id
	skills, err := NewSkills(params)
	if err != nil {
		return err
	}
	r.skills = skills
	return nil
}

func (r *Resume) Skills() *Skills {
	if r == nil || r.skills == nil {
		return nil
	}
	copy := *r.skills
	return &copy
}
