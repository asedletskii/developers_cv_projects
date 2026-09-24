package queries

import (
	"context"
	"testing"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"

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

type stubProfileReader struct {
	view ResumeProfileView
	got  GetResumeProfileQuery
}

func (s *stubProfileReader) GetResumeProfile(_ context.Context, q GetResumeProfileQuery) (ResumeProfileView, error) {
	s.got = q
	return s.view, nil
}

type stubContactsReader struct {
	view ResumeContactsView
	got  GetResumeContactsQuery
}

func (s *stubContactsReader) GetResumeContacts(_ context.Context, q GetResumeContactsQuery) (ResumeContactsView, error) {
	s.got = q
	return s.view, nil
}

type stubSkillsReader struct {
	view SkillsView
	got  GetSkillsQuery
}

func (s *stubSkillsReader) GetSkills(_ context.Context, q GetSkillsQuery) (SkillsView, error) {
	s.got = q
	return s.view, nil
}

func TestPrivateGetQueryHandlersDelegateToReaders(t *testing.T) {
	uid, rid := useridentity.NewID(), valueobjects.NewID()
	rr := &stubResumeRootReader{view: ResumeRootView{ID: rid.String()}}
	pv := &stubProfileReader{view: ResumeProfileView{ID: "profile"}}
	cv := &stubContactsReader{view: ResumeContactsView{ID: "contacts"}}
	sv := &stubSkillsReader{view: SkillsView{ID: "skills"}}
	rv, err := NewGetResumeRootQueryHandler(rr).Handle(context.Background(), GetResumeRootQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, err)
	require.Equal(t, rr.view, rv)
	require.Equal(t, rid, rr.got.ResumeID)
	p, err := NewGetResumeProfileQueryHandler(pv).Handle(context.Background(), GetResumeProfileQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, err)
	require.Equal(t, pv.view, p)
	c, err := NewGetResumeContactsQueryHandler(cv).Handle(context.Background(), GetResumeContactsQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, err)
	require.Equal(t, cv.view, c)
	s, err := NewGetSkillsQueryHandler(sv).Handle(context.Background(), GetSkillsQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, err)
	require.Equal(t, sv.view, s)
}

func TestPrivateGetQueryHandlersRejectMissingIDs(t *testing.T) {
	_, err := NewGetResumeRootQueryHandler(&stubResumeRootReader{}).Handle(context.Background(), GetResumeRootQuery{})
	require.Error(t, err)
	_, err = NewGetResumeProfileQueryHandler(&stubProfileReader{}).Handle(context.Background(), GetResumeProfileQuery{})
	require.Error(t, err)
	_, err = NewGetResumeContactsQueryHandler(&stubContactsReader{}).Handle(context.Background(), GetResumeContactsQuery{})
	require.Error(t, err)
	_, err = NewGetSkillsQueryHandler(&stubSkillsReader{}).Handle(context.Background(), GetSkillsQuery{})
	require.Error(t, err)
}

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

func TestLanguagesAndTechnologiesQueryHandlersDelegateToReaders(t *testing.T) {
	uid, rid, lid, tid := useridentity.NewID(), valueobjects.NewID(), valueobjects.NewID(), valueobjects.NewID()
	lr := &stubLanguagesReader{views: []LanguageView{{ID: lid.String()}}, view: LanguageView{ID: lid.String()}}
	tr := &stubTechnologiesReader{views: []TechnologyView{{ID: tid.String()}}, view: TechnologyView{ID: tid.String()}}

	languages, err := NewListLanguagesQueryHandler(lr).Handle(context.Background(), ListLanguagesQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, err)
	require.Equal(t, lr.views, languages)
	require.Equal(t, rid, lr.list.ResumeID)

	language, err := NewGetLanguageQueryHandler(lr).Handle(context.Background(), GetLanguageQuery{ResumeID: rid, LanguageID: lid, UserID: uid})
	require.NoError(t, err)
	require.Equal(t, lr.view, language)
	require.Equal(t, lid, lr.get.LanguageID)

	technologies, err := NewListTechnologiesQueryHandler(tr).Handle(context.Background(), ListTechnologiesQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, err)
	require.Equal(t, tr.views, technologies)
	require.Equal(t, rid, tr.list.ResumeID)

	technology, err := NewGetTechnologyQueryHandler(tr).Handle(context.Background(), GetTechnologyQuery{ResumeID: rid, TechnologyID: tid, UserID: uid})
	require.NoError(t, err)
	require.Equal(t, tr.view, technology)
	require.Equal(t, tid, tr.get.TechnologyID)
}

func TestLanguagesAndTechnologiesQueryHandlersRejectMissingIDs(t *testing.T) {
	_, err := NewListLanguagesQueryHandler(&stubLanguagesReader{}).Handle(context.Background(), ListLanguagesQuery{})
	require.Error(t, err)
	_, err = NewGetLanguageQueryHandler(&stubLanguagesReader{}).Handle(context.Background(), GetLanguageQuery{})
	require.Error(t, err)
	_, err = NewListTechnologiesQueryHandler(&stubTechnologiesReader{}).Handle(context.Background(), ListTechnologiesQuery{})
	require.Error(t, err)
	_, err = NewGetTechnologyQueryHandler(&stubTechnologiesReader{}).Handle(context.Background(), GetTechnologyQuery{})
	require.Error(t, err)
}
