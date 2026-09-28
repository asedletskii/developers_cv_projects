package domain

import (
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestResumeForAggregate(t *testing.T) *Resume {
	t.Helper()
	r, err := NewResume(ResumeParams{ID: valueobjects.NewID(), UserID: useridentity.NewID(), Title: "Resume", Language: enums.ResumeLanguageEN, Version: 1})
	require.NoError(t, err)
	return r
}

func newReplacementProfileParams(t *testing.T, id, resumeID valueobjects.ID, name string) ProfileParams {
	t.Helper()
	photo, err := valueobjects.NewExternalURL("https://example.com/profile.jpg")
	require.NoError(t, err)
	dob, err := valueobjects.ParseDate("1998-04-12")
	require.NoError(t, err)
	return ProfileParams{ID: id, ResumeID: resumeID, Name: name, Surname: "Doe", Position: "Developer", PhotoPath: photo, DateOfBirth: dob, Location: "Paris", WorkFormats: []enums.WorkFormat{enums.WorkFormatRemote}}
}

func TestResumeReplaceProfileCreatesAndGeneratesID(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	err := resume.ReplaceProfile(newReplacementProfileParams(t, valueobjects.ID{}, resume.ID(), "John"))
	require.NoError(t, err)
	require.NotNil(t, resume.Profile())
	require.False(t, resume.Profile().ID().IsZero())
	require.Equal(t, resume.ID(), resume.Profile().ResumeID())
}

func TestResumeReplaceProfilePreservesExistingID(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceProfile(newReplacementProfileParams(t, valueobjects.ID{}, resume.ID(), "John")))
	firstID := resume.Profile().ID()
	require.NoError(t, resume.ReplaceProfile(newReplacementProfileParams(t, valueobjects.ID{}, resume.ID(), "Jane")))
	require.Equal(t, firstID, resume.Profile().ID())
	require.Equal(t, "Jane", resume.Profile().Name())
}

func TestResumeReplaceProfileRejectsUnknownID(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceProfile(newReplacementProfileParams(t, valueobjects.ID{}, resume.ID(), "John")))
	err := resume.ReplaceProfile(newReplacementProfileParams(t, valueobjects.NewID(), resume.ID(), "Jane"))
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceProfileRejectsForeignExistingID(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceProfile(newReplacementProfileParams(t, valueobjects.ID{}, resume.ID(), "John")))
	err := resume.ReplaceProfile(newReplacementProfileParams(t, valueobjects.NewID(), resume.ID(), "Jane"))
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceProfileRejectsInvalidProfile(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	p := newReplacementProfileParams(t, valueobjects.ID{}, resume.ID(), "John")
	p.WorkFormats = []enums.WorkFormat{enums.WorkFormatRemote, enums.WorkFormatRemote}
	err := resume.ReplaceProfile(p)
	require.ErrorIs(t, err, ErrDuplicateEntity)
}
