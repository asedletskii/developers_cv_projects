package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"testing"

	"github.com/stretchr/testify/require"
)

type stubLanguagesReader struct {
	views []LanguageView
	view  LanguageView
	list  ListLanguagesQuery
	get   GetLanguageQuery
}

func (s *stubLanguagesReader) ListLanguages(_ context.Context, q ListLanguagesQuery) ([]LanguageView, error) {
	s.list = q
	return s.views, nil
}
func (s *stubLanguagesReader) GetLanguage(_ context.Context, q GetLanguageQuery) (LanguageView, error) {
	s.get = q
	return s.view, nil
}
func TestLanguagesQueryHandlers(t *testing.T) {
	uid, rid, lid := useridentity.NewID(), valueobjects.NewID(), valueobjects.NewID()
	r := &stubLanguagesReader{views: []LanguageView{{ID: lid.String()}}, view: LanguageView{ID: lid.String()}}
	v, e := NewListLanguagesQueryHandler(r).Handle(context.Background(), ListLanguagesQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, e)
	require.Equal(t, r.views, v)
	_, e = NewGetLanguageQueryHandler(r).Handle(context.Background(), GetLanguageQuery{ResumeID: rid, LanguageID: lid, UserID: uid})
	require.NoError(t, e)
	require.Equal(t, lid, r.get.LanguageID)
}
func TestLanguagesQueryHandlersRejectMissingIDs(t *testing.T) {
	_, e := NewListLanguagesQueryHandler(&stubLanguagesReader{}).Handle(context.Background(), ListLanguagesQuery{})
	require.Error(t, e)
	_, e = NewGetLanguageQueryHandler(&stubLanguagesReader{}).Handle(context.Background(), GetLanguageQuery{})
	require.Error(t, e)
}
