package readers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type compositeReaderStub struct{}

func (compositeReaderStub) GetResumeMetadata(context.Context, queries.GetResumeMetadataQuery) (queries.ResumeMetadataView, error) {
	return queries.ResumeMetadataView{ID: "resume-1", Title: "Backend Developer", Language: "en"}, nil
}

func (compositeReaderStub) ReadResumeMetadata(context.Context, valueobjects.ID) (queries.ResumeMetadataView, error) {
	return queries.ResumeMetadataView{ID: "resume-1", Title: "Backend Developer", Language: "en"}, nil
}

func (compositeReaderStub) ReadOverview(context.Context, valueobjects.ID) (queries.ResumeOverviewView, error) {
	return queries.ResumeOverviewView{}, nil
}
func (compositeReaderStub) ReadResumeProfile(context.Context, valueobjects.ID) (queries.ResumeProfileView, error) {
	return queries.ResumeProfileView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) ReadResumeContacts(context.Context, valueobjects.ID) (queries.ResumeContactsView, error) {
	return queries.ResumeContactsView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) ReadSkills(context.Context, valueobjects.ID) (queries.SkillsView, error) {
	return queries.SkillsView{}, queries.ErrResumeNotFound
}
func (compositeReaderStub) ReadLanguages(context.Context, valueobjects.ID) ([]queries.LanguageView, error) {
	return []queries.LanguageView{}, nil
}
func (compositeReaderStub) ReadTechnologies(context.Context, valueobjects.ID) ([]queries.TechnologyView, error) {
	return []queries.TechnologyView{}, nil
}
func (compositeReaderStub) ReadExperiences(context.Context, valueobjects.ID) ([]queries.ExperienceView, error) {
	return []queries.ExperienceView{}, nil
}
func (compositeReaderStub) ReadEducations(context.Context, valueobjects.ID) ([]queries.EducationView, error) {
	return []queries.EducationView{}, nil
}
func (compositeReaderStub) ReadCourses(context.Context, valueobjects.ID) ([]queries.CourseView, error) {
	return []queries.CourseView{}, nil
}
func (compositeReaderStub) ReadActivities(context.Context, valueobjects.ID) ([]queries.ActivityView, error) {
	return []queries.ActivityView{}, nil
}
func (compositeReaderStub) ReadProjects(context.Context, valueobjects.ID) ([]queries.ProjectView, error) {
	return []queries.ProjectView{}, nil
}

func TestResumeReaderRejectsNilDependencies(t *testing.T) {
	reader := NewResumeReader(ResumeReaderDependencies{})
	_, err := reader.GetOwned(context.Background(), valueobjects.NewID(), useridentity.NewID())
	require.Error(t, err)
}

func TestResumeReaderReturnsConsistentEmptyCollections(t *testing.T) {
	stub := compositeReaderStub{}
	reader := NewResumeReader(ResumeReaderDependencies{Metadata: stub, Overview: stub, Profile: stub, Contacts: stub, Skills: stub, Languages: stub, Technologies: stub, Experiences: stub, Educations: stub, Courses: stub, Activities: stub, Projects: stub})

	view, err := reader.GetOwned(context.Background(), valueobjects.NewID(), useridentity.NewID())
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

func TestResumeReaderPublicDoesNotRequireUserID(t *testing.T) {
	stub := compositeReaderStub{}
	reader := NewResumeReader(ResumeReaderDependencies{Metadata: stub, Overview: stub, Profile: stub, Contacts: stub, Skills: stub, Languages: stub, Technologies: stub, Experiences: stub, Educations: stub, Courses: stub, Activities: stub, Projects: stub})

	view, err := reader.GetPublic(context.Background(), valueobjects.NewID())
	require.NoError(t, err)
	require.Equal(t, "resume-1", view.ID)
}
