package readers

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type ResumeReader struct {
	metadata     queries.ResumeMetadataReader
	overview     queries.OverviewReader
	profile      queries.ResumeProfileReader
	contacts     queries.ResumeContactsReader
	skills       queries.SkillsReader
	languages    queries.LanguagesReader
	technologies queries.TechnologiesReader
	experiences  queries.ExperiencesReader
	educations   queries.EducationsReader
	courses      queries.CoursesReader
	activities   queries.ActivitiesReader
	projects     queries.ProjectsReader
}

func NewResumeReader(
	metadata queries.ResumeMetadataReader,
	overview queries.OverviewReader,
	profile queries.ResumeProfileReader,
	contacts queries.ResumeContactsReader,
	skills queries.SkillsReader,
	languages queries.LanguagesReader,
	technologies queries.TechnologiesReader,
	experiences queries.ExperiencesReader,
	educations queries.EducationsReader,
	courses queries.CoursesReader,
	activities queries.ActivitiesReader,
	projects queries.ProjectsReader,
) *ResumeReader {
	return &ResumeReader{metadata: metadata, overview: overview, profile: profile, contacts: contacts, skills: skills, languages: languages, technologies: technologies, experiences: experiences, educations: educations, courses: courses, activities: activities, projects: projects}
}

func (r *ResumeReader) GetResume(ctx context.Context, q queries.GetResumeQuery) (queries.ResumeView, error) {
	if r == nil || r.metadata == nil || r.overview == nil || r.profile == nil || r.contacts == nil || r.skills == nil || r.languages == nil || r.technologies == nil || r.experiences == nil || r.educations == nil || r.courses == nil || r.activities == nil || r.projects == nil {
		return queries.ResumeView{}, fmt.Errorf("get resume: reader dependencies are nil")
	}

	metadata, err := r.metadata.GetResumeMetadata(ctx, queries.GetResumeMetadataQuery(q))
	if err != nil {
		return queries.ResumeView{}, err
	}
	overview, err := r.overview.GetOverview(ctx, queries.GetOverviewQuery(q))
	if err != nil {
		return queries.ResumeView{}, err
	}

	v := queries.ResumeView{
		ID: metadata.ID, Title: metadata.Title, Language: metadata.Language,
		CreatedAt: metadata.CreatedAt, UpdatedAt: metadata.UpdatedAt,
		Overview:  overview.Presentation,
		Languages: make([]queries.LanguageView, 0), Technologies: make([]queries.TechnologyView, 0),
		Experiences: make([]queries.ExperienceView, 0), Educations: make([]queries.EducationView, 0),
		Courses: make([]queries.CourseView, 0), Activities: make([]queries.ActivityView, 0), Projects: make([]queries.ProjectView, 0),
	}

	if x, e := r.profile.GetResumeProfile(ctx, queries.GetResumeProfileQuery(q)); e == nil {
		v.Profile = &x
	} else if e != queries.ErrResumeNotFound {
		return queries.ResumeView{}, e
	}
	if x, e := r.contacts.GetResumeContacts(ctx, queries.GetResumeContactsQuery(q)); e == nil {
		v.Contacts = &x
	} else if e != queries.ErrResumeNotFound {
		return queries.ResumeView{}, e
	}
	if x, e := r.skills.GetSkills(ctx, queries.GetSkillsQuery(q)); e == nil {
		v.Skills = &x
	} else if e != queries.ErrResumeNotFound {
		return queries.ResumeView{}, e
	}

	if v.Languages, err = r.languages.ListLanguages(ctx, queries.ListLanguagesQuery(q)); err != nil {
		return queries.ResumeView{}, err
	}
	if v.Technologies, err = r.technologies.ListTechnologies(ctx, queries.ListTechnologiesQuery(q)); err != nil {
		return queries.ResumeView{}, err
	}
	if v.Experiences, err = r.experiences.ListExperiences(ctx, queries.ListExperiencesQuery(q)); err != nil {
		return queries.ResumeView{}, err
	}
	if v.Educations, err = r.educations.ListEducations(ctx, queries.ListEducationsQuery(q)); err != nil {
		return queries.ResumeView{}, err
	}
	if v.Courses, err = r.courses.ListCourses(ctx, queries.ListCoursesQuery(q)); err != nil {
		return queries.ResumeView{}, err
	}
	if v.Activities, err = r.activities.ListActivities(ctx, queries.ListActivitiesQuery(q)); err != nil {
		return queries.ResumeView{}, err
	}
	if v.Projects, err = r.projects.ListProjects(ctx, queries.ListProjectsQuery(q)); err != nil {
		return queries.ResumeView{}, err
	}

	return v, nil
}
