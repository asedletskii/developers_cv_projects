package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func newProjectTestResume(t *testing.T) *Resume {
	t.Helper()
	resume, err := NewResume(ResumeParams{
		ID: valueobjects.NewID(), UserID: useridentity.NewID(),
		Title: "Resume", Language: enums.ResumeLanguageEN,
	})
	require.NoError(t, err)
	return resume
}

func TestReplaceProjectsCreatesNewAndPreservesExistingID(t *testing.T) {
	resume := newProjectTestResume(t)

	require.NoError(t, resume.ReplaceProjects([]ProjectParams{{Name: "First"}}))
	first := resume.Projects()[0]
	require.False(t, first.ID().IsZero())

	require.NoError(t, resume.ReplaceProjects([]ProjectParams{
		{ID: first.ID(), Name: "Updated"},
		{Name: "Second"},
	}))
	items := resume.Projects()
	require.Len(t, items, 2)
	require.Equal(t, first.ID(), items[0].ID())
	require.Equal(t, "Updated", items[0].Name())
	require.False(t, items[1].ID().IsZero())
}

func TestReplaceProjectsRejectsUnknownAndDuplicateIDs(t *testing.T) {
	resume := newProjectTestResume(t)
	require.NoError(t, resume.ReplaceProjects([]ProjectParams{{Name: "Existing"}}))
	existing := resume.Projects()[0].ID()

	err := resume.ReplaceProjects([]ProjectParams{{ID: valueobjects.NewID(), Name: "Unknown"}})
	require.ErrorIs(t, err, ErrResumeInvariant)

	err = resume.ReplaceProjects([]ProjectParams{
		{ID: existing, Name: "One"},
		{ID: existing, Name: "Two"},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestReplaceProjectsRemovesMissingAndPreservesOrder(t *testing.T) {
	resume := newProjectTestResume(t)
	require.NoError(t, resume.ReplaceProjects([]ProjectParams{{Name: "First"}, {Name: "Second"}}))
	first, second := resume.Projects()[0], resume.Projects()[1]

	require.NoError(t, resume.ReplaceProjects([]ProjectParams{
		{ID: second.ID(), Name: "Second updated"},
	}))
	items := resume.Projects()
	require.Len(t, items, 1)
	require.Equal(t, second.ID(), items[0].ID())
	require.Equal(t, "Second updated", items[0].Name())
	require.NotEqual(t, first.ID(), items[0].ID())
}
