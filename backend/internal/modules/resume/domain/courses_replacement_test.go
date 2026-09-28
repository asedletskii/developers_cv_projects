package domain

import (
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"testing"

	"github.com/stretchr/testify/require"
)

func courseParams(t *testing.T, id valueobjects.ID, name string) CourseParams {
	t.Helper()
	end, err := valueobjects.ParseYearMonth("2024-06")
	require.NoError(t, err)
	return CourseParams{ID: id, Name: name, Organization: "Academy", EndDate: end}
}

func TestResumeReplaceCoursesCreatesAndPreservesIDs(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceCourses([]CourseParams{courseParams(t, valueobjects.ID{}, "Go")}))
	first := resume.Courses()[0].ID()
	require.NoError(t, resume.ReplaceCourses([]CourseParams{courseParams(t, first, "Advanced Go")}))
	require.Equal(t, first, resume.Courses()[0].ID())
	require.Equal(t, "Advanced Go", resume.Courses()[0].Name())
}

func TestResumeReplaceCoursesRejectsUnknownAndDuplicateIDs(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceCourses([]CourseParams{courseParams(t, valueobjects.ID{}, "Go")}))
	_, existing := resume.Courses()[0], resume.Courses()[0].ID()
	err := resume.ReplaceCourses([]CourseParams{courseParams(t, valueobjects.NewID(), "Unknown")})
	require.ErrorIs(t, err, ErrResumeInvariant)
	err = resume.ReplaceCourses([]CourseParams{courseParams(t, existing, "A"), courseParams(t, existing, "B")})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceCoursesRemovesMissingItems(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceCourses([]CourseParams{courseParams(t, valueobjects.ID{}, "Go"), courseParams(t, valueobjects.ID{}, "Docker")}))
	keep := resume.Courses()[1].ID()
	require.NoError(t, resume.ReplaceCourses([]CourseParams{courseParams(t, keep, "Docker")}))
	require.Len(t, resume.Courses(), 1)
	require.Equal(t, keep, resume.Courses()[0].ID())
}
