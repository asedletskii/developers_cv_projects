package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func newTestResume(t *testing.T) *Resume {
	t.Helper()
	resume, err := NewResume(ResumeParams{
		ID:       sharedvalueobjects.NewID(),
		UserID:   sharedvalueobjects.NewID(),
		Title:    "Resume",
		Language: ResumeLanguageEN,
		Version:  1,
	})
	require.NoError(t, err)
	return resume
}

func validResumeParams() ResumeParams {
	return ResumeParams{
		ID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID(), Title: "Frontend Developer",
		Language: ResumeLanguageEN,
	}
}

func TestNewResumeRequiresValidLanguage(t *testing.T) {
	params := validResumeParams()
	resume, err := NewResume(params)
	require.NoError(t, err)
	require.Equal(t, ResumeLanguageEN, resume.Language())

	params.Language = ""
	_, err = NewResume(params)
	require.Error(t, err)

	params.Language = ResumeLanguage("xx")
	_, err = NewResume(params)
	require.Error(t, err)
}

func TestResumeReplaceTitleKeepsLanguageImmutable(t *testing.T) {
	params := validResumeParams()
	resume, err := NewResume(params)
	require.NoError(t, err)

	require.NoError(t, resume.ReplaceTitle("Backend Developer"))
	require.Equal(t, ResumeLanguageEN, resume.Language())
	require.Equal(t, "Backend Developer", resume.Title())
}

func TestNewResumeGeneratesIDAndStartsAtVersionOne(t *testing.T) {
	resume, err := NewResume(ResumeParams{UserID: sharedvalueobjects.NewID(), Title: "Backend Developer", Language: ResumeLanguageEN})
	require.NoError(t, err)
	require.False(t, resume.ID().IsZero())
	require.Equal(t, int64(1), resume.Version())
}

func TestNewResumeRejectsInvalidUserAndTitle(t *testing.T) {
	tests := []ResumeParams{
		{Title: "Backend Developer", Language: ResumeLanguageEN},
		{UserID: sharedvalueobjects.NewID(), Language: ResumeLanguageEN},
	}
	for _, params := range tests {
		_, err := NewResume(params)
		require.Error(t, err)
	}
}

func TestNewResumeRejectsInvalidVersion(t *testing.T) {
	_, err := NewResume(ResumeParams{UserID: sharedvalueobjects.NewID(), Title: "Backend Developer", Language: ResumeLanguageEN, Version: -1})
	require.ErrorIs(t, err, ErrInvalidResume)
}

func TestReconstructResumeAllowsAbsentOptionalChildren(t *testing.T) {
	resume, err := ReconstructResume(ResumeSnapshot{
		ID:       sharedvalueobjects.NewID(),
		UserID:   sharedvalueobjects.NewID(),
		Title:    "Backend",
		Language: ResumeLanguageEN,
		Status:   ResumeStatusDraft,
		Version:  1,
	})
	require.NoError(t, err)
	require.NotNil(t, resume)
	require.Nil(t, resume.Profile())
	require.Nil(t, resume.Contacts())
	require.Nil(t, resume.Skills())
	require.Nil(t, resume.Languages())
	require.Nil(t, resume.Technologies())
}

func TestReconstructResumeRejectsForeignOwnedChildren(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	foreignResumeID := sharedvalueobjects.NewID()
	photoPath, err := NewExternalURL("https://example.com/photo.jpg")
	require.NoError(t, err)
	dob, err := NewDate(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	profile, err := NewProfile(ProfileParams{
		ID: sharedvalueobjects.NewID(), ResumeID: foreignResumeID, Name: "John", Surname: "Doe", Position: "Developer",
		PhotoPath: photoPath, DateOfBirth: dob, Location: "Paris",
	})
	require.NoError(t, err)
	_, err = ReconstructResume(ResumeSnapshot{ID: resumeID, UserID: sharedvalueobjects.NewID(), Title: "Backend", Language: ResumeLanguageEN, Status: ResumeStatusDraft, Version: 1, Profile: profile})
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceTitleTrimsAndValidates(t *testing.T) {
	resume, err := NewResume(ResumeParams{UserID: sharedvalueobjects.NewID(), Title: "Old", Language: ResumeLanguage("en")})
	require.NoError(t, err)

	require.NoError(t, resume.ReplaceTitle("  New title  "))
	require.Equal(t, "New title", resume.Title())

	err = resume.ReplaceTitle("   ")
	require.ErrorIs(t, err, ErrInvalidResume)
}

func TestNewResumeStartsAsDraft(t *testing.T) {
	resume, err := NewResume(ResumeParams{
		UserID: sharedvalueobjects.NewID(), Title: "Backend", Language: ResumeLanguageEN,
	})
	require.NoError(t, err)
	require.Equal(t, ResumeStatusDraft, resume.Status())
}

func TestResumeStatusTransitions(t *testing.T) {
	resume, err := NewResume(ResumeParams{
		UserID: sharedvalueobjects.NewID(), Title: "Backend", Language: ResumeLanguageEN,
	})
	require.NoError(t, err)

	require.NoError(t, resume.Publish())
	require.Equal(t, ResumeStatusOpened, resume.Status())
	require.ErrorIs(t, resume.Publish(), ErrInvalidResumeStatusTransition)
	require.NoError(t, resume.Close())
	require.Equal(t, ResumeStatusClosed, resume.Status())
	require.NoError(t, resume.Publish())
	require.Equal(t, ResumeStatusOpened, resume.Status())
}

func TestResumeCannotCloseDraft(t *testing.T) {
	resume, err := NewResume(ResumeParams{
		UserID: sharedvalueobjects.NewID(), Title: "Backend", Language: ResumeLanguageEN,
	})
	require.NoError(t, err)
	require.ErrorIs(t, resume.Close(), ErrInvalidResumeStatusTransition)
	require.Equal(t, ResumeStatusDraft, resume.Status())
}

func TestNormalizeOptionalStringTrimsWhitespaceAndTurnsEmptyIntoNil(t *testing.T) {
	value := "  optional value  "
	got := cloneString(&value)
	require.NotNil(t, got)
	require.Equal(t, "optional value", *got)

	empty := "   "
	require.Nil(t, cloneString(&empty))
	require.Nil(t, cloneString(nil))
}
