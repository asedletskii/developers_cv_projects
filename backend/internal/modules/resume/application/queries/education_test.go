package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type stubEducationQueryReader struct{}

func (stubEducationQueryReader) ListEducations(context.Context, ListEducationsQuery) ([]EducationView, error) {
	return []EducationView{{ID: uuid.NewString()}}, nil
}
func (stubEducationQueryReader) GetEducation(context.Context, GetEducationQuery) (EducationView, error) {
	return EducationView{ID: uuid.NewString()}, nil
}

func TestEducationQueryHandlersDelegateToReader(t *testing.T) {
	rid, uid := valueobjects.NewID(), useridentity.NewID()
	r := stubEducationQueryReader{}
	items, err := NewListEducationsQueryHandler(r).Handle(context.Background(), ListEducationsQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, err)
	require.Len(t, items, 1)
	_, err = NewGetEducationQueryHandler(r).Handle(context.Background(), GetEducationQuery{ResumeID: rid, EducationID: valueobjects.NewID(), UserID: uid})
	require.NoError(t, err)
}

func TestEducationQueryHandlersRejectMissingIDs(t *testing.T) {
	rid, uid := valueobjects.NewID(), useridentity.NewID()
	_, err := NewListEducationsQueryHandler(stubEducationQueryReader{}).Handle(context.Background(), ListEducationsQuery{ResumeID: valueobjects.ID{}, UserID: uid})
	require.Error(t, err)
	_, err = NewGetEducationQueryHandler(stubEducationQueryReader{}).Handle(context.Background(), GetEducationQuery{ResumeID: rid, EducationID: valueobjects.ID{}, UserID: uid})
	require.Error(t, err)
}
