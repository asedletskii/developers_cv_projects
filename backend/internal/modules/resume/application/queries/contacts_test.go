package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"testing"

	"github.com/stretchr/testify/require"
)

type stubContactsReader struct{ view ResumeContactsView }

func (s *stubContactsReader) GetResumeContacts(_ context.Context, q GetResumeContactsQuery) (ResumeContactsView, error) {
	return s.view, nil
}
func TestGetResumeContactsQueryHandler(t *testing.T) {
	uid, rid := useridentity.NewID(), valueobjects.NewID()
	r := &stubContactsReader{view: ResumeContactsView{ID: "contacts"}}
	v, e := NewGetResumeContactsQueryHandler(r).Handle(context.Background(), GetResumeContactsQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, e)
	require.Equal(t, r.view, v)
}
func TestGetResumeContactsQueryHandlerRejectsMissingIDs(t *testing.T) {
	_, e := NewGetResumeContactsQueryHandler(&stubContactsReader{}).Handle(context.Background(), GetResumeContactsQuery{})
	require.Error(t, e)
}
