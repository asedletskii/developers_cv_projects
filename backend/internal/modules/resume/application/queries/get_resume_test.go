package queries

import (
	"context"
	"testing"

	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeResumeReader struct {
	query  GetResumeQuery
	resume Resume
	err    error
}

func (f *fakeResumeReader) GetResume(_ context.Context, query GetResumeQuery) (Resume, error) {
	f.query = query
	return f.resume, f.err
}

func TestGetResumeHandlerPassesQueryObjectToReader(t *testing.T) {
	reader := &fakeResumeReader{resume: Resume{ID: "resume-1"}}
	handler := NewGetResumeHandler(reader)
	id := uuid.New()
	query := GetResumeQuery{UserID: useridentity.NewID(), ID: id, Lang: LanguageEN}

	resume, err := handler.Handle(context.Background(), query)

	require.NoError(t, err)
	require.Equal(t, reader.resume, resume)
	require.Equal(t, query, reader.query)
}

func TestGetResumeHandlerRejectsInvalidQuery(t *testing.T) {
	reader := &fakeResumeReader{}
	handler := NewGetResumeHandler(reader)

	_, err := handler.Handle(context.Background(), GetResumeQuery{Lang: LanguageEN})
	require.Error(t, err)

	_, err = handler.Handle(context.Background(), GetResumeQuery{ID: uuid.New(), Lang: "de"})
	require.Error(t, err)
}
