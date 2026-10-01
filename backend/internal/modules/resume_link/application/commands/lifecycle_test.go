package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	valueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type lifecycleCommandRepositoryStub struct {
	link       *domain.ResumeLink
	getErr     error
	replaced   bool
	active     bool
	version    int64
	replaceErr error
}

func (s *lifecycleCommandRepositoryStub) Create(context.Context, *domain.ResumeLink) error {
	return nil
}
func (s *lifecycleCommandRepositoryStub) GetByID(context.Context, valueobjects.ID) (*domain.ResumeLink, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.link, nil
}
func (s *lifecycleCommandRepositoryStub) Delete(context.Context, valueobjects.ID) error { return nil }
func (s *lifecycleCommandRepositoryStub) ReplaceStatus(_ context.Context, _ valueobjects.ID, active bool, version int64) error {
	s.replaced, s.active, s.version = true, active, version
	return s.replaceErr
}
func (s *lifecycleCommandRepositoryStub) RecordEvent(context.Context, *domain.ResumeLinkEvent) error {
	return nil
}

var _ ports.ResumeLinkRepository = (*lifecycleCommandRepositoryStub)(nil)

type lifecycleCommandOwnershipStub struct{ err error }

func (s lifecycleCommandOwnershipStub) GetOwned(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	if s.err != nil {
		return resumequeries.ResumeView{}, s.err
	}
	return resumequeries.ResumeView{}, nil
}

type lifecycleCommandUOW struct{}

func (lifecycleCommandUOW) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestReplaceStatusResumeLinkDeactivatesOwnedLink(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: valueobjects.NewID(), ResumeID: valueobjects.NewID(), Name: "HH", Version: 2})
	require.NoError(t, err)
	repo := &lifecycleCommandRepositoryStub{link: link}
	err = NewReplaceStatusResumeLinkHandler(repo, lifecycleCommandOwnershipStub{}, lifecycleCommandUOW{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: valueobjects.NewID(), LinkID: link.ID(), Active: false, Version: 2,
	})
	require.NoError(t, err)
	require.True(t, repo.replaced)
	require.False(t, repo.active)
	require.Equal(t, int64(2), repo.version)
	require.False(t, link.IsActive())
}

func TestReplaceStatusResumeLinkReactivatesLink(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: valueobjects.NewID(), ResumeID: valueobjects.NewID(), Name: "HH", Version: 5})
	require.NoError(t, err)
	require.NoError(t, link.Deactivate())
	repo := &lifecycleCommandRepositoryStub{link: link}
	err = NewReplaceStatusResumeLinkHandler(repo, lifecycleCommandOwnershipStub{}, lifecycleCommandUOW{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: valueobjects.NewID(), LinkID: link.ID(), Active: true, Version: 5,
	})
	require.NoError(t, err)
	require.True(t, repo.active)
	require.True(t, link.IsActive())
}

func TestReplaceStatusResumeLinkChecksOwnership(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: valueobjects.NewID(), ResumeID: valueobjects.NewID(), Name: "HH"})
	require.NoError(t, err)
	repo := &lifecycleCommandRepositoryStub{link: link}
	err = NewReplaceStatusResumeLinkHandler(repo, lifecycleCommandOwnershipStub{err: resumequeries.ErrResumeNotFound}, lifecycleCommandUOW{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: valueobjects.NewID(), LinkID: link.ID(), Active: false, Version: 1,
	})
	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	require.False(t, repo.replaced)
}

func TestReplaceStatusResumeLinkRejectsStaleVersion(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: valueobjects.NewID(), ResumeID: valueobjects.NewID(), Name: "HH", Version: 2})
	require.NoError(t, err)
	repo := &lifecycleCommandRepositoryStub{link: link}
	err = NewReplaceStatusResumeLinkHandler(repo, lifecycleCommandOwnershipStub{}, lifecycleCommandUOW{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: valueobjects.NewID(), LinkID: link.ID(), Active: false, Version: 1,
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.replaced)
}

func TestReplaceStatusResumeLinkRejectsMissingLink(t *testing.T) {
	repo := &lifecycleCommandRepositoryStub{getErr: domain.ErrEntityNotFound}
	err := NewReplaceStatusResumeLinkHandler(repo, lifecycleCommandOwnershipStub{}, lifecycleCommandUOW{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: valueobjects.NewID(), LinkID: valueobjects.NewID(), Active: false, Version: 1,
	})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}

func TestReplaceStatusResumeLinkNoOpKeepsVersionAndSkipsPersistence(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: valueobjects.NewID(), ResumeID: valueobjects.NewID(), Name: "HH", Version: 7})
	require.NoError(t, err)
	repo := &lifecycleCommandRepositoryStub{link: link}

	err = NewReplaceStatusResumeLinkHandler(repo, lifecycleCommandOwnershipStub{}, lifecycleCommandUOW{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: valueobjects.NewID(), LinkID: link.ID(), Active: true, Version: 7,
	})

	require.NoError(t, err)
	require.False(t, repo.replaced)
	require.True(t, link.IsActive())
	require.Equal(t, int64(7), link.Version())
}

func TestReplaceStatusResumeLinkNoOpStillRejectsStaleVersion(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: valueobjects.NewID(), ResumeID: valueobjects.NewID(), Name: "HH", Version: 7})
	require.NoError(t, err)
	repo := &lifecycleCommandRepositoryStub{link: link}

	err = NewReplaceStatusResumeLinkHandler(repo, lifecycleCommandOwnershipStub{}, lifecycleCommandUOW{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: valueobjects.NewID(), LinkID: link.ID(), Active: true, Version: 6,
	})

	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.replaced)
}
