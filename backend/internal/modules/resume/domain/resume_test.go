package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func newTestResume(t *testing.T) *Resume {
	t.Helper()
	resume, err := NewResume(ResumeParams{
		UserID:   sharedvalueobjects.NewID(),
		Title:    "Resume",
		Language: ResumeLanguageEN,
	})
	require.NoError(t, err)
	return resume
}

func validResumeParams() ResumeParams {
	return ResumeParams{
		UserID: sharedvalueobjects.NewID(), Title: "Frontend Developer",
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

	another, err := NewResume(ResumeParams{UserID: sharedvalueobjects.NewID(), Title: "Backend Developer", Language: ResumeLanguageEN})
	require.NoError(t, err)
	require.NotEqual(t, resume.ID(), another.ID())
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

func TestRestoreResumeAllowsAbsentOptionalChildren(t *testing.T) {
	resume := RestoreResume(ResumeSnapshot{
		ID:       sharedvalueobjects.NewID(),
		UserID:   sharedvalueobjects.NewID(),
		Title:    "Backend",
		Language: ResumeLanguageEN,
		Status:   ResumeStatusDraft,
		Version:  1,
	})
	require.NotNil(t, resume)
	require.Nil(t, resume.Profile())
	require.Nil(t, resume.Contacts())
	require.Nil(t, resume.Skills())
	require.Nil(t, resume.SpokenLanguages())
	require.Nil(t, resume.Technologies())
}

func TestRestoreResumeDetachesChildrenFromSnapshot(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	photoURL, err := NewExternalURL("https://example.com/photo.jpg")
	require.NoError(t, err)
	dob, err := NewDate(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	profile := RestoreProfile(ProfileSnapshot{
		ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "Original", Surname: "Person",
		Position: "Developer", PhotoURL: photoURL, DateOfBirth: dob, Location: "Paris",
	})
	resume := RestoreResume(ResumeSnapshot{
		ID: resumeID, UserID: sharedvalueobjects.NewID(), Title: "Backend",
		Language: ResumeLanguageEN, Status: ResumeStatusDraft, Version: 1, Profile: profile,
	})

	profile.name = "Mutated outside aggregate"
	require.Equal(t, "Original", resume.Profile().Name())
}

func TestRestoreResumeDoesNotRevalidatePersistedChildFields(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	profile := RestoreProfile(ProfileSnapshot{
		ID: sharedvalueobjects.NewID(), ResumeID: resumeID,
	})
	resume := RestoreResume(ResumeSnapshot{
		ID: resumeID, UserID: sharedvalueobjects.NewID(), Title: "Backend",
		Language: ResumeLanguageEN, Status: ResumeStatusDraft, Version: 1, Profile: profile,
	})
	require.NotNil(t, resume.Profile())
	require.Equal(t, profile.ID(), resume.Profile().ID())
}

func TestRestoreResumePreservesPersistedChildOwnership(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	foreignResumeID := sharedvalueobjects.NewID()
	profile := RestoreProfile(ProfileSnapshot{
		ID: sharedvalueobjects.NewID(), ResumeID: foreignResumeID,
	})

	resume := RestoreResume(ResumeSnapshot{
		ID: resumeID, UserID: sharedvalueobjects.NewID(), Profile: profile,
	})

	require.Equal(t, foreignResumeID, resume.Profile().ResumeID())
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

	require.NoError(t, resume.ReplaceStatus(ResumeStatusOpened))
	require.Equal(t, ResumeStatusOpened, resume.Status())
	require.ErrorIs(t, resume.ReplaceStatus(ResumeStatusOpened), ErrInvalidResumeStatusTransition)
	require.NoError(t, resume.ReplaceStatus(ResumeStatusClosed))
	require.Equal(t, ResumeStatusClosed, resume.Status())
	require.NoError(t, resume.ReplaceStatus(ResumeStatusOpened))
	require.Equal(t, ResumeStatusOpened, resume.Status())
}

func TestResumeCannotCloseDraft(t *testing.T) {
	resume, err := NewResume(ResumeParams{
		UserID: sharedvalueobjects.NewID(), Title: "Backend", Language: ResumeLanguageEN,
	})
	require.NoError(t, err)
	require.ErrorIs(t, resume.ReplaceStatus(ResumeStatusClosed), ErrInvalidResumeStatusTransition)
	require.Equal(t, ResumeStatusDraft, resume.Status())
}

func TestNormalizeOptionalStringTrimsWhitespaceAndTurnsEmptyIntoNil(t *testing.T) {
	value := "  optional value  "
	got := normalizeOptionalString(&value)
	require.NotNil(t, got)
	require.Equal(t, "optional value", *got)

	empty := "   "
	require.Nil(t, normalizeOptionalString(&empty))
	require.Nil(t, normalizeOptionalString(nil))
}

func TestResumeTitleLengthLimit(t *testing.T) {
	valid := strings.Repeat("я", MaxResumeTitleLength)
	require.NoError(t, validateResumeTitle(valid))
	require.Error(t, validateResumeTitle(valid+"я"))
}

func TestRestoreResumePreservesSnapshotWithoutValidation(t *testing.T) {
	resume := RestoreResume(ResumeSnapshot{
		ID: sharedvalueobjects.ID{}, UserID: sharedvalueobjects.ID{},
		Title: "", Language: ResumeLanguage("legacy"), Status: ResumeStatus("legacy"), Version: 0,
	})

	require.True(t, resume.ID().IsZero())
	require.True(t, resume.UserID().IsZero())
	require.Empty(t, resume.Title())
	require.Equal(t, ResumeLanguage("legacy"), resume.Language())
	require.Equal(t, ResumeStatus("legacy"), resume.Status())
	require.Zero(t, resume.Version())
}

func TestRestoreResumePreservesDuplicateChildIDs(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	experienceID := sharedvalueobjects.NewID()
	start, err := ParseYearMonth("2024-01")
	require.NoError(t, err)
	first := RestoreExperience(ExperienceSnapshot{ID: experienceID, ResumeID: resumeID, Company: "Acme", Position: "Developer", StartDate: start})
	second := RestoreExperience(ExperienceSnapshot{ID: experienceID, ResumeID: resumeID, Company: "Other", Position: "Engineer", StartDate: start})

	resume := RestoreResume(ResumeSnapshot{ID: resumeID, UserID: sharedvalueobjects.NewID(), Experiences: []*Experience{first, second}})

	require.Len(t, resume.Experiences(), 2)
	require.Equal(t, resume.Experiences()[0].ID(), resume.Experiences()[1].ID())
}

func TestRestoreResumePreservesRootSnapshotFields(t *testing.T) {
	snapshot := ResumeSnapshot{
		ID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID(), Title: " Backend ",
		Language: ResumeLanguage("legacy"), Status: ResumeStatus("legacy"), Version: 0,
	}

	resume := RestoreResume(snapshot)
	require.Equal(t, snapshot.Title, resume.Title())
	require.Equal(t, snapshot.Language, resume.Language())
	require.Equal(t, snapshot.Status, resume.Status())
	require.Equal(t, snapshot.Version, resume.Version())
}

func TestRestoreResumePreservesStoredCollectionPositions(t *testing.T) {
	technologies := RestoreResumeTechnologies(ResumeTechnologiesParams{Entries: []ResumeTechnologyEntryParams{
		{ID: sharedvalueobjects.NewID(), Name: "Go", Position: 0},
		{ID: sharedvalueobjects.NewID(), Name: "PostgreSQL", Position: 0},
	}})

	resume := RestoreResume(ResumeSnapshot{ID: sharedvalueobjects.NewID(), Technologies: technologies})

	entries := resume.Technologies().Entries()
	require.Len(t, entries, 2)
	require.Equal(t, 0, entries[0].Position())
	require.Equal(t, 0, entries[1].Position())
}

func mustGreeting(t *testing.T, params GreetingParams) Greeting {
	t.Helper()
	value, err := NewGreeting(params)
	require.NoError(t, err)
	return value
}

func TestRestoreResumePreservesChildEntitySnapshotIdentity(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	userID := sharedvalueobjects.NewID()
	experienceID := sharedvalueobjects.NewID()
	start, err := ParseYearMonth("2024-01")
	require.NoError(t, err)
	experience := RestoreExperience(ExperienceSnapshot{
		ID: experienceID, ResumeID: resumeID, Company: "Acme", Position: "Developer",
		StartDate: start, Achievements: []string{"Delivered API"},
	})

	resume := RestoreResume(ResumeSnapshot{
		ID: resumeID, UserID: userID, Title: "Backend", Language: ResumeLanguageEN,
		Status: ResumeStatusDraft, Version: 1, Experiences: []*Experience{experience},
	})
	require.Len(t, resume.Experiences(), 1)
	require.Equal(t, experienceID, resume.Experiences()[0].ID())
	require.Equal(t, resumeID, resume.Experiences()[0].ResumeID())
}

func TestRestoreResumePreservesEmptyTechnologyCollections(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	resume := RestoreResume(ResumeSnapshot{
		ID: resumeID, UserID: sharedvalueobjects.NewID(), Title: "Backend",
		Language: ResumeLanguageEN, Status: ResumeStatusDraft, Version: 1,
		Technologies:    RestoreResumeTechnologies(ResumeTechnologiesParams{}),
		SpokenLanguages: RestoreResumeSpokenLanguages(ResumeSpokenLanguagesParams{}),
	})
	require.NotNil(t, resume.Technologies())
	require.Empty(t, resume.Technologies().Entries())
	require.NotNil(t, resume.SpokenLanguages())
	require.Empty(t, resume.SpokenLanguages().Entries())
}

func TestRestoreResumePreservesPersistedSkillsSnapshot(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	skills := RestoreSkills(SkillsSnapshot{
		ID: sharedvalueobjects.NewID(), ResumeID: resumeID,
		Hard: []string{" Go "},
	})

	resume := RestoreResume(ResumeSnapshot{
		ID: resumeID, UserID: sharedvalueobjects.NewID(), Title: "Backend",
		Language: ResumeLanguageEN, Status: ResumeStatusDraft, Version: 1, Skills: skills,
	})
	require.Equal(t, []string{" Go "}, resume.Skills().HardSkills())
}

func TestRestoreResumePreservesCollectionOrderAndStoredPositions(t *testing.T) {
	technologies := RestoreResumeTechnologies(ResumeTechnologiesParams{Entries: []ResumeTechnologyEntryParams{
		{ID: sharedvalueobjects.NewID(), Name: "Go", Position: 1},
		{ID: sharedvalueobjects.NewID(), Name: "Rust", Position: 0},
	}})

	resume := RestoreResume(ResumeSnapshot{ID: sharedvalueobjects.NewID(), Technologies: technologies})

	entries := resume.Technologies().Entries()
	require.Equal(t, "Go", entries[0].Name())
	require.Equal(t, 1, entries[0].Position())
	require.Equal(t, "Rust", entries[1].Name())
	require.Equal(t, 0, entries[1].Position())
}
