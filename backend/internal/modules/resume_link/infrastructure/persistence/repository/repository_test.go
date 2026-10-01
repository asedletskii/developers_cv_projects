package repository

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestStatusOf(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(), Name: "GitHub"})
	require.NoError(t, err)
	require.Equal(t, "active", statusOf(link))
	require.NoError(t, link.Deactivate())
	require.Equal(t, "inactive", statusOf(link))
}
