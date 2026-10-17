package requests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProfileRequestAcceptsExternalPhotoURL(t *testing.T) {
	input, err := (ProfileRequest{
		Name:        "Антон",
		Surname:     "Седлецкий",
		Position:    "Backend Developer",
		PhotoURL:    "https://example.com/avatar.jpg",
		DateOfBirth: "1998-04-12",
		Location:    "Paris",
	}).ToInput()

	require.NoError(t, err)
	require.Equal(t, "https://example.com/avatar.jpg", input.PhotoURL.String())
}

func TestProfileRequestRejectsRelativePhotoPath(t *testing.T) {
	_, err := (ProfileRequest{
		Name:        "Антон",
		Surname:     "Седлецкий",
		Position:    "Backend Developer",
		PhotoURL:    "/assets/avatar.jpg",
		DateOfBirth: "1998-04-12",
		Location:    "Paris",
	}).ToInput()

	require.Error(t, err)
}
