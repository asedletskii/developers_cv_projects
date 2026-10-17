package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func newReplacementProfileParams(t *testing.T, id, _ sharedvalueobjects.ID, name string) ProfileParams {
	t.Helper()
	photo, err := NewExternalURL("https://example.com/profile.jpg")
	require.NoError(t, err)
	dob, err := ParseDate("1998-04-12")
	require.NoError(t, err)
	return ProfileParams{ID: id, Name: name, Surname: "Doe", Position: "Developer", PhotoURL: photo, DateOfBirth: dob, Location: "Paris", WorkFormats: []WorkFormat{WorkFormatRemote}}
}

func TestResumeReplaceProfileCreatesAndGeneratesID(t *testing.T) {
	resume := newTestResume(t)
	err := resume.ReplaceProfile(newReplacementProfileParams(t, sharedvalueobjects.ID{}, resume.ID(), "John"))
	require.NoError(t, err)
	require.NotNil(t, resume.Profile())
	require.False(t, resume.Profile().ID().IsZero())
	require.Equal(t, resume.ID(), resume.Profile().ResumeID())
}

func TestResumeReplaceProfilePreservesExistingID(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceProfile(newReplacementProfileParams(t, sharedvalueobjects.ID{}, resume.ID(), "John")))
	firstID := resume.Profile().ID()
	require.NoError(t, resume.ReplaceProfile(newReplacementProfileParams(t, sharedvalueobjects.ID{}, resume.ID(), "Jane")))
	require.Equal(t, firstID, resume.Profile().ID())
	require.Equal(t, "Jane", resume.Profile().Name())
}

func TestResumeReplaceProfileRejectsUnknownID(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceProfile(newReplacementProfileParams(t, sharedvalueobjects.ID{}, resume.ID(), "John")))
	err := resume.ReplaceProfile(newReplacementProfileParams(t, sharedvalueobjects.NewID(), resume.ID(), "Jane"))
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceProfileRejectsInvalidProfile(t *testing.T) {
	resume := newTestResume(t)
	p := newReplacementProfileParams(t, sharedvalueobjects.ID{}, resume.ID(), "John")
	p.WorkFormats = []WorkFormat{WorkFormatRemote, WorkFormatRemote}
	err := resume.ReplaceProfile(p)
	require.ErrorIs(t, err, ErrDuplicateEntity)
}

func TestNewProfileNormalizesFieldsAndValidates(t *testing.T) {
	photo, err := NewExternalURL("https://example.com/photo")
	require.NoError(t, err)
	birth, err := ParseDate("1990-01-01")
	require.NoError(t, err)
	id, resumeID := sharedvalueobjects.NewID(), sharedvalueobjects.NewID()

	profile, err := newProfile(ProfileParams{ID: id, Name: " John ", Surname: " Doe ", Position: " Developer ", PhotoURL: photo, DateOfBirth: birth, Location: " Paris ", WorkFormats: []WorkFormat{WorkFormatRemote}}, resumeID)
	require.NoError(t, err)
	require.Equal(t, "John", profile.Name())
	require.Equal(t, []WorkFormat{WorkFormatRemote}, profile.WorkFormats())

	_, err = newProfile(ProfileParams{ID: id, Name: "John", Surname: "Doe", Position: "Developer", PhotoURL: photo, DateOfBirth: birth, Location: "Paris", WorkFormats: []WorkFormat{WorkFormatRemote, WorkFormatRemote}}, resumeID)
	require.ErrorIs(t, err, ErrDuplicateEntity)

	future, err := ParseDate("2999-01-01")
	require.NoError(t, err)
	_, err = newProfile(ProfileParams{ID: id, Name: "John", Surname: "Doe", Position: "Developer", PhotoURL: photo, DateOfBirth: future, Location: "Paris"}, resumeID)
	require.Error(t, err)

}

func TestNewProfileRejectsNameAboveConfiguredLimit(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	params := newReplacementProfileParams(t, sharedvalueobjects.NewID(), resumeID, strings.Repeat("a", MaxProfileNameLength+1))
	_, err := newProfile(params, resumeID)
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestRestoreProfilePreservesSnapshotText(t *testing.T) {
	photo, err := NewExternalURL("https://example.com/profile.jpg")
	require.NoError(t, err)
	dob, err := ParseDate("1998-04-12")
	require.NoError(t, err)
	profile := RestoreProfile(ProfileSnapshot{
		ID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(),
		Name: " John ", Surname: "Doe", Position: "Developer", PhotoURL: photo,
		DateOfBirth: dob, Location: "Paris",
	})

	err = profile.validate()

	require.NoError(t, err)
	require.Equal(t, " John ", profile.Name())
}

func TestResumeReplaceProfileDoesNotMutateOnError(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceProfile(newReplacementProfileParams(t, sharedvalueobjects.ID{}, resume.ID(), "John")))
	existing := resume.Profile()

	invalid := newReplacementProfileParams(t, sharedvalueobjects.ID{}, resume.ID(), "Jane")
	invalid.WorkFormats = []WorkFormat{WorkFormatRemote, WorkFormatRemote}
	err := resume.ReplaceProfile(invalid)

	require.ErrorIs(t, err, ErrDuplicateEntity)
	require.Equal(t, existing.ID(), resume.Profile().ID())
	require.Equal(t, existing.Name(), resume.Profile().Name())
}
