package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeOverviewAggregateRepository struct {
	resume  *domain.Resume
	saved   bool
	version int64
}

func (r *fakeOverviewAggregateRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}
func (r *fakeOverviewAggregateRepository) Save(_ context.Context, resume *domain.Resume) error {
	r.saved = true
	r.resume = resume
	r.version = resume.Version()
	return nil
}
func (r *fakeOverviewAggregateRepository) Create(context.Context, *domain.Resume) error { return nil }
func (r *fakeOverviewAggregateRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

type overviewUnitOfWork struct{}

func (overviewUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestReplaceOverviewHandlerLoadsReplacesAndSavesAggregate(t *testing.T) {
	resume, err := domain.NewResume(domain.ResumeParams{UserID: sharedvalueobjects.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	overview := domain.ResumeOverviewParams{Greeting: domain.NewGreeting(domain.GreetingParams{First: "Old"}), Sections: []domain.SectionSettings{}}
	require.NoError(t, resume.ReplaceOverview(overview))
	repo := &fakeOverviewAggregateRepository{resume: resume}
	h := NewReplaceOverviewHandler(repo, overviewUnitOfWork{})
	err = h.Handle(context.Background(), ReplaceOverviewCommand{ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1, Input: OverviewInput{Greeting: GreetingInput{First: "New"}, Sections: []SectionInput{}}})
	require.NoError(t, err)
	require.True(t, repo.saved)
	require.Equal(t, int64(1), repo.version)
	require.Equal(t, "New", repo.resume.Overview().Greeting().First())
}

func TestReplaceOverviewHandlerRejectsStaleVersionBeforeSave(t *testing.T) {
	resume, err := domain.NewResume(domain.ResumeParams{UserID: sharedvalueobjects.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	overview := domain.ResumeOverviewParams{Greeting: domain.NewGreeting(domain.GreetingParams{}), Sections: []domain.SectionSettings{}}
	require.NoError(t, resume.ReplaceOverview(overview))
	repo := &fakeOverviewAggregateRepository{resume: resume}
	h := NewReplaceOverviewHandler(repo, overviewUnitOfWork{})
	err = h.Handle(context.Background(), ReplaceOverviewCommand{ResumeID: resume.ID(), UserID: resume.UserID(), Version: 2, Input: OverviewInput{Sections: []SectionInput{}}})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saved)
}
