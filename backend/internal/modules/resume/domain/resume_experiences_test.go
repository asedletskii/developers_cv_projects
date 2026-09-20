package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestNewExperienceRejectsEndBeforeStart(t *testing.T) {
	id := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	start, err := ParseYearMonth("2024-01")
	require.NoError(t, err)
	end, err := ParseYearMonth("2023-12")
	require.NoError(t, err)

	_, err = newExperience(id, resumeID, ExperienceParams{
		ID: id, Company: "Acme", Position: "Developer",
		StartDate: start, EndDate: &end,
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestNewExperienceTrimsValuesAndBuildsChildren(t *testing.T) {
	id := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	start, err := ParseYearMonth("2024-01")
	require.NoError(t, err)
	item, err := newExperience(id, resumeID, ExperienceParams{
		ID: id, Company: " Acme ", Position: " Developer ",
		StartDate:    start,
		Achievements: []string{" shipped v2 ", "reduced latency"},
		Tasks:        []string{" own backend ", "review PRs "},
		Stack:        []ExperienceStackParams{{Name: " Go "}},
	})
	require.NoError(t, err)
	require.Equal(t, "Acme", item.Company())
	require.Equal(t, "Developer", item.Position())
	require.Equal(t, []string{"shipped v2", "reduced latency"}, item.Achievements())
	require.Equal(t, []string{"own backend", "review PRs"}, item.Tasks())
	require.Len(t, item.Stack(), 1)
	require.Equal(t, "Go", item.Stack()[0].Name())
}

func TestNewExperienceRejectsZeroStartDate(t *testing.T) {
	_, err := newExperience(sharedvalueobjects.NewID(), sharedvalueobjects.NewID(), ExperienceParams{
		ID:        sharedvalueobjects.NewID(),
		Company:   "Acme",
		Position:  "Developer",
		StartDate: YearMonth{},
	})

	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestNewExperienceRejectsEmptyTask(t *testing.T) {
	start, err := ParseYearMonth("2024-01")
	require.NoError(t, err)
	_, err = newExperience(sharedvalueobjects.NewID(), sharedvalueobjects.NewID(), ExperienceParams{
		ID: sharedvalueobjects.NewID(), Company: "Acme", Position: "Developer", StartDate: start,
		Tasks: []string{""},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func experienceParamsForTest(t *testing.T, id sharedvalueobjects.ID, company string) ExperienceParams {
	t.Helper()
	start, err := ParseYearMonth("2024-01")
	require.NoError(t, err)
	return ExperienceParams{ID: id, Company: company, Position: "Developer", StartDate: start, Achievements: []string{"Built API"}}
}

func TestResumeReplaceExperiencesGeneratesIDAndReplacesSnapshot(t *testing.T) {
	resume := newTestResume(t)
	params := experienceParamsForTest(t, sharedvalueobjects.ID{}, "Acme")
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{params}))
	require.Len(t, resume.Experiences(), 1)
	firstID := resume.Experiences()[0].ID()
	require.False(t, firstID.IsZero())

	params = experienceParamsForTest(t, firstID, "Acme Updated")
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{params}))
	require.Equal(t, firstID, resume.Experiences()[0].ID())
	require.Equal(t, "Acme Updated", resume.Experiences()[0].Company())
	require.Equal(t, resume.ID(), resume.Experiences()[0].ResumeID())
}

func TestResumeReplaceExperiencesRejectsUnknownAndDuplicateIDs(t *testing.T) {
	resume := newTestResume(t)
	unknown := experienceParamsForTest(t, sharedvalueobjects.NewID(), "Unknown")
	err := resume.ReplaceExperiences([]ExperienceParams{unknown})
	require.ErrorIs(t, err, ErrResumeInvariant)
	require.ErrorIs(t, err, ErrEntityNotFound)

	first := experienceParamsForTest(t, sharedvalueobjects.ID{}, "First")
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{first}))
	id := resume.Experiences()[0].ID()
	one := experienceParamsForTest(t, id, "One")
	two := experienceParamsForTest(t, id, "Two")
	require.ErrorIs(t, resume.ReplaceExperiences([]ExperienceParams{one, two}), ErrInvalidEntity)
}

func TestResumeReplaceExperiencesRemovesMissingAndPreservesOrder(t *testing.T) {
	resume := newTestResume(t)
	first := experienceParamsForTest(t, sharedvalueobjects.ID{}, "First")
	second := experienceParamsForTest(t, sharedvalueobjects.ID{}, "Second")
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{first, second}))
	secondID := resume.Experiences()[1].ID()
	firstID := resume.Experiences()[0].ID()

	second = experienceParamsForTest(t, secondID, "Second Updated")
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{second}))
	require.Len(t, resume.Experiences(), 1)
	require.Equal(t, secondID, resume.Experiences()[0].ID())
	require.NotEqual(t, firstID, resume.Experiences()[0].ID())
}

func TestResumeExperiencesReturnsNilForNilReceiver(t *testing.T) {
	var resume *Resume

	require.Nil(t, resume.Experiences())
}

func TestResumeExperiencesReturnsDetachedEntities(t *testing.T) {
	resume := newTestResume(t)
	start, err := ParseYearMonth("2024-01")
	require.NoError(t, err)
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{{Company: "Original", Position: "Engineer", StartDate: start}}))
	returned := resume.Experiences()
	returned[0].company = "Changed"
	require.Equal(t, "Original", resume.Experiences()[0].Company())
}

func TestResumeReplaceExperiencesDoesNotPartiallyMutateOnError(t *testing.T) {
	resume := newTestResume(t)
	first := experienceParamsForTest(t, sharedvalueobjects.ID{}, "Existing")
	require.NoError(t, resume.ReplaceExperiences([]ExperienceParams{first}))
	existing := resume.Experiences()[0]

	invalid := experienceParamsForTest(t, sharedvalueobjects.ID{}, "Valid replacement")
	invalid.StartDate = YearMonth{}
	err := resume.ReplaceExperiences([]ExperienceParams{invalid})
	require.ErrorIs(t, err, ErrInvalidEntity)

	items := resume.Experiences()
	require.Len(t, items, 1)
	require.Equal(t, existing.ID(), items[0].ID())
	require.Equal(t, existing.Company(), items[0].Company())
}

func TestResumeReplaceExperiencesRejectsDuplicateStackTechnologies(t *testing.T) {
	resume := newTestResume(t)
	startDate, parseErr := ParseYearMonth("2024-01")
	require.NoError(t, parseErr)
	_, err := newExperience(sharedvalueobjects.NewID(), resume.ID(), ExperienceParams{
		Company: "Acme", Position: "Developer", StartDate: startDate,
		Stack: []ExperienceStackParams{{Name: "Go"}, {Name: " Go "}},
	})
	require.ErrorIs(t, err, ErrDuplicateEntity)
}
