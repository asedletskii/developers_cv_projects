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
