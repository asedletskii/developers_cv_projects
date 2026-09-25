package readers

import (
	"context"
	"fmt"
	"sync"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type ResumeReaderDependencies struct {
	Metadata     queries.ResumeMetadataReader
	Overview     queries.OverviewReader
	Profile      queries.ResumeProfileReader
	Contacts     queries.ResumeContactsReader
	Skills       queries.SkillsReader
	Languages    queries.LanguagesReader
	Technologies queries.TechnologiesReader
	Experiences  queries.ExperiencesReader
	Educations   queries.EducationsReader
	Courses      queries.CoursesReader
	Activities   queries.ActivitiesReader
	Projects     queries.ProjectsReader
}

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

func NewResumeReader(deps ResumeReaderDependencies) *ResumeReader {
	return &ResumeReader{
		metadata: deps.Metadata, overview: deps.Overview, profile: deps.Profile,
		contacts: deps.Contacts, skills: deps.Skills, languages: deps.Languages,
		technologies: deps.Technologies, experiences: deps.Experiences,
		educations: deps.Educations, courses: deps.Courses,
		activities: deps.Activities, projects: deps.Projects,
	}
}

func (r *ResumeReader) GetResume(ctx context.Context, q queries.GetResumeQuery) (queries.ResumeView, error) {
	if r == nil || r.metadata == nil || r.overview == nil || r.profile == nil || r.contacts == nil || r.skills == nil || r.languages == nil || r.technologies == nil || r.experiences == nil || r.educations == nil || r.courses == nil || r.activities == nil || r.projects == nil {
		return queries.ResumeView{}, fmt.Errorf("get resume: reader dependencies are nil")
	}

	metadata, err := r.metadata.GetResumeMetadata(ctx, queries.GetResumeMetadataQuery(q))
	if err != nil {
		return queries.ResumeView{}, err
	}

	v := queries.ResumeView{
		ID: metadata.ID, Title: metadata.Title, Language: metadata.Language,
		CreatedAt: metadata.CreatedAt, UpdatedAt: metadata.UpdatedAt,
		Languages: make([]queries.LanguageView, 0), Technologies: make([]queries.TechnologyView, 0),
		Experiences: make([]queries.ExperienceView, 0), Educations: make([]queries.EducationView, 0),
		Courses: make([]queries.CourseView, 0), Activities: make([]queries.ActivityView, 0), Projects: make([]queries.ProjectView, 0),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error
	setErr := func(err error) {
		if err == nil {
			return
		}
		mu.Lock()
		if firstErr == nil {
			firstErr = err
		}
		mu.Unlock()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		x, err := r.overview.GetOverview(ctx, queries.GetOverviewQuery(q))
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		v.Overview = x.Presentation
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		x, err := r.profile.GetResumeProfile(ctx, queries.GetResumeProfileQuery(q))
		if err == nil {
			mu.Lock()
			v.Profile = &x
			mu.Unlock()
			return
		}
		if err != queries.ErrResumeNotFound {
			setErr(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		x, err := r.contacts.GetResumeContacts(ctx, queries.GetResumeContactsQuery(q))
		if err == nil {
			mu.Lock()
			v.Contacts = &x
			mu.Unlock()
			return
		}
		if err != queries.ErrResumeNotFound {
			setErr(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		x, err := r.skills.GetSkills(ctx, queries.GetSkillsQuery(q))
		if err == nil {
			mu.Lock()
			v.Skills = &x
			mu.Unlock()
			return
		}
		if err != queries.ErrResumeNotFound {
			setErr(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		x, err := r.languages.ListLanguages(ctx, queries.ListLanguagesQuery(q))
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		v.Languages = x
		mu.Unlock()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		x, err := r.technologies.ListTechnologies(ctx, queries.ListTechnologiesQuery(q))
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		v.Technologies = x
		mu.Unlock()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		x, err := r.experiences.ListExperiences(ctx, queries.ListExperiencesQuery(q))
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		v.Experiences = x
		mu.Unlock()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		x, err := r.educations.ListEducations(ctx, queries.ListEducationsQuery(q))
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		v.Educations = x
		mu.Unlock()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		x, err := r.courses.ListCourses(ctx, queries.ListCoursesQuery(q))
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		v.Courses = x
		mu.Unlock()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		x, err := r.activities.ListActivities(ctx, queries.ListActivitiesQuery(q))
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		v.Activities = x
		mu.Unlock()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		x, err := r.projects.ListProjects(ctx, queries.ListProjectsQuery(q))
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		v.Projects = x
		mu.Unlock()
	}()

	wg.Wait()
	if firstErr != nil {
		return queries.ResumeView{}, firstErr
	}
	return v, nil
}
