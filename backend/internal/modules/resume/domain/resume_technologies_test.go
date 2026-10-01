package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestResumeTechnologyEntryTrimsName(t *testing.T) {
	entry, err := NewResumeTechnologyEntry(ResumeTechnologyEntryParams{
		ID:       identity.NewID(),
		Name:     "  Go  ",
		Position: 0,
	})
	require.NoError(t, err)
	require.Equal(t, "Go", entry.Name())
}

func TestResumeTechnologyEntryRejectsInvalidGrade(t *testing.T) {
	grade := TechnologyGrade("invalid")
	_, err := NewResumeTechnologyEntry(ResumeTechnologyEntryParams{
		ID:       identity.NewID(),
		Name:     "Go",
		Grade:    &grade,
		Position: 0,
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeTechnologiesRejectsDuplicateIDs(t *testing.T) {
	id := identity.NewID()
	_, err := NewResumeTechnologies(ResumeTechnologiesParams{
		ID:       identity.NewID(),
		ResumeID: identity.NewID(),
		Entries: []ResumeTechnologyEntryParams{
			{ID: id, Name: "Go", Position: 0},
			{ID: id, Name: "Rust", Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeTechnologiesRejectsDuplicatePositions(t *testing.T) {
	_, err := NewResumeTechnologies(ResumeTechnologiesParams{
		ID:       identity.NewID(),
		ResumeID: identity.NewID(),
		Entries: []ResumeTechnologyEntryParams{
			{ID: identity.NewID(), Name: "Go", Position: 0},
			{ID: identity.NewID(), Name: "Rust", Position: 0},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceTechnologiesCreatesAndReplacesSnapshot(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{
			{Name: " Go ", Position: 0},
			{Name: "Rust", Position: 1},
		},
	}))
	first := resume.Technologies().Entries()
	require.Len(t, first, 2)
	firstID := first[0].ID()

	require.NoError(t, resume.ReplaceTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{
			{ID: firstID, Name: "Golang", Position: 0},
		},
	}))
	current := resume.Technologies().Entries()
	require.Len(t, current, 1)
	require.Equal(t, firstID, current[0].ID())
	require.Equal(t, "Golang", current[0].Name())
}

func TestResumeReplaceTechnologiesRejectsUnknownID(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{{Name: "Go", Position: 0}},
	}))
	err := resume.ReplaceTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{{ID: identity.NewID(), Name: "Rust", Position: 0}},
	})
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceTechnologiesRejectsNonContiguousPositions(t *testing.T) {
	resume := newTestResume(t)
	err := resume.ReplaceTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{
			{Name: "Go", Position: 0},
			{Name: "Rust", Position: 2},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeTechnologiesReplaceAndCloneEntries(t *testing.T) {
	id, resumeID := identity.NewID(), identity.NewID()
	entryID := identity.NewID()
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
	_, err := NewResumeLanguages(ResumeLanguagesParams{ResumeID: identity.NewID()})
	require.Error(t, err)
	_, err = NewResumeTechnologies(ResumeTechnologiesParams{ID: identity.NewID()})
	require.Error(t, err)
}
