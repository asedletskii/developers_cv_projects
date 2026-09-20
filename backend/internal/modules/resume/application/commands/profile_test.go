package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestReplaceProfileHandlerUsesVersionAndAggregate(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{}
	photo := testURL(t)
	birth := testDate(t)
	err := NewReplaceProfileHandler(repo, uow).Handle(context.Background(), ReplaceProfileCommand{ResumeID: resume.ID(), UserID: uid, Version: 1, Input: ProfileInput{Name: "John", Surname: "Doe", Position: "Developer", PhotoURL: photo, DateOfBirth: birth, Location: "Paris", WorkFormats: []domain.WorkFormat{domain.WorkFormatRemote}}})
	require.NoError(t, err)
	require.True(t, repo.saveCalled)
	require.Equal(t, int64(1), repo.version)
	require.NotNil(t, repo.resume.Profile())
}
func TestReplaceProfileHandlerRejectsStaleVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{}
	err := NewReplaceProfileHandler(repo, uow).Handle(context.Background(), ReplaceProfileCommand{ResumeID: resume.ID(), UserID: uid, Version: 2})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saveCalled)
}
