package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ResumeTechnologyEntry struct {
	id       valueobjects.ID
	name     string
	grade    *enums.TechnologyGrade
	position int
}
type ResumeTechnologyEntryParams struct {
	ID       valueobjects.ID
	Name     string
	Grade    *enums.TechnologyGrade
	Position int
}
type ResumeTechnologiesParams struct {
	ID, ResumeID valueobjects.ID
	Entries      []ResumeTechnologyEntryParams
}
type ResumeTechnologies struct {
	id, resumeID valueobjects.ID
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
	seenIDs := make(map[valueobjects.ID]struct{}, len(p.Entries))
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

func (t ResumeTechnologyEntry) ID() valueobjects.ID { return t.id }
func (t ResumeTechnologyEntry) Name() string        { return t.name }
func (t ResumeTechnologyEntry) Grade() *enums.TechnologyGrade {
	if t.grade == nil {
		return nil
	}
	x := *t.grade
	return &x
}
func (t ResumeTechnologyEntry) Position() int           { return t.position }
func (t *ResumeTechnologies) ID() valueobjects.ID       { return t.id }
func (t *ResumeTechnologies) ResumeID() valueobjects.ID { return t.resumeID }
func (t *ResumeTechnologies) Entries() []ResumeTechnologyEntry {
	return append([]ResumeTechnologyEntry(nil), t.entries...)
}
