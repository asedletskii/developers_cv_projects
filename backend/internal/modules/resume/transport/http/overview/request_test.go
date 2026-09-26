package overview

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestOverviewRequestAcceptsExternalResumeURL(t *testing.T) {
	command, err := (overviewRequest{
		Presentation: presentationRequest{
			Greeting:   greetingRequest{First: "A", Second: "B", SecondHighlight: "C", Third: "D"},
			ResumeLink: stringPtr("https://example.com/cv.pdf"),
			Sections:   []sectionRequest{},
		},
	}).ToCommand(useridentity.NewID(), valueobjects.NewID())

	require.NoError(t, err)
	require.NotNil(t, command.Presentation.ResumeLink)
	require.Equal(t, "https://example.com/cv.pdf", command.Presentation.ResumeLink.String())
}

func TestOverviewRequestRejectsRelativeResumePath(t *testing.T) {
	_, err := (overviewRequest{
		Presentation: presentationRequest{
			Greeting:   greetingRequest{First: "A", Second: "B", SecondHighlight: "C", Third: "D"},
			ResumeLink: stringPtr("/cv.pdf"),
			Sections:   []sectionRequest{},
		},
	}).ToCommand(useridentity.NewID(), valueobjects.NewID())

	require.Error(t, err)
}

func stringPtr(value string) *string {
	return &value
}
