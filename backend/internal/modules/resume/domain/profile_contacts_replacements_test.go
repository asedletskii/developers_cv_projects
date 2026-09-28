package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestProfileValidationAndReplacement(t *testing.T) {
	photo, err := valueobjects.NewExternalURL("https://example.com/photo")
	require.NoError(t, err)
	birth, err := valueobjects.ParseDate("1990-01-01")
	require.NoError(t, err)
	id, resumeID := valueobjects.NewID(), valueobjects.NewID()

	profile, err := NewProfile(ProfileParams{ID: id, ResumeID: resumeID, Name: " John ", Surname: " Doe ", Position: " Developer ", PhotoPath: photo, DateOfBirth: birth, Location: " Paris ", WorkFormats: []enums.WorkFormat{enums.WorkFormatRemote}})
	require.NoError(t, err)
	require.Equal(t, "John", profile.Name())
	require.Equal(t, []enums.WorkFormat{enums.WorkFormatRemote}, profile.WorkFormats())

	invalid := *profile
	_ = invalid
	_, err = NewProfile(ProfileParams{ID: id, ResumeID: resumeID, Name: "John", Surname: "Doe", Position: "Developer", PhotoPath: photo, DateOfBirth: birth, Location: "Paris", WorkFormats: []enums.WorkFormat{enums.WorkFormatRemote, enums.WorkFormatRemote}})
	require.ErrorIs(t, err, ErrDuplicateEntity)

	future, err := valueobjects.ParseDate("2999-01-01")
	require.NoError(t, err)
	_, err = NewProfile(ProfileParams{ID: id, ResumeID: resumeID, Name: "John", Surname: "Doe", Position: "Developer", PhotoPath: photo, DateOfBirth: future, Location: "Paris"})
	require.Error(t, err)

	err = profile.Replace(ProfileParams{Name: " Jane ", Surname: " Doe ", Position: " Lead ", PhotoPath: photo, DateOfBirth: birth, Location: " Lyon "})
	require.NoError(t, err)
	require.Equal(t, id, profile.ID())
	require.Equal(t, resumeID, profile.ResumeID())
	require.Equal(t, "Jane", profile.Name())
}

func TestContactsValidationAndReplacement(t *testing.T) {
	id, resumeID := valueobjects.NewID(), valueobjects.NewID()
	contactsEntity, err := entities.NewContacts(entities.ContactsParams{Phone: "+33123456789"})
	require.NoError(t, err)
	contacts, err := NewContacts(ContactsParams{ID: id, ResumeID: resumeID, Contacts: contactsEntity, HiddenContacts: []enums.ContactKey{enums.ContactKeyPhone}})
	require.NoError(t, err)

	email, err := valueobjects.NewEmail("person@example.com")
	require.NoError(t, err)
	replacementEntity, err := entities.NewContacts(entities.ContactsParams{
		Phone: "+33123456789",
		Email: &email,
	})
	require.NoError(t, err)

	err = contacts.Replace(ContactsParams{Contacts: replacementEntity, HiddenContacts: []enums.ContactKey{enums.ContactKeyEmail}})
	require.NoError(t, err)
	require.Equal(t, id, contacts.ID())
	require.Equal(t, resumeID, contacts.ResumeID())
	require.Equal(t, []enums.ContactKey{enums.ContactKeyEmail}, contacts.HiddenContacts())

	_, err = NewContacts(ContactsParams{ID: id, ResumeID: resumeID, Contacts: contactsEntity, HiddenContacts: []enums.ContactKey{enums.ContactKeyEmail, enums.ContactKeyEmail}})
	require.Error(t, err)
}
