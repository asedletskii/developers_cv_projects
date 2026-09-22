package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func newTestResume(t *testing.T) *Resume {
	t.Helper()
	id, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	date, err := valueobjects.ParseDate("1998-04-12")
	require.NoError(t, err)
	photoPath, err := valueobjects.NewAssetPath("/mock/avatar.svg")
	require.NoError(t, err)
	contacts, err := entities.NewContacts(entities.ContactsParams{Phone: "+79999999999"})
	require.NoError(t, err)
	section, err := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionAbout, Enabled: true})
	require.NoError(t, err)

	resume, err := NewResume(ResumeParams{
		ID:          id,
		Name:        "Anton",
		Surname:     "Sedletskii",
		Position:    "Frontend",
		PhotoPath:   photoPath,
		DateOfBirth: date,
		Location:    "Novosibirsk",
		Sections:    []valueobjects.SectionSettings{section},
		Contacts:    contacts,
		WorkFormats: []enums.WorkFormat{enums.WorkFormatRemote},
	})
	require.NoError(t, err)
	return resume
}

func TestResumeRejectsDuplicateSections(t *testing.T) {
	id, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	date, err := valueobjects.ParseDate("1998-04-12")
	require.NoError(t, err)
	photoPath, err := valueobjects.NewAssetPath("/mock/avatar.svg")
	require.NoError(t, err)
	section, err := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionAbout, Enabled: true})
	require.NoError(t, err)
	contacts, err := entities.NewContacts(entities.ContactsParams{Phone: "+79999999999"})
	require.NoError(t, err)

	_, err = NewResume(ResumeParams{
		ID: id, Name: "Anton", Surname: "Sedletskii", Position: "Frontend", PhotoPath: photoPath,
		DateOfBirth: date, Location: "Novosibirsk", Sections: []valueobjects.SectionSettings{section, section}, Contacts: contacts,
	})
	require.ErrorIs(t, err, ErrDuplicateEntity)
}

func TestResumeAddAndRemoveExperience(t *testing.T) {
	resume := newTestResume(t)
	id, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	start, err := valueobjects.ParseYearMonth("2025-01")
	require.NoError(t, err)
	experience, err := entities.NewExperience(entities.ExperienceParams{ID: id, Company: "Acme", Position: "Developer", StartDate: start})
	require.NoError(t, err)

	require.NoError(t, resume.AddExperience(experience))
	require.Len(t, resume.Experience(), 1)
	require.ErrorIs(t, resume.AddExperience(experience), ErrDuplicateEntity)
	require.NoError(t, resume.RemoveExperience(id))
	require.Empty(t, resume.Experience())
	require.ErrorIs(t, resume.RemoveExperience(id), ErrEntityNotFound)
}

func TestResumeRejectsFutureDateOfBirth(t *testing.T) {
	id, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	date, err := valueobjects.NewDate(time.Now().UTC().AddDate(1, 0, 0))
	require.NoError(t, err)
	photoPath, err := valueobjects.NewAssetPath("/mock/avatar.svg")
	require.NoError(t, err)
	contacts, err := entities.NewContacts(entities.ContactsParams{Phone: "+79999999999"})
	require.NoError(t, err)

	_, err = NewResume(ResumeParams{
		ID: id, Name: "Anton", Surname: "Sedletskii", Position: "Frontend", PhotoPath: photoPath,
		DateOfBirth: date, Location: "Novosibirsk", Contacts: contacts,
	})
	require.Error(t, err)
}
