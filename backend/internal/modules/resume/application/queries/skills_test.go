package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"testing"

	"github.com/stretchr/testify/require"
)

type stubSkillsReader struct{ view SkillsView }

func (s *stubSkillsReader) GetSkills(_ context.Context, q GetSkillsQuery) (SkillsView, error) {
	return s.view, nil
}
func TestGetSkillsQueryHandler(t *testing.T) {
	uid, rid := useridentity.NewID(), valueobjects.NewID()
	r := &stubSkillsReader{view: SkillsView{ID: "skills"}}
	v, e := NewGetSkillsQueryHandler(r).Handle(context.Background(), GetSkillsQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, e)
	require.Equal(t, r.view, v)
}
func TestGetSkillsQueryHandlerRejectsMissingIDs(t *testing.T) {
	_, e := NewGetSkillsQueryHandler(&stubSkillsReader{}).Handle(context.Background(), GetSkillsQuery{})
	require.Error(t, e)
}
