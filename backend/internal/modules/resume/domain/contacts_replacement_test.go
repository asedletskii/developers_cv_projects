package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func newContactsParamsForReplacement(t *testing.T, id valueobjects.ID, resumeID valueobjects.ID, email string) ContactsParams {
	t.Helper()
	emailValue, err := valueobjects.NewEmail(email)
	require.NoError(t, err)
	entity, err := entities.NewContacts(entities.ContactsParams{Email: &emailValue})
	require.NoError(t, err)
	return ContactsParams{ID: id, ResumeID: resumeID, Contacts: entity, HiddenContacts: []enums.ContactKey{enums.ContactKeyEmail}}
}

func TestResumeReplaceContactsCreatesAndPreservesID(t *testing.T) {
	resume := newTestResumeForAggregate(t)

	require.NoError(t, resume.ReplaceContacts(newContactsParamsForReplacement(t, valueobjects.ID{}, resume.ID(), "first@example.com")))
	firstID := resume.Contacts().ID()

	require.NoError(t, resume.ReplaceContacts(newContactsParamsForReplacement(t, valueobjects.ID{}, resume.ID(), "second@example.com")))
	require.Equal(t, firstID, resume.Contacts().ID())
	require.Equal(t, "second@example.com", resume.Contacts().Contacts().Email().String())
}

func TestResumeReplaceContactsRejectsUnknownID(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceContacts(newContactsParamsForReplacement(t, valueobjects.ID{}, resume.ID(), "first@example.com")))

	err := resume.ReplaceContacts(newContactsParamsForReplacement(t, valueobjects.NewID(), resume.ID(), "second@example.com"))
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceContactsOwnsResumeID(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	params := newContactsParamsForReplacement(t, valueobjects.ID{}, valueobjects.NewID(), "first@example.com")

	require.NoError(t, resume.ReplaceContacts(params))
	require.Equal(t, resume.ID(), resume.Contacts().ResumeID())
}
