package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestResumeLanguageEntryReplacePreservesID(t *testing.T) {
	id := valueobjects.NewID()
	entry, err := NewResumeLanguageEntry(ResumeLanguageEntryParams{ID: id, Language: enums.SpokenLanguageEN, Position: 0})
	require.NoError(t, err)
	err = entry.Replace(ResumeLanguageEntryParams{Language: enums.SpokenLanguageFR, Position: 2})
	require.NoError(t, err)
	require.Equal(t, id, entry.ID())
	require.Equal(t, enums.SpokenLanguageFR, entry.Language())
	require.Equal(t, 2, entry.Position())
}

func TestResumeTechnologiesReplaceAndCloneEntries(t *testing.T) {
	id, resumeID := valueobjects.NewID(), valueobjects.NewID()
	entryID := valueobjects.NewID()
	collection, err := NewResumeTechnologies(ResumeTechnologiesParams{ID: id, ResumeID: resumeID, Entries: []ResumeTechnologyEntryParams{{ID: entryID, Name: " Go ", Position: 0}}})
	require.NoError(t, err)
	require.Equal(t, "Go", collection.Entries()[0].Name())

	entry := collection.Entries()[0]
	err = entry.Replace(ResumeTechnologyEntryParams{Name: "Rust", Position: 1})
	require.NoError(t, err)
	require.Equal(t, entryID, entry.ID())
	require.Equal(t, "Rust", entry.Name())
}

func TestResumeLanguagesAndTechnologiesRejectInvalidCollectionIdentity(t *testing.T) {
	_, err := NewResumeLanguages(ResumeLanguagesParams{ResumeID: valueobjects.NewID()})
	require.Error(t, err)
	_, err = NewResumeTechnologies(ResumeTechnologiesParams{ID: valueobjects.NewID()})
	require.Error(t, err)
}
