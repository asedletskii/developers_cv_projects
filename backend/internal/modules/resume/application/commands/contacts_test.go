package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeContactsAggregateRepository struct {
	resume  *domain.Resume
	version int64
	saved   bool
}

func (r *fakeContactsAggregateRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}

func (r *fakeContactsAggregateRepository) Save(_ context.Context, resume *domain.Resume) error {
	r.resume, r.version, r.saved = resume, resume.Version(), true
	return nil
}

func (r *fakeContactsAggregateRepository) SaveSection(ctx context.Context, resume *domain.Resume, _ ports.ResumeMutationSection) error {
	return r.Save(ctx, resume)
}
func (r *fakeContactsAggregateRepository) Create(context.Context, *domain.Resume) error { return nil }
func (r *fakeContactsAggregateRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

func TestReplaceContactsHandlerUsesVersionAndAggregate(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeContactsAggregateRepository{resume: resume}
	uow := &fakeCommandUnitOfWork{}

	err := NewReplaceContactsHandler(repo, uow).Handle(context.Background(), ReplaceContactsCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 1,
		Input: ContactsInput{Email: func() *domain.Email {
			v, e := domain.NewEmail("john@example.com")
			require.NoError(t, e)
			return &v
		}(), HiddenContacts: []domain.ContactKey{domain.ContactKeyEmail}},
	})
	require.NoError(t, err)
	require.True(t, repo.saved)
	require.Equal(t, int64(1), repo.version)
	require.NotNil(t, repo.resume.Contacts())
}

func TestReplaceContactsHandlerRejectsStaleVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeContactsAggregateRepository{resume: resume}
	uow := &fakeCommandUnitOfWork{}

	email, err := domain.NewEmail("stale@example.com")
	require.NoError(t, err)

	err = NewReplaceContactsHandler(repo, uow).Handle(context.Background(), ReplaceContactsCommand{
		ResumeID: resume.ID(),
		UserID:   uid,
		Version:  2,
		Input:    ContactsInput{Email: &email},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saved)
}
