package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestReplaceOverviewHandlerLoadsReplacesAndSavesAggregate(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "Resume", 1)
	overview := domain.ResumeOverviewParams{Greeting: mustGreeting(t, domain.GreetingParams{First: "Old"}), Sections: []domain.SectionSettings{}}
	require.NoError(t, resume.ReplaceOverview(overview))
	repo := &commandRepositoryStub{resume: resume}
	h := NewReplaceOverviewHandler(repo, &unitOfWorkStub{})
	err := h.Handle(context.Background(), ReplaceOverviewCommand{ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1, Input: OverviewInput{Greeting: GreetingInput{First: "New"}, Sections: []SectionInput{}}})
	require.NoError(t, err)
	require.True(t, repo.saveCalled)
	require.Equal(t, int64(1), repo.version)
	require.Equal(t, "New", repo.resume.Overview().Greeting().First())
}

func TestReplaceOverviewHandlerRejectsStaleVersionBeforeSave(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "Resume", 1)
	overview := domain.ResumeOverviewParams{Greeting: mustGreeting(t, domain.GreetingParams{}), Sections: []domain.SectionSettings{}}
	require.NoError(t, resume.ReplaceOverview(overview))
	repo := &commandRepositoryStub{resume: resume}
	h := NewReplaceOverviewHandler(repo, &unitOfWorkStub{})
	err := h.Handle(context.Background(), ReplaceOverviewCommand{ResumeID: resume.ID(), UserID: resume.UserID(), Version: 2, Input: OverviewInput{Sections: []SectionInput{}}})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saveCalled)
}
