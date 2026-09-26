package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type stubActivityQueryReader struct{}

func (stubActivityQueryReader) ListActivities(context.Context, ListActivitiesQuery) ([]ActivityView, error) {
	return []ActivityView{{ID: uuid.NewString()}}, nil
}
func (stubActivityQueryReader) GetActivity(context.Context, GetActivityQuery) (ActivityView, error) {
	return ActivityView{ID: uuid.NewString()}, nil
}

func TestActivityQueryHandlersDelegateToReader(t *testing.T) {
	rid, uid := valueobjects.NewID(), useridentity.NewID()
	r := stubActivityQueryReader{}
	items, err := NewListActivitiesQueryHandler(r).Handle(context.Background(), ListActivitiesQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, err)
	require.Len(t, items, 1)
	_, err = NewGetActivityQueryHandler(r).Handle(context.Background(), GetActivityQuery{ResumeID: rid, ActivityID: valueobjects.NewID(), UserID: uid})
	require.NoError(t, err)
}

func TestActivityQueryHandlersRejectMissingIDs(t *testing.T) {
	rid, uid := valueobjects.NewID(), useridentity.NewID()
	_, err := NewListActivitiesQueryHandler(stubActivityQueryReader{}).Handle(context.Background(), ListActivitiesQuery{ResumeID: valueobjects.ID{}, UserID: uid})
	require.Error(t, err)
	_, err = NewGetActivityQueryHandler(stubActivityQueryReader{}).Handle(context.Background(), GetActivityQuery{ResumeID: rid, ActivityID: valueobjects.ID{}, UserID: uid})
	require.Error(t, err)
}
