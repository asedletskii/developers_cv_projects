package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	application "developers-cv/backend/internal/modules/resume_link/application"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestListResumeLinksQueryHandlerChecksOwnershipBeforeReading(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
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
	ownership := &ownershipReaderStub{viewID: sharedvalueobjects.NewID().String()}

	_, err := NewListResumeLinksQueryHandler(reader, ownership).Handle(context.Background(), ListResumeLinksQuery{
		UserID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(),
	})
	require.ErrorIs(t, err, application.ErrResumeNotFound)
	require.True(t, reader.query.ResumeID.IsZero())
}

func TestListResumeLinksQueryHandlerRejectsInvalidQuery(t *testing.T) {
	h := NewListResumeLinksQueryHandler(&listReaderStub{}, &ownershipReaderStub{})
	_, err := h.Handle(context.Background(), ListResumeLinksQuery{UserID: sharedvalueobjects.NewID()})
	require.Error(t, err)
}
