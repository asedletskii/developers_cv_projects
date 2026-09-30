package repository

import (
	"testing"

	"github.com/stretchr/testify/require"

	resumevalueobjects "developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/domain/valueobjects"
)

func TestStatusOf(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: valueobjects.NewID(), ResumeID: resumevalueobjects.NewID(), Name: "GitHub"})
	require.NoError(t, err)
	require.Equal(t, "active", statusOf(link))
	require.NoError(t, link.Deactivate())
	require.Equal(t, "inactive", statusOf(link))
}
