package profile

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProfileRequestAcceptsExternalPhotoURL(t *testing.T) {
	input, err := (profileRequest{
		Name:        "Антон",
		Surname:     "Седлецкий",
		Position:    "Backend Developer",
		PhotoPath:   "https://example.com/avatar.jpg",
		DateOfBirth: "1998-04-12",
		Location:    "Paris",
	}).input()

	require.NoError(t, err)
	require.Equal(t, "https://example.com/avatar.jpg", input.PhotoPath.String())
}

func TestProfileRequestRejectsRelativePhotoPath(t *testing.T) {
	_, err := (profileRequest{
		Name:        "Антон",
		Surname:     "Седлецкий",
		Position:    "Backend Developer",
		PhotoPath:   "/assets/avatar.jpg",
		DateOfBirth: "1998-04-12",
		Location:    "Paris",
	}).input()

	require.Error(t, err)
}
