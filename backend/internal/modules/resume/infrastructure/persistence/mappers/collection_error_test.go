package mappers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestLanguagesToDomainAddsRowContext(t *testing.T) {
	_, err := LanguagesToDomain([]gormmodel.ResumeSpokenLanguage{
		{ID: uuid.New()},
		{},
	})

	require.ErrorIs(t, err, sharedvalueobjects.ErrInvalidID)
	require.ErrorContains(t, err, "spoken language at index 1")
}

func TestTechnologiesToDomainAddsRowContext(t *testing.T) {
	_, err := TechnologiesToDomain([]gormmodel.ResumeTechnology{
		{ID: uuid.New()},
		{},
	})

	require.ErrorIs(t, err, sharedvalueobjects.ErrInvalidID)
	require.ErrorContains(t, err, "technology at index 1")
}
