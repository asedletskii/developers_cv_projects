package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestResumeLanguagesRejectsDuplicateLanguages(t *testing.T) {
	language := SpokenLanguage("en")
	_, err := NewResumeLanguages(ResumeLanguagesParams{
		ID:       sharedvalueobjects.NewID(),
		ResumeID: sharedvalueobjects.NewID(),
		Entries: []ResumeLanguageEntryParams{
			{ID: sharedvalueobjects.NewID(), Language: language, Position: 0},
			{ID: sharedvalueobjects.NewID(), Language: language, Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeLanguageEntryRejectsInvalidPosition(t *testing.T) {
	_, err := NewResumeLanguageEntry(ResumeLanguageEntryParams{
		ID:       sharedvalueobjects.NewID(),
		Language: SpokenLanguage("en"),
		Position: -1,
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeLanguagesRejectsDuplicateIDs(t *testing.T) {
	id := sharedvalueobjects.NewID()
	_, err := NewResumeLanguages(ResumeLanguagesParams{
		ID:       sharedvalueobjects.NewID(),
		ResumeID: sharedvalueobjects.NewID(),
		Entries: []ResumeLanguageEntryParams{
			{ID: id, Language: SpokenLanguageEN, Position: 0},
			{ID: id, Language: SpokenLanguageRU, Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeLanguagesRejectsDuplicatePositions(t *testing.T) {
	_, err := NewResumeLanguages(ResumeLanguagesParams{
		ID:       sharedvalueobjects.NewID(),
		ResumeID: sharedvalueobjects.NewID(),
		Entries: []ResumeLanguageEntryParams{
			{ID: sharedvalueobjects.NewID(), Language: SpokenLanguageEN, Position: 0},
			{ID: sharedvalueobjects.NewID(), Language: SpokenLanguageRU, Position: 0},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceLanguagesCreatesAndReplacesOrderedSnapshot(t *testing.T) {
	resume := newTestResume(t)

	require.NoError(t, resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{Language: SpokenLanguageEN, Position: 0},
		{Language: SpokenLanguageRU, Position: 1},
	}}))
	first := resume.Languages().Entries()
	require.Len(t, first, 2)
	require.NotEqual(t, sharedvalueobjects.ID{}, first[0].ID())

	firstID := first[0].ID()
	secondID := first[1].ID()
	require.NoError(t, resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{ID: secondID, Language: SpokenLanguageRU, Position: 0},
		{ID: firstID, Language: SpokenLanguageEN, Position: 1},
	}}))

	entries := resume.Languages().Entries()
	require.Equal(t, []SpokenLanguage{SpokenLanguageRU, SpokenLanguageEN}, []SpokenLanguage{entries[0].Language(), entries[1].Language()})
	require.Equal(t, secondID, entries[0].ID())
	require.Equal(t, firstID, entries[1].ID())
}

func TestResumeReplaceLanguagesRejectsUnknownAndDuplicateIDs(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{Language: SpokenLanguageEN, Position: 0},
	}}))
	known := resume.Languages().Entries()[0].ID()

	err := resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{ID: sharedvalueobjects.NewID(), Language: SpokenLanguageRU, Position: 0},
	}})
	require.ErrorIs(t, err, ErrResumeInvariant)

	err = resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{ID: known, Language: SpokenLanguageRU, Position: 0},
		{ID: known, Language: SpokenLanguageEN, Position: 1},
	}})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceLanguagesRejectsNonContiguousPositions(t *testing.T) {
	resume := newTestResume(t)
	err := resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{Language: SpokenLanguageEN, Position: 0},
		{Language: SpokenLanguageRU, Position: 2},
	}})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceLanguagesRejectsCollectionIdentityFromAnotherState(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{Language: SpokenLanguageEN, Position: 0},
	}}))
	collectionID := resume.Languages().ID()
	err := resume.ReplaceLanguages(ResumeLanguagesParams{ID: sharedvalueobjects.NewID(), Entries: []ResumeLanguageEntryParams{
		{ID: resume.Languages().Entries()[0].ID(), Language: SpokenLanguageRU, Position: 0},
	}})
	require.ErrorIs(t, err, ErrResumeInvariant)
	require.Equal(t, collectionID, resume.Languages().ID())
}

func TestResumeLanguageEntryReplacePreservesID(t *testing.T) {
	id := sharedvalueobjects.NewID()
	entry, err := NewResumeLanguageEntry(ResumeLanguageEntryParams{ID: id, Language: SpokenLanguageEN, Position: 0})
	require.NoError(t, err)
	err = entry.Replace(ResumeLanguageEntryParams{Language: SpokenLanguageRU, Position: 2})
	require.NoError(t, err)
	require.Equal(t, id, entry.ID())
	require.Equal(t, SpokenLanguageRU, entry.Language())
	require.Equal(t, 2, entry.Position())
}
