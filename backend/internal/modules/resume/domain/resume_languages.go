package domain

import (
	"fmt"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LanguageGrade string

const (
	LanguageGradeA1 LanguageGrade = "A1"
	LanguageGradeA2 LanguageGrade = "A2"
	LanguageGradeB1 LanguageGrade = "B1"
	LanguageGradeB2 LanguageGrade = "B2"
	LanguageGradeC1 LanguageGrade = "C1"
	LanguageGradeC2 LanguageGrade = "C2"
)

func (v LanguageGrade) IsValid() bool {
	switch v {
	case LanguageGradeA1, LanguageGradeA2, LanguageGradeB1, LanguageGradeB2, LanguageGradeC1, LanguageGradeC2:
		return true
	default:
		return false
	}
}

type SpokenLanguage string

const (
	SpokenLanguageRU SpokenLanguage = "ru"
	SpokenLanguageEN SpokenLanguage = "en"
)

func (v SpokenLanguage) IsValid() bool {
	switch v {
	case SpokenLanguageRU, SpokenLanguageEN:
		return true
	default:
		return false
	}
}

type ResumeLanguageEntry struct {
	id       sharedvalueobjects.ID
	language SpokenLanguage
	grade    *LanguageGrade
	position int
}
type ResumeLanguagesParams struct {
	ID, ResumeID sharedvalueobjects.ID
	Entries      []ResumeLanguageEntryParams
}
type ResumeLanguageEntryParams struct {
	ID       sharedvalueobjects.ID
	Language SpokenLanguage
	Grade    *LanguageGrade
	Position int
}
type ResumeLanguages struct {
	id, resumeID sharedvalueobjects.ID
	entries      []ResumeLanguageEntry
}

func NewResumeLanguageEntry(p ResumeLanguageEntryParams) (ResumeLanguageEntry, error) {
	if p.ID.IsZero() || !p.Language.IsValid() || p.Position < 0 {
		return ResumeLanguageEntry{}, fmt.Errorf("%w: invalid language entry", ErrInvalidEntity)
	}
	if p.Grade != nil && !p.Grade.IsValid() {
		return ResumeLanguageEntry{}, fmt.Errorf("%w: invalid language grade", ErrInvalidEntity)
	}
	g := p.Grade
	if g != nil {
		x := *g
		g = &x
	}
	return ResumeLanguageEntry{id: p.ID, language: p.Language, grade: g, position: p.Position}, nil
}
func NewResumeLanguages(p ResumeLanguagesParams) (*ResumeLanguages, error) {
	if p.ID.IsZero() || p.ResumeID.IsZero() {
		return nil, fmt.Errorf("%w: language collection ids are required", ErrInvalidEntity)
	}
	out := make([]ResumeLanguageEntry, 0, len(p.Entries))
	seenLanguages := make(map[string]struct{}, len(p.Entries))
	seenIDs := make(map[sharedvalueobjects.ID]struct{}, len(p.Entries))
	seenPositions := make(map[int]struct{}, len(p.Entries))
	for _, x := range p.Entries {
		e, err := NewResumeLanguageEntry(x)
		if err != nil {
			return nil, err
		}
		if _, ok := seenLanguages[string(e.language)]; ok {
			return nil, fmt.Errorf("%w: duplicate language", ErrInvalidEntity)
		}
		if _, ok := seenIDs[e.id]; ok {
			return nil, fmt.Errorf("%w: duplicate language id", ErrInvalidEntity)
		}
		if _, ok := seenPositions[e.position]; ok {
			return nil, fmt.Errorf("%w: duplicate language position", ErrInvalidEntity)
		}
		seenLanguages[string(e.language)] = struct{}{}
		seenIDs[e.id] = struct{}{}
		seenPositions[e.position] = struct{}{}
		out = append(out, e)
	}
	for position := range seenPositions {
		if position < 0 || position >= len(out) {
			return nil, fmt.Errorf("%w: language positions must be contiguous", ErrInvalidEntity)
		}
	}
	return &ResumeLanguages{id: p.ID, resumeID: p.ResumeID, entries: out}, nil
}
func (l *ResumeLanguageEntry) Replace(params ResumeLanguageEntryParams) error {
	params.ID = l.id
	updated, err := NewResumeLanguageEntry(params)
	if err != nil {
		return err
	}
	*l = updated
	return nil
}

func (l ResumeLanguageEntry) ID() sharedvalueobjects.ID { return l.id }
func (l ResumeLanguageEntry) Language() SpokenLanguage  { return l.language }
func (l ResumeLanguageEntry) Grade() *LanguageGrade {
	if l.grade == nil {
		return nil
	}
	x := *l.grade
	return &x
}
func (l ResumeLanguageEntry) Position() int                { return l.position }
func (l *ResumeLanguages) ID() sharedvalueobjects.ID       { return l.id }
func (l *ResumeLanguages) ResumeID() sharedvalueobjects.ID { return l.resumeID }
func (l *ResumeLanguages) Entries() []ResumeLanguageEntry {
	return append([]ResumeLanguageEntry(nil), l.entries...)
}

func (r *Resume) ReplaceLanguages(params ResumeLanguagesParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if params.ID.IsZero() {
		if r.languages != nil {
			params.ID = r.languages.ID()
		} else {
			params.ID = sharedvalueobjects.NewID()
		}
	} else if r.languages == nil || params.ID != r.languages.ID() {
		return fmt.Errorf("%w: languages do not belong to resume", ErrResumeInvariant)
	}

	existing := make(map[sharedvalueobjects.ID]struct{})
	if r.languages != nil {
		for _, entry := range r.languages.Entries() {
			existing[entry.ID()] = struct{}{}
		}
	}
	seen := make(map[sharedvalueobjects.ID]struct{}, len(params.Entries))
	entries := make([]ResumeLanguageEntryParams, 0, len(params.Entries))
	for _, entry := range params.Entries {
		if entry.ID.IsZero() {
			entry.ID = sharedvalueobjects.NewID()
		} else {
			if _, ok := existing[entry.ID]; !ok {
				return fmt.Errorf("%w: language does not belong to resume", ErrResumeInvariant)
			}
		}
		if _, ok := seen[entry.ID]; ok {
			return fmt.Errorf("%w: duplicate language id", ErrInvalidEntity)
		}
		seen[entry.ID] = struct{}{}
		entries = append(entries, entry)
	}

	params.ResumeID = r.id
	params.Entries = entries
	languages, err := NewResumeLanguages(params)
	if err != nil {
		return err
	}
	r.languages = languages
	return nil
}

func (r *Resume) Languages() *ResumeLanguages {
	if r == nil || r.languages == nil {
		return nil
	}
	copy := *r.languages
	copy.entries = append([]ResumeLanguageEntry(nil), r.languages.entries...)
	return &copy
}
