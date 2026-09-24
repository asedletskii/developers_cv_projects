package commands

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type updateRepositoryStub struct {
	resume       *domain.Resume
	saved        *domain.Resume
	err          error
	title        string
	updateCalls  int
	lastUserID   useridentity.ID
	lastResumeID valueobjects.ID
}

func (r *updateRepositoryStub) Get(_ context.Context, query GetResumeRepositoryQuery) (*domain.Resume, error) {
	r.lastUserID = query.UserID
	r.lastResumeID = query.ResumeID
	if r.err != nil {
		return nil, r.err
	}
	return r.resume, nil
}

func (r *updateRepositoryStub) Update(_ context.Context, command UpdateResumeRepositoryCommand) error {
	r.saved = command.Resume
	return nil
}

func (r *updateRepositoryStub) UpdateTitle(_ context.Context, command UpdateResumeTitleRepositoryCommand) error {
	r.updateCalls++
	r.lastUserID = command.UserID
	r.lastResumeID = command.ResumeID
	r.title = command.Title
	return r.err
}

type updateUnitOfWorkStub struct{}

func (updateUnitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestUpdateResumeHandlerPreservesTitleAndUsesOwnerScope(t *testing.T) {
	owner := useridentity.NewID()
	resume := testUpdateResume(t, owner, "Manually renamed")
	repo := &updateRepositoryStub{resume: resume}
	handler := NewUpdateResumeHandler(repo, repo, updateUnitOfWorkStub{})

	id := resume.ID()
	command := UpdateResumeCommand{
		ResumeID:    id,
		UserID:      owner,
		Name:        "Updated",
		Surname:     "User",
		Position:    "Staff Go Developer",
		Greeting:    resume.Greeting(),
		PhotoPath:   resume.PhotoPath(),
		DateOfBirth: resume.DateOfBirth(),
		Location:    resume.Location(),
		Contacts:    resume.Contacts(),
	}

	updatedID, err := handler.Handle(context.Background(), command)
	require.NoError(t, err)
	require.Equal(t, id, updatedID)
	require.NotNil(t, repo.saved)
	require.Equal(t, "Updated", repo.saved.Name())
	require.Equal(t, "Staff Go Developer", repo.saved.Position())
	require.Equal(t, "Manually renamed", repo.saved.Title())
	require.Equal(t, owner.UUID(), repo.saved.UserID().UUID())
}

func TestUpdateResumeHandlerReturnsNotFoundForForeignResume(t *testing.T) {
	owner := useridentity.NewID()
	otherUser := useridentity.NewID()
	resume := testUpdateResume(t, owner, "Manual title")
	repo := &updateRepositoryStub{err: domain.ErrEntityNotFound}
	handler := NewUpdateResumeHandler(repo, repo, updateUnitOfWorkStub{})

	_, err := handler.Handle(context.Background(), UpdateResumeCommand{
		ResumeID: resume.ID(),
		UserID:   otherUser,
	})

	require.ErrorIs(t, err, domain.ErrEntityNotFound)
	require.Equal(t, otherUser, repo.lastUserID)
	require.Equal(t, resume.ID(), repo.lastResumeID)
}

func TestPatchResumeHandlerUpdatesOnlyTitleWithinOwnerScope(t *testing.T) {
	owner := useridentity.NewID()
	resume := testUpdateResume(t, owner, "Old title")
	repo := &updateRepositoryStub{resume: resume}
	handler := NewPatchResumeHandler(repo, updateUnitOfWorkStub{})

	result, err := handler.Handle(context.Background(), PatchResumeCommand{
		ResumeID: resume.ID(),
		UserID:   owner,
		Title:    " New title ",
	})
	require.NoError(t, err)
	require.Equal(t, resume.ID(), result.ResumeID)
	require.Equal(t, "New title", result.Title)
	require.Equal(t, "New title", repo.title)
	require.Equal(t, 1, repo.updateCalls)
	require.Equal(t, owner, repo.lastUserID)
	require.Equal(t, resume.ID(), repo.lastResumeID)
	require.Nil(t, repo.saved)
}

func TestPatchResumeHandlerRejectsEmptyTitleWithoutRepositoryCall(t *testing.T) {
	repo := &updateRepositoryStub{}
	handler := NewPatchResumeHandler(repo, updateUnitOfWorkStub{})

	_, err := handler.Handle(context.Background(), PatchResumeCommand{
		ResumeID: valueobjects.NewID(),
		UserID:   useridentity.NewID(),
		Title:    "   ",
	})

	require.Error(t, err)
	require.Zero(t, repo.updateCalls)
}

func testUpdateResume(t *testing.T, userID useridentity.ID, title string) *domain.Resume {
	t.Helper()
	id, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	date, err := valueobjects.ParseDate("1998-04-12")
	require.NoError(t, err)
	photo, err := valueobjects.NewAssetPath("/assets/photo.svg")
	require.NoError(t, err)
	contacts, err := entities.NewContacts(entities.ContactsParams{Phone: "+79999999999"})
	require.NoError(t, err)
	resume, err := domain.NewResume(domain.ResumeParams{
		ID: id, UserID: userID, Name: "User", Surname: "Test", Position: "Senior Go Developer", Title: title,
		PhotoPath: photo, DateOfBirth: date, Location: "Remote", Contacts: contacts,
	})
	require.NoError(t, err)
	return resume
}
