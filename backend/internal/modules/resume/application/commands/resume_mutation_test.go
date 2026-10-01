package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type mutationRepositoryStub struct {
	resume     *domain.Resume
	getErr     error
	saveErr    error
	getCalled  bool
	saveCalled bool
}

func (r *mutationRepositoryStub) Get(_ context.Context, _ ports.GetResumeParams) (*domain.Resume, error) {
	r.getCalled = true
	if r.getErr != nil {
		return nil, r.getErr
	}
	return r.resume, nil
}

func (r *mutationRepositoryStub) Save(_ context.Context, _ *domain.Resume) error {
	r.saveCalled = true
	return r.saveErr
}

func (r *mutationRepositoryStub) Create(context.Context, *domain.Resume) error           { return nil }
func (r *mutationRepositoryStub) Delete(context.Context, ports.DeleteResumeParams) error { return nil }

type mutationUnitOfWorkStub struct {
	called bool
	err    error
}

func (u *mutationUnitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	u.called = true
	if u.err != nil {
		return u.err
	}
	return fn(ctx)
}

func newMutationTestResume(t *testing.T) *domain.Resume {
	t.Helper()
	resume, err := domain.NewResume(domain.ResumeParams{
		ID:       identity.NewID(),
		UserID:   identity.NewID(),
		Title:    "Original",
		Language: domain.ResumeLanguage("en"),
		Version:  3,
	})
	require.NoError(t, err)
	return resume
}

func TestExecuteResumeMutationRejectsMissingDependencies(t *testing.T) {
	resume := newMutationTestResume(t)
	uow := &mutationUnitOfWorkStub{}

	err := executeResumeMutation(context.Background(), nil, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", func(*domain.Resume) error {
		return nil
	})
	require.EqualError(t, err, "replace title handler dependencies are nil")
	require.False(t, uow.called)

	repo := &mutationRepositoryStub{resume: resume}
	err = executeResumeMutation(context.Background(), repo, nil, resume.ID(), resume.UserID(), resume.Version(), "replace title", func(*domain.Resume) error {
		return nil
	})
	require.EqualError(t, err, "replace title handler dependencies are nil")
	require.False(t, repo.getCalled)
}

func TestExecuteResumeMutationRejectsInvalidCommandBeforeTransaction(t *testing.T) {
	resume := newMutationTestResume(t)
	repo := &mutationRepositoryStub{resume: resume}
	uow := &mutationUnitOfWorkStub{}

	cases := []struct {
		name     string
		resumeID identity.ID
		userID   identity.ID
		version  int64
	}{
		{name: "missing resume id", userID: resume.UserID(), version: 3},
		{name: "missing user id", resumeID: resume.ID(), version: 3},
		{name: "zero version", resumeID: resume.ID(), userID: resume.UserID()},
		{name: "negative version", resumeID: resume.ID(), userID: resume.UserID(), version: -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := executeResumeMutation(context.Background(), repo, uow, tc.resumeID, tc.userID, tc.version, "replace title", func(*domain.Resume) error {
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
	repo := &mutationRepositoryStub{resume: resume}
	uow := &mutationUnitOfWorkStub{err: repositoryErr}

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", func(*domain.Resume) error {
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
	repo := &mutationRepositoryStub{resume: resume, getErr: getErr}
	uow := &mutationUnitOfWorkStub{}

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", func(*domain.Resume) error {
		t.Fatal("mutation must not run")
		return nil
	})

	require.ErrorIs(t, err, getErr)
	require.True(t, repo.getCalled)
	require.False(t, repo.saveCalled)
}

func TestExecuteResumeMutationRejectsStaleVersion(t *testing.T) {
	resume := newMutationTestResume(t)
	repo := &mutationRepositoryStub{resume: resume}
	uow := &mutationUnitOfWorkStub{}

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), resume.Version()-1, "replace title", func(*domain.Resume) error {
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
	repo := &mutationRepositoryStub{resume: resume}
	uow := &mutationUnitOfWorkStub{}

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", func(*domain.Resume) error {
		return mutationErr
	})

	require.ErrorIs(t, err, mutationErr)
	require.False(t, repo.saveCalled)
}

func TestExecuteResumeMutationPropagatesSaveError(t *testing.T) {
	saveErr := errors.New("save failure")
	resume := newMutationTestResume(t)
	repo := &mutationRepositoryStub{resume: resume, saveErr: saveErr}
	uow := &mutationUnitOfWorkStub{}

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", func(r *domain.Resume) error {
		return r.ReplaceTitle("Updated")
	})

	require.ErrorIs(t, err, saveErr)
	require.True(t, repo.saveCalled)
	require.Equal(t, "Updated", resume.Title())
}

func TestExecuteResumeMutationLoadsMutatesAndSavesAggregate(t *testing.T) {
	resume := newMutationTestResume(t)
	repo := &mutationRepositoryStub{resume: resume}
	uow := &mutationUnitOfWorkStub{}
	var received *domain.Resume

	err := executeResumeMutation(context.Background(), repo, uow, resume.ID(), resume.UserID(), resume.Version(), "replace title", func(r *domain.Resume) error {
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
