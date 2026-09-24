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
