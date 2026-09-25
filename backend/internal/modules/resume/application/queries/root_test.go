package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"testing"

	"github.com/stretchr/testify/require"
)

type stubResumeRootReader struct {
	view ResumeRootView
	got  GetResumeRootQuery
}

func (s *stubResumeRootReader) GetResumeRoot(_ context.Context, q GetResumeRootQuery) (ResumeRootView, error) {
	s.got = q
	return s.view, nil
}
func TestGetResumeRootQueryHandler(t *testing.T) {
	uid, rid := useridentity.NewID(), valueobjects.NewID()
	r := &stubResumeRootReader{view: ResumeRootView{ID: rid.String()}}
	got, err := NewGetResumeRootQueryHandler(r).Handle(context.Background(), GetResumeRootQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, err)
	require.Equal(t, r.view, got)
	require.Equal(t, rid, r.got.ResumeID)
}
func TestGetResumeRootQueryHandlerRejectsMissingIDs(t *testing.T) {
	_, err := NewGetResumeRootQueryHandler(&stubResumeRootReader{}).Handle(context.Background(), GetResumeRootQuery{})
	require.Error(t, err)
}
