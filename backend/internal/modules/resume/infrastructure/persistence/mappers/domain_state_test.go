package mappers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	contactkey "developers-cv/backend/internal/shared/domain/contact"
)

func TestContactsToDomainRestoresPersistedStateWithoutBusinessValidation(t *testing.T) {
	contactsID := uuid.New()
	contacts, err := ContactsToDomain(
		gormmodel.ResumeContacts{ID: contactsID, ResumeID: uuid.New(), Phone: stringPtr("+123456")},
		[]gormmodel.ResumeContactHidden{{ContactsID: contactsID, Value: "invalid", Position: 0}},
		nil,
	)

	require.NoError(t, err)
	require.Equal(t, []contactkey.Key{"invalid"}, contacts.HiddenContacts())
}

func TestOverviewToDomainCopiesPersistedFieldsWithoutValidation(t *testing.T) {
	date := time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)
	note := "Persisted note"
	status := "open"
	model := gormmodel.ResumeOverview{AvailabilityStatus: &status, AvailabilityFrom: &date, AvailabilityNote: &note}

	got := OverviewToDomain(model, nil)

	require.NotNil(t, got.Availability())
	require.Equal(t, domain.AvailabilityStatusOpen, got.Availability().Status())
	require.Equal(t, date, got.Availability().From().Time())
	require.Equal(t, note, *got.Availability().Note())
}

func TestSkillsToDomainCopiesModelWithoutOwnershipValidation(t *testing.T) {
	resumeID := uuid.New()
	skillsID := uuid.New()
	skills, err := SkillsToDomain(
		gormmodel.ResumeSkills{ID: skillsID, ResumeID: resumeID},
		[]gormmodel.ResumeSkillHard{{SkillsID: skillsID, Value: " Go ", Position: 0}},
		nil,
	)

	require.NoError(t, err)
	require.Equal(t, " Go ", skills.HardSkills()[0])
}
