package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func validResumeLinkParams() ResumeLinkParams {
	return ResumeLinkParams{
		ResumeID: sharedvalueobjects.NewID(),
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

func TestNewResumeLinkRejectsInvalidResumeID(t *testing.T) {
	params := validResumeLinkParams()
	params.ResumeID = sharedvalueobjects.ID{}

	_, err := NewResumeLink(params)
	require.ErrorIs(t, err, ErrInvalidResumeLink)
}

func TestNewResumeLinkGeneratesIDAndStartsAtVersionOne(t *testing.T) {
	link, err := NewResumeLink(validResumeLinkParams())
	require.NoError(t, err)
	require.False(t, link.ID().IsZero())
	require.Equal(t, int64(1), link.Version())
}

func TestEnsureResumeLinkCapacity(t *testing.T) {
	require.NoError(t, EnsureResumeLinkCapacity(MaxResumeLinksPerResume-1))
	require.ErrorIs(t, EnsureResumeLinkCapacity(MaxResumeLinksPerResume), ErrResumeLinkLimitExceeded)
	require.ErrorIs(t, EnsureResumeLinkCapacity(MaxResumeLinksPerResume+1), ErrResumeLinkLimitExceeded)
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
		_, err := RestoreResumeLink(ResumeLinkSnapshot{
			ID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(),
			Name: "HH", Status: ResumeLinkStatusActive, Version: version,
		})
		require.ErrorIs(t, err, ErrInvalidResumeLink)
	}
}

func TestNewResumeLinkAcceptsMaximumNameLength(t *testing.T) {
	params := validResumeLinkParams()
	params.Name = string(make([]rune, MaxResumeLinkNameLength))
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

func TestNewResumeLinkRejectsInvalidStatus(t *testing.T) {
	_, err := RestoreResumeLink(ResumeLinkSnapshot{
		ID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(),
		Name: "HH", Status: ResumeLinkStatus("broken"), Version: 1,
	})
	require.ErrorIs(t, err, ErrInvalidResumeLink)
}

func TestNewResumeLinkDefaultsStatusToActive(t *testing.T) {
	link, err := NewResumeLink(validResumeLinkParams())

	require.NoError(t, err)
	require.Equal(t, ResumeLinkStatusActive, link.Status())
	require.True(t, link.IsActive())
}

func TestReconstructResumeLinkRestoresSnapshot(t *testing.T) {
	snapshot := ResumeLinkSnapshot{
		ID:       sharedvalueobjects.NewID(),
		ResumeID: sharedvalueobjects.NewID(),
		Name:     "  Telegram  ",
		Status:   ResumeLinkStatusInactive,
		Version:  7,
	}

	link, err := ReconstructResumeLink(snapshot)
	require.NoError(t, err)
	require.Equal(t, snapshot.ID, link.ID())
	require.Equal(t, snapshot.ResumeID, link.ResumeID())
	require.Equal(t, "  Telegram  ", link.Name())
	require.Equal(t, snapshot.Status, link.Status())
	require.Equal(t, snapshot.Version, link.Version())
	require.False(t, link.IsActive())
}

func TestReconstructResumeLinkRejectsInvalidSnapshot(t *testing.T) {
	_, err := ReconstructResumeLink(ResumeLinkSnapshot{
		ID:       sharedvalueobjects.NewID(),
		ResumeID: sharedvalueobjects.NewID(),
		Name:     "HH",
		Status:   ResumeLinkStatus("broken"),
		Version:  1,
	})
	require.ErrorIs(t, err, ErrInvalidResumeLink)
}

func TestResumeLinkRenameTrimsAndUpdatesName(t *testing.T) {
	link, err := NewResumeLink(validResumeLinkParams())
	require.NoError(t, err)

	require.NoError(t, link.Rename("  New Name  "))
	require.Equal(t, "New Name", link.Name())
	require.Equal(t, int64(1), link.Version())
}

func TestResumeLinkRenameRejectsEmptyName(t *testing.T) {
	link, err := NewResumeLink(validResumeLinkParams())
	require.NoError(t, err)

	for _, name := range []string{"", "   ", "\t\n"} {
		require.ErrorIs(t, link.Rename(name), ErrInvalidResumeLinkName)
	}
	require.Equal(t, "HH", link.Name())
}

func TestResumeLinkRenameRejectsTooLongName(t *testing.T) {
	link, err := NewResumeLink(validResumeLinkParams())
	require.NoError(t, err)

	err = link.Rename(strings.Repeat("ж", MaxResumeLinkNameLength+1))
	require.ErrorIs(t, err, ErrInvalidResumeLinkName)
	require.Equal(t, "HH", link.Name())
}

func TestResumeLinkRenameRejectsNilLink(t *testing.T) {
	var link *ResumeLink
	require.ErrorIs(t, link.Rename("Name"), ErrInvalidResumeLink)
}

func TestResumeLinkStateTransitions(t *testing.T) {
	link, err := NewResumeLink(validResumeLinkParams())
	require.NoError(t, err)

	require.NoError(t, link.Deactivate())
	require.False(t, link.IsActive())
	require.NoError(t, link.Activate())
	require.True(t, link.IsActive())
}
