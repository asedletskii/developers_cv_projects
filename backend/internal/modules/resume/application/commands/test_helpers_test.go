package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type unitOfWorkStub struct {
	called bool
	err    error
}

func (u *unitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	u.called = true
	if u.err != nil {
		return u.err
	}
	return fn(ctx)
}

type resumeLinkCleanupStub struct {
	called    bool
	resumeIDs []sharedvalueobjects.ID
	err       error
}

func (s *resumeLinkCleanupStub) DeleteByResumeIDs(_ context.Context, resumeIDs []sharedvalueobjects.ID) error {
	s.called = true
	s.resumeIDs = append([]sharedvalueobjects.ID(nil), resumeIDs...)
	return s.err
}

type userExistenceCheckerStub struct {
	userIDs []sharedvalueobjects.ID
	err     error
}

func (s *userExistenceCheckerStub) EnsureUserExists(_ context.Context, userID sharedvalueobjects.ID) error {
	s.userIDs = append(s.userIDs, userID)
	return s.err
}

type commandRepositoryStub struct {
	resume *domain.Resume

	err       error
	getErr    error
	saveErr   error
	createErr error
	deleteErr error

	called       bool
	getCalled    bool
	saveCalled   bool
	createCalled bool
	deleteCalled bool
	saved        *domain.Resume
	version      int64
	deleted      ports.DeleteResumeParams

	advanceSectionVersion bool
}

func (r *commandRepositoryStub) Get(_ context.Context, _ ports.GetResumeParams) (*domain.Resume, error) {
	r.getCalled = true
	if r.getErr != nil {
		return nil, r.getErr
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.resume, nil
}

func (r *commandRepositoryStub) Save(_ context.Context, resume *domain.Resume) error {
	r.called = true
	r.saveCalled = true
	if r.saveErr != nil {
		return r.saveErr
	}
	if r.err != nil {
		return r.err
	}
	r.saved = resume
	r.resume = resume
	r.version = resume.Version()
	return nil
}

func (r *commandRepositoryStub) SaveSection(ctx context.Context, resume *domain.Resume, _ ports.ResumeMutationSection) error {
	if err := r.Save(ctx, resume); err != nil {
		return err
	}
	if r.advanceSectionVersion {
		r.resume = restoreCommandTestResumeVersion(resume, resume.Version()+1)
	}
	return nil
}

func (r *commandRepositoryStub) Create(_ context.Context, resume *domain.Resume) error {
	r.called = true
	r.createCalled = true
	if r.createErr != nil {
		return r.createErr
	}
	if r.err != nil {
		return r.err
	}
	r.resume = resume
	return nil
}

func (r *commandRepositoryStub) LockForDelete(_ context.Context, _ ports.DeleteResumeParams) error {
	r.called = true
	return r.err
}

func (r *commandRepositoryStub) Delete(_ context.Context, params ports.DeleteResumeParams) error {
	r.called = true
	r.deleteCalled = true
	r.deleted = params
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return r.err
}

func newCommandTestResume(t *testing.T, userID sharedvalueobjects.ID, title string, version int64) *domain.Resume {
	t.Helper()
	resume, err := domain.NewResume(domain.ResumeParams{
		UserID: userID, Title: title, Language: domain.ResumeLanguageEN,
	})
	require.NoError(t, err)
	if version == resume.Version() {
		return resume
	}
	return restoreCommandTestResumeVersion(resume, version)
}

func newMutationTestResume(t *testing.T) *domain.Resume {
	t.Helper()
	return newCommandTestResume(t, sharedvalueobjects.NewID(), "Original", 1)
}

func restoreCommandTestResumeVersion(resume *domain.Resume, version int64) *domain.Resume {
	return domain.RestoreResume(domain.ResumeSnapshot{
		ID: resume.ID(), UserID: resume.UserID(), Title: resume.Title(),
		Language: resume.Language(), Status: resume.Status(), Version: version,
		Overview: resume.Overview(), Profile: resume.Profile(), Contacts: resume.Contacts(),
		Skills: resume.Skills(), SpokenLanguages: resume.SpokenLanguages(), Technologies: resume.Technologies(),
		Experiences: resume.Experiences(), Educations: resume.Educations(), Courses: resume.Courses(),
		Activities: resume.Activities(), Projects: resume.Projects(),
	})
}

func testDate(t *testing.T) domain.Date {
	t.Helper()
	v, err := domain.ParseDate("2020-01-02")
	require.NoError(t, err)
	return v
}

func testURL(t *testing.T) domain.ExternalURL {
	t.Helper()
	v, err := domain.NewExternalURL("https://example.com/item")
	require.NoError(t, err)
	return v
}

func mustGreeting(t *testing.T, params domain.GreetingParams) domain.Greeting {
	t.Helper()
	value, err := domain.NewGreeting(params)
	require.NoError(t, err)
	return value
}
