package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestCourseValidation(t *testing.T) {
	endDate, err := ParseYearMonth("2024-06")
	require.NoError(t, err)
	certificate, err := NewExternalURL("https://example.com/certificate")
	require.NoError(t, err)
	hours := 40
	item, err := NewCourse(CourseParams{ID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(), Name: "Go", Organization: "Academy", EndDate: endDate, Hours: &hours, CertificateLink: &certificate})
	require.NoError(t, err)
	require.Equal(t, "2024-06", item.EndDate().String())
	require.Equal(t, certificate.String(), item.CertificateLink().String())
}

func courseParams(t *testing.T, id sharedvalueobjects.ID, name string) CourseParams {
	t.Helper()
	end, err := ParseYearMonth("2024-06")
	require.NoError(t, err)
	return CourseParams{ID: id, Name: name, Organization: "Academy", EndDate: end}
}

func TestResumeReplaceCoursesCreatesAndPreservesIDs(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceCourses([]CourseParams{courseParams(t, sharedvalueobjects.ID{}, "Go")}))
	first := resume.Courses()[0].ID()
	require.NoError(t, resume.ReplaceCourses([]CourseParams{courseParams(t, first, "Advanced Go")}))
	require.Equal(t, first, resume.Courses()[0].ID())
	require.Equal(t, "Advanced Go", resume.Courses()[0].Name())
}

func TestResumeReplaceCoursesRejectsUnknownAndDuplicateIDs(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceCourses([]CourseParams{courseParams(t, sharedvalueobjects.ID{}, "Go")}))
	_, existing := resume.Courses()[0], resume.Courses()[0].ID()
	err := resume.ReplaceCourses([]CourseParams{courseParams(t, sharedvalueobjects.NewID(), "Unknown")})
	require.ErrorIs(t, err, ErrResumeInvariant)
	err = resume.ReplaceCourses([]CourseParams{courseParams(t, existing, "A"), courseParams(t, existing, "B")})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceCoursesRemovesMissingItems(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceCourses([]CourseParams{courseParams(t, sharedvalueobjects.ID{}, "Go"), courseParams(t, sharedvalueobjects.ID{}, "Docker")}))
	keep := resume.Courses()[1].ID()
	require.NoError(t, resume.ReplaceCourses([]CourseParams{courseParams(t, keep, "Docker")}))
	require.Len(t, resume.Courses(), 1)
	require.Equal(t, keep, resume.Courses()[0].ID())
}
