package domain

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ResumeLanguageEntry struct {
	id       valueobjects.ID
	language enums.SpokenLanguage
	grade    *enums.LanguageGrade
	position int
}
type ResumeLanguagesParams struct {
	ID, ResumeID valueobjects.ID
	Entries      []ResumeLanguageEntryParams
}
type ResumeLanguageEntryParams struct {
	ID       valueobjects.ID
	Language enums.SpokenLanguage
	Grade    *enums.LanguageGrade
	Position int
}
type ResumeLanguages struct {
	id, resumeID valueobjects.ID
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
	seen := map[string]bool{}
	for _, x := range p.Entries {
		e, err := NewResumeLanguageEntry(x)
		if err != nil {
			return nil, err
		}
		if seen[string(e.language)] {
			return nil, fmt.Errorf("%w: duplicate language", ErrInvalidEntity)
		}
		seen[string(e.language)] = true
		out = append(out, e)
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

func (l ResumeLanguageEntry) ID() valueobjects.ID            { return l.id }
func (l ResumeLanguageEntry) Language() enums.SpokenLanguage { return l.language }
func (l ResumeLanguageEntry) Grade() *enums.LanguageGrade {
	if l.grade == nil {
		return nil
	}
	x := *l.grade
	return &x
}
func (l ResumeLanguageEntry) Position() int          { return l.position }
func (l *ResumeLanguages) ID() valueobjects.ID       { return l.id }
func (l *ResumeLanguages) ResumeID() valueobjects.ID { return l.resumeID }
func (l *ResumeLanguages) Entries() []ResumeLanguageEntry {
	return append([]ResumeLanguageEntry(nil), l.entries...)
}
