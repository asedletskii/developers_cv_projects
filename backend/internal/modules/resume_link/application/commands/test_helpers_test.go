package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
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
	link          *domain.ResumeLink
	err           error
	getErr        error
	count         int64
	countErr      error
	createErr     error
	deleteErr     error
	delErr        error
	replaceErr    error
	renameErr     error
	returnNil     bool
	countCalled   bool
	called        bool
	createCalled  bool
	deleted       sharedvalueobjects.ID
	replaced      bool
	renamed       bool
	active        bool
	version       int64
	renameName    string
	renameVersion int64
}

func validCreateCommand() CreateResumeLinkCommand {
	return CreateResumeLinkCommand{UserID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(), Name: "  HH  "}
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
func (r *resumeLinkRepositoryStub) GetByIDForUpdate(context.Context, sharedvalueobjects.ID) (*domain.ResumeLink, error) {
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
func (r *resumeLinkRepositoryStub) Rename(_ context.Context, _ sharedvalueobjects.ID, name string, version int64) error {
	r.renamed = true
	r.renameName = name
	r.renameVersion = version
	if r.renameErr != nil {
		return r.renameErr
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

type analyticsCleanupStub struct {
	called  bool
	linkIDs []sharedvalueobjects.ID
	err     error
}

func (s *analyticsCleanupStub) DeleteByLinkIDs(_ context.Context, linkIDs []sharedvalueobjects.ID) error {
	s.called = true
	s.linkIDs = append([]sharedvalueobjects.ID(nil), linkIDs...)
	return s.err
}

type resumeOwnershipReaderStub struct {
	err        error
	query      resumecontracts.OwnedResumeLock
	lockCalled bool
}

func (r *resumeOwnershipReaderStub) LockOwnedResumeForUpdate(_ context.Context, query resumecontracts.OwnedResumeLock) (bool, error) {
	r.lockCalled = true
	r.query = query
	if r.err != nil {
		if errors.Is(r.err, resumecontracts.ErrResumeNotFound) {
			return false, nil
		}
		return false, r.err
	}
	return true, nil
}

type resumeReaderStub struct {
	err       error
	called    bool
	notOwned  bool
	lockQuery resumecontracts.OwnedResumeLock
}

func (r *resumeReaderStub) LockOwnedResumeForUpdate(_ context.Context, query resumecontracts.OwnedResumeLock) (bool, error) {
	r.called = true
	r.lockQuery = query
	if r.err != nil {
		if errors.Is(r.err, resumecontracts.ErrResumeNotFound) {
			return false, nil
		}
		return false, r.err
	}
	return !r.notOwned, nil
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
