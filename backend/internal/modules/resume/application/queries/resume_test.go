package queries

import (
	"context"
	"testing"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"

	"github.com/stretchr/testify/require"
)

type stubResumeReader struct{ view ResumeView }

func (s stubResumeReader) GetOwned(context.Context, valueobjects.ID, useridentity.ID) (ResumeView, error) {
	return s.view, nil
}

func TestGetResumeQueryHandlerDelegatesToOwnedReader(t *testing.T) {
	expected := ResumeView{ID: "resume-1", Title: "Backend Developer", Language: "en"}
	got, err := NewGetResumeQueryHandler(stubResumeReader{view: expected}).Handle(context.Background(), GetResumeQuery{ResumeID: valueobjects.NewID(), UserID: useridentity.NewID()})
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestGetResumeQueryHandlerRejectsMissingUserID(t *testing.T) {
	_, err := NewGetResumeQueryHandler(stubResumeReader{}).Handle(context.Background(), GetResumeQuery{ResumeID: valueobjects.NewID()})
	require.Error(t, err)
}

func TestGetResumeQueryHandlerRejectsMissingResumeID(t *testing.T) {
	_, err := NewGetResumeQueryHandler(stubResumeReader{}).Handle(context.Background(), GetResumeQuery{UserID: useridentity.NewID()})
	require.Error(t, err)
}

func TestGetResumeQueryHandlerRejectsNilReader(t *testing.T) {
	_, err := NewGetResumeQueryHandler(nil).Handle(context.Background(), GetResumeQuery{ResumeID: valueobjects.NewID(), UserID: useridentity.NewID()})
	require.Error(t, err)
}

func TestGetResumeQueryHandlerPropagatesReaderError(t *testing.T) {
	want := context.Canceled
	reader := resumeReaderErrorStub{err: want}
	_, err := NewGetResumeQueryHandler(reader).Handle(context.Background(), GetResumeQuery{ResumeID: valueobjects.NewID(), UserID: useridentity.NewID()})
	require.ErrorIs(t, err, want)
}

type resumeReaderErrorStub struct{ err error }

func (s resumeReaderErrorStub) GetOwned(context.Context, valueobjects.ID, useridentity.ID) (ResumeView, error) {
	return ResumeView{}, s.err
}
