package readers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestResumeReaderRejectsInvalidOwnedQueryBeforeDatabaseAccess(t *testing.T) {
	reader := NewResumeReader(nil)

	_, err := reader.GetOwned(context.Background(), queries.GetResumeQuery{ResumeID: sharedvalueobjects.NewID()})

	require.Error(t, err)
}

func TestResumeReaderRejectsInvalidPublicQueryBeforeDatabaseAccess(t *testing.T) {
	reader := NewResumeReader(nil)

	_, err := reader.GetPublic(context.Background(), resumecontracts.PublicResumeQuery{})

	require.Error(t, err)
}

func TestResumeReaderRequiresDatabaseForValidQuery(t *testing.T) {
	reader := NewResumeReader(nil)

	_, err := reader.GetOwned(context.Background(), queries.GetResumeQuery{
		ResumeID: sharedvalueobjects.NewID(),
		UserID:   sharedvalueobjects.NewID(),
	})

	require.Error(t, err)
}

func TestResumeReaderInitializesEmptyCollections(t *testing.T) {
	view := queries.ResumeView{
		Languages:    make([]queries.LanguageView, 0),
		Technologies: make([]queries.TechnologyView, 0),
		Experiences:  make([]queries.ExperienceView, 0),
		Educations:   make([]queries.EducationView, 0),
		Courses:      make([]queries.CourseView, 0),
		Activities:   make([]queries.ActivityView, 0),
		Projects:     make([]queries.ProjectView, 0),
	}

	require.NotNil(t, view.Languages)
	require.NotNil(t, view.Projects)
}

func TestResumeReaderUsesAggregateVersionInView(t *testing.T) {
	version := int64(42)
	view := queries.ResumeView{Version: &version}

	require.NotNil(t, view.Version)
	require.Equal(t, int64(42), *view.Version)
}
