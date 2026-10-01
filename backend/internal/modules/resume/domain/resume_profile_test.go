package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func newReplacementProfileParams(t *testing.T, id, resumeID sharedvalueobjects.ID, name string) ProfileParams {
	t.Helper()
	photo, err := NewExternalURL("https://example.com/profile.jpg")
	require.NoError(t, err)
	dob, err := ParseDate("1998-04-12")
	require.NoError(t, err)
	return ProfileParams{ID: id, ResumeID: resumeID, Name: name, Surname: "Doe", Position: "Developer", PhotoPath: photo, DateOfBirth: dob, Location: "Paris", WorkFormats: []WorkFormat{WorkFormatRemote}}
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

func TestProfileReplacePreservesIdentityAndNormalizesFields(t *testing.T) {
	photo, err := NewExternalURL("https://example.com/photo")
	require.NoError(t, err)
	birth, err := ParseDate("1990-01-01")
	require.NoError(t, err)
	id, resumeID := sharedvalueobjects.NewID(), sharedvalueobjects.NewID()

	profile, err := NewProfile(ProfileParams{ID: id, ResumeID: resumeID, Name: " John ", Surname: " Doe ", Position: " Developer ", PhotoPath: photo, DateOfBirth: birth, Location: " Paris ", WorkFormats: []WorkFormat{WorkFormatRemote}})
	require.NoError(t, err)
	require.Equal(t, "John", profile.Name())
	require.Equal(t, []WorkFormat{WorkFormatRemote}, profile.WorkFormats())

	invalid := *profile
	_ = invalid
	_, err = NewProfile(ProfileParams{ID: id, ResumeID: resumeID, Name: "John", Surname: "Doe", Position: "Developer", PhotoPath: photo, DateOfBirth: birth, Location: "Paris", WorkFormats: []WorkFormat{WorkFormatRemote, WorkFormatRemote}})
	require.ErrorIs(t, err, ErrDuplicateEntity)

	future, err := ParseDate("2999-01-01")
	require.NoError(t, err)
	_, err = NewProfile(ProfileParams{ID: id, ResumeID: resumeID, Name: "John", Surname: "Doe", Position: "Developer", PhotoPath: photo, DateOfBirth: future, Location: "Paris"})
	require.Error(t, err)

	err = profile.Replace(ProfileParams{Name: " Jane ", Surname: " Doe ", Position: " Lead ", PhotoPath: photo, DateOfBirth: birth, Location: " Lyon "})
	require.NoError(t, err)
	require.Equal(t, id, profile.ID())
	require.Equal(t, resumeID, profile.ResumeID())
	require.Equal(t, "Jane", profile.Name())
}
