package mappers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func TestExperienceToDomainRejectsInvalidTechnologyGrade(t *testing.T) {
	resumeID := uuid.New()
	_, err := ExperienceToDomain(gormmodel.ResumeExperience{
		ID:        uuid.New(),
		ResumeID:  resumeID,
		Company:   "Company",
		Position:  "Developer",
		StartDate: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
	}, nil, nil, []gormmodel.ResumeExperienceStack{{
		ExperienceID: uuid.New(),
		Position:     0,
		Name:         "Go",
		Grade:        stringPtr("invalid"),
	}})

	require.ErrorIs(t, err, domain.ErrInvalidEntity)
}

func TestProjectToDomainRejectsInvalidTechnologyGrade(t *testing.T) {
	_, err := ProjectToDomain(gormmodel.ResumeProject{
		ID:       uuid.New(),
		ResumeID: uuid.New(),
		Name:     "Project",
	}, nil, []gormmodel.ResumeProjectTechnology{{
		ProjectID: uuid.New(),
		Position:  0,
		Name:      "Go",
		Grade:     stringPtr("invalid"),
	}}, nil)

	require.ErrorIs(t, err, domain.ErrInvalidEntity)
}

func TestContactsToDomainRejectsInvalidHiddenContact(t *testing.T) {
	contactsID := uuid.New()
	_, err := ContactsToDomain(
		gormmodel.ResumeContacts{ID: contactsID, ResumeID: uuid.New(), Phone: stringPtr("+123456")},
		[]gormmodel.ResumeContactHidden{{ContactsID: contactsID, Value: "invalid", Position: 0}},
		nil,
	)

	require.ErrorIs(t, err, domain.ErrInvalidResume)
}
