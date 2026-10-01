package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type contactViewReaderStub struct {
	view  ContactView
	err   error
	query ContactViewQuery
}

func (s *contactViewReaderStub) GetContactForView(_ context.Context, query ContactViewQuery) (ContactView, error) {
	s.query = query
	if s.err != nil {
		return ContactView{}, s.err
	}
	return s.view, nil
}

func TestContactViewQueryHandlerReturnsContact(t *testing.T) {
	reader := &contactViewReaderStub{view: ContactView{Contact: domain.ContactKeyEmail, Value: "person@example.com"}}
	resumeID := sharedvalueobjects.NewID()

	result, err := NewContactViewQueryHandler(reader).Handle(context.Background(), ContactViewQuery{
		ResumeID: resumeID,
		Contact:  domain.ContactKeyEmail,
	})

	require.NoError(t, err)
	require.Equal(t, reader.view, result)
	require.Equal(t, resumeID, reader.query.ResumeID)
	require.Equal(t, domain.ContactKeyEmail, reader.query.Contact)
}

func TestContactViewQueryHandlerRejectsInvalidQuery(t *testing.T) {
	reader := &contactViewReaderStub{}

	_, err := NewContactViewQueryHandler(reader).Handle(context.Background(), ContactViewQuery{
		ResumeID: sharedvalueobjects.NewID(),
		Contact:  "invalid",
	})

	require.Error(t, err)
	require.Equal(t, ContactViewQuery{}, reader.query)
}

func TestContactViewQueryHandlerPropagatesReaderError(t *testing.T) {
	want := ErrContactNotHidden
	reader := &contactViewReaderStub{err: want}

	_, err := NewContactViewQueryHandler(reader).Handle(context.Background(), ContactViewQuery{
		ResumeID: sharedvalueobjects.NewID(),
		Contact:  domain.ContactKeyEmail,
	})

	require.ErrorIs(t, err, want)
}
