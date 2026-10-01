package domain

import (
	"fmt"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Skills struct {
	id, resumeID sharedvalueobjects.ID
	hard         []string
	soft         []string
}
type SkillsParams struct {
	ID, ResumeID sharedvalueobjects.ID
	Hard         []string
	Soft         []string
}

func NewSkills(p SkillsParams) (*Skills, error) {
	if p.ID.IsZero() || p.ResumeID.IsZero() {
		return nil, fmt.Errorf("%w: skills ids are required", ErrInvalidEntity)
	}
	h := make([]string, 0, len(p.Hard))
	seen := map[string]bool{}
	for _, v := range p.Hard {
		s := strings.TrimSpace(v)
		if s == "" || seen[s] {
			return nil, fmt.Errorf("%w: invalid or duplicate hard skill", ErrInvalidEntity)
		}
		seen[s] = true
		h = append(h, s)
	}
	sft := make([]string, 0, len(p.Soft))
	seen = map[string]bool{}
	for _, v := range p.Soft {
		s := strings.TrimSpace(v)
		if s == "" || seen[s] {
			return nil, fmt.Errorf("%w: invalid or duplicate soft skill", ErrInvalidEntity)
		}
		seen[s] = true
		sft = append(sft, s)
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

func (s *Skills) ID() sharedvalueobjects.ID       { return s.id }
func (s *Skills) ResumeID() sharedvalueobjects.ID { return s.resumeID }
func (s *Skills) HardSkills() []string            { return append([]string(nil), s.hard...) }
func (s *Skills) SoftSkills() []string            { return append([]string(nil), s.soft...) }

func (r *Resume) ReplaceSkills(params SkillsParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if params.ID.IsZero() {
		if r.skills != nil {
			params.ID = r.skills.ID()
		} else {
			params.ID = sharedvalueobjects.NewID()
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
