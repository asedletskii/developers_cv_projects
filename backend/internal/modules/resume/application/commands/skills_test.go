package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeSkillsAggregateRepository struct {
	resume  *domain.Resume
	version int64
	saved   bool
}

func (r *fakeSkillsAggregateRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}

func (r *fakeSkillsAggregateRepository) Save(_ context.Context, resume *domain.Resume) error {
	r.resume, r.version, r.saved = resume, resume.Version(), true
	return nil
}

func (r *fakeSkillsAggregateRepository) SaveSection(ctx context.Context, resume *domain.Resume, _ ports.ResumeMutationSection) error {
	return r.Save(ctx, resume)
}
func (r *fakeSkillsAggregateRepository) Create(context.Context, *domain.Resume) error { return nil }
func (r *fakeSkillsAggregateRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

func TestReplaceSkillsHandlerUsesVersionAndAggregate(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeSkillsAggregateRepository{resume: resume}
	uow := &fakeCommandUnitOfWork{}
	err := NewReplaceSkillsHandler(repo, uow).Handle(context.Background(), ReplaceSkillsCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 1,
		Input: SkillsInput{Hard: []string{"testing"}, Soft: []string{"teamwork"}},
	})
	require.NoError(t, err)
	require.True(t, repo.saved)
	require.Equal(t, int64(1), repo.version)
	require.NotNil(t, repo.resume.Skills())
	require.Equal(t, []string{"testing"}, repo.resume.Skills().HardSkills())
}

func TestReplaceSkillsHandlerRejectsStaleVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeSkillsAggregateRepository{resume: resume}
	uow := &fakeCommandUnitOfWork{}
	err := NewReplaceSkillsHandler(repo, uow).Handle(context.Background(), ReplaceSkillsCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 2,
		Input: SkillsInput{Hard: []string{"testing"}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saved)
}
