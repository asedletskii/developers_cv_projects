package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type unitOfWorkStub struct {
	called bool
	calls  int
	err    error
}

func (u *unitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	u.called = true
	u.calls++
	if u.err != nil {
		return u.err
	}
	return fn(ctx)
}

type resumeLinkRepositoryStub struct {
	link         *domain.ResumeLink
	err          error
	getErr       error
	count        int64
	countErr     error
	createErr    error
	deleteErr    error
	delErr       error
	replaceErr   error
	returnNil    bool
	lockCalled   bool
	countCalled  bool
	called       bool
	createCalled bool
	deleted      sharedvalueobjects.ID
	replaced     bool
	active       bool
	version      int64
}

func validCreateCommand() CreateResumeLinkCommand {
	return CreateResumeLinkCommand{UserID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(), Name: "  HH  "}
}

func (r *resumeLinkRepositoryStub) LockResumeForLinkCreation(context.Context, sharedvalueobjects.ID) error {
	r.lockCalled = true
	return nil
}
func (r *resumeLinkRepositoryStub) CountByResumeID(context.Context, sharedvalueobjects.ID) (int64, error) {
	r.countCalled = true
	return r.count, r.countErr
}
func (r *resumeLinkRepositoryStub) Create(_ context.Context, link *domain.ResumeLink) error {
	r.called = true
	r.createCalled = true
	if r.createErr != nil {
		return r.createErr
	}
	if r.err != nil {
		return r.err
	}
	r.link = link
	return nil
}
func (r *resumeLinkRepositoryStub) GetByID(context.Context, sharedvalueobjects.ID) (*domain.ResumeLink, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	if r.err != nil {
		return nil, r.err
	}
	if r.returnNil {
		return nil, nil
	}
	return r.link, nil
}
func (r *resumeLinkRepositoryStub) Delete(_ context.Context, id sharedvalueobjects.ID) error {
	r.deleted = id
	if r.deleteErr != nil {
		return r.deleteErr
	}
	if r.delErr != nil {
		return r.delErr
	}
	return r.err
}
func (r *resumeLinkRepositoryStub) ReplaceStatus(_ context.Context, _ sharedvalueobjects.ID, active bool, version int64) error {
	r.replaced, r.active, r.version = true, active, version
	if r.replaceErr != nil {
		return r.replaceErr
	}
	return r.err
}

type resumeOwnershipReaderStub struct {
	err   error
	query resumequeries.GetResumeQuery
}

func (r *resumeOwnershipReaderStub) GetOwned(_ context.Context, query resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	r.query = query
	if r.err != nil {
		return resumequeries.ResumeView{}, r.err
	}
	return resumequeries.ResumeView{ID: query.ResumeID.String()}, nil
}
func (r *resumeOwnershipReaderStub) OwnsResume(_ context.Context, query resumequeries.GetResumeQuery) (bool, error) {
	r.query = query
	return r.err == nil, r.err
}

type resumeReaderStub struct {
	view   resumequeries.ResumeView
	err    error
	called bool
	query  resumequeries.GetResumeQuery
}

func (r *resumeReaderStub) GetOwned(_ context.Context, query resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	r.called = true
	r.query = query
	if r.err != nil {
		return resumequeries.ResumeView{}, r.err
	}
	return r.view, nil
}
func (r *resumeReaderStub) OwnsResume(_ context.Context, query resumequeries.GetResumeQuery) (bool, error) {
	r.called = true
	r.query = query
	if r.err != nil {
		return false, r.err
	}
	return true, nil
}

func restoreCommandTestLink(t *testing.T, status domain.ResumeLinkStatus, version int64) *domain.ResumeLink {
	t.Helper()
	link, err := domain.RestoreResumeLink(domain.ResumeLinkSnapshot{
		ID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(),
		Name: "HH", Status: status, Version: version,
	})
	require.NoError(t, err)
	return link
}
