package queries

import (
	"context"
	"testing"

	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeGetResumeReader struct {
	query  GetResumeQuery
	resume Resume
	err    error
}

func (f *fakeGetResumeReader) GetResume(_ context.Context, query GetResumeQuery) (Resume, error) {
	f.query = query
	return f.resume, f.err
}

func TestGetResumeHandlerPassesQueryObjectToReader(t *testing.T) {
	reader := &fakeGetResumeReader{resume: Resume{ID: "resume-1"}}
	handler := NewGetResumeHandler(reader)
	id := uuid.New()
	userID := useridentity.NewID()
	query := GetResumeQuery{UserID: &userID, ID: id, Lang: LanguageEN}

	resume, err := handler.Handle(context.Background(), query)

	require.NoError(t, err)
	require.Equal(t, reader.resume, resume)
	require.Equal(t, query, reader.query)
}

func TestGetResumeHandlerAllowsPublicQuery(t *testing.T) {
	reader := &fakeGetResumeReader{resume: Resume{ID: "resume-1"}}
	handler := NewGetResumeHandler(reader)
	query := GetResumeQuery{ID: uuid.New(), Lang: LanguageRU}

	resume, err := handler.Handle(context.Background(), query)

	require.NoError(t, err)
	require.Equal(t, reader.resume, resume)
	require.Nil(t, reader.query.UserID)
}

func TestGetResumeHandlerRejectsInvalidQuery(t *testing.T) {
	reader := &fakeGetResumeReader{}
	handler := NewGetResumeHandler(reader)

	zeroUserID := useridentity.ID{}
	_, err := handler.Handle(context.Background(), GetResumeQuery{UserID: &zeroUserID, ID: uuid.New(), Lang: LanguageEN})
	require.Error(t, err)

	_, err = handler.Handle(context.Background(), GetResumeQuery{Lang: LanguageEN})
	require.Error(t, err)

	_, err = handler.Handle(context.Background(), GetResumeQuery{ID: uuid.New(), Lang: "de"})
	require.Error(t, err)
}
