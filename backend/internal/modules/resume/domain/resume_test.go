package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

func validResumeParams() ResumeParams {
	return ResumeParams{
		ID: valueobjects.NewID(), UserID: useridentity.NewID(), Title: "Frontend Developer",
		Language: enums.ResumeLanguageEN,
	}
}

func TestNewResumeRequiresValidLanguage(t *testing.T) {
	params := validResumeParams()
	resume, err := NewResume(params)
	require.NoError(t, err)
	require.Equal(t, enums.ResumeLanguageEN, resume.Language())

	params.Language = ""
	_, err = NewResume(params)
	require.Error(t, err)

	params.Language = enums.ResumeLanguage("xx")
	_, err = NewResume(params)
	require.Error(t, err)
}

func TestResumeRenameKeepsLanguageImmutable(t *testing.T) {
	params := validResumeParams()
	resume, err := NewResume(params)
	require.NoError(t, err)

	require.NoError(t, resume.Rename("Backend Developer"))
	require.Equal(t, enums.ResumeLanguageEN, resume.Language())
	require.Equal(t, "Backend Developer", resume.Title())
}

func TestNewResumeRejectsInvalidIdentityAndTitle(t *testing.T) {
	tests := []ResumeParams{
		{UserID: useridentity.NewID(), Title: "Backend Developer", Language: enums.ResumeLanguageEN},
		{ID: valueobjects.NewID(), Title: "Backend Developer", Language: enums.ResumeLanguageEN},
		{ID: valueobjects.NewID(), UserID: useridentity.NewID(), Language: enums.ResumeLanguageEN},
	}
	for _, params := range tests {
		_, err := NewResume(params)
		require.Error(t, err)
	}
}

func TestResumeRenameTrimsTitle(t *testing.T) {
	resume, err := NewResume(validResumeParams())
	require.NoError(t, err)
	require.NoError(t, resume.Rename("  Senior Backend Developer  "))
	require.Equal(t, "Senior Backend Developer", resume.Title())
}
