package ownership

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	application "developers-cv/backend/internal/modules/resume_link/application"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ownershipReaderStub struct {
	owned       bool
	accessErr   error
	accessCalls int
}

func (s *ownershipReaderStub) OwnsResume(_ context.Context, _ resumecontracts.ResumeOwnershipQuery) (bool, error) {
	s.accessCalls++
	return s.owned, s.accessErr
}

func TestEnsureResumeOwnedUsesLightweightAccess(t *testing.T) {
	reader := &ownershipReaderStub{owned: true}
	query := resumecontracts.ResumeOwnershipQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	err := EnsureResumeOwned(context.Background(), reader, query)

	require.NoError(t, err)
	require.Equal(t, 1, reader.accessCalls)
}

func TestEnsureResumeOwnedRejectsForeignResume(t *testing.T) {
	reader := &ownershipReaderStub{owned: false}
	query := resumecontracts.ResumeOwnershipQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	err := EnsureResumeOwned(context.Background(), reader, query)

	require.ErrorIs(t, err, application.ErrResumeNotFound)
	require.Equal(t, 1, reader.accessCalls)
}

func TestEnsureResumeOwnedNormalizesProviderNotFound(t *testing.T) {
	reader := &ownershipReaderStub{accessErr: resumecontracts.ErrResumeNotFound}
	query := resumecontracts.ResumeOwnershipQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	err := EnsureResumeOwned(context.Background(), reader, query)

	require.ErrorIs(t, err, application.ErrResumeNotFound)
	require.ErrorIs(t, err, resumecontracts.ErrResumeNotFound)
}

func TestEnsureResumeOwnedRejectsInvalidQuery(t *testing.T) {
	reader := &ownershipReaderStub{owned: true}

	err := EnsureResumeOwned(context.Background(), reader, resumecontracts.ResumeOwnershipQuery{UserID: sharedvalueobjects.NewID()})

	require.Error(t, err)
	require.Zero(t, reader.accessCalls)
}

func TestEnsureResumeOwnedForUpdateLocksOwnedResume(t *testing.T) {
	guard := &resumeLockGuardStub{owned: true}
	resumeID, userID := sharedvalueobjects.NewID(), sharedvalueobjects.NewID()

	err := EnsureResumeOwnedForUpdate(context.Background(), guard, resumeID, userID)

	require.NoError(t, err)
	require.Equal(t, resumecontracts.OwnedResumeLock{ResumeID: resumeID, UserID: userID}, guard.query)
}

func TestEnsureResumeOwnedForUpdateRejectsUnownedResume(t *testing.T) {
	guard := &resumeLockGuardStub{}

	err := EnsureResumeOwnedForUpdate(context.Background(), guard, sharedvalueobjects.NewID(), sharedvalueobjects.NewID())

	require.ErrorIs(t, err, application.ErrResumeNotFound)
}

func TestEnsureResumeOwnedForUpdatePropagatesLockError(t *testing.T) {
	want := errors.New("database unavailable")
	guard := &resumeLockGuardStub{err: want}

	err := EnsureResumeOwnedForUpdate(context.Background(), guard, sharedvalueobjects.NewID(), sharedvalueobjects.NewID())

	require.ErrorIs(t, err, want)
}

type resumeLockGuardStub struct {
	owned bool
	err   error
	query resumecontracts.OwnedResumeLock
}

func (s *resumeLockGuardStub) LockOwnedResumeForUpdate(_ context.Context, query resumecontracts.OwnedResumeLock) (bool, error) {
	s.query = query
	return s.owned, s.err
}
