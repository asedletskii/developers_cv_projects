package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestResumeReplaceLanguagesCreatesAndReplacesOrderedSnapshot(t *testing.T) {
	resume := newTestResumeForAggregate(t)

	require.NoError(t, resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{Language: enums.SpokenLanguageEN, Position: 0},
		{Language: enums.SpokenLanguageFR, Position: 1},
	}}))
	first := resume.Languages().Entries()
	require.Len(t, first, 2)
	require.NotEqual(t, valueobjects.ID{}, first[0].ID())

	firstID := first[0].ID()
	secondID := first[1].ID()
	require.NoError(t, resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{ID: secondID, Language: enums.SpokenLanguageDE, Position: 0},
		{ID: firstID, Language: enums.SpokenLanguageEN, Position: 1},
	}}))

	entries := resume.Languages().Entries()
	require.Equal(t, []enums.SpokenLanguage{enums.SpokenLanguageDE, enums.SpokenLanguageEN}, []enums.SpokenLanguage{entries[0].Language(), entries[1].Language()})
	require.Equal(t, secondID, entries[0].ID())
	require.Equal(t, firstID, entries[1].ID())
}

func TestResumeReplaceLanguagesRejectsUnknownAndDuplicateIDs(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{Language: enums.SpokenLanguageEN, Position: 0},
	}}))
	known := resume.Languages().Entries()[0].ID()

	err := resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{ID: valueobjects.NewID(), Language: enums.SpokenLanguageFR, Position: 0},
	}})
	require.ErrorIs(t, err, ErrResumeInvariant)

	err = resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{ID: known, Language: enums.SpokenLanguageFR, Position: 0},
		{ID: known, Language: enums.SpokenLanguageDE, Position: 1},
	}})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceLanguagesRejectsNonContiguousPositions(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	err := resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{Language: enums.SpokenLanguageEN, Position: 0},
		{Language: enums.SpokenLanguageFR, Position: 2},
	}})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceLanguagesRejectsCollectionIdentityFromAnotherState(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{Language: enums.SpokenLanguageEN, Position: 0},
	}}))
	collectionID := resume.Languages().ID()
	err := resume.ReplaceLanguages(ResumeLanguagesParams{ID: valueobjects.NewID(), Entries: []ResumeLanguageEntryParams{
		{ID: resume.Languages().Entries()[0].ID(), Language: enums.SpokenLanguageFR, Position: 0},
	}})
	require.ErrorIs(t, err, ErrResumeInvariant)
	require.Equal(t, collectionID, resume.Languages().ID())
}
