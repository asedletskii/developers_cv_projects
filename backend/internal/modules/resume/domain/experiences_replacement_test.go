package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func experienceParamsForTest(t *testing.T, id valueobjects.ID, company string) ExperienceParams {
	t.Helper()
	start, err := valueobjects.ParseYearMonth("2024-01")
	require.NoError(t, err)
	return ExperienceParams{ID: id, ResumeID: valueobjects.NewID(), Company: company, Position: "Developer", StartDate: start, Achievements: []string{"Built API"}}
}

func TestResumeReplaceExperiencesGeneratesIDAndReplacesSnapshot(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	params := experienceParamsForTest(t, valueobjects.ID{}, "Acme")
	params.ResumeID = resume.ID()
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{params}))
	require.Len(t, resume.Experiences(), 1)
	firstID := resume.Experiences()[0].ID()
	require.False(t, firstID.IsZero())

	params = experienceParamsForTest(t, firstID, "Acme Updated")
	params.ResumeID = valueobjects.NewID()
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{params}))
	require.Equal(t, firstID, resume.Experiences()[0].ID())
	require.Equal(t, "Acme Updated", resume.Experiences()[0].Company())
	require.Equal(t, resume.ID(), resume.Experiences()[0].ResumeID())
}

func TestResumeReplaceExperiencesRejectsUnknownAndDuplicateIDs(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	unknown := experienceParamsForTest(t, valueobjects.NewID(), "Unknown")
	unknown.ResumeID = resume.ID()
	require.ErrorIs(t, resume.ReplaceExperiences([]ExperienceParams{unknown}), ErrResumeInvariant)

	first := experienceParamsForTest(t, valueobjects.ID{}, "First")
	first.ResumeID = resume.ID()
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{first}))
	id := resume.Experiences()[0].ID()
	one := experienceParamsForTest(t, id, "One")
	two := experienceParamsForTest(t, id, "Two")
	require.ErrorIs(t, resume.ReplaceExperiences([]ExperienceParams{one, two}), ErrInvalidEntity)
}

func TestResumeReplaceExperiencesRemovesMissingAndPreservesOrder(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	first := experienceParamsForTest(t, valueobjects.ID{}, "First")
	second := experienceParamsForTest(t, valueobjects.ID{}, "Second")
	first.ResumeID, second.ResumeID = resume.ID(), resume.ID()
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{first, second}))
	secondID := resume.Experiences()[1].ID()
	firstID := resume.Experiences()[0].ID()

	second = experienceParamsForTest(t, secondID, "Second Updated")
	second.ResumeID = valueobjects.NewID()
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{second}))
	require.Len(t, resume.Experiences(), 1)
	require.Equal(t, secondID, resume.Experiences()[0].ID())
	require.NotEqual(t, firstID, resume.Experiences()[0].ID())
}
