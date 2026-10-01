package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestNewProjectBuildsOwnedChildren(t *testing.T) {
	id := identity.NewID()
	resumeID := identity.NewID()
	grade := TechnologyGradeSenior
	technology, err := NewTechnology(TechnologyParams{Name: " Go ", Grade: &grade})
	require.NoError(t, err)
	email, err := NewEmail("jane@example.com")
	require.NoError(t, err)
	contacts, err := NewContactMethods(ContactMethodsParams{Email: &email})
	require.NoError(t, err)
	contributor, err := NewContributor(ContributorParams{Name: " Jane ", Surname: " Doe ", Role: " Developer ", Contacts: contacts})
	require.NoError(t, err)
	screenshot, err := NewExternalURL("https://example.com/assets/project.png")
	require.NoError(t, err)

	project, err := NewProject(ProjectParams{
		ID: id, ResumeID: resumeID, Name: " My project ",
		Screenshots: []ExternalURL{screenshot},
		Stack:       []Technology{technology}, Contributors: []Contributor{contributor},
	})
	require.NoError(t, err)
	require.Equal(t, "My project", project.Name())
	require.Len(t, project.Screenshots(), 1)
	require.Equal(t, "Go", project.Stack()[0].Name())
	require.Equal(t, "Jane", project.Contributors()[0].Name())
}

func TestNewProjectRejectsMissingIdentity(t *testing.T) {
	_, err := NewProject(ProjectParams{Name: "Project", ResumeID: identity.NewID()})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceProjectsCreatesNewAndPreservesExistingID(t *testing.T) {
	resume := newTestResume(t)

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

func TestResumeReplaceProjectsRejectsUnknownAndDuplicateIDs(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceProjects([]ProjectParams{{Name: "Existing"}}))
	existing := resume.Projects()[0].ID()

	err := resume.ReplaceProjects([]ProjectParams{{ID: identity.NewID(), Name: "Unknown"}})
	require.ErrorIs(t, err, ErrResumeInvariant)

	err = resume.ReplaceProjects([]ProjectParams{
		{ID: existing, Name: "One"},
		{ID: existing, Name: "Two"},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceProjectsRemovesMissingAndPreservesOrder(t *testing.T) {
	resume := newTestResume(t)
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
