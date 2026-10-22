package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type deleteUserRepositoryStub struct {
	exists    bool
	lockErr   error
	deleteErr error
	calls     []string
}

func (s *deleteUserRepositoryStub) LockForDelete(context.Context, sharedvalueobjects.ID) (bool, error) {
	s.calls = append(s.calls, "lock")
	return s.exists, s.lockErr
}

func (s *deleteUserRepositoryStub) Delete(context.Context, sharedvalueobjects.ID) error {
	s.calls = append(s.calls, "delete")
	if s.deleteErr == nil {
		s.exists = false
	}
	return s.deleteErr
}

type deleteUserCleanupStub struct {
	name string
	err  error
	ids  []sharedvalueobjects.ID
}

func (s *deleteUserCleanupStub) DeleteByUserID(_ context.Context, userID sharedvalueobjects.ID) error {
	s.ids = append(s.ids, userID)
	return s.err
}

func TestDeleteUserHandlerDeletesDependentsBeforeUser(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	users := &deleteUserRepositoryStub{exists: true}
	authCleanup := &deleteUserCleanupStub{name: "auth"}
	resumeCleanup := &deleteUserCleanupStub{name: "resume"}
	handler := NewDeleteUserHandler(users, authCleanup, resumeCleanup, &unitOfWorkStub{})

	require.NoError(t, handler.Handle(context.Background(), DeleteUserCommand{UserID: userID}))

	require.Equal(t, []string{"lock", "delete"}, users.calls)
	require.Equal(t, []sharedvalueobjects.ID{userID}, authCleanup.ids)
	require.Equal(t, []sharedvalueobjects.ID{userID}, resumeCleanup.ids)
}

func TestDeleteUserHandlerIsIdempotentForMissingUser(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	users := &deleteUserRepositoryStub{exists: true}
	authCleanup := &deleteUserCleanupStub{name: "auth"}
	resumeCleanup := &deleteUserCleanupStub{name: "resume"}
	handler := NewDeleteUserHandler(users, authCleanup, resumeCleanup, &unitOfWorkStub{})

	require.NoError(t, handler.Handle(context.Background(), DeleteUserCommand{UserID: userID}))
	require.NoError(t, handler.Handle(context.Background(), DeleteUserCommand{UserID: userID}))
	require.Equal(t, []string{"lock", "delete", "lock"}, users.calls)
	require.Equal(t, []sharedvalueobjects.ID{userID}, authCleanup.ids)
	require.Equal(t, []sharedvalueobjects.ID{userID}, resumeCleanup.ids)
}

func TestDeleteUserHandlerStopsWhenAuthCleanupFails(t *testing.T) {
	cleanupErr := errors.New("auth cleanup failed")
	users := &deleteUserRepositoryStub{exists: true}
	authCleanup := &deleteUserCleanupStub{err: cleanupErr}
	resumeCleanup := &deleteUserCleanupStub{name: "resume"}
	handler := NewDeleteUserHandler(users, authCleanup, resumeCleanup, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), DeleteUserCommand{UserID: sharedvalueobjects.NewID()})

	require.ErrorIs(t, err, cleanupErr)
	require.Equal(t, []string{"lock"}, users.calls)
	require.Empty(t, resumeCleanup.ids)
}

func TestDeleteUserHandlerStopsWhenResumeCleanupFails(t *testing.T) {
	cleanupErr := errors.New("resume cleanup failed")
	users := &deleteUserRepositoryStub{exists: true}
	authCleanup := &deleteUserCleanupStub{name: "auth"}
	resumeCleanup := &deleteUserCleanupStub{err: cleanupErr}
	handler := NewDeleteUserHandler(users, authCleanup, resumeCleanup, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), DeleteUserCommand{UserID: sharedvalueobjects.NewID()})

	require.ErrorIs(t, err, cleanupErr)
	require.Equal(t, []string{"lock"}, users.calls)
}

func TestDeleteUserHandlerRejectsEmptyUserID(t *testing.T) {
	users := &deleteUserRepositoryStub{exists: true}
	handler := NewDeleteUserHandler(users, &deleteUserCleanupStub{}, &deleteUserCleanupStub{}, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), DeleteUserCommand{})

	require.Error(t, err)
	require.Empty(t, users.calls)
}
