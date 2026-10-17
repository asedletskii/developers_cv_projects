package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestDownloadContactsQueryHandlerReturnsAllContacts(t *testing.T) {
	phone, email := "+33123456789", "person@example.com"
	reader := &downloadContactsReaderStub{view: DownloadContactsView{Phone: &phone, Email: &email}}
	resumeID := sharedvalueobjects.NewID()

	result, err := NewDownloadContactsQueryHandler(reader).Handle(context.Background(), DownloadContactsQuery{
		ResumeID: resumeID,
	})

	require.NoError(t, err)
	require.Equal(t, reader.view, result)
	require.Equal(t, resumeID, reader.query.ResumeID)
}

func TestDownloadContactsQueryHandlerRejectsInvalidQuery(t *testing.T) {
	reader := &downloadContactsReaderStub{}

	_, err := NewDownloadContactsQueryHandler(reader).Handle(context.Background(), DownloadContactsQuery{})

	require.ErrorContains(t, err, "invalid download contacts query")
	require.Equal(t, DownloadContactsQuery{}, reader.query)
}
