package mappers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestActivitiesToDomainPreservesOrderAndAddsRowContext(t *testing.T) {
	resumeID := uuid.New()
	firstID := uuid.New()
	secondID := uuid.New()

	activities, err := ActivitiesToDomain([]gormmodel.ResumeActivity{
		{ID: firstID, ResumeID: resumeID, Title: "First"},
		{ID: secondID, ResumeID: resumeID, Title: "Second"},
	})

	require.NoError(t, err)
	require.Len(t, activities, 2)
	require.Equal(t, firstID, activities[0].ID().UUID())
	require.Equal(t, "First", activities[0].Title())
	require.Equal(t, secondID, activities[1].ID().UUID())
	require.Equal(t, "Second", activities[1].Title())

	_, err = ActivitiesToDomain([]gormmodel.ResumeActivity{{ID: firstID, ResumeID: resumeID}, {}})
	require.ErrorIs(t, err, sharedvalueobjects.ErrInvalidID)
	require.ErrorContains(t, err, "activity at index 1")
}

func TestCoursesToDomainPreservesOrderAndAddsRowContext(t *testing.T) {
	resumeID := uuid.New()
	firstID := uuid.New()
	secondID := uuid.New()

	courses, err := CoursesToDomain([]gormmodel.ResumeCourse{
		{ID: firstID, ResumeID: resumeID, Name: "First"},
		{ID: secondID, ResumeID: resumeID, Name: "Second"},
	})

	require.NoError(t, err)
	require.Len(t, courses, 2)
	require.Equal(t, firstID, courses[0].ID().UUID())
	require.Equal(t, "First", courses[0].Name())
	require.Equal(t, secondID, courses[1].ID().UUID())
	require.Equal(t, "Second", courses[1].Name())

	_, err = CoursesToDomain([]gormmodel.ResumeCourse{{ID: firstID, ResumeID: resumeID}, {}})
	require.ErrorIs(t, err, sharedvalueobjects.ErrInvalidID)
	require.ErrorContains(t, err, "course at index 1")
}

func TestExperiencesToDomainMapsGroupedRowsAndAddsRowContext(t *testing.T) {
	resumeID := uuid.New()
	firstID := uuid.New()
	secondID := uuid.New()
	achievements := map[uuid.UUID][]gormmodel.ResumeExperienceAchievement{
		firstID: {{ExperienceID: firstID, Position: 0, Value: "First achievement"}},
	}

	experiences, err := ExperiencesToDomain([]gormmodel.ResumeExperience{
		{ID: firstID, ResumeID: resumeID, Company: "First"},
		{ID: secondID, ResumeID: resumeID, Company: "Second"},
	}, achievements, nil, nil)

	require.NoError(t, err)
	require.Len(t, experiences, 2)
	require.Equal(t, firstID, experiences[0].ID().UUID())
	require.Equal(t, "First", experiences[0].Company())
	require.Equal(t, []string{"First achievement"}, experiences[0].Achievements())
	require.Equal(t, secondID, experiences[1].ID().UUID())
	require.Equal(t, "Second", experiences[1].Company())

	_, err = ExperiencesToDomain([]gormmodel.ResumeExperience{{ID: firstID, ResumeID: resumeID}, {}}, nil, nil, nil)
	require.ErrorIs(t, err, sharedvalueobjects.ErrInvalidID)
	require.ErrorContains(t, err, "experience at index 1")
}

func TestProjectsToDomainMapsGroupedRowsAndAddsRowContext(t *testing.T) {
	resumeID := uuid.New()
	firstID := uuid.New()
	secondID := uuid.New()
	screenshots := map[uuid.UUID][]gormmodel.ResumeProjectScreenshot{
		firstID: {{ProjectID: firstID, Position: 0, Path: "https://example.com/image.png"}},
	}

	projects, err := ProjectsToDomain([]gormmodel.ResumeProject{
		{ID: firstID, ResumeID: resumeID, Name: "First"},
		{ID: secondID, ResumeID: resumeID, Name: "Second"},
	}, screenshots, nil, nil)

	require.NoError(t, err)
	require.Len(t, projects, 2)
	require.Equal(t, firstID, projects[0].ID().UUID())
	require.Equal(t, "First", projects[0].Name())
	require.Len(t, projects[0].Screenshots(), 1)
	require.Equal(t, "https://example.com/image.png", projects[0].Screenshots()[0].String())
	require.Equal(t, secondID, projects[1].ID().UUID())
	require.Equal(t, "Second", projects[1].Name())

	_, err = ProjectsToDomain([]gormmodel.ResumeProject{{ID: firstID, ResumeID: resumeID}, {}}, nil, nil, nil)
	require.ErrorIs(t, err, sharedvalueobjects.ErrInvalidID)
	require.ErrorContains(t, err, "project at index 1")
}
