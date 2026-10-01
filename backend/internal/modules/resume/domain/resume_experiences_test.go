package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestNewExperienceRejectsEndBeforeStart(t *testing.T) {
	id := identity.NewID()
	resumeID := identity.NewID()
	start, err := ParseYearMonth("2024-01")
	require.NoError(t, err)
	end, err := ParseYearMonth("2023-12")
	require.NoError(t, err)

	_, err = NewExperience(ExperienceParams{
		ID: id, ResumeID: resumeID, Company: "Acme", Position: "Developer",
		StartDate: start, EndDate: &end,
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestNewExperienceTrimsValuesAndBuildsChildren(t *testing.T) {
	id := identity.NewID()
	resumeID := identity.NewID()
	start, err := ParseYearMonth("2024-01")
	require.NoError(t, err)
	grade := TechnologyGradeSenior
	item, err := NewExperience(ExperienceParams{
		ID: id, ResumeID: resumeID, Company: " Acme ", Position: " Developer ",
		StartDate:    start,
		Achievements: []string{" shipped v2 ", "reduced latency"},
		Tasks:        []string{" own backend ", "review PRs "},
		Stack:        []ExperienceStackParams{{Name: " Go ", Grade: &grade}},
	})
	require.NoError(t, err)
	require.Equal(t, "Acme", item.Company())
	require.Equal(t, "Developer", item.Position())
	require.Equal(t, []string{"shipped v2", "reduced latency"}, item.Achievements())
	require.Equal(t, []string{"own backend", "review PRs"}, item.Tasks())
	require.Len(t, item.Stack(), 1)
	require.Equal(t, "Go", item.Stack()[0].Name())
	require.Equal(t, &grade, item.Stack()[0].Grade())
}

func TestNewExperienceRejectsZeroStartDate(t *testing.T) {
	_, err := NewExperience(ExperienceParams{
		ID:        identity.NewID(),
		ResumeID:  identity.NewID(),
		Company:   "Acme",
		Position:  "Developer",
		StartDate: YearMonth{},
	})

	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestNewExperienceRejectsEmptyTask(t *testing.T) {
	start, err := ParseYearMonth("2024-01")
	require.NoError(t, err)
	_, err = NewExperience(ExperienceParams{
		ID: identity.NewID(), ResumeID: identity.NewID(), Company: "Acme", Position: "Developer", StartDate: start,
		Tasks: []string{""},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func experienceParamsForTest(t *testing.T, id identity.ID, company string) ExperienceParams {
	t.Helper()
	start, err := ParseYearMonth("2024-01")
	require.NoError(t, err)
	return ExperienceParams{ID: id, ResumeID: identity.NewID(), Company: company, Position: "Developer", StartDate: start, Achievements: []string{"Built API"}}
}

func TestResumeReplaceExperiencesGeneratesIDAndReplacesSnapshot(t *testing.T) {
	resume := newTestResume(t)
	params := experienceParamsForTest(t, identity.ID{}, "Acme")
	params.ResumeID = resume.ID()
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{params}))
	require.Len(t, resume.Experiences(), 1)
	firstID := resume.Experiences()[0].ID()
	require.False(t, firstID.IsZero())

	params = experienceParamsForTest(t, firstID, "Acme Updated")
	params.ResumeID = identity.NewID()
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{params}))
	require.Equal(t, firstID, resume.Experiences()[0].ID())
	require.Equal(t, "Acme Updated", resume.Experiences()[0].Company())
	require.Equal(t, resume.ID(), resume.Experiences()[0].ResumeID())
}

func TestResumeReplaceExperiencesRejectsUnknownAndDuplicateIDs(t *testing.T) {
	resume := newTestResume(t)
	unknown := experienceParamsForTest(t, identity.NewID(), "Unknown")
	unknown.ResumeID = resume.ID()
	require.ErrorIs(t, resume.ReplaceExperiences([]ExperienceParams{unknown}), ErrResumeInvariant)

	first := experienceParamsForTest(t, identity.ID{}, "First")
	first.ResumeID = resume.ID()
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{first}))
	id := resume.Experiences()[0].ID()
	one := experienceParamsForTest(t, id, "One")
	two := experienceParamsForTest(t, id, "Two")
	require.ErrorIs(t, resume.ReplaceExperiences([]ExperienceParams{one, two}), ErrInvalidEntity)
}

func TestResumeReplaceExperiencesRemovesMissingAndPreservesOrder(t *testing.T) {
	resume := newTestResume(t)
	first := experienceParamsForTest(t, identity.ID{}, "First")
	second := experienceParamsForTest(t, identity.ID{}, "Second")
	first.ResumeID, second.ResumeID = resume.ID(), resume.ID()
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{first, second}))
	secondID := resume.Experiences()[1].ID()
	firstID := resume.Experiences()[0].ID()

	second = experienceParamsForTest(t, secondID, "Second Updated")
	second.ResumeID = identity.NewID()
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{second}))
	require.Len(t, resume.Experiences(), 1)
	require.Equal(t, secondID, resume.Experiences()[0].ID())
	require.NotEqual(t, firstID, resume.Experiences()[0].ID())
}

func TestResumeExperiencesReturnsNilForNilReceiver(t *testing.T) {
	var resume *Resume

	require.Nil(t, resume.Experiences())
}
