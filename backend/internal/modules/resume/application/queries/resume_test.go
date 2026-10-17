package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestGetResumeQueryHandlerDelegatesToOwnedReader(t *testing.T) {
	expected := ResumeView{ID: "resume-1", Title: "Backend Developer", Language: "en"}
	got, err := NewGetResumeQueryHandler(resumeReaderStub{view: expected}).Handle(context.Background(), GetResumeQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()})
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestGetResumeQueryHandlerRejectsMissingUserID(t *testing.T) {
	_, err := NewGetResumeQueryHandler(resumeReaderStub{}).Handle(context.Background(), GetResumeQuery{ResumeID: sharedvalueobjects.NewID()})
	require.Error(t, err)
}

func TestGetResumeQueryHandlerRejectsMissingResumeID(t *testing.T) {
	_, err := NewGetResumeQueryHandler(resumeReaderStub{}).Handle(context.Background(), GetResumeQuery{UserID: sharedvalueobjects.NewID()})
	require.Error(t, err)
}

func TestGetResumeQueryHandlerRejectsNilReader(t *testing.T) {
	_, err := NewGetResumeQueryHandler(nil).Handle(context.Background(), GetResumeQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()})
	require.Error(t, err)
}

func TestGetResumeQueryHandlerPropagatesReaderError(t *testing.T) {
	want := context.Canceled
	reader := resumeReaderErrorStub{err: want}
	_, err := NewGetResumeQueryHandler(reader).Handle(context.Background(), GetResumeQuery{ResumeID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID()})
	require.ErrorIs(t, err, want)
}
