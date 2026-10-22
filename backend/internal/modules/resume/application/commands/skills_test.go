package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestReplaceSkillsHandlerUsesVersionAndAggregate(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{}
	err := NewReplaceSkillsHandler(repo, uow).Handle(context.Background(), ReplaceSkillsCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 1,
		Input: SkillsInput{Hard: []string{"testing"}, Soft: []string{"teamwork"}},
	})
	require.NoError(t, err)
	require.True(t, repo.saveCalled)
	require.Equal(t, int64(1), repo.version)
	require.NotNil(t, repo.resume.Skills())
	require.Equal(t, []string{"testing"}, repo.resume.Skills().HardSkills())
}

func TestReplaceSkillsHandlerRejectsStaleVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{}
	err := NewReplaceSkillsHandler(repo, uow).Handle(context.Background(), ReplaceSkillsCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 2,
		Input: SkillsInput{Hard: []string{"testing"}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saveCalled)
}
