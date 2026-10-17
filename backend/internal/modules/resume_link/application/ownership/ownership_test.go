package ownership

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ownershipReaderStub struct {
	owned         bool
	accessErr     error
	getOwnedErr   error
	accessCalls   int
	getOwnedCalls int
}

func (s *ownershipReaderStub) OwnsResume(_ context.Context, _ ports.ResumeOwnershipQuery) (bool, error) {
	s.accessCalls++
	return s.owned, s.accessErr
}

func (s *ownershipReaderStub) GetOwned(_ context.Context, _ ports.ResumeOwnershipQuery) (resumequeries.ResumeView, error) {
	s.getOwnedCalls++
	return resumequeries.ResumeView{}, s.getOwnedErr
}

func TestEnsureResumeOwnedUsesLightweightAccess(t *testing.T) {
	reader := &ownershipReaderStub{owned: true}
	query := ports.ResumeOwnershipQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	err := EnsureResumeOwned(context.Background(), reader, query)

	require.NoError(t, err)
	require.Equal(t, 1, reader.accessCalls)
	require.Zero(t, reader.getOwnedCalls)
}

func TestEnsureResumeOwnedRejectsForeignResume(t *testing.T) {
	reader := &ownershipReaderStub{owned: false}
	query := ports.ResumeOwnershipQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	err := EnsureResumeOwned(context.Background(), reader, query)

	require.ErrorIs(t, err, resumedomain.ErrEntityNotFound)
	require.Equal(t, 1, reader.accessCalls)
	require.Zero(t, reader.getOwnedCalls)
}

func TestEnsureResumeOwnedRejectsInvalidQuery(t *testing.T) {
	reader := &ownershipReaderStub{owned: true}

	err := EnsureResumeOwned(context.Background(), reader, ports.ResumeOwnershipQuery{UserID: sharedvalueobjects.NewID()})

	require.Error(t, err)
	require.Zero(t, reader.accessCalls)
	require.Zero(t, reader.getOwnedCalls)
}
