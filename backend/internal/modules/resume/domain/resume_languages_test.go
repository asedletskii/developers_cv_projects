package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestResumeLanguagesRejectsDuplicateLanguages(t *testing.T) {
	language := SpokenLanguage("en")
	_, err := NewResumeLanguages(ResumeLanguagesParams{
		ID:       identity.NewID(),
		ResumeID: identity.NewID(),
		Entries: []ResumeLanguageEntryParams{
			{ID: identity.NewID(), Language: language, Position: 0},
			{ID: identity.NewID(), Language: language, Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeLanguageEntryRejectsInvalidPosition(t *testing.T) {
	_, err := NewResumeLanguageEntry(ResumeLanguageEntryParams{
		ID:       identity.NewID(),
		Language: SpokenLanguage("en"),
		Position: -1,
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeLanguagesRejectsDuplicateIDs(t *testing.T) {
	id := identity.NewID()
	_, err := NewResumeLanguages(ResumeLanguagesParams{
		ID:       identity.NewID(),
		ResumeID: identity.NewID(),
		Entries: []ResumeLanguageEntryParams{
			{ID: id, Language: SpokenLanguageEN, Position: 0},
			{ID: id, Language: SpokenLanguageFR, Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeLanguagesRejectsDuplicatePositions(t *testing.T) {
	_, err := NewResumeLanguages(ResumeLanguagesParams{
		ID:       identity.NewID(),
		ResumeID: identity.NewID(),
		Entries: []ResumeLanguageEntryParams{
			{ID: identity.NewID(), Language: SpokenLanguageEN, Position: 0},
			{ID: identity.NewID(), Language: SpokenLanguageFR, Position: 0},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceLanguagesCreatesAndReplacesOrderedSnapshot(t *testing.T) {
	resume := newTestResume(t)

	require.NoError(t, resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{Language: SpokenLanguageEN, Position: 0},
		{Language: SpokenLanguageFR, Position: 1},
	}}))
	first := resume.Languages().Entries()
	require.Len(t, first, 2)
	require.NotEqual(t, identity.ID{}, first[0].ID())

	firstID := first[0].ID()
	secondID := first[1].ID()
	require.NoError(t, resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{ID: secondID, Language: SpokenLanguageDE, Position: 0},
		{ID: firstID, Language: SpokenLanguageEN, Position: 1},
	}}))

	entries := resume.Languages().Entries()
	require.Equal(t, []SpokenLanguage{SpokenLanguageDE, SpokenLanguageEN}, []SpokenLanguage{entries[0].Language(), entries[1].Language()})
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
		{ID: identity.NewID(), Language: SpokenLanguageFR, Position: 0},
	}})
	require.ErrorIs(t, err, ErrResumeInvariant)

	err = resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{ID: known, Language: SpokenLanguageFR, Position: 0},
		{ID: known, Language: SpokenLanguageDE, Position: 1},
	}})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceLanguagesRejectsNonContiguousPositions(t *testing.T) {
	resume := newTestResume(t)
	err := resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{Language: SpokenLanguageEN, Position: 0},
		{Language: SpokenLanguageFR, Position: 2},
	}})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceLanguagesRejectsCollectionIdentityFromAnotherState(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceLanguages(ResumeLanguagesParams{Entries: []ResumeLanguageEntryParams{
		{Language: SpokenLanguageEN, Position: 0},
	}}))
	collectionID := resume.Languages().ID()
	err := resume.ReplaceLanguages(ResumeLanguagesParams{ID: identity.NewID(), Entries: []ResumeLanguageEntryParams{
		{ID: resume.Languages().Entries()[0].ID(), Language: SpokenLanguageFR, Position: 0},
	}})
	require.ErrorIs(t, err, ErrResumeInvariant)
	require.Equal(t, collectionID, resume.Languages().ID())
}

func TestResumeLanguageEntryReplacePreservesID(t *testing.T) {
	id := identity.NewID()
	entry, err := NewResumeLanguageEntry(ResumeLanguageEntryParams{ID: id, Language: SpokenLanguageEN, Position: 0})
	require.NoError(t, err)
	err = entry.Replace(ResumeLanguageEntryParams{Language: SpokenLanguageFR, Position: 2})
	require.NoError(t, err)
	require.Equal(t, id, entry.ID())
	require.Equal(t, SpokenLanguageFR, entry.Language())
	require.Equal(t, 2, entry.Position())
}
