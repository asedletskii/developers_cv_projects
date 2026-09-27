package entities

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestTechnologyTrimsAndDefensivelyCopiesGrade(t *testing.T) {
	grade := enums.TechnologyGradeSenior
	technology, err := NewTechnology(TechnologyParams{Name: " Go ", Grade: &grade})
	require.NoError(t, err)
	require.Equal(t, "Go", technology.Name())
	got := technology.Grade()
	require.NotNil(t, got)
	*got = enums.TechnologyGradeJunior
	require.Equal(t, enums.TechnologyGradeSenior, *technology.Grade())

	_, err = NewTechnology(TechnologyParams{Name: "   "})
	require.Error(t, err)
}

func TestContributorValidation(t *testing.T) {
	email, err := valueobjects.NewEmail("user@example.com")
	require.NoError(t, err)
	contacts, err := NewContacts(ContactsParams{Email: &email})
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

	contacts, err := NewContacts(ContactsParams{Phone: " +33123456789 "})
	require.NoError(t, err)
	phone := contacts.Phone()
	require.NotNil(t, phone)
	require.Equal(t, "+33123456789", *phone)
}

func TestContactsUseValidatedEmailAndTelegram(t *testing.T) {
	email, err := valueobjects.NewEmail("user@example.com")
	require.NoError(t, err)
	telegram, err := valueobjects.NewTelegram("@developer")
	require.NoError(t, err)
	contacts, err := NewContacts(ContactsParams{Email: &email, Telegram: &telegram})
	require.NoError(t, err)
	require.Equal(t, "user@example.com", contacts.Email().String())
	require.Equal(t, "developer", contacts.Telegram().String())
}
