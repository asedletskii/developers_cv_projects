package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestContactsRejectsDuplicateHiddenFields(t *testing.T) {
	tg, err := valueobjects.NewTelegram("telegram")
	require.NoError(t, err)
	c, err := entities.NewContacts(entities.ContactsParams{Telegram: &tg})
	require.NoError(t, err)
	_, err = domain.NewContacts(domain.ContactsParams{ID: valueobjects.NewID(), ResumeID: valueobjects.NewID(), Contacts: c, HiddenContacts: []enums.ContactKey{enums.ContactKeyTelegram, enums.ContactKeyTelegram}})
	require.Error(t, err)
}

func TestContactsKeepsOwnedStateSeparateFromPresentation(t *testing.T) {
	email, err := valueobjects.NewEmail("person@example.com")
	require.NoError(t, err)
	c, err := entities.NewContacts(entities.ContactsParams{Email: &email})
	require.NoError(t, err)
	aggregate, err := domain.NewContacts(domain.ContactsParams{ID: valueobjects.NewID(), ResumeID: valueobjects.NewID(), Contacts: c})
	require.NoError(t, err)
	require.Equal(t, "person@example.com", aggregate.Contacts().Email().String())
}

func TestContactsTextRejectsWhitespaceOnlyValues(t *testing.T) {
	value := "   "
	_, err := valueobjects.NewContactsText(valueobjects.ContactsTextParams{Title: &value})
	require.ErrorIs(t, err, valueobjects.ErrInvalidValue)
}

func TestContactsClonePreservesContactsText(t *testing.T) {
	title := "  Contacts  "
	text, err := valueobjects.NewContactsText(valueobjects.ContactsTextParams{Title: &title})
	require.NoError(t, err)

	aggregate, err := domain.NewContacts(domain.ContactsParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Contacts: func() entities.Contacts {
			contacts, createErr := entities.NewContacts(entities.ContactsParams{Phone: "+33123456789"})
			require.NoError(t, createErr)
			return contacts
		}(),
		Text: &text,
	})
	require.NoError(t, err)
	require.NotNil(t, aggregate.Text())
	require.Equal(t, "Contacts", *aggregate.Text().Title())
}

func TestContactsAllowsHiddenContactThatIsNotPresent(t *testing.T) {
	email, err := valueobjects.NewEmail("person@example.com")
	require.NoError(t, err)
	contacts, err := entities.NewContacts(entities.ContactsParams{Email: &email})
	require.NoError(t, err)

	aggregate, err := domain.NewContacts(domain.ContactsParams{
		ID:             valueobjects.NewID(),
		ResumeID:       valueobjects.NewID(),
		Contacts:       contacts,
		HiddenContacts: []enums.ContactKey{enums.ContactKeyPhone},
	})

	require.NoError(t, err)
	require.Equal(t, []enums.ContactKey{enums.ContactKeyPhone}, aggregate.HiddenContacts())
}

func TestContactsRejectsEmptyContactSet(t *testing.T) {
	_, err := domain.NewContacts(domain.ContactsParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
	})

	require.ErrorIs(t, err, domain.ErrInvalidResume)
}
