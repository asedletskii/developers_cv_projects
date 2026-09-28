package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func validResumeParams() ResumeParams {
	return ResumeParams{
		ID: valueobjects.NewID(), UserID: useridentity.NewID(), Title: "Frontend Developer",
		Language: enums.ResumeLanguageEN,
	}
}

func TestNewResumeRequiresValidLanguage(t *testing.T) {
	params := validResumeParams()
	resume, err := NewResume(params)
	require.NoError(t, err)
	require.Equal(t, enums.ResumeLanguageEN, resume.Language())

	params.Language = ""
	_, err = NewResume(params)
	require.Error(t, err)

	params.Language = enums.ResumeLanguage("xx")
	_, err = NewResume(params)
	require.Error(t, err)
}

func TestResumeReplaceTitleKeepsLanguageImmutable(t *testing.T) {
	params := validResumeParams()
	resume, err := NewResume(params)
	require.NoError(t, err)

	require.NoError(t, resume.ReplaceTitle("Backend Developer"))
	require.Equal(t, enums.ResumeLanguageEN, resume.Language())
	require.Equal(t, "Backend Developer", resume.Title())
}

func TestNewResumeGeneratesIDAndStartsAtVersionOne(t *testing.T) {
	resume, err := NewResume(ResumeParams{UserID: useridentity.NewID(), Title: "Backend Developer", Language: enums.ResumeLanguageEN})
	require.NoError(t, err)
	require.False(t, resume.ID().IsZero())
	require.Equal(t, int64(1), resume.Version())
}

func TestNewResumeRejectsInvalidUserAndTitle(t *testing.T) {
	tests := []ResumeParams{
		{Title: "Backend Developer", Language: enums.ResumeLanguageEN},
		{UserID: useridentity.NewID(), Language: enums.ResumeLanguageEN},
	}
	for _, params := range tests {
		_, err := NewResume(params)
		require.Error(t, err)
	}
}

func TestNewResumeRejectsInvalidVersion(t *testing.T) {
	_, err := NewResume(ResumeParams{UserID: useridentity.NewID(), Title: "Backend Developer", Language: enums.ResumeLanguageEN, Version: -1})
	require.ErrorIs(t, err, ErrInvalidResume)
}

func TestReconstructResumeAllowsAbsentOptionalChildren(t *testing.T) {
	resume, err := ReconstructResume(ResumeSnapshot{
		ID:       valueobjects.NewID(),
		UserID:   useridentity.NewID(),
		Title:    "Backend",
		Language: enums.ResumeLanguageEN,
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
	resumeID := valueobjects.NewID()
	foreignResumeID := valueobjects.NewID()
	photoPath, err := valueobjects.NewExternalURL("https://example.com/photo.jpg")
	require.NoError(t, err)
	dob, err := valueobjects.NewDate(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	profile, err := NewProfile(ProfileParams{
		ID: valueobjects.NewID(), ResumeID: foreignResumeID, Name: "John", Surname: "Doe", Position: "Developer",
		PhotoPath: photoPath, DateOfBirth: dob, Location: "Paris",
	})
	require.NoError(t, err)
	_, err = ReconstructResume(ResumeSnapshot{ID: resumeID, UserID: useridentity.NewID(), Title: "Backend", Language: enums.ResumeLanguageEN, Version: 1, Profile: profile})
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceTitleTrimsTitle(t *testing.T) {
	resume, err := NewResume(validResumeParams())
	require.NoError(t, err)
	require.NoError(t, resume.ReplaceTitle("  Senior Backend Developer  "))
	require.Equal(t, "Senior Backend Developer", resume.Title())
}

func TestResumeReplaceTitleRejectsEmptyTitle(t *testing.T) {
	resume, err := NewResume(validResumeParams())
	require.NoError(t, err)

	err = resume.ReplaceTitle("   ")
	require.ErrorIs(t, err, ErrInvalidResume)
	require.Equal(t, "Frontend Developer", resume.Title())
}

func TestResumeReplaceOverviewReplacesWholeSection(t *testing.T) {
	resume, err := NewResume(validResumeParams())
	require.NoError(t, err)

	first := ResumeOverviewParams{Greeting: valueobjects.NewGreeting(valueobjects.GreetingParams{First: "First"}), Sections: []valueobjects.SectionSettings{}}
	require.NoError(t, resume.ReplaceOverview(first))

	second := ResumeOverviewParams{Greeting: valueobjects.NewGreeting(valueobjects.GreetingParams{First: "Second"}), Sections: []valueobjects.SectionSettings{}}
	require.NoError(t, resume.ReplaceOverview(second))

	require.Equal(t, "Second", resume.Overview().Greeting().First())
}

func TestResumeReplaceOverviewRejectsInvalidSection(t *testing.T) {
	resume, err := NewResume(validResumeParams())
	require.NoError(t, err)

	err = resume.ReplaceOverview(ResumeOverviewParams{Sections: []valueobjects.SectionSettings{{}}})
	require.ErrorIs(t, err, valueobjects.ErrInvalidValue)
	require.Nil(t, resume.Overview())
}

func TestResumeReplaceTitleTrimsAndValidates(t *testing.T) {
	resume, err := NewResume(ResumeParams{UserID: useridentity.NewID(), Title: "Old", Language: enums.ResumeLanguage("en")})
	require.NoError(t, err)

	require.NoError(t, resume.ReplaceTitle("  New title  "))
	require.Equal(t, "New title", resume.Title())

	err = resume.ReplaceTitle("   ")
	require.ErrorIs(t, err, ErrInvalidResume)
}
