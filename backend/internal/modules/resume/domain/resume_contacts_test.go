package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestContactsRejectsDuplicateHiddenFields(t *testing.T) {
	tg, err := NewTelegram("telegram")
	require.NoError(t, err)
	c, err := NewContactMethods(ContactMethodsParams{Telegram: &tg})
	require.NoError(t, err)
	_, err = NewContacts(ContactsParams{ID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(), Contacts: c, HiddenContacts: []ContactKey{ContactKeyTelegram, ContactKeyTelegram}})
	require.Error(t, err)
}

func TestContactsKeepsOwnedStateSeparateFromPresentation(t *testing.T) {
	email, err := NewEmail("person@example.com")
	require.NoError(t, err)
	c, err := NewContactMethods(ContactMethodsParams{Email: &email})
	require.NoError(t, err)
	aggregate, err := NewContacts(ContactsParams{ID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(), Contacts: c})
	require.NoError(t, err)
	require.Equal(t, "person@example.com", aggregate.Contacts().Email().String())
}

func TestContactsTextRejectsWhitespaceOnlyValues(t *testing.T) {
	value := "   "
	_, err := NewContactsText(ContactsTextParams{Title: &value})
	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestContactsClonePreservesContactsText(t *testing.T) {
	title := "  Contacts  "
	text, err := NewContactsText(ContactsTextParams{Title: &title})
	require.NoError(t, err)

	aggregate, err := NewContacts(ContactsParams{
		ID:       sharedvalueobjects.NewID(),
		ResumeID: sharedvalueobjects.NewID(),
		Contacts: func() ContactMethods {
			contacts, createErr := NewContactMethods(ContactMethodsParams{Phone: "+33123456789"})
			require.NoError(t, createErr)
			return contacts
		}(),
		Text: &text,
	})
	require.NoError(t, err)
	require.NotNil(t, aggregate.Text())
	require.Equal(t, "Contacts", *aggregate.Text().Title())
}

func TestContactsRejectsHiddenContactThatIsNotPresent(t *testing.T) {
	email, err := NewEmail("person@example.com")
	require.NoError(t, err)
	contacts, err := NewContactMethods(ContactMethodsParams{Email: &email})
	require.NoError(t, err)

	_, err = NewContacts(ContactsParams{
		ID:             sharedvalueobjects.NewID(),
		ResumeID:       sharedvalueobjects.NewID(),
		Contacts:       contacts,
		HiddenContacts: []ContactKey{ContactKeyPhone},
	})

	require.ErrorIs(t, err, ErrResumeInvariant)
	require.ErrorContains(t, err, "hidden contact is not filled")
}

func TestContactsRejectsEmptyContactSet(t *testing.T) {
	_, err := NewContacts(ContactsParams{
		ID:       sharedvalueobjects.NewID(),
		ResumeID: sharedvalueobjects.NewID(),
	})

	require.ErrorIs(t, err, ErrInvalidResume)
}

func newContactsParamsForReplacement(t *testing.T, id sharedvalueobjects.ID, resumeID sharedvalueobjects.ID, email string) ContactsParams {
	t.Helper()
	emailValue, err := NewEmail(email)
	require.NoError(t, err)
	entity, err := NewContactMethods(ContactMethodsParams{Email: &emailValue})
	require.NoError(t, err)
	return ContactsParams{ID: id, ResumeID: resumeID, Contacts: entity, HiddenContacts: []ContactKey{ContactKeyEmail}}
}

func TestResumeReplaceContactsCreatesAndPreservesID(t *testing.T) {
	resume := newTestResume(t)

	require.NoError(t, resume.ReplaceContacts(newContactsParamsForReplacement(t, sharedvalueobjects.ID{}, resume.ID(), "first@example.com")))
	firstID := resume.Contacts().ID()

	require.NoError(t, resume.ReplaceContacts(newContactsParamsForReplacement(t, sharedvalueobjects.ID{}, resume.ID(), "second@example.com")))
	require.Equal(t, firstID, resume.Contacts().ID())
	require.Equal(t, "second@example.com", resume.Contacts().Contacts().Email().String())
}

func TestResumeReplaceContactsRejectsUnknownID(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceContacts(newContactsParamsForReplacement(t, sharedvalueobjects.ID{}, resume.ID(), "first@example.com")))

	err := resume.ReplaceContacts(newContactsParamsForReplacement(t, sharedvalueobjects.NewID(), resume.ID(), "second@example.com"))
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceContactsOwnsResumeID(t *testing.T) {
	resume := newTestResume(t)
	params := newContactsParamsForReplacement(t, sharedvalueobjects.ID{}, sharedvalueobjects.NewID(), "first@example.com")

	require.NoError(t, resume.ReplaceContacts(params))
	require.Equal(t, resume.ID(), resume.Contacts().ResumeID())
}

func TestContactsValidationAndReplacement(t *testing.T) {
	id, resumeID := sharedvalueobjects.NewID(), sharedvalueobjects.NewID()
	contactsEntity, err := NewContactMethods(ContactMethodsParams{Phone: "+33123456789"})
	require.NoError(t, err)
	contacts, err := NewContacts(ContactsParams{ID: id, ResumeID: resumeID, Contacts: contactsEntity, HiddenContacts: []ContactKey{ContactKeyPhone}})
	require.NoError(t, err)

	email, err := NewEmail("person@example.com")
	require.NoError(t, err)
	replacementEntity, err := NewContactMethods(ContactMethodsParams{
		Phone: "+33123456789",
		Email: &email,
	})
	require.NoError(t, err)

	err = contacts.Replace(ContactsParams{Contacts: replacementEntity, HiddenContacts: []ContactKey{ContactKeyEmail}})
	require.NoError(t, err)
	require.Equal(t, id, contacts.ID())
	require.Equal(t, resumeID, contacts.ResumeID())
	require.Equal(t, []ContactKey{ContactKeyEmail}, contacts.HiddenContacts())

	_, err = NewContacts(ContactsParams{ID: id, ResumeID: resumeID, Contacts: contactsEntity, HiddenContacts: []ContactKey{ContactKeyEmail, ContactKeyEmail}})
	require.Error(t, err)
}

func TestTechnologyTrimsAndDefensivelyCopiesGrade(t *testing.T) {
	grade := TechnologyGradeSenior
	technology, err := NewTechnology(TechnologyParams{Name: " Go ", Grade: &grade})
	require.NoError(t, err)
	require.Equal(t, "Go", technology.Name())
	got := technology.Grade()
	require.NotNil(t, got)
	*got = TechnologyGradeJunior
	require.Equal(t, TechnologyGradeSenior, *technology.Grade())

	_, err = NewTechnology(TechnologyParams{Name: "   "})
	require.Error(t, err)
}

func TestContributorValidation(t *testing.T) {
	email, err := NewEmail("user@example.com")
	require.NoError(t, err)
	contacts, err := NewContactMethods(ContactMethodsParams{Email: &email})
	require.NoError(t, err)

	contributor, err := NewContributor(ContributorParams{Name: " John ", Surname: " Doe ", Role: " Developer ", Contacts: contacts})
	require.NoError(t, err)
	require.Equal(t, "John", contributor.Name())
	require.Equal(t, "Doe", contributor.Surname())
	require.Equal(t, "Developer", contributor.Role())

	_, err = NewContributor(ContributorParams{Name: "John", Surname: "Doe", Role: "Developer"})
	require.Error(t, err)
	require.ErrorContains(t, err, "contributor requires at least one contact")
}

func TestContactsRequireAtLeastOneValue(t *testing.T) {
	_, err := NewContacts(ContactsParams{})
	require.Error(t, err)

	contacts, err := NewContactMethods(ContactMethodsParams{Phone: " +33123456789 "})
	require.NoError(t, err)
	phone := contacts.Phone()
	require.NotNil(t, phone)
	require.Equal(t, "+33123456789", *phone)
}

func TestContactsUseValidatedEmailAndTelegram(t *testing.T) {
	email, err := NewEmail("user@example.com")
	require.NoError(t, err)
	telegram, err := NewTelegram("@developer")
	require.NoError(t, err)
	contacts, err := NewContactMethods(ContactMethodsParams{Email: &email, Telegram: &telegram})
	require.NoError(t, err)
	require.Equal(t, "user@example.com", contacts.Email().String())
	require.Equal(t, "developer", contacts.Telegram().String())
}

func TestNewContactsTextRejectsValuesThatNormalizeToEmpty(t *testing.T) {
	empty := "   "

	_, err := NewContactsText(ContactsTextParams{Title: &empty})

	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestNewContactsTextNormalizesAndClonesValues(t *testing.T) {
	title := "  Contacts  "
	subtitle := "  Reach me  "

	value, err := NewContactsText(ContactsTextParams{Title: &title, Subtitle: &subtitle})
	require.NoError(t, err)
	require.Equal(t, "Contacts", *value.Title())
	require.Equal(t, "Reach me", *value.Subtitle())

	clone := value.Clone()
	require.Equal(t, "Contacts", *clone.Title())
	require.Equal(t, "Reach me", *clone.Subtitle())
}
