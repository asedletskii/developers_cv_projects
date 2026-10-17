package requests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectInputAcceptsExternalScreenshotURLs(t *testing.T) {
	input, err := (ProjectRequest{
		Name:        "Project",
		Screenshots: []string{"https://example.com/project.png"},
	}).ToInput()

	require.NoError(t, err)
	require.Len(t, input.Screenshots, 1)
	require.Equal(t, "https://example.com/project.png", input.Screenshots[0].String())
}

func TestProjectInputRejectsRelativeScreenshotPaths(t *testing.T) {
	_, err := (ProjectRequest{
		Name:        "Project",
		Screenshots: []string{"/assets/project.png"},
	}).ToInput()

	require.Error(t, err)
}
