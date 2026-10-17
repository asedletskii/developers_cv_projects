package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func experienceCommandItem(t *testing.T, company string) ExperienceInput {
	t.Helper()
	start, err := domain.ParseYearMonth("2024-01")
	require.NoError(t, err)
	return ExperienceInput{
		Company: company, Position: "Go Developer", StartDate: start,
		Achievements: []string{"Built API"},
		Tasks:        []string{"Implement backend", "Review PRs"},
		Stack:        []ExperienceStackItemInput{{Name: "Go"}},
	}
}

func TestReplaceExperiencesHandlerUsesVersionAndAggregate(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}

	err := NewReplaceExperiencesHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceExperiencesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 1,
		Input: ExperiencesInput{Items: []ExperienceInput{experienceCommandItem(t, "Acme")}},
	})
	require.NoError(t, err)
	require.True(t, repo.saveCalled)
	require.Equal(t, int64(1), repo.version)
	require.Len(t, repo.resume.Experiences(), 1)
	require.Equal(t, "Acme", repo.resume.Experiences()[0].Company())
}

func TestReplaceExperiencesHandlerRejectsStaleVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}

	err := NewReplaceExperiencesHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceExperiencesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 2,
		Input: ExperiencesInput{Items: []ExperienceInput{experienceCommandItem(t, "Acme")}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saveCalled)
}
