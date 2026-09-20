package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ownedSectionsReaderStub struct {
	sections []ResumeSectionSettings
	err      error
	query    resumecontracts.OwnedResumeSectionsQuery
}

func (s *ownedSectionsReaderStub) GetOwnedSections(_ context.Context, query resumecontracts.OwnedResumeSectionsQuery) ([]ResumeSectionSettings, error) {
	s.query = query
	return s.sections, s.err
}

func TestOwnedResumeSectionsQueryHandlerReturnsSections(t *testing.T) {
	style := "second"
	reader := &ownedSectionsReaderStub{sections: []ResumeSectionSettings{{ID: "about", Enabled: true, Style: &style}}}
	query := resumecontracts.OwnedResumeSectionsQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()}

	got, err := NewOwnedResumeSectionsQueryHandler(reader).Handle(context.Background(), query)

	require.NoError(t, err)
	require.Equal(t, []resumecontracts.ResumeSectionSettings{{ID: "about", Enabled: true, Style: &style}}, got)
	require.Equal(t, query, reader.query)
}

func TestOwnedResumeSectionsQueryHandlerReturnsNonNilEmptySlice(t *testing.T) {
	reader := &ownedSectionsReaderStub{}

	got, err := NewOwnedResumeSectionsQueryHandler(reader).Handle(context.Background(), resumecontracts.OwnedResumeSectionsQuery{
		ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID(),
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Empty(t, got)
}

func TestOwnedResumeSectionsQueryHandlerRejectsMissingIdentifiers(t *testing.T) {
	reader := &ownedSectionsReaderStub{}

	_, err := NewOwnedResumeSectionsQueryHandler(reader).Handle(context.Background(), resumecontracts.OwnedResumeSectionsQuery{ResumeID: sharedvalueobjects.NewID()})

	require.Error(t, err)
	require.Equal(t, resumecontracts.OwnedResumeSectionsQuery{}, reader.query)
}

func TestOwnedResumeSectionsQueryHandlerPropagatesReaderError(t *testing.T) {
	want := errors.New("database unavailable")
	reader := &ownedSectionsReaderStub{err: want}

	_, err := NewOwnedResumeSectionsQueryHandler(reader).Handle(context.Background(), resumecontracts.OwnedResumeSectionsQuery{
		ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID(),
	})

	require.ErrorIs(t, err, want)
}
