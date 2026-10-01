package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/shared/domain/valueobjects"
)

func validResumeLinkParams() ResumeLinkParams {
	return ResumeLinkParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Name:     "HH",
	}
}

func TestNewResumeLink(t *testing.T) {
	link, err := NewResumeLink(validResumeLinkParams())
	require.NoError(t, err)
	require.False(t, link.ID().IsZero())
	require.False(t, link.ResumeID().IsZero())
	require.Equal(t, "HH", link.Name())
	require.True(t, link.IsActive())
	require.Equal(t, int64(1), link.Version())
}

func TestNewResumeLinkRejectsInvalidIdentity(t *testing.T) {
	params := validResumeLinkParams()
	params.ID = valueobjects.ID{}
	_, err := NewResumeLink(params)
	require.ErrorIs(t, err, ErrInvalidResumeLink)

	params = validResumeLinkParams()
	params.ResumeID = valueobjects.ID{}
	_, err = NewResumeLink(params)
	require.ErrorIs(t, err, ErrInvalidResumeLink)
}

func TestNewResumeLinkRejectsEmptyName(t *testing.T) {
	for _, name := range []string{"", "   ", "\t\n"} {
		params := validResumeLinkParams()
		params.Name = name

		_, err := NewResumeLink(params)
		require.ErrorIs(t, err, ErrInvalidResumeLinkName)
	}
}

func TestNewResumeLinkTrimsName(t *testing.T) {
	params := validResumeLinkParams()
	params.Name = "  Telegram  "

	link, err := NewResumeLink(params)
	require.NoError(t, err)
	require.Equal(t, "Telegram", link.Name())
}

func TestNewResumeLinkRejectsNonPositiveVersion(t *testing.T) {
	for _, version := range []int64{-1, -2} {
		params := validResumeLinkParams()
		params.Version = version
		_, err := NewResumeLink(params)
		require.ErrorIs(t, err, ErrInvalidResumeLink)
	}
}

func TestResumeLinkStateTransitions(t *testing.T) {
	link, err := NewResumeLink(validResumeLinkParams())
	require.NoError(t, err)

	require.NoError(t, link.Deactivate())
	require.False(t, link.IsActive())
	require.NoError(t, link.Activate())
	require.True(t, link.IsActive())
}

func TestNewResumeLinkAcceptsMaximumNameLength(t *testing.T) {
	params := validResumeLinkParams()
	params.Name = string([]rune{'ж'})
	params.Name += string(make([]rune, MaxResumeLinkNameLength-1))
	runes := []rune(params.Name)
	for i := range runes {
		runes[i] = 'ж'
	}
	params.Name = string(runes)

	link, err := NewResumeLink(params)
	require.NoError(t, err)
	require.Len(t, []rune(link.Name()), MaxResumeLinkNameLength)
}

func TestNewResumeLinkRejectsNameLongerThanMaximum(t *testing.T) {
	params := validResumeLinkParams()
	params.Name = string(make([]rune, MaxResumeLinkNameLength+1))

	_, err := NewResumeLink(params)
	require.ErrorIs(t, err, ErrInvalidResumeLinkName)
}
