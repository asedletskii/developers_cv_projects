package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestReplaceCoursesHandlerUsesAggregateRepository(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}
	end, err := domain.ParseYearMonth("2024-06")
	require.NoError(t, err)
	err = NewReplaceCoursesHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceCoursesCommand{ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1, Input: CoursesInput{Items: []CourseInput{{Name: "Go", Organization: "Academy", EndDate: end}}}})
	require.NoError(t, err)
	require.True(t, repo.saveCalled)
	require.Len(t, resume.Courses(), 1)
}
func TestReplaceCoursesHandlerRejectsStaleVersion(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}
	end, err := domain.ParseYearMonth("2024-06")
	require.NoError(t, err)
	err = NewReplaceCoursesHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceCoursesCommand{ResumeID: resume.ID(), UserID: resume.UserID(), Version: 2, Input: CoursesInput{Items: []CourseInput{{Name: "Go", Organization: "Academy", EndDate: end}}}})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saveCalled)
}
