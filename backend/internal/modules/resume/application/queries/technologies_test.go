package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"testing"

	"github.com/stretchr/testify/require"
)

type stubTechnologiesReader struct {
	views []TechnologyView
	view  TechnologyView
	list  ListTechnologiesQuery
	get   GetTechnologyQuery
}

func (s *stubTechnologiesReader) ListTechnologies(_ context.Context, q ListTechnologiesQuery) ([]TechnologyView, error) {
	s.list = q
	return s.views, nil
}
func (s *stubTechnologiesReader) GetTechnology(_ context.Context, q GetTechnologyQuery) (TechnologyView, error) {
	s.get = q
	return s.view, nil
}
func TestTechnologiesQueryHandlers(t *testing.T) {
	uid, rid, tid := useridentity.NewID(), valueobjects.NewID(), valueobjects.NewID()
	r := &stubTechnologiesReader{views: []TechnologyView{{ID: tid.String()}}, view: TechnologyView{ID: tid.String()}}
	v, e := NewListTechnologiesQueryHandler(r).Handle(context.Background(), ListTechnologiesQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, e)
	require.Equal(t, r.views, v)
	_, e = NewGetTechnologyQueryHandler(r).Handle(context.Background(), GetTechnologyQuery{ResumeID: rid, TechnologyID: tid, UserID: uid})
	require.NoError(t, e)
	require.Equal(t, tid, r.get.TechnologyID)
}
func TestTechnologiesQueryHandlersRejectMissingIDs(t *testing.T) {
	_, e := NewListTechnologiesQueryHandler(&stubTechnologiesReader{}).Handle(context.Background(), ListTechnologiesQuery{})
	require.Error(t, e)
	_, e = NewGetTechnologyQueryHandler(&stubTechnologiesReader{}).Handle(context.Background(), GetTechnologyQuery{})
	require.Error(t, e)
}
