package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

func validResumeRootParams() ResumeRootParams {
	return ResumeRootParams{
		ID: valueobjects.NewID(), UserID: useridentity.NewID(), Title: "Frontend Developer",
		Language:     enums.ResumeLanguageEN,
		Presentation: ResumePresentation{Greeting: valueobjects.NewGreeting(valueobjects.GreetingParams{First: "Hello", Second: "I build", SecondHighlight: "web products", Third: "for users"})},
	}
}

func TestNewResumeRootRequiresValidLanguage(t *testing.T) {
	params := validResumeRootParams()
	resume, err := NewResumeRoot(params)
	require.NoError(t, err)
	require.Equal(t, enums.ResumeLanguageEN, resume.Language())

	params.Language = ""
	_, err = NewResumeRoot(params)
	require.Error(t, err)

	params.Language = enums.ResumeLanguage("xx")
	_, err = NewResumeRoot(params)
	require.Error(t, err)
}

func TestResumeRootReplaceKeepsLanguageImmutable(t *testing.T) {
	params := validResumeRootParams()
	resume, err := NewResumeRoot(params)
	require.NoError(t, err)

	params.Language = enums.ResumeLanguageRU
	params.Title = "Backend Developer"
	require.NoError(t, resume.Replace(params))
	require.Equal(t, enums.ResumeLanguageEN, resume.Language())
	require.Equal(t, "Backend Developer", resume.Title())
}
