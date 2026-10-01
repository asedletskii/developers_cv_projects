package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeLanguagesAggregateRepository struct {
	resume  *domain.Resume
	saved   bool
	version int64
}

func (r *fakeLanguagesAggregateRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}
func (r *fakeLanguagesAggregateRepository) Save(_ context.Context, resume *domain.Resume) error {
	r.resume, r.version, r.saved = resume, resume.Version(), true
	return nil
}
func (r *fakeLanguagesAggregateRepository) Create(context.Context, *domain.Resume) error { return nil }
func (r *fakeLanguagesAggregateRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

func TestReplaceLanguagesHandlerUsesAggregateAndVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeLanguagesAggregateRepository{resume: resume}
	err := NewReplaceLanguagesHandler(repo, &fakeCommandUnitOfWork{}).Handle(context.Background(), ReplaceLanguagesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 1,
		Input: LanguagesInput{Entries: []LanguageInput{{Language: domain.SpokenLanguageEN, Position: 0}}},
	})
	require.NoError(t, err)
	require.True(t, repo.saved)
	require.Equal(t, int64(1), repo.version)
	require.Len(t, repo.resume.Languages().Entries(), 1)
}

func TestReplaceLanguagesHandlerRejectsStaleVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeLanguagesAggregateRepository{resume: resume}
	err := NewReplaceLanguagesHandler(repo, &fakeCommandUnitOfWork{}).Handle(context.Background(), ReplaceLanguagesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 2,
		Input: LanguagesInput{Entries: []LanguageInput{{Language: domain.SpokenLanguageEN, Position: 0}}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saved)
}
