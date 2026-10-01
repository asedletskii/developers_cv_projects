package domain

import (
	"fmt"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeTechnologyEntry struct {
	id       sharedvalueobjects.ID
	name     string
	grade    *TechnologyGrade
	position int
}
type ResumeTechnologyEntryParams struct {
	ID       sharedvalueobjects.ID
	Name     string
	Grade    *TechnologyGrade
	Position int
}
type ResumeTechnologiesParams struct {
	ID, ResumeID sharedvalueobjects.ID
	Entries      []ResumeTechnologyEntryParams
}
type ResumeTechnologies struct {
	id, resumeID sharedvalueobjects.ID
	entries      []ResumeTechnologyEntry
}

func NewResumeTechnologyEntry(p ResumeTechnologyEntryParams) (ResumeTechnologyEntry, error) {
	n := strings.TrimSpace(p.Name)
	if p.ID.IsZero() || n == "" || p.Position < 0 {
		return ResumeTechnologyEntry{}, fmt.Errorf("%w: invalid technology entry", ErrInvalidEntity)
	}
	if p.Grade != nil && !p.Grade.IsValid() {
		return ResumeTechnologyEntry{}, fmt.Errorf("%w: invalid technology grade", ErrInvalidEntity)
	}
	g := p.Grade
	if g != nil {
		x := *g
		g = &x
	}
	return ResumeTechnologyEntry{id: p.ID, name: n, grade: g, position: p.Position}, nil
}
func NewResumeTechnologies(p ResumeTechnologiesParams) (*ResumeTechnologies, error) {
	if p.ID.IsZero() || p.ResumeID.IsZero() {
		return nil, fmt.Errorf("%w: technology collection ids are required", ErrInvalidEntity)
	}
	out := make([]ResumeTechnologyEntry, 0, len(p.Entries))
	seenIDs := make(map[sharedvalueobjects.ID]struct{}, len(p.Entries))
	seenPositions := make(map[int]struct{}, len(p.Entries))
	for _, x := range p.Entries {
		e, err := NewResumeTechnologyEntry(x)
		if err != nil {
			return nil, err
		}
		if _, ok := seenIDs[e.id]; ok {
			return nil, fmt.Errorf("%w: duplicate technology id", ErrInvalidEntity)
		}
		if _, ok := seenPositions[e.position]; ok {
			return nil, fmt.Errorf("%w: duplicate technology position", ErrInvalidEntity)
		}
		seenIDs[e.id] = struct{}{}
		seenPositions[e.position] = struct{}{}
		out = append(out, e)
	}
	for position := range seenPositions {
		if position < 0 || position >= len(out) {
			return nil, fmt.Errorf("%w: technology positions must be contiguous", ErrInvalidEntity)
		}
	}
	return &ResumeTechnologies{id: p.ID, resumeID: p.ResumeID, entries: out}, nil
}
func (t *ResumeTechnologyEntry) Replace(params ResumeTechnologyEntryParams) error {
	params.ID = t.id
	updated, err := NewResumeTechnologyEntry(params)
	if err != nil {
		return err
	}
	*t = updated
	return nil
}

func (t ResumeTechnologyEntry) ID() sharedvalueobjects.ID { return t.id }
func (t ResumeTechnologyEntry) Name() string              { return t.name }
func (t ResumeTechnologyEntry) Grade() *TechnologyGrade {
	if t.grade == nil {
		return nil
	}
	x := *t.grade
	return &x
}
func (t ResumeTechnologyEntry) Position() int                 { return t.position }
func (t *ResumeTechnologies) ID() sharedvalueobjects.ID       { return t.id }
func (t *ResumeTechnologies) ResumeID() sharedvalueobjects.ID { return t.resumeID }
func (t *ResumeTechnologies) Entries() []ResumeTechnologyEntry {
	return append([]ResumeTechnologyEntry(nil), t.entries...)
}

func (r *Resume) ReplaceTechnologies(params ResumeTechnologiesParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if params.ID.IsZero() {
		if r.technologies != nil {
			params.ID = r.technologies.ID()
		} else {
			params.ID = sharedvalueobjects.NewID()
		}
	} else if r.technologies == nil || params.ID != r.technologies.ID() {
		return fmt.Errorf("%w: technologies do not belong to resume", ErrResumeInvariant)
	}

	existing := make(map[sharedvalueobjects.ID]struct{})
	if r.technologies != nil {
		for _, entry := range r.technologies.Entries() {
			existing[entry.ID()] = struct{}{}
		}
	}
	seen := make(map[sharedvalueobjects.ID]struct{}, len(params.Entries))
	entries := make([]ResumeTechnologyEntryParams, 0, len(params.Entries))
	for _, entry := range params.Entries {
		if entry.ID.IsZero() {
			entry.ID = sharedvalueobjects.NewID()
		} else {
			if _, ok := existing[entry.ID]; !ok {
				return fmt.Errorf("%w: technology does not belong to resume", ErrResumeInvariant)
			}
		}
		if _, ok := seen[entry.ID]; ok {
			return fmt.Errorf("%w: duplicate technology id", ErrInvalidEntity)
		}
		seen[entry.ID] = struct{}{}
		entries = append(entries, entry)
	}

	params.ResumeID = r.id
	params.Entries = entries
	technologies, err := NewResumeTechnologies(params)
	if err != nil {
		return err
	}
	r.technologies = technologies
	return nil
}

func (r *Resume) Technologies() *ResumeTechnologies {
	if r == nil || r.technologies == nil {
		return nil
	}
	copy := *r.technologies
	copy.entries = append([]ResumeTechnologyEntry(nil), r.technologies.entries...)
	return &copy
}

type TechnologyGrade string

const (
	TechnologyGradeIntern TechnologyGrade = "intern"
	TechnologyGradeJunior TechnologyGrade = "junior"
	TechnologyGradeMiddle TechnologyGrade = "middle"
	TechnologyGradeSenior TechnologyGrade = "senior"
)

func (v TechnologyGrade) IsValid() bool {
	switch v {
	case TechnologyGradeIntern, TechnologyGradeJunior, TechnologyGradeMiddle, TechnologyGradeSenior:
		return true
	default:
		return false
	}
}

type TechnologyParams struct {
	Name  string
	Grade *TechnologyGrade
}

type Technology struct {
	name  string
	grade *TechnologyGrade
}

func NewTechnology(params TechnologyParams) (Technology, error) {
	name := strings.TrimSpace(params.Name)
	if name == "" {
		return Technology{}, fmt.Errorf("%w: technology name is empty", ErrInvalidValue)
	}
	if params.Grade != nil && !params.Grade.IsValid() {
		return Technology{}, fmt.Errorf("%w: invalid technology grade", ErrInvalidValue)
	}
	return Technology{name: name, grade: cloneTechnologyGrade(params.Grade)}, nil
}

func (t Technology) Name() string            { return t.name }
func (t Technology) Grade() *TechnologyGrade { return cloneTechnologyGrade(t.grade) }

func cloneTechnologyGrade(value *TechnologyGrade) *TechnologyGrade {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
