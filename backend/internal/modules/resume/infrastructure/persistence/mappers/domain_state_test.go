package mappers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func TestContactsToDomainRejectsInvalidHiddenContact(t *testing.T) {
	contactsID := uuid.New()
	_, err := ContactsToDomain(
		gormmodel.ResumeContacts{ID: contactsID, ResumeID: uuid.New(), Phone: stringPtr("+123456")},
		[]gormmodel.ResumeContactHidden{{ContactsID: contactsID, Value: "invalid", Position: 0}},
		nil,
	)

	require.ErrorIs(t, err, domain.ErrInvalidResume)
}
