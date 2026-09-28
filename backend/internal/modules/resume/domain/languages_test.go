package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestResumeLanguagesRejectsDuplicateLanguages(t *testing.T) {
	language := enums.SpokenLanguage("en")
	_, err := NewResumeLanguages(ResumeLanguagesParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Entries: []ResumeLanguageEntryParams{
			{ID: valueobjects.NewID(), Language: language, Position: 0},
			{ID: valueobjects.NewID(), Language: language, Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeLanguageEntryRejectsInvalidPosition(t *testing.T) {
	_, err := NewResumeLanguageEntry(ResumeLanguageEntryParams{
		ID:       valueobjects.NewID(),
		Language: enums.SpokenLanguage("en"),
		Position: -1,
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeLanguagesRejectsDuplicateIDs(t *testing.T) {
	id := valueobjects.NewID()
	_, err := NewResumeLanguages(ResumeLanguagesParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Entries: []ResumeLanguageEntryParams{
			{ID: id, Language: enums.SpokenLanguageEN, Position: 0},
			{ID: id, Language: enums.SpokenLanguageFR, Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeLanguagesRejectsDuplicatePositions(t *testing.T) {
	_, err := NewResumeLanguages(ResumeLanguagesParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Entries: []ResumeLanguageEntryParams{
			{ID: valueobjects.NewID(), Language: enums.SpokenLanguageEN, Position: 0},
			{ID: valueobjects.NewID(), Language: enums.SpokenLanguageFR, Position: 0},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}
