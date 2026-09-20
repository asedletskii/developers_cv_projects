package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestContactViewQueryHandlerReturnsContact(t *testing.T) {
	reader := &contactViewReaderStub{view: resumecontracts.ContactView{Contact: resumecontracts.ContactKeyEmail, Value: "person@example.com"}}
	resumeID := sharedvalueobjects.NewID()

	result, err := NewContactViewQueryHandler(reader).Handle(context.Background(), resumecontracts.ContactViewQuery{
		ResumeID: resumeID,
		Contact:  resumecontracts.ContactKeyEmail,
	})

	require.NoError(t, err)
	require.Equal(t, reader.view, result)
	require.Equal(t, resumeID, reader.query.ResumeID)
	require.Equal(t, resumecontracts.ContactKeyEmail, reader.query.Contact)
}

func TestContactViewQueryHandlerRejectsInvalidQuery(t *testing.T) {
	reader := &contactViewReaderStub{}

	_, err := NewContactViewQueryHandler(reader).Handle(context.Background(), resumecontracts.ContactViewQuery{
		ResumeID: sharedvalueobjects.NewID(),
		Contact:  "invalid",
	})

	require.Error(t, err)
	require.Equal(t, resumecontracts.ContactViewQuery{}, reader.query)
}

func TestContactViewQueryHandlerPropagatesReaderError(t *testing.T) {
	want := resumecontracts.ErrContactNotHidden
	reader := &contactViewReaderStub{err: want}

	_, err := NewContactViewQueryHandler(reader).Handle(context.Background(), resumecontracts.ContactViewQuery{
		ResumeID: sharedvalueobjects.NewID(),
		Contact:  resumecontracts.ContactKeyEmail,
	})

	require.ErrorIs(t, err, want)
}
