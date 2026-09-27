package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestNewProjectBuildsOwnedChildren(t *testing.T) {
	id := valueobjects.NewID()
	resumeID := valueobjects.NewID()
	grade := enums.TechnologyGradeSenior
	technology, err := entities.NewTechnology(entities.TechnologyParams{Name: " Go ", Grade: &grade})
	require.NoError(t, err)
	email, err := valueobjects.NewEmail("jane@example.com")
	require.NoError(t, err)
	contacts, err := entities.NewContacts(entities.ContactsParams{Email: &email})
	require.NoError(t, err)
	contributor, err := entities.NewContributor(entities.ContributorParams{Name: " Jane ", Surname: " Doe ", Role: " Developer ", Contacts: contacts})
	require.NoError(t, err)
	screenshot, err := valueobjects.NewExternalURL("https://example.com/assets/project.png")
	require.NoError(t, err)

	project, err := NewProject(ProjectParams{
		ID: id, ResumeID: resumeID, Name: " My project ",
		Screenshots: []valueobjects.ExternalURL{screenshot},
		Stack:       []entities.Technology{technology}, Contributors: []entities.Contributor{contributor},
	})
	require.NoError(t, err)
	require.Equal(t, "My project", project.Name())
	require.Len(t, project.Screenshots(), 1)
	require.Equal(t, "Go", project.Stack()[0].Name())
	require.Equal(t, "Jane", project.Contributors()[0].Name())
}

func TestNewProjectRejectsMissingIdentity(t *testing.T) {
	_, err := NewProject(ProjectParams{Name: "Project", ResumeID: valueobjects.NewID()})
	require.ErrorIs(t, err, ErrInvalidEntity)
}
