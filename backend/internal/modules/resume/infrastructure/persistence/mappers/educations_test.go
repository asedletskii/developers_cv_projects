package mappers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func TestEducationsToDomainPreservesRowOrder(t *testing.T) {
	resumeID := uuid.New()
	firstID := uuid.New()
	secondID := uuid.New()

	educations, err := EducationsToDomain([]gormmodel.ResumeEducation{
		{ID: firstID, ResumeID: resumeID, University: "First"},
		{ID: secondID, ResumeID: resumeID, University: "Second"},
	})

	require.NoError(t, err)
	require.Len(t, educations, 2)
	require.Equal(t, firstID, educations[0].ID().UUID())
	require.Equal(t, "First", educations[0].University())
	require.Equal(t, secondID, educations[1].ID().UUID())
	require.Equal(t, "Second", educations[1].University())
}

func TestEducationsToDomainReturnsEmptySliceForNoRows(t *testing.T) {
	educations, err := EducationsToDomain(nil)

	require.NoError(t, err)
	require.NotNil(t, educations)
	require.Empty(t, educations)
}
