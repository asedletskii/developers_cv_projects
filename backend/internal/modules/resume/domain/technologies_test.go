package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestResumeTechnologyEntryTrimsName(t *testing.T) {
	entry, err := NewResumeTechnologyEntry(ResumeTechnologyEntryParams{
		ID:       valueobjects.NewID(),
		Name:     "  Go  ",
		Position: 0,
	})
	require.NoError(t, err)
	require.Equal(t, "Go", entry.Name())
}

func TestResumeTechnologyEntryRejectsInvalidGrade(t *testing.T) {
	grade := enums.TechnologyGrade("invalid")
	_, err := NewResumeTechnologyEntry(ResumeTechnologyEntryParams{
		ID:       valueobjects.NewID(),
		Name:     "Go",
		Grade:    &grade,
		Position: 0,
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeTechnologiesRejectsDuplicateIDs(t *testing.T) {
	id := valueobjects.NewID()
	_, err := NewResumeTechnologies(ResumeTechnologiesParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Entries: []ResumeTechnologyEntryParams{
			{ID: id, Name: "Go", Position: 0},
			{ID: id, Name: "Rust", Position: 1},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeTechnologiesRejectsDuplicatePositions(t *testing.T) {
	_, err := NewResumeTechnologies(ResumeTechnologiesParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Entries: []ResumeTechnologyEntryParams{
			{ID: valueobjects.NewID(), Name: "Go", Position: 0},
			{ID: valueobjects.NewID(), Name: "Rust", Position: 0},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceTechnologiesCreatesAndReplacesSnapshot(t *testing.T) {
	resume := newTestResumeForAggregate(t)
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
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{{Name: "Go", Position: 0}},
	}))
	err := resume.ReplaceTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{{ID: valueobjects.NewID(), Name: "Rust", Position: 0}},
	})
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceTechnologiesRejectsNonContiguousPositions(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	err := resume.ReplaceTechnologies(ResumeTechnologiesParams{
		Entries: []ResumeTechnologyEntryParams{
			{Name: "Go", Position: 0},
			{Name: "Rust", Position: 2},
		},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}
