package queries

import (
	"context"
	"errors"
	"testing"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"

	"github.com/stretchr/testify/require"
)

type fastResumeAccessStub struct {
	owned        bool
	accessErr    error
	accessCalled bool
	getCalled    bool
}

func (s *fastResumeAccessStub) OwnsResume(context.Context, resumequeries.GetResumeQuery) (bool, error) {
	s.accessCalled = true
	return s.owned, s.accessErr
}

func (s *fastResumeAccessStub) GetOwned(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	s.getCalled = true
	return resumequeries.ResumeView{}, nil
}

func TestEnsureResumeOwnedUsesLightweightAccessPort(t *testing.T) {
	reader := &fastResumeAccessStub{owned: true}
	query := resumequeries.GetResumeQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	err := EnsureResumeOwned(context.Background(), reader, query)

	require.NoError(t, err)
	require.True(t, reader.accessCalled)
	require.False(t, reader.getCalled)
}

func TestEnsureResumeOwnedReturnsNotFoundWhenAccessDeniesOwnership(t *testing.T) {
	reader := &fastResumeAccessStub{owned: false}
	query := resumequeries.GetResumeQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	err := EnsureResumeOwned(context.Background(), reader, query)

	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	require.True(t, reader.accessCalled)
	require.False(t, reader.getCalled)
}

func TestEnsureResumeOwnedPropagatesAccessError(t *testing.T) {
	want := errors.New("database unavailable")
	reader := &fastResumeAccessStub{accessErr: want}
	query := resumequeries.GetResumeQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	err := EnsureResumeOwned(context.Background(), reader, query)

	require.ErrorIs(t, err, want)
	require.False(t, reader.getCalled)
}
