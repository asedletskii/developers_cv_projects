package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestReplaceLanguagesHandlerUsesAggregateAndVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}
	err := NewReplaceLanguagesHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceLanguagesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 1,
		Input: LanguagesInput{Entries: []LanguageInput{{Language: domain.SpokenLanguageEN, Position: 0}}},
	})
	require.NoError(t, err)
	require.True(t, repo.saveCalled)
	require.Equal(t, int64(1), repo.version)
	require.Len(t, repo.resume.SpokenLanguages().Entries(), 1)
}

func TestReplaceLanguagesHandlerRejectsStaleVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}
	err := NewReplaceLanguagesHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceLanguagesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 2,
		Input: LanguagesInput{Entries: []LanguageInput{{Language: domain.SpokenLanguageEN, Position: 0}}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saveCalled)
}
