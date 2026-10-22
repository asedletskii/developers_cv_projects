package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	contactkey "developers-cv/backend/internal/shared/domain/contact"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestContactsRejectsDuplicateHiddenFields(t *testing.T) {
	tg, err := NewTelegram("telegram")
	require.NoError(t, err)
	c, err := NewContactMethods(ContactMethodsParams{Telegram: &tg})
	require.NoError(t, err)
	_, err = newContacts(ContactsParams{ID: sharedvalueobjects.NewID(), Contacts: c, HiddenContacts: []contactkey.Key{contactkey.Telegram, contactkey.Telegram}}, sharedvalueobjects.NewID())
	require.Error(t, err)
}

func TestContactsKeepsOwnedStateSeparateFromPresentation(t *testing.T) {
	email, err := NewEmail("person@example.com")
	require.NoError(t, err)
	c, err := NewContactMethods(ContactMethodsParams{Email: &email})
	require.NoError(t, err)
	aggregate, err := newContacts(ContactsParams{ID: sharedvalueobjects.NewID(), Contacts: c}, sharedvalueobjects.NewID())
	require.NoError(t, err)
	require.Equal(t, "person@example.com", aggregate.Contacts().Email().String())
}

func TestContactsTextTreatsWhitespaceOnlyValueAsAbsent(t *testing.T) {
	value := "   "
	text, err := NewContactsText(ContactsTextParams{Title: &value})
	require.NoError(t, err)
	require.Nil(t, text)
}

func TestContactsClonePreservesContactsText(t *testing.T) {
	title := "  Contacts  "
	text, err := NewContactsText(ContactsTextParams{Title: &title})
	require.NoError(t, err)

	aggregate, err := newContacts(ContactsParams{
		ID: sharedvalueobjects.NewID(),
		Contacts: func() ContactMethods {
			contacts, createErr := NewContactMethods(ContactMethodsParams{Phone: "+33123456789"})
			require.NoError(t, createErr)
			return contacts
		}(),
		Text: text,
	}, sharedvalueobjects.NewID())
	require.NoError(t, err)
	require.NotNil(t, aggregate.Text())
	require.Equal(t, "Contacts", *aggregate.Text().Title())
}

func TestContactsRejectsHiddenContactThatIsNotPresent(t *testing.T) {
	email, err := NewEmail("person@example.com")
	require.NoError(t, err)
	contacts, err := NewContactMethods(ContactMethodsParams{Email: &email})
	require.NoError(t, err)

	_, err = newContacts(ContactsParams{
		ID:             sharedvalueobjects.NewID(),
		Contacts:       contacts,
		HiddenContacts: []contactkey.Key{contactkey.Phone},
	}, sharedvalueobjects.NewID())

	require.ErrorIs(t, err, ErrResumeInvariant)
	require.ErrorContains(t, err, "hidden contact is not filled")
}

func TestContactsRejectsEmptyContactSet(t *testing.T) {
	_, err := newContacts(ContactsParams{ID: sharedvalueobjects.NewID()}, sharedvalueobjects.NewID())

	require.ErrorIs(t, err, ErrInvalidResume)
}

func newContactsParamsForReplacement(t *testing.T, id sharedvalueobjects.ID, _ sharedvalueobjects.ID, email string) ContactsParams {
	t.Helper()
	emailValue, err := NewEmail(email)
	require.NoError(t, err)
	entity, err := NewContactMethods(ContactMethodsParams{Email: &emailValue})
	require.NoError(t, err)
	return ContactsParams{ID: id, Contacts: entity, HiddenContacts: []contactkey.Key{contactkey.Email}}
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

func TestContactsValidation(t *testing.T) {
	id, resumeID := sharedvalueobjects.NewID(), sharedvalueobjects.NewID()
	contactsEntity, err := NewContactMethods(ContactMethodsParams{Phone: "+33123456789"})
	require.NoError(t, err)
	_, err = newContacts(ContactsParams{ID: id, Contacts: contactsEntity, HiddenContacts: []contactkey.Key{contactkey.Phone}}, resumeID)
	require.NoError(t, err)

	_, err = newContacts(ContactsParams{ID: id, Contacts: contactsEntity, HiddenContacts: []contactkey.Key{contactkey.Email, contactkey.Email}}, resumeID)
	require.Error(t, err)
}

func TestTechnologyTrimsName(t *testing.T) {
	technology, err := NewTechnology(TechnologyParams{Name: " Go "})
	require.NoError(t, err)
	require.Equal(t, "Go", technology.Name())

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
	_, err := newContacts(ContactsParams{ID: sharedvalueobjects.NewID()}, sharedvalueobjects.NewID())
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

func TestNewContactsTextTreatsEmptyValuesAsMissing(t *testing.T) {
	empty := ""
	whitespace := "   "

	text, err := NewContactsText(ContactsTextParams{Title: &empty, Subtitle: &whitespace})
	require.NoError(t, err)
	require.Nil(t, text)
}

func TestNewContactsTextTreatsMissingValuesAsAbsent(t *testing.T) {
	text, err := NewContactsText(ContactsTextParams{})
	require.NoError(t, err)
	require.Nil(t, text)
}

func TestNewContactsTextDropsEmptyFields(t *testing.T) {
	empty := ""
	whitespace := "  "
	subtitle := "  Reach me  "

	text, err := NewContactsText(ContactsTextParams{Title: &empty, Subtitle: &subtitle})
	require.NoError(t, err)
	require.NotNil(t, text)
	require.Nil(t, text.Title())
	require.Equal(t, "Reach me", *text.Subtitle())

	text, err = NewContactsText(ContactsTextParams{Title: &whitespace, Subtitle: &empty})
	require.NoError(t, err)
	require.Nil(t, text)
}

func TestNewContactsTextNormalizesAndClonesValues(t *testing.T) {
	title := "  Contacts  "
	subtitle := "  Reach me  "

	value, err := NewContactsText(ContactsTextParams{Title: &title, Subtitle: &subtitle})
	require.NoError(t, err)
	require.Equal(t, "Contacts", *value.Title())
	require.Equal(t, "Reach me", *value.Subtitle())

	clone := value.clone()
	require.Equal(t, "Contacts", *clone.Title())
	require.Equal(t, "Reach me", *clone.Subtitle())
}

func TestRestoreContactsTextPreservesSnapshotValues(t *testing.T) {
	title := " Contacts "
	text := RestoreContactsText(ContactsTextParams{Title: &title})

	err := text.validate()

	require.NoError(t, err)
	require.Equal(t, " Contacts ", *text.Title())
}

func TestResumeReplaceContactsDoesNotMutateOnError(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceContacts(newContactsParamsForReplacement(t, sharedvalueobjects.ID{}, resume.ID(), "first@example.com")))
	existing := resume.Contacts()

	invalid := ContactsParams{ID: existing.ID()}
	err := resume.ReplaceContacts(invalid)

	require.ErrorIs(t, err, ErrInvalidResume)
	require.Equal(t, existing.ID(), resume.Contacts().ID())
	require.Equal(t, existing.Contacts().Email().String(), resume.Contacts().Contacts().Email().String())
}
