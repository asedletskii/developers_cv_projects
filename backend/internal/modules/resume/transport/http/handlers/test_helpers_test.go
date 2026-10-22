package handlers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/resumefixture"
)

func newHandlerTestResume(t *testing.T) (*domain.Resume, *resumefixture.Repository) {
	t.Helper()
	resume, err := domain.NewResume(domain.ResumeParams{
		UserID: sharedvalueobjects.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN,
	})
	require.NoError(t, err)
	return resume, &resumefixture.Repository{Resume: resume}
}

func newOverviewTestResume(t *testing.T) (*domain.Resume, *resumefixture.Repository) {
	t.Helper()
	resume, repo := newHandlerTestResume(t)
	overview := domain.ResumeOverviewParams{Greeting: mustGreeting(t, domain.GreetingParams{}), Sections: []domain.SectionSettings{}}
	require.NoError(t, resume.ReplaceOverview(overview))
	return resume, repo
}

func mustGreeting(t *testing.T, params domain.GreetingParams) domain.Greeting {
	t.Helper()
	greeting, err := domain.NewGreeting(params)
	require.NoError(t, err)
	return greeting
}
