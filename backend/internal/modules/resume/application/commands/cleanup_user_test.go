package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type userDataCleanupRepositoryStub struct {
	resumeIDs []sharedvalueobjects.ID
	deleted   []sharedvalueobjects.ID
	listErr   error
	deleteErr error
}

func (s *userDataCleanupRepositoryStub) ListIDsForUpdateByUserID(_ context.Context, userID sharedvalueobjects.ID) ([]sharedvalueobjects.ID, error) {
	s.deleted = append(s.deleted, userID)
	return s.resumeIDs, s.listErr
}

func (s *userDataCleanupRepositoryStub) DeleteByUserID(_ context.Context, userID sharedvalueobjects.ID) error {
	s.deleted = append(s.deleted, userID)
	return s.deleteErr
}

func TestDeleteUserResumesHandlerCleansLinksBeforeResumes(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	resumeIDs := []sharedvalueobjects.ID{sharedvalueobjects.NewID(), sharedvalueobjects.NewID()}
	repository := &userDataCleanupRepositoryStub{resumeIDs: resumeIDs}
	links := &resumeLinkCleanupStub{}
	handler := NewDeleteUserResumesHandler(repository, links)

	require.NoError(t, handler.DeleteByUserID(context.Background(), userID))

	require.True(t, links.called)
	require.Equal(t, resumeIDs, links.resumeIDs)
	require.Equal(t, []sharedvalueobjects.ID{userID, userID}, repository.deleted)
}

func TestDeleteUserResumesHandlerIsNoOpWithoutResumes(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	repository := &userDataCleanupRepositoryStub{}
	links := &resumeLinkCleanupStub{}
	handler := NewDeleteUserResumesHandler(repository, links)

	require.NoError(t, handler.DeleteByUserID(context.Background(), userID))

	require.False(t, links.called)
	require.Equal(t, []sharedvalueobjects.ID{userID, userID}, repository.deleted)
}

func TestDeleteUserResumesHandlerStopsWhenLinkCleanupFails(t *testing.T) {
	cleanupErr := errors.New("links cleanup failed")
	repository := &userDataCleanupRepositoryStub{resumeIDs: []sharedvalueobjects.ID{sharedvalueobjects.NewID()}}
	handler := NewDeleteUserResumesHandler(repository, &resumeLinkCleanupStub{err: cleanupErr})

	err := handler.DeleteByUserID(context.Background(), sharedvalueobjects.NewID())

	require.ErrorIs(t, err, cleanupErr)
	require.Len(t, repository.deleted, 1, "resume delete must not run after link cleanup fails")
}

func TestDeleteUserResumesHandlerRejectsEmptyUserID(t *testing.T) {
	repository := &userDataCleanupRepositoryStub{}
	handler := NewDeleteUserResumesHandler(repository, &resumeLinkCleanupStub{})

	require.Error(t, handler.DeleteByUserID(context.Background(), sharedvalueobjects.ID{}))
	require.Empty(t, repository.deleted)
}
