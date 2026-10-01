package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type courseRepoStub struct {
	resume *domain.Resume
	saved  bool
}

func (s *courseRepoStub) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return s.resume, nil
}
func (s *courseRepoStub) Save(context.Context, *domain.Resume) error {
	s.saved = true
	return nil
}
func (s *courseRepoStub) Create(context.Context, *domain.Resume) error           { return nil }
func (s *courseRepoStub) Delete(context.Context, ports.DeleteResumeParams) error { return nil }

type courseUOWStub struct{}

func (courseUOWStub) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

func TestReplaceCoursesHandlerUsesAggregateRepository(t *testing.T) {
	resume, err := domain.NewResume(domain.ResumeParams{ID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &courseRepoStub{resume: resume}
	end, err := domain.ParseYearMonth("2024-06")
	require.NoError(t, err)
	err = NewReplaceCoursesHandler(repo, courseUOWStub{}).Handle(context.Background(), ReplaceCoursesCommand{ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1, Input: CoursesInput{Items: []CourseInput{{Name: "Go", Organization: "Academy", EndDate: end}}}})
	require.NoError(t, err)
	require.True(t, repo.saved)
	require.Len(t, resume.Courses(), 1)
}
func TestReplaceCoursesHandlerRejectsStaleVersion(t *testing.T) {
	resume, err := domain.NewResume(domain.ResumeParams{ID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &courseRepoStub{resume: resume}
	end, err := domain.ParseYearMonth("2024-06")
	require.NoError(t, err)
	err = NewReplaceCoursesHandler(repo, courseUOWStub{}).Handle(context.Background(), ReplaceCoursesCommand{ResumeID: resume.ID(), UserID: resume.UserID(), Version: 2, Input: CoursesInput{Items: []CourseInput{{Name: "Go", Organization: "Academy", EndDate: end}}}})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saved)
}
