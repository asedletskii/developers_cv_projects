package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type stubExperienceQueryReader struct{}

func (stubExperienceQueryReader) ListExperiences(context.Context, ListExperiencesQuery) ([]ExperienceView, error) {
	return []ExperienceView{{ID: uuid.NewString()}}, nil
}
func (stubExperienceQueryReader) GetExperience(context.Context, GetExperienceQuery) (ExperienceView, error) {
	return ExperienceView{ID: uuid.NewString()}, nil
}

func TestExperienceQueryHandlersDelegateToReader(t *testing.T) {
	rid, uid := valueobjects.NewID(), useridentity.NewID()
	r := stubExperienceQueryReader{}
	items, err := NewListExperiencesQueryHandler(r).Handle(context.Background(), ListExperiencesQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, err)
	require.Len(t, items, 1)
	_, err = NewGetExperienceQueryHandler(r).Handle(context.Background(), GetExperienceQuery{ResumeID: rid, ExperienceID: valueobjects.NewID(), UserID: uid})
	require.NoError(t, err)
}

func TestExperienceQueryHandlersRejectMissingIDs(t *testing.T) {
	rid, uid := valueobjects.NewID(), useridentity.NewID()
	_, err := NewListExperiencesQueryHandler(stubExperienceQueryReader{}).Handle(context.Background(), ListExperiencesQuery{ResumeID: valueobjects.ID{}, UserID: uid})
	require.Error(t, err)
	_, err = NewGetExperienceQueryHandler(stubExperienceQueryReader{}).Handle(context.Background(), GetExperienceQuery{ResumeID: rid, ExperienceID: valueobjects.ID{}, UserID: uid})
	require.Error(t, err)
}
