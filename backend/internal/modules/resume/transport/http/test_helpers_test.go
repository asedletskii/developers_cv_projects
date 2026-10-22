package http

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type stubListReader struct {
	items  []queries.ResumeListItem
	userID sharedvalueobjects.ID
}

func (r *stubListReader) ListResumes(_ context.Context, query queries.ListResumesQuery) ([]queries.ResumeListItem, error) {
	r.userID = query.UserID
	return r.items, nil
}

type stubOwnedResumeReader struct {
	view queries.ResumeView
}

func (r stubOwnedResumeReader) GetOwned(context.Context, queries.GetResumeQuery) (queries.ResumeView, error) {
	return r.view, nil
}

type stubResumeRepository struct {
	resume    *domain.Resume
	saved     *domain.Resume
	version   int64
	deleteErr error
}

func (r *stubResumeRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}

func (r *stubResumeRepository) Save(_ context.Context, resume *domain.Resume) error {
	r.saved = resume
	r.version = resume.Version()
	return nil
}

func (r *stubResumeRepository) SaveSection(ctx context.Context, resume *domain.Resume, _ ports.ResumeMutationSection) error {
	return r.Save(ctx, resume)
}

func (r *stubResumeRepository) Create(context.Context, *domain.Resume) error { return nil }

func (r *stubResumeRepository) LockForDelete(context.Context, ports.DeleteResumeParams) error {
	return r.deleteErr
}

func (r *stubResumeRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return r.deleteErr
}

type stubResumeLinkCleanup struct{}

func (stubResumeLinkCleanup) DeleteByResumeIDs(context.Context, []sharedvalueobjects.ID) error {
	return nil
}

type stubResumeUnitOfWork struct{}

func (stubResumeUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func newStubResume(t *testing.T, resumeID, userID uuid.UUID, version int64) *domain.Resume {
	t.Helper()
	id, err := sharedvalueobjects.IDFromUUID(resumeID)
	require.NoError(t, err)
	uid, err := sharedvalueobjects.IDFromUUID(userID)
	require.NoError(t, err)
	return domain.RestoreResume(domain.ResumeSnapshot{
		ID: id, UserID: uid, Title: "Old title", Language: domain.ResumeLanguageEN, Status: domain.ResumeStatusDraft, Version: version,
	})
}

func int64Ptr(value int64) *int64 {
	return &value
}
