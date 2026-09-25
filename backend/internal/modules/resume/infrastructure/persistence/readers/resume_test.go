package readers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type compositeReaderStub struct{}

func (compositeReaderStub) GetResumeMetadata(context.Context, queries.GetResumeMetadataQuery) (queries.ResumeMetadataView, error) {
	return queries.ResumeMetadataView{ID: "resume-1", Title: "Backend Developer", Language: "en"}, nil
}

func (compositeReaderStub) GetOverview(context.Context, queries.GetOverviewQuery) (queries.ResumeOverviewView, error) {
	return queries.ResumeOverviewView{}, nil
}
func (compositeReaderStub) GetResumeProfile(context.Context, queries.GetResumeProfileQuery) (queries.ResumeProfileView, error) {
	return queries.ResumeProfileView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) GetResumeContacts(context.Context, queries.GetResumeContactsQuery) (queries.ResumeContactsView, error) {
	return queries.ResumeContactsView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) GetSkills(context.Context, queries.GetSkillsQuery) (queries.SkillsView, error) {
	return queries.SkillsView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) ListLanguages(context.Context, queries.ListLanguagesQuery) ([]queries.LanguageView, error) {
	return []queries.LanguageView{}, nil
}
func (compositeReaderStub) GetLanguage(context.Context, queries.GetLanguageQuery) (queries.LanguageView, error) {
	return queries.LanguageView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) ListTechnologies(context.Context, queries.ListTechnologiesQuery) ([]queries.TechnologyView, error) {
	return []queries.TechnologyView{}, nil
}
func (compositeReaderStub) GetTechnology(context.Context, queries.GetTechnologyQuery) (queries.TechnologyView, error) {
	return queries.TechnologyView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) ListExperiences(context.Context, queries.ListExperiencesQuery) ([]queries.ExperienceView, error) {
	return []queries.ExperienceView{}, nil
}
func (compositeReaderStub) GetExperience(context.Context, queries.GetExperienceQuery) (queries.ExperienceView, error) {
	return queries.ExperienceView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) ListEducations(context.Context, queries.ListEducationsQuery) ([]queries.EducationView, error) {
	return []queries.EducationView{}, nil
}
func (compositeReaderStub) GetEducation(context.Context, queries.GetEducationQuery) (queries.EducationView, error) {
	return queries.EducationView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) ListCourses(context.Context, queries.ListCoursesQuery) ([]queries.CourseView, error) {
	return []queries.CourseView{}, nil
}
func (compositeReaderStub) GetCourse(context.Context, queries.GetCourseQuery) (queries.CourseView, error) {
	return queries.CourseView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) ListActivities(context.Context, queries.ListActivitiesQuery) ([]queries.ActivityView, error) {
	return []queries.ActivityView{}, nil
}
func (compositeReaderStub) GetActivity(context.Context, queries.GetActivityQuery) (queries.ActivityView, error) {
	return queries.ActivityView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) ListProjects(context.Context, queries.ListProjectsQuery) ([]queries.ProjectView, error) {
	return []queries.ProjectView{}, nil
}
func (compositeReaderStub) GetProject(context.Context, queries.GetProjectQuery) (queries.ProjectView, error) {
	return queries.ProjectView{}, queries.ErrResumeNotFound
}

func TestResumeReaderRejectsNilDependencies(t *testing.T) {
	reader := NewResumeReader(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	_, err := reader.GetResume(context.Background(), queries.GetResumeQuery{ResumeID: valueobjects.NewID(), UserID: useridentity.NewID()})
	require.Error(t, err)
}

func TestResumeReaderReturnsConsistentEmptyCollections(t *testing.T) {
	stub := compositeReaderStub{}
	reader := NewResumeReader(stub, stub, stub, stub, stub, stub, stub, stub, stub, stub, stub, stub)

	view, err := reader.GetResume(context.Background(), queries.GetResumeQuery{ResumeID: valueobjects.NewID(), UserID: useridentity.NewID()})
	require.NoError(t, err)
	require.Equal(t, "resume-1", view.ID)
	require.Nil(t, view.Profile)
	require.Nil(t, view.Contacts)
	require.Nil(t, view.Skills)
	require.NotNil(t, view.Languages)
	require.NotNil(t, view.Technologies)
	require.NotNil(t, view.Experiences)
	require.NotNil(t, view.Educations)
	require.NotNil(t, view.Courses)
	require.NotNil(t, view.Activities)
	require.NotNil(t, view.Projects)
	require.Empty(t, view.Languages)
}
