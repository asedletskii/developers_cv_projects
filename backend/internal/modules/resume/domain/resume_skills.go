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
	ID   sharedvalueobjects.ID
	Hard []string
	Soft []string
}

type SkillsSnapshot struct {
	ID, ResumeID sharedvalueobjects.ID
	Hard         []string
	Soft         []string
}

func newSkills(params SkillsParams, resumeID sharedvalueobjects.ID) (*Skills, error) {
	skills := &Skills{id: params.ID, resumeID: resumeID}
	skills.hard = normalizeSkillValues(params.Hard)
	skills.soft = normalizeSkillValues(params.Soft)
	if err := skills.validate(); err != nil {
		return nil, err
	}
	return skills, nil
}

func RestoreSkills(snapshot SkillsSnapshot) *Skills {
	return &Skills{
		id: snapshot.ID, resumeID: snapshot.ResumeID,
		hard: append([]string(nil), snapshot.Hard...), soft: append([]string(nil), snapshot.Soft...),
	}
}

func normalizeSkillValues(values []string) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = strings.TrimSpace(value)
	}
	return result
}

func (s *Skills) validate() error {
	if s == nil || s.id.IsZero() || s.resumeID.IsZero() {
		return fmt.Errorf("%w: skills ids are required", ErrInvalidEntity)
	}
	if len(s.hard) > MaxSkillsPerResume || len(s.soft) > MaxSkillsPerResume {
		return fmt.Errorf("%w: skills exceed the per-category limit of %d", ErrInvalidEntity, MaxSkillsPerResume)
	}
	if err := validateSkillValues(s.hard); err != nil {
		return fmt.Errorf("%w: invalid or duplicate hard skill", ErrInvalidEntity)
	}
	if err := validateSkillValues(s.soft); err != nil {
		return fmt.Errorf("%w: invalid or duplicate soft skill", ErrInvalidEntity)
	}
	return nil
}

func validateSkillValues(values []string) error {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value == "" || value != strings.TrimSpace(value) || !textLengthWithin(value, MaxSkillLength) {
			return ErrInvalidEntity
		}
		if _, exists := seen[value]; exists {
			return ErrDuplicateEntity
		}
		seen[value] = struct{}{}
	}
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
	var currentID sharedvalueobjects.ID
	if r.skills != nil {
		currentID = r.skills.ID()
	}
	resolvedID, err := resolveSingletonID(params.ID, currentID, "skills")
	if err != nil {
		return err
	}
	params.ID = resolvedID
	skills, err := newSkills(params, r.id)
	if err != nil {
		return err
	}
	r.skills = skills
	return nil
}

func (r *Resume) Skills() *Skills {
	if r == nil {
		return nil
	}
	return cloneSkills(r.skills)
}

func cloneSkills(value *Skills) *Skills {
	if value == nil {
		return nil
	}
	copy := *value
	copy.hard = append([]string(nil), value.hard...)
	copy.soft = append([]string(nil), value.soft...)
	return &copy
}
