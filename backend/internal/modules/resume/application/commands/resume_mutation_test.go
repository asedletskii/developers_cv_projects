package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestExecuteResumeMutationRejectsMissingDependencies(t *testing.T) {
	resume := newMutationTestResume(t)
	uow := &unitOfWorkStub{}

	err := executeResumeMutation(context.Background(), nil, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", ports.ResumeMutationTitle, func(*domain.Resume) error {
		return nil
	})
	require.EqualError(t, err, "replace title handler dependencies are nil")
	require.False(t, uow.called)

	repo := &commandRepositoryStub{resume: resume}
	err = executeResumeMutation(context.Background(), repo, nil, resume.ID(), resume.UserID(), resume.Version(), "replace title", ports.ResumeMutationTitle, func(*domain.Resume) error {
		return nil
	})
	require.EqualError(t, err, "replace title handler dependencies are nil")
	require.False(t, repo.getCalled)
}

func TestExecuteResumeMutationRejectsInvalidCommandBeforeTransaction(t *testing.T) {
	resume := newMutationTestResume(t)
	repo := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{}

	cases := []struct {
		name     string
		resumeID sharedvalueobjects.ID
		userID   sharedvalueobjects.ID
		version  int64
	}{
		{name: "missing resume id", userID: resume.UserID(), version: 3},
		{name: "missing user id", resumeID: resume.ID(), version: 3},
		{name: "zero version", resumeID: resume.ID(), userID: resume.UserID()},
		{name: "negative version", resumeID: resume.ID(), userID: resume.UserID(), version: -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := executeResumeMutation(context.Background(), repo, uow, tc.resumeID, tc.userID, tc.version, "replace title", ports.ResumeMutationTitle, func(*domain.Resume) error {
				t.Fatal("mutation must not run")
				return nil
			})
			require.Error(t, err)
		})
	}

	require.False(t, uow.called)
	require.False(t, repo.getCalled)
}

func TestExecuteResumeMutationPropagatesUnitOfWorkError(t *testing.T) {
	resume := newMutationTestResume(t)
	repositoryErr := errors.New("transaction failure")
	repo := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{err: repositoryErr}

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", ports.ResumeMutationTitle, func(*domain.Resume) error {
		t.Fatal("mutation must not run")
		return nil
	})

	require.ErrorIs(t, err, repositoryErr)
	require.True(t, uow.called)
	require.False(t, repo.getCalled)
}

func TestExecuteResumeMutationPropagatesRepositoryGetError(t *testing.T) {
	getErr := errors.New("get failure")
	resume := newMutationTestResume(t)
	repo := &commandRepositoryStub{resume: resume, getErr: getErr}
	uow := &unitOfWorkStub{}

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", ports.ResumeMutationTitle, func(*domain.Resume) error {
		t.Fatal("mutation must not run")
		return nil
	})

	require.ErrorIs(t, err, getErr)
	require.True(t, repo.getCalled)
	require.False(t, repo.saveCalled)
}

func TestExecuteResumeMutationRejectsNilResumeFromRepository(t *testing.T) {
	repo := &commandRepositoryStub{}
	uow := &unitOfWorkStub{}
	resumeID := sharedvalueobjects.NewID()
	userID := sharedvalueobjects.NewID()

	err := executeResumeMutation(
		context.Background(), repo, uow, resumeID, userID, 1,
		"replace title", ports.ResumeMutationTitle,
		func(*domain.Resume) error {
			t.Fatal("mutation must not run")
			return nil
		},
	)

	require.EqualError(t, err, "replace title: repository returned nil resume")
	require.True(t, uow.called)
	require.True(t, repo.getCalled)
	require.False(t, repo.saveCalled)
}

func TestExecuteResumeMutationRejectsStaleVersion(t *testing.T) {
	created := newMutationTestResume(t)
	resume := domain.RestoreResume(domain.ResumeSnapshot{
		ID: created.ID(), UserID: created.UserID(), Title: created.Title(),
		Language: created.Language(), Status: created.Status(), Version: 2,
	})
	repo := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{}

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), 1, "replace title", ports.ResumeMutationTitle, func(*domain.Resume) error {
		t.Fatal("mutation must not run")
		return nil
	})

	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.True(t, uow.called)
	require.True(t, repo.getCalled)
	require.False(t, repo.saveCalled)
}

func TestExecuteResumeMutationPropagatesMutationErrorWithoutSaving(t *testing.T) {
	mutationErr := errors.New("mutation failure")
	resume := newMutationTestResume(t)
	repo := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{}

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", ports.ResumeMutationTitle, func(*domain.Resume) error {
		return mutationErr
	})

	require.ErrorIs(t, err, mutationErr)
	require.False(t, repo.saveCalled)
}

func TestExecuteResumeMutationPropagatesSaveError(t *testing.T) {
	saveErr := errors.New("save failure")
	resume := newMutationTestResume(t)
	repo := &commandRepositoryStub{resume: resume, saveErr: saveErr}
	uow := &unitOfWorkStub{}

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", ports.ResumeMutationTitle, func(r *domain.Resume) error {
		return r.ReplaceTitle("Updated")
	})

	require.ErrorIs(t, err, saveErr)
	require.True(t, repo.saveCalled)
	require.Equal(t, "Updated", resume.Title())
}

func TestExecuteResumeMutationLoadsMutatesAndSavesAggregate(t *testing.T) {
	resume := newMutationTestResume(t)
	repo := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{}
	var received *domain.Resume

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", ports.ResumeMutationTitle, func(r *domain.Resume) error {
		received = r
		return r.ReplaceTitle("Updated")
	})

	require.NoError(t, err)
	require.True(t, uow.called)
	require.True(t, repo.getCalled)
	require.True(t, repo.saveCalled)
	require.Same(t, resume, received)
	require.Equal(t, "Updated", resume.Title())
}
