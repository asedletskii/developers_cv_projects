package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"testing"

	"github.com/stretchr/testify/require"
)

type stubProfileReader struct{ view ResumeProfileView }

func (s *stubProfileReader) GetResumeProfile(_ context.Context, q GetResumeProfileQuery) (ResumeProfileView, error) {
	return s.view, nil
}
func TestGetResumeProfileQueryHandler(t *testing.T) {
	uid, rid := useridentity.NewID(), valueobjects.NewID()
	r := &stubProfileReader{view: ResumeProfileView{ID: "profile"}}
	v, e := NewGetResumeProfileQueryHandler(r).Handle(context.Background(), GetResumeProfileQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, e)
	require.Equal(t, r.view, v)
}
func TestGetResumeProfileQueryHandlerRejectsMissingIDs(t *testing.T) {
	_, e := NewGetResumeProfileQueryHandler(&stubProfileReader{}).Handle(context.Background(), GetResumeProfileQuery{})
	require.Error(t, e)
}
