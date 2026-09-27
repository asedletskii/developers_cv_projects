package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestResumeTechnologyEntryTrimsName(t *testing.T) {
	entry, err := domain.NewResumeTechnologyEntry(domain.ResumeTechnologyEntryParams{
		ID:       valueobjects.NewID(),
		Name:     "  Go  ",
		Position: 0,
	})
	require.NoError(t, err)
	require.Equal(t, "Go", entry.Name())
}

func TestResumeTechnologyEntryRejectsInvalidGrade(t *testing.T) {
	grade := enums.TechnologyGrade("invalid")
	_, err := domain.NewResumeTechnologyEntry(domain.ResumeTechnologyEntryParams{
		ID:       valueobjects.NewID(),
		Name:     "Go",
		Grade:    &grade,
		Position: 0,
	})
	require.ErrorIs(t, err, domain.ErrInvalidEntity)
}

func TestResumeTechnologiesRejectsDuplicateIDs(t *testing.T) {
	id := valueobjects.NewID()
	_, err := domain.NewResumeTechnologies(domain.ResumeTechnologiesParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Entries: []domain.ResumeTechnologyEntryParams{
			{ID: id, Name: "Go", Position: 0},
			{ID: id, Name: "Rust", Position: 1},
		},
	})
	require.ErrorIs(t, err, domain.ErrInvalidEntity)
}

func TestResumeTechnologiesRejectsDuplicatePositions(t *testing.T) {
	_, err := domain.NewResumeTechnologies(domain.ResumeTechnologiesParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Entries: []domain.ResumeTechnologyEntryParams{
			{ID: valueobjects.NewID(), Name: "Go", Position: 0},
			{ID: valueobjects.NewID(), Name: "Rust", Position: 0},
		},
	})
	require.ErrorIs(t, err, domain.ErrInvalidEntity)
}
