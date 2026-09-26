package readers

import (
	"context"
	"fmt"
	"sync"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeReaderDependencies struct {
	Metadata     queries.ResumeMetadataReader
	Overview     queries.OverviewSource
	Profile      queries.ResumeProfileSource
	Contacts     queries.ResumeContactsSource
	Skills       queries.SkillsSource
	Languages    queries.LanguagesSource
	Technologies queries.TechnologiesSource
	Experiences  queries.ExperiencesSource
	Educations   queries.EducationsSource
	Courses      queries.CoursesSource
	Activities   queries.ActivitiesSource
	Projects     queries.ProjectsSource
}

type ResumeReader struct {
	metadata     queries.ResumeMetadataReader
	overview     queries.OverviewSource
	profile      queries.ResumeProfileSource
	contacts     queries.ResumeContactsSource
	skills       queries.SkillsSource
	languages    queries.LanguagesSource
	technologies queries.TechnologiesSource
	experiences  queries.ExperiencesSource
	educations   queries.EducationsSource
	courses      queries.CoursesSource
	activities   queries.ActivitiesSource
	projects     queries.ProjectsSource
}

func NewResumeReader(deps ResumeReaderDependencies) *ResumeReader {
	return &ResumeReader{
		metadata:     deps.Metadata,
		overview:     deps.Overview,
		profile:      deps.Profile,
		contacts:     deps.Contacts,
		skills:       deps.Skills,
		languages:    deps.Languages,
		technologies: deps.Technologies,
		experiences:  deps.Experiences,
		educations:   deps.Educations,
		courses:      deps.Courses,
		activities:   deps.Activities,
		projects:     deps.Projects,
	}
}

func (r *ResumeReader) GetOwned(ctx context.Context, resumeID valueobjects.ID, userID useridentity.ID) (queries.ResumeView, error) {
	if r == nil || r.metadata == nil {
		return queries.ResumeView{}, fmt.Errorf("get owned resume: reader dependencies are nil")
	}
	if resumeID.IsZero() || userID.IsZero() {
		return queries.ResumeView{}, fmt.Errorf("get owned resume: invalid identifiers")
	}
	metadata, err := r.metadata.GetResumeMetadata(ctx, queries.GetResumeMetadataQuery{ResumeID: resumeID, UserID: userID})
	if err != nil {
		return queries.ResumeView{}, err
	}
	return r.read(ctx, resumeID, metadata)
}

func (r *ResumeReader) GetPublic(ctx context.Context, resumeID valueobjects.ID) (queries.ResumeView, error) {
	if r == nil || r.metadata == nil {
		return queries.ResumeView{}, fmt.Errorf("get public resume: reader dependencies are nil")
	}
	if resumeID.IsZero() {
		return queries.ResumeView{}, fmt.Errorf("get public resume: invalid resume id")
	}
	metadata, err := r.metadata.ReadResumeMetadata(ctx, resumeID)
	if err != nil {
		return queries.ResumeView{}, err
	}
	return r.read(ctx, resumeID, metadata)
}

func (r *ResumeReader) read(ctx context.Context, resumeID valueobjects.ID, metadata queries.ResumeMetadataView) (queries.ResumeView, error) {
	if r == nil || r.metadata == nil || r.overview == nil || r.profile == nil || r.contacts == nil || r.skills == nil || r.languages == nil || r.technologies == nil || r.experiences == nil || r.educations == nil || r.courses == nil || r.activities == nil || r.projects == nil {
		return queries.ResumeView{}, fmt.Errorf("get resume: reader dependencies are nil")
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
		x, err := r.overview.ReadOverview(ctx, resumeID)
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
		x, err := r.profile.ReadResumeProfile(ctx, resumeID)
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
		x, err := r.contacts.ReadResumeContacts(ctx, resumeID)
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
		x, err := r.skills.ReadSkills(ctx, resumeID)
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
		x, err := r.languages.ReadLanguages(ctx, resumeID)
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
		x, err := r.technologies.ReadTechnologies(ctx, resumeID)
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
		x, err := r.experiences.ReadExperiences(ctx, resumeID)
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
		x, err := r.educations.ReadEducations(ctx, resumeID)
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
		x, err := r.courses.ReadCourses(ctx, resumeID)
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
		x, err := r.activities.ReadActivities(ctx, resumeID)
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
		x, err := r.projects.ReadProjects(ctx, resumeID)
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
