package queries

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type stubResumeReader struct{}

func (stubResumeReader) ListExperiences(context.Context, ListExperiencesQuery) ([]ExperienceView, error) {
	return []ExperienceView{{ID: uuid.NewString()}}, nil
}
func (stubResumeReader) GetExperience(context.Context, GetExperienceQuery) (ExperienceView, error) {
	return ExperienceView{ID: uuid.NewString()}, nil
}
func (stubResumeReader) ListEducations(context.Context, ListEducationsQuery) ([]EducationView, error) {
	return []EducationView{{ID: uuid.NewString()}}, nil
}
func (stubResumeReader) GetEducation(context.Context, GetEducationQuery) (EducationView, error) {
	return EducationView{ID: uuid.NewString()}, nil
}
func (stubResumeReader) ListCourses(context.Context, ListCoursesQuery) ([]CourseView, error) {
	return []CourseView{{ID: uuid.NewString()}}, nil
}
func (stubResumeReader) GetCourse(context.Context, GetCourseQuery) (CourseView, error) {
	return CourseView{ID: uuid.NewString()}, nil
}
func (stubResumeReader) ListActivities(context.Context, ListActivitiesQuery) ([]ActivityView, error) {
	return []ActivityView{{ID: uuid.NewString()}}, nil
}
func (stubResumeReader) GetActivity(context.Context, GetActivityQuery) (ActivityView, error) {
	return ActivityView{ID: uuid.NewString()}, nil
}

func testResumeIDs() (valueobjects.ID, useridentity.ID) {
	return valueobjects.NewID(), useridentity.NewID()
}

func TestExperienceEducationCourseActivityQueryHandlersDelegateToReaders(t *testing.T) {
	resumeID, userID := testResumeIDs()
	reader := stubResumeReader{}

	experiences, err := NewListExperiencesQueryHandler(reader).Handle(context.Background(), ListExperiencesQuery{ResumeID: resumeID, UserID: userID})
	require.NoError(t, err)
	require.Len(t, experiences, 1)
	_, err = NewGetExperienceQueryHandler(reader).Handle(context.Background(), GetExperienceQuery{ResumeID: resumeID, ExperienceID: valueobjects.NewID(), UserID: userID})
	require.NoError(t, err)

	educations, err := NewListEducationsQueryHandler(reader).Handle(context.Background(), ListEducationsQuery{ResumeID: resumeID, UserID: userID})
	require.NoError(t, err)
	require.Len(t, educations, 1)
	_, err = NewGetEducationQueryHandler(reader).Handle(context.Background(), GetEducationQuery{ResumeID: resumeID, EducationID: valueobjects.NewID(), UserID: userID})
	require.NoError(t, err)

	courses, err := NewListCoursesQueryHandler(reader).Handle(context.Background(), ListCoursesQuery{ResumeID: resumeID, UserID: userID})
	require.NoError(t, err)
	require.Len(t, courses, 1)
	_, err = NewGetCourseQueryHandler(reader).Handle(context.Background(), GetCourseQuery{ResumeID: resumeID, CourseID: valueobjects.NewID(), UserID: userID})
	require.NoError(t, err)

	activities, err := NewListActivitiesQueryHandler(reader).Handle(context.Background(), ListActivitiesQuery{ResumeID: resumeID, UserID: userID})
	require.NoError(t, err)
	require.Len(t, activities, 1)
	_, err = NewGetActivityQueryHandler(reader).Handle(context.Background(), GetActivityQuery{ResumeID: resumeID, ActivityID: valueobjects.NewID(), UserID: userID})
	require.NoError(t, err)
}

func TestExperienceEducationCourseActivityQueryHandlersRejectMissingIDs(t *testing.T) {
	_, userID := testResumeIDs()
	zero := valueobjects.ID{}
	_, err := NewListExperiencesQueryHandler(stubResumeReader{}).Handle(context.Background(), ListExperiencesQuery{ResumeID: zero, UserID: userID})
	require.Error(t, err)
	_, err = NewGetEducationQueryHandler(stubResumeReader{}).Handle(context.Background(), GetEducationQuery{ResumeID: valueobjects.NewID(), EducationID: zero, UserID: userID})
	require.Error(t, err)
	_, err = NewListCoursesQueryHandler(stubResumeReader{}).Handle(context.Background(), ListCoursesQuery{ResumeID: valueobjects.NewID()})
	require.Error(t, err)
	_, err = NewGetActivityQueryHandler(stubResumeReader{}).Handle(context.Background(), GetActivityQuery{ResumeID: valueobjects.NewID(), ActivityID: valueobjects.NewID()})
	require.Error(t, err)
}
