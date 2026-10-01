package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeProfileAggregateRepository struct {
	resume  *domain.Resume
	version int64
	saved   bool
}

func (r *fakeProfileAggregateRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}
func (r *fakeProfileAggregateRepository) Save(_ context.Context, resume *domain.Resume) error {
	r.resume = resume
	r.version = resume.Version()
	r.saved = true
	return nil
}
func (r *fakeProfileAggregateRepository) Create(context.Context, *domain.Resume) error { return nil }
func (r *fakeProfileAggregateRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}
func newTestResumeForCommand(t *testing.T, uid identity.ID) *domain.Resume {
	t.Helper()
	r, err := domain.NewResume(domain.ResumeParams{ID: identity.NewID(), UserID: uid, Title: "Resume", Language: domain.ResumeLanguageEN, Version: 1})
	require.NoError(t, err)
	return r
}

func TestReplaceProfileHandlerUsesVersionAndAggregate(t *testing.T) {
	uid := identity.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeProfileAggregateRepository{resume: resume}
	uow := &fakeCommandUnitOfWork{}
	photo := testURL(t)
	birth := testDate(t)
	err := NewReplaceProfileHandler(repo, uow).Handle(context.Background(), ReplaceProfileCommand{ResumeID: resume.ID(), UserID: uid, Version: 1, Input: ProfileInput{Name: "John", Surname: "Doe", Position: "Developer", PhotoPath: photo, DateOfBirth: birth, Location: "Paris", WorkFormats: []domain.WorkFormat{domain.WorkFormatRemote}}})
	require.NoError(t, err)
	require.True(t, repo.saved)
	require.Equal(t, int64(1), repo.version)
	require.NotNil(t, repo.resume.Profile())
}
func TestReplaceProfileHandlerRejectsStaleVersion(t *testing.T) {
	uid := identity.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeProfileAggregateRepository{resume: resume}
	uow := &fakeCommandUnitOfWork{}
	err := NewReplaceProfileHandler(repo, uow).Handle(context.Background(), ReplaceProfileCommand{ResumeID: resume.ID(), UserID: uid, Version: 2})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saved)
}
