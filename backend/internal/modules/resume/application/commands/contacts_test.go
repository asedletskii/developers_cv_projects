package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	contactkey "developers-cv/backend/internal/shared/domain/contact"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestReplaceContactsHandlerUsesVersionAndAggregate(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{}

	err := NewReplaceContactsHandler(repo, uow).Handle(context.Background(), ReplaceContactsCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 1,
		Input: ContactsInput{Email: func() *domain.Email {
			v, e := domain.NewEmail("john@example.com")
			require.NoError(t, e)
			return &v
		}(), HiddenContacts: []contactkey.Key{contactkey.Email}},
	})
	require.NoError(t, err)
	require.True(t, repo.saveCalled)
	require.Equal(t, int64(1), repo.version)
	require.NotNil(t, repo.resume.Contacts())
}

func TestReplaceContactsHandlerRejectsStaleVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{}

	email, err := domain.NewEmail("stale@example.com")
	require.NoError(t, err)

	err = NewReplaceContactsHandler(repo, uow).Handle(context.Background(), ReplaceContactsCommand{
		ResumeID: resume.ID(),
		UserID:   uid,
		Version:  2,
		Input:    ContactsInput{Email: &email},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saveCalled)
}
