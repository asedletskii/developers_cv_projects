package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type resumeLinkCleanupRepositoryStub struct {
	linkIDs      []sharedvalueobjects.ID
	listErr      error
	deleteErr    error
	listCalled   bool
	deleteCalled bool
	deleted      []sharedvalueobjects.ID
}

func (s *resumeLinkCleanupRepositoryStub) ListIDsForUpdateByResumeIDs(context.Context, []sharedvalueobjects.ID) ([]sharedvalueobjects.ID, error) {
	s.listCalled = true
	return s.linkIDs, s.listErr
}

func (s *resumeLinkCleanupRepositoryStub) DeleteByResumeIDs(_ context.Context, resumeIDs []sharedvalueobjects.ID) error {
	s.deleteCalled = true
	s.deleted = append([]sharedvalueobjects.ID(nil), resumeIDs...)
	return s.deleteErr
}

func TestDeleteResumeLinksCleanupRemovesAnalyticsBeforeLinks(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	repository := &resumeLinkCleanupRepositoryStub{linkIDs: []sharedvalueobjects.ID{linkID}}
	analytics := &analyticsCleanupStub{}

	handler := NewDeleteResumeLinksHandler(repository, analytics)
	require.NoError(t, handler.DeleteByResumeIDs(context.Background(), []sharedvalueobjects.ID{resumeID}))
	require.True(t, repository.listCalled)
	require.True(t, analytics.called)
	require.Equal(t, []sharedvalueobjects.ID{linkID}, analytics.linkIDs)
	require.True(t, repository.deleteCalled)
	require.Equal(t, []sharedvalueobjects.ID{resumeID}, repository.deleted)
}

func TestDeleteResumeLinksCleanupIsNoOpWithoutResumes(t *testing.T) {
	repository := &resumeLinkCleanupRepositoryStub{}
	analytics := &analyticsCleanupStub{}

	handler := NewDeleteResumeLinksHandler(repository, analytics)
	require.NoError(t, handler.DeleteByResumeIDs(context.Background(), nil))
	require.False(t, repository.listCalled)
	require.False(t, analytics.called)
}

func TestDeleteResumeLinksCleanupRollsBackWhenAnalyticsFails(t *testing.T) {
	repository := &resumeLinkCleanupRepositoryStub{linkIDs: []sharedvalueobjects.ID{sharedvalueobjects.NewID()}}
	analytics := &analyticsCleanupStub{err: errors.New("cleanup failed")}

	handler := NewDeleteResumeLinksHandler(repository, analytics)
	err := handler.DeleteByResumeIDs(context.Background(), []sharedvalueobjects.ID{sharedvalueobjects.NewID()})
	require.Error(t, err)
	require.True(t, analytics.called)
	require.False(t, repository.deleteCalled)
}
