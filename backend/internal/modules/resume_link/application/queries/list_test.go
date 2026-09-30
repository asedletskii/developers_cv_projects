package queries

import (
	"context"
	"testing"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumevalueobjects "developers-cv/backend/internal/modules/resume/domain/valueobjects"
	uservalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"

	"github.com/stretchr/testify/require"
)

type listReaderStub struct {
	items []ResumeLinkListItem
	err   error
	query ListResumeLinksQuery
}

func (s *listReaderStub) ListByResumeID(_ context.Context, query ListResumeLinksQuery) ([]ResumeLinkListItem, error) {
	s.query = query
	return s.items, s.err
}

type ownershipReaderStub struct {
	err   error
	query resumequeries.GetResumeQuery
}

func (s *ownershipReaderStub) GetOwned(_ context.Context, query resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	s.query = query
	if s.err != nil {
		return resumequeries.ResumeView{}, s.err
	}
	return resumequeries.ResumeView{ID: query.ResumeID.String()}, nil
}

func TestListResumeLinksQueryHandlerChecksOwnershipBeforeReading(t *testing.T) {
	userID := uservalueobjects.NewID()
	resumeID := resumevalueobjects.NewID()
	reader := &listReaderStub{items: []ResumeLinkListItem{{ID: "link", Name: "HH"}}}
	ownership := &ownershipReaderStub{}

	items, err := NewListResumeLinksQueryHandler(reader, ownership).Handle(context.Background(), ListResumeLinksQuery{UserID: userID, ResumeID: resumeID})
	require.NoError(t, err)
	require.Equal(t, reader.items, items)
	require.Equal(t, resumeID, ownership.query.ResumeID)
	require.Equal(t, userID, ownership.query.UserID)
	require.Equal(t, resumeID, reader.query.ResumeID)
}

func TestListResumeLinksQueryHandlerRejectsForeignResume(t *testing.T) {
	reader := &listReaderStub{}
	ownership := &ownershipReaderStub{err: resumequeries.ErrResumeNotFound}

	_, err := NewListResumeLinksQueryHandler(reader, ownership).Handle(context.Background(), ListResumeLinksQuery{
		UserID: uservalueobjects.NewID(), ResumeID: resumevalueobjects.NewID(),
	})
	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	require.True(t, reader.query.ResumeID.IsZero())
}

func TestListResumeLinksQueryHandlerRejectsInvalidQuery(t *testing.T) {
	h := NewListResumeLinksQueryHandler(&listReaderStub{}, &ownershipReaderStub{})
	_, err := h.Handle(context.Background(), ListResumeLinksQuery{UserID: uservalueobjects.NewID()})
	require.Error(t, err)
}
