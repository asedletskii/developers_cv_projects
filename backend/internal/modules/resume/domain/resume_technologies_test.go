package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestResumeTechnologyEntryTrimsName(t *testing.T) {
	entry, err := newResumeTechnologyEntry(ResumeTechnologyEntryParams{
		ID:       sharedvalueobjects.NewID(),
		Name:     "  Go  ",
		Position: 0,
	})
	require.NoError(t, err)
	require.Equal(t, "Go", entry.Name())
}

func TestResumeTechnologiesRejectsDuplicateIDs(t *testing.T) {
	id := sharedvalueobjects.NewID()
	_, err := newResumeTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{
			{ID: id, Name: "Go", Position: 0},
			{ID: id, Name: "Rust", Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeTechnologiesRejectsDuplicatePositions(t *testing.T) {
	_, err := newResumeTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{
			{ID: sharedvalueobjects.NewID(), Name: "Go", Position: 0},
			{ID: sharedvalueobjects.NewID(), Name: "Rust", Position: 0},
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

func TestResumeReplaceTechnologiesClearsCollectionWhenEmpty(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceTechnologies(ResumeTechnologiesParams{Entries: []ResumeTechnologyEntryParams{
		{Name: "Go", Position: 0},
	}}))
	require.NotNil(t, resume.Technologies())
	require.NoError(t, resume.ReplaceTechnologies(ResumeTechnologiesParams{}))
	require.Nil(t, resume.Technologies())
}

func TestResumeReplaceTechnologiesRejectsUnknownID(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{{Name: "Go", Position: 0}},
	}))
	err := resume.ReplaceTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{{ID: sharedvalueobjects.NewID(), Name: "Rust", Position: 0}},
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

func TestTechnologyCollectionAllowsEmptyValue(t *testing.T) {
	_, err := newResumeTechnologies(ResumeTechnologiesParams{})
	require.NoError(t, err)
}

func TestResumeTechnologiesRejectsPositionsThatDisagreeWithSliceOrder(t *testing.T) {
	_, err := newResumeTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{
			{ID: sharedvalueobjects.NewID(), Name: "Go", Position: 1},
			{ID: sharedvalueobjects.NewID(), Name: "Rust", Position: 0},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeTechnologiesRejectsDuplicateNames(t *testing.T) {
	_, err := newResumeTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{
			{ID: sharedvalueobjects.NewID(), Name: " Go ", Position: 0},
			{ID: sharedvalueobjects.NewID(), Name: "Go", Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrDuplicateEntity)
}
