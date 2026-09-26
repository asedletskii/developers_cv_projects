package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type stubCourseQueryReader struct{}

func (stubCourseQueryReader) ListCourses(context.Context, ListCoursesQuery) ([]CourseView, error) {
	return []CourseView{{ID: uuid.NewString()}}, nil
}
func (stubCourseQueryReader) GetCourse(context.Context, GetCourseQuery) (CourseView, error) {
	return CourseView{ID: uuid.NewString()}, nil
}

func TestCourseQueryHandlersDelegateToReader(t *testing.T) {
	rid, uid := valueobjects.NewID(), useridentity.NewID()
	r := stubCourseQueryReader{}
	items, err := NewListCoursesQueryHandler(r).Handle(context.Background(), ListCoursesQuery{ResumeID: rid, UserID: uid})
	require.NoError(t, err)
	require.Len(t, items, 1)
	_, err = NewGetCourseQueryHandler(r).Handle(context.Background(), GetCourseQuery{ResumeID: rid, CourseID: valueobjects.NewID(), UserID: uid})
	require.NoError(t, err)
}

func TestCourseQueryHandlersRejectMissingIDs(t *testing.T) {
	rid, uid := valueobjects.NewID(), useridentity.NewID()
	_, err := NewListCoursesQueryHandler(stubCourseQueryReader{}).Handle(context.Background(), ListCoursesQuery{ResumeID: valueobjects.ID{}, UserID: uid})
	require.Error(t, err)
	_, err = NewGetCourseQueryHandler(stubCourseQueryReader{}).Handle(context.Background(), GetCourseQuery{ResumeID: rid, CourseID: valueobjects.ID{}, UserID: uid})
	require.Error(t, err)
}
