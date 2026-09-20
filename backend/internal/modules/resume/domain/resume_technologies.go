package domain

import (
	"fmt"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeTechnologyEntry struct {
	id       sharedvalueobjects.ID
	name     string
	position int
}
type ResumeTechnologyEntryParams struct {
	ID       sharedvalueobjects.ID
	Name     string
	Position int
}
type ResumeTechnologiesParams struct {
	Entries []ResumeTechnologyEntryParams
}
type ResumeTechnologies struct {
	entries []ResumeTechnologyEntry
}

func newResumeTechnologyEntry(params ResumeTechnologyEntryParams) (ResumeTechnologyEntry, error) {
	entry := ResumeTechnologyEntry{id: params.ID, name: strings.TrimSpace(params.Name), position: params.Position}
	if err := entry.validate(); err != nil {
		return ResumeTechnologyEntry{}, err
	}
	return entry, nil
}

func (t ResumeTechnologyEntry) validate() error {
	if t.id.IsZero() || t.name == "" || !textLengthWithin(t.name, MaxSkillLength) || t.position < 0 {
		return fmt.Errorf("%w: invalid technology entry", ErrInvalidEntity)
	}
	return nil
}

func newResumeTechnologies(params ResumeTechnologiesParams) (*ResumeTechnologies, error) {
	entries := make([]ResumeTechnologyEntry, 0, len(params.Entries))
	for _, value := range params.Entries {
		entries = append(entries, ResumeTechnologyEntry{
			id: value.ID, name: strings.TrimSpace(value.Name), position: value.Position,
		})
	}
	technologies := &ResumeTechnologies{entries: entries}
	if err := technologies.validate(); err != nil {
		return nil, err
	}
	return technologies, nil
}

func RestoreResumeTechnologies(params ResumeTechnologiesParams) *ResumeTechnologies {
	entries := make([]ResumeTechnologyEntry, 0, len(params.Entries))
	for _, value := range params.Entries {
		entries = append(entries, ResumeTechnologyEntry{id: value.ID, name: value.Name, position: value.Position})
	}
	return &ResumeTechnologies{entries: entries}
}

func (t *ResumeTechnologies) validate() error {
	if t == nil {
		return fmt.Errorf("%w: technology collection is nil", ErrInvalidEntity)
	}
	if len(t.entries) > MaxTechnologiesPerResume {
		return fmt.Errorf("%w: technologies exceed limit of %d", ErrInvalidEntity, MaxTechnologiesPerResume)
	}
	seenIDs := make(map[sharedvalueobjects.ID]struct{}, len(t.entries))
	seenNames := make(map[string]struct{}, len(t.entries))
	for index, entry := range t.entries {
		if err := entry.validate(); err != nil {
			return err
		}
		if _, exists := seenIDs[entry.id]; exists {
			return fmt.Errorf("%w: duplicate technology id", ErrInvalidEntity)
		}
		if entry.position != index {
			return fmt.Errorf("%w: technology positions must match collection order", ErrInvalidEntity)
		}
		if _, exists := seenNames[entry.name]; exists {
			return fmt.Errorf("%w: duplicate technology name", ErrDuplicateEntity)
		}
		seenIDs[entry.id] = struct{}{}
		seenNames[entry.name] = struct{}{}
	}
	return nil
}

func (t ResumeTechnologyEntry) ID() sharedvalueobjects.ID { return t.id }
func (t ResumeTechnologyEntry) Name() string              { return t.name }
func (t ResumeTechnologyEntry) Position() int             { return t.position }
func (t *ResumeTechnologies) Entries() []ResumeTechnologyEntry {
	return append([]ResumeTechnologyEntry(nil), t.entries...)
}

func (r *Resume) ReplaceTechnologies(params ResumeTechnologiesParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
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
		resolvedID, err := resolveCollectionItemID(entry.ID, existing, "technology")
		if err != nil {
			return err
		}
		entry.ID = resolvedID
		if _, ok := seen[entry.ID]; ok {
			return fmt.Errorf("%w: duplicate technology id", ErrInvalidEntity)
		}
		seen[entry.ID] = struct{}{}
		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		r.technologies = nil
		return nil
	}
	params.Entries = entries
	technologies, err := newResumeTechnologies(params)
	if err != nil {
		return err
	}
	r.technologies = technologies
	return nil
}

func (r *Resume) Technologies() *ResumeTechnologies {
	if r == nil {
		return nil
	}
	return cloneResumeTechnologies(r.technologies)
}

type TechnologyParams struct {
	Name string
}

type Technology struct {
	name string
}

func NewTechnology(params TechnologyParams) (Technology, error) {
	technology := Technology{name: strings.TrimSpace(params.Name)}
	if err := technology.validate(); err != nil {
		return Technology{}, err
	}
	return technology, nil
}

func RestoreTechnology(params TechnologyParams) Technology {
	return Technology{name: params.Name}
}

func (t Technology) validate() error {
	if t.name == "" || !textLengthWithin(t.name, MaxSkillLength) {
		return fmt.Errorf("%w: technology name is empty", ErrInvalidValue)
	}
	return nil
}

func (t Technology) Name() string { return t.name }

func cloneResumeTechnologies(value *ResumeTechnologies) *ResumeTechnologies {
	if value == nil {
		return nil
	}
	copy := *value
	copy.entries = append([]ResumeTechnologyEntry(nil), value.entries...)
	return &copy
}
