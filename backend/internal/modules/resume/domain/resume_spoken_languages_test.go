package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestSpokenLanguageCollectionAllowsEmptyValue(t *testing.T) {
	_, err := newResumeSpokenLanguages(ResumeSpokenLanguagesParams{})
	require.NoError(t, err)
}

func TestResumeSpokenLanguagesRejectsDuplicateLanguages(t *testing.T) {
	language := SpokenLanguage("en")
	_, err := newResumeSpokenLanguages(ResumeSpokenLanguagesParams{
		Entries: []SpokenLanguageEntryParams{
			{ID: sharedvalueobjects.NewID(), Language: language, Position: 0},
			{ID: sharedvalueobjects.NewID(), Language: language, Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestSpokenLanguageEntryRejectsInvalidPosition(t *testing.T) {
	_, err := newSpokenLanguageEntry(SpokenLanguageEntryParams{
		ID:       sharedvalueobjects.NewID(),
		Language: SpokenLanguage("en"),
		Position: -1,
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeSpokenLanguagesRejectsDuplicateIDs(t *testing.T) {
	id := sharedvalueobjects.NewID()
	_, err := newResumeSpokenLanguages(ResumeSpokenLanguagesParams{
		Entries: []SpokenLanguageEntryParams{
			{ID: id, Language: SpokenLanguageEN, Position: 0},
			{ID: id, Language: SpokenLanguageRU, Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeSpokenLanguagesRejectsDuplicatePositions(t *testing.T) {
	_, err := newResumeSpokenLanguages(ResumeSpokenLanguagesParams{
		Entries: []SpokenLanguageEntryParams{
			{ID: sharedvalueobjects.NewID(), Language: SpokenLanguageEN, Position: 0},
			{ID: sharedvalueobjects.NewID(), Language: SpokenLanguageRU, Position: 0},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceLanguagesCreatesAndReplacesOrderedSnapshot(t *testing.T) {
	resume := newTestResume(t)

	require.NoError(t, resume.ReplaceSpokenLanguages(ResumeSpokenLanguagesParams{Entries: []SpokenLanguageEntryParams{
		{Language: SpokenLanguageEN, Position: 0},
		{Language: SpokenLanguageRU, Position: 1},
	}}))
	first := resume.SpokenLanguages().Entries()
	require.Len(t, first, 2)
	require.NotEqual(t, sharedvalueobjects.ID{}, first[0].ID())

	firstID := first[0].ID()
	secondID := first[1].ID()
	require.NoError(t, resume.ReplaceSpokenLanguages(ResumeSpokenLanguagesParams{Entries: []SpokenLanguageEntryParams{
		{ID: secondID, Language: SpokenLanguageRU, Position: 0},
		{ID: firstID, Language: SpokenLanguageEN, Position: 1},
	}}))

	entries := resume.SpokenLanguages().Entries()
	require.Equal(t, []SpokenLanguage{SpokenLanguageRU, SpokenLanguageEN}, []SpokenLanguage{entries[0].Language(), entries[1].Language()})
	require.Equal(t, secondID, entries[0].ID())
	require.Equal(t, firstID, entries[1].ID())
}

func TestResumeReplaceLanguagesRejectsUnknownAndDuplicateIDs(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceSpokenLanguages(ResumeSpokenLanguagesParams{Entries: []SpokenLanguageEntryParams{
		{Language: SpokenLanguageEN, Position: 0},
	}}))
	known := resume.SpokenLanguages().Entries()[0].ID()

	err := resume.ReplaceSpokenLanguages(ResumeSpokenLanguagesParams{Entries: []SpokenLanguageEntryParams{
		{ID: sharedvalueobjects.NewID(), Language: SpokenLanguageRU, Position: 0},
	}})
	require.ErrorIs(t, err, ErrResumeInvariant)

	err = resume.ReplaceSpokenLanguages(ResumeSpokenLanguagesParams{Entries: []SpokenLanguageEntryParams{
		{ID: known, Language: SpokenLanguageRU, Position: 0},
		{ID: known, Language: SpokenLanguageEN, Position: 1},
	}})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceLanguagesRejectsNonContiguousPositions(t *testing.T) {
	resume := newTestResume(t)
	err := resume.ReplaceSpokenLanguages(ResumeSpokenLanguagesParams{Entries: []SpokenLanguageEntryParams{
		{Language: SpokenLanguageEN, Position: 0},
		{Language: SpokenLanguageRU, Position: 2},
	}})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceLanguagesClearsCollectionWhenEmpty(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceSpokenLanguages(ResumeSpokenLanguagesParams{Entries: []SpokenLanguageEntryParams{
		{Language: SpokenLanguageEN, Position: 0},
	}}))
	require.NotNil(t, resume.SpokenLanguages())
	require.NoError(t, resume.ReplaceSpokenLanguages(ResumeSpokenLanguagesParams{}))
	require.Nil(t, resume.SpokenLanguages())
}

func TestSpokenLanguageCatalogueIsDistinctFromResumeLanguage(t *testing.T) {
	for _, language := range []SpokenLanguage{
		SpokenLanguageRU, SpokenLanguageEN, SpokenLanguageFR, SpokenLanguageDE,
		SpokenLanguageKK, SpokenLanguageZH, SpokenLanguageES,
	} {
		require.True(t, language.isValid(), string(language))
	}
	require.False(t, SpokenLanguage("kz").isValid())
	require.False(t, SpokenLanguage("xx").isValid())
}

func TestResumeSpokenLanguagesRejectsPositionsThatDisagreeWithSliceOrder(t *testing.T) {
	_, err := newResumeSpokenLanguages(ResumeSpokenLanguagesParams{
		Entries: []SpokenLanguageEntryParams{
			{ID: sharedvalueobjects.NewID(), Language: SpokenLanguageEN, Position: 1},
			{ID: sharedvalueobjects.NewID(), Language: SpokenLanguageRU, Position: 0},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}
