package queries

import (
	"context"
	"errors"
	"testing"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"

	"github.com/stretchr/testify/require"
)

func TestEnsureResumeOwnedUsesLightweightAccessPort(t *testing.T) {
	reader := &ownershipReaderStub{}
	query := resumequeries.GetResumeQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	err := EnsureResumeOwned(context.Background(), reader, query)

	require.NoError(t, err)
	require.Equal(t, 1, reader.ownershipCalls)
	require.Zero(t, reader.getCalls)
}

func TestEnsureResumeOwnedReturnsNotFoundWhenAccessDeniesOwnership(t *testing.T) {
	reader := &ownershipReaderStub{viewID: sharedvalueobjects.NewID().String()}
	query := resumequeries.GetResumeQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	err := EnsureResumeOwned(context.Background(), reader, query)

	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	require.Equal(t, 1, reader.ownershipCalls)
	require.Zero(t, reader.getCalls)
}

func TestEnsureResumeOwnedPropagatesAccessError(t *testing.T) {
	want := errors.New("database unavailable")
	reader := &ownershipReaderStub{err: want}
	query := resumequeries.GetResumeQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	err := EnsureResumeOwned(context.Background(), reader, query)

	require.ErrorIs(t, err, want)
	require.Zero(t, reader.getCalls)
}
