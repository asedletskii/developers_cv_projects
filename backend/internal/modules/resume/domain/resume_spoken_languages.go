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

func (v LanguageGrade) isValid() bool {
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
	SpokenLanguageFR SpokenLanguage = "fr"
	SpokenLanguageDE SpokenLanguage = "de"
	SpokenLanguageKK SpokenLanguage = "kk"
	SpokenLanguageZH SpokenLanguage = "zh"
	SpokenLanguageES SpokenLanguage = "es"
)

func (v SpokenLanguage) isValid() bool {
	switch v {
	case SpokenLanguageRU, SpokenLanguageEN, SpokenLanguageFR, SpokenLanguageDE, SpokenLanguageKK, SpokenLanguageZH, SpokenLanguageES:
		return true
	default:
		return false
	}
}

type SpokenLanguageEntry struct {
	id       sharedvalueobjects.ID
	language SpokenLanguage
	grade    *LanguageGrade
	position int
}
type ResumeSpokenLanguagesParams struct {
	Entries []SpokenLanguageEntryParams
}
type SpokenLanguageEntryParams struct {
	ID       sharedvalueobjects.ID
	Language SpokenLanguage
	Grade    *LanguageGrade
	Position int
}
type ResumeSpokenLanguages struct {
	entries []SpokenLanguageEntry
}

func newSpokenLanguageEntry(params SpokenLanguageEntryParams) (SpokenLanguageEntry, error) {
	entry := SpokenLanguageEntry{id: params.ID, language: params.Language, grade: cloneLanguageGrade(params.Grade), position: params.Position}
	if err := entry.validate(); err != nil {
		return SpokenLanguageEntry{}, err
	}
	return entry, nil
}

func restoreSpokenLanguageEntry(params SpokenLanguageEntryParams) SpokenLanguageEntry {
	return SpokenLanguageEntry{id: params.ID, language: params.Language, grade: cloneLanguageGrade(params.Grade), position: params.Position}
}

func (l SpokenLanguageEntry) validate() error {
	if l.id.IsZero() || !l.language.isValid() || l.position < 0 {
		return fmt.Errorf("%w: invalid language entry", ErrInvalidEntity)
	}
	if l.grade != nil && !l.grade.isValid() {
		return fmt.Errorf("%w: invalid language grade", ErrInvalidEntity)
	}
	return nil
}

func cloneLanguageGrade(value *LanguageGrade) *LanguageGrade {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func newResumeSpokenLanguages(params ResumeSpokenLanguagesParams) (*ResumeSpokenLanguages, error) {
	entries := make([]SpokenLanguageEntry, 0, len(params.Entries))
	for _, value := range params.Entries {
		entries = append(entries, restoreSpokenLanguageEntry(value))
	}
	languages := &ResumeSpokenLanguages{entries: entries}
	if err := languages.validate(); err != nil {
		return nil, err
	}
	return languages, nil
}

func RestoreResumeSpokenLanguages(params ResumeSpokenLanguagesParams) *ResumeSpokenLanguages {
	entries := make([]SpokenLanguageEntry, 0, len(params.Entries))
	for _, value := range params.Entries {
		entries = append(entries, restoreSpokenLanguageEntry(value))
	}
	return &ResumeSpokenLanguages{entries: entries}
}

func (l *ResumeSpokenLanguages) validate() error {
	if l == nil {
		return fmt.Errorf("%w: language collection is nil", ErrInvalidEntity)
	}
	if len(l.entries) > MaxLanguagesPerResume {
		return fmt.Errorf("%w: languages exceed limit of %d", ErrInvalidEntity, MaxLanguagesPerResume)
	}
	seenLanguages := make(map[SpokenLanguage]struct{}, len(l.entries))
	seenIDs := make(map[sharedvalueobjects.ID]struct{}, len(l.entries))
	for index, entry := range l.entries {
		if err := entry.validate(); err != nil {
			return err
		}
		if _, exists := seenLanguages[entry.language]; exists {
			return fmt.Errorf("%w: duplicate language", ErrInvalidEntity)
		}
		if _, exists := seenIDs[entry.id]; exists {
			return fmt.Errorf("%w: duplicate language id", ErrInvalidEntity)
		}
		if entry.position != index {
			return fmt.Errorf("%w: language positions must match collection order", ErrInvalidEntity)
		}
		seenLanguages[entry.language] = struct{}{}
		seenIDs[entry.id] = struct{}{}
	}
	return nil
}

func (l SpokenLanguageEntry) ID() sharedvalueobjects.ID { return l.id }
func (l SpokenLanguageEntry) Language() SpokenLanguage  { return l.language }
func (l SpokenLanguageEntry) Grade() *LanguageGrade {
	if l.grade == nil {
		return nil
	}
	x := *l.grade
	return &x
}
func (l SpokenLanguageEntry) Position() int { return l.position }
func (l *ResumeSpokenLanguages) Entries() []SpokenLanguageEntry {
	entries := append([]SpokenLanguageEntry(nil), l.entries...)
	for i := range entries {
		entries[i].grade = cloneLanguageGrade(entries[i].grade)
	}
	return entries
}

func (r *Resume) ReplaceSpokenLanguages(params ResumeSpokenLanguagesParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	existing := make(map[sharedvalueobjects.ID]struct{})
	if r.spokenLanguages != nil {
		for _, entry := range r.spokenLanguages.Entries() {
			existing[entry.ID()] = struct{}{}
		}
	}
	seen := make(map[sharedvalueobjects.ID]struct{}, len(params.Entries))
	entries := make([]SpokenLanguageEntryParams, 0, len(params.Entries))
	for _, entry := range params.Entries {
		resolvedID, err := resolveCollectionItemID(entry.ID, existing, "language")
		if err != nil {
			return err
		}
		entry.ID = resolvedID
		if _, ok := seen[entry.ID]; ok {
			return fmt.Errorf("%w: duplicate language id", ErrInvalidEntity)
		}
		seen[entry.ID] = struct{}{}
		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		r.spokenLanguages = nil
		return nil
	}
	params.Entries = entries
	languages, err := newResumeSpokenLanguages(params)
	if err != nil {
		return err
	}
	r.spokenLanguages = languages
	return nil
}

func (r *Resume) SpokenLanguages() *ResumeSpokenLanguages {
	if r == nil {
		return nil
	}
	return cloneResumeSpokenLanguages(r.spokenLanguages)
}

func cloneResumeSpokenLanguages(value *ResumeSpokenLanguages) *ResumeSpokenLanguages {
	if value == nil {
		return nil
	}
	copy := *value
	copy.entries = append([]SpokenLanguageEntry(nil), value.entries...)
	for i := range copy.entries {
		copy.entries[i].grade = cloneLanguageGrade(copy.entries[i].grade)
	}
	return &copy
}
