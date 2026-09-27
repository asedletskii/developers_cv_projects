//go:build integration

package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type phase6PostgresFixture struct {
	ctx    context.Context
	dsn    string
	gorm   *gorm.DB
	owner  useridentity.ID
	resume *domain.Resume
	uow    *database.UnitOfWork
}

func newPhase6PostgresFixture(t *testing.T) *phase6PostgresFixture {
	t.Helper()
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx := context.Background()
	container, err := testpostgres.Run(
		ctx,
		"postgres:17-alpine",
		testpostgres.WithDatabase("developers_cv"),
		testpostgres.WithUsername("developers_cv"),
		testpostgres.WithPassword("test-password"),
		testpostgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, testcontainers.TerminateContainer(container)) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyMigrations(t, ctx, dsn)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	owner := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	resume := newTestResume(t, owner)
	uow := database.NewUnitOfWork(gormDB)
	resumeRepo := repository.NewResumeRepository(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: resume})
	}))

	return &phase6PostgresFixture{ctx: ctx, dsn: dsn, gorm: gormDB, owner: owner, resume: resume, uow: uow}
}

func phase6YearMonth(t *testing.T, value string) valueobjects.YearMonth {
	t.Helper()
	result, err := valueobjects.ParseYearMonth(value)
	require.NoError(t, err)
	return result
}

func phase6Year(t *testing.T, value string) valueobjects.Year {
	t.Helper()
	result, err := valueobjects.ParseYear(value)
	require.NoError(t, err)
	return result
}

func phase6Technology(t *testing.T, name string) entities.Technology {
	t.Helper()
	technology, err := entities.NewTechnology(entities.TechnologyParams{Name: name})
	require.NoError(t, err)
	return technology
}

func TestPostgresOverviewRepositoryCRUD(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeID := fixture.resume.ID()
	owner := fixture.owner
	overviewRepo := repository.NewOverviewRepository(fixture.gorm)

	primary, err := valueobjects.NewColor("#112233")
	require.NoError(t, err)
	secondary, err := valueobjects.NewColor("#abc")
	require.NoError(t, err)
	theme, err := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{Primary: &primary, Secondary: &secondary})
	require.NoError(t, err)
	availability, err := valueobjects.NewAvailability(valueobjects.AvailabilityParams{Status: enums.AvailabilityStatusOpen})
	require.NoError(t, err)
	section, err := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionAbout, Enabled: true, Collapsed: false})
	require.NoError(t, err)
	link, err := valueobjects.NewExternalURL("https://example.com/cv.pdf")
	require.NoError(t, err)
	overview, err := domain.NewResumeOverview(domain.ResumeOverviewParams{
		Theme:        &theme,
		Greeting:     valueobjects.NewGreeting(valueobjects.GreetingParams{First: "Hello", Second: "I am", SecondHighlight: "Developer", Third: "Welcome"}),
		Availability: &availability,
		ResumeLink:   &link,
		Sections:     []valueobjects.SectionSettings{section},
	})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return overviewRepo.CreateOverview(tx, ports.CreateOverviewRepositoryCommand{ResumeID: resumeID, UserID: owner, Overview: overview})
	}))

	loaded, err := overviewRepo.GetOverview(ctx, ports.GetOverviewRepositoryCommand{ResumeID: resumeID, UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Hello", loaded.Greeting().First())
	require.Equal(t, "#112233", loaded.Theme().Primary.String())
	require.Len(t, loaded.Sections(), 1)
	require.Equal(t, enums.SectionAbout, loaded.Sections()[0].ID())

	updatedSection, err := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionProjects, Enabled: true, Collapsed: true})
	require.NoError(t, err)
	updatedAvailability, err := valueobjects.NewAvailability(valueobjects.AvailabilityParams{Status: enums.AvailabilityStatusBusy, Note: stringPtr("Not available")})
	require.NoError(t, err)
	updatedOverview, err := domain.NewResumeOverview(domain.ResumeOverviewParams{
		Greeting:     valueobjects.NewGreeting(valueobjects.GreetingParams{First: "Updated", Second: "I build", SecondHighlight: "systems", Third: ""}),
		Availability: &updatedAvailability,
		Sections:     []valueobjects.SectionSettings{updatedSection},
	})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return overviewRepo.UpdateOverview(tx, ports.UpdateOverviewRepositoryCommand{ResumeID: resumeID, UserID: owner, Overview: updatedOverview})
	}))

	loaded, err = overviewRepo.GetOverview(ctx, ports.GetOverviewRepositoryCommand{ResumeID: resumeID, UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Updated", loaded.Greeting().First())
	require.Nil(t, loaded.Theme())
	require.Nil(t, loaded.ResumeLink())
	require.Equal(t, enums.AvailabilityStatusBusy, loaded.Availability().Status())
	require.Len(t, loaded.Sections(), 1)
	require.Equal(t, enums.SectionProjects, loaded.Sections()[0].ID())
	require.True(t, loaded.Sections()[0].Collapsed())
}

func TestPostgresNestedRepositoriesCRUD(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeID := fixture.resume.ID()
	owner := fixture.owner

	experienceRepo := repository.NewExperiencesRepository(fixture.gorm)
	experience, err := domain.NewExperience(domain.ExperienceParams{
		ID:           valueobjects.NewID(),
		ResumeID:     resumeID,
		Company:      "Acme",
		Position:     "Go Developer",
		Location:     stringPtr("Paris"),
		StartDate:    phase6YearMonth(t, "2024-01"),
		Description:  stringPtr("Initial description"),
		Achievements: []string{"Built API"},
		Stack:        []entities.Technology{phase6Technology(t, "Go")},
	})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return experienceRepo.CreateExperience(tx, ports.CreateExperienceRepositoryCommand{Experience: experience, ResumeID: resumeID, UserID: owner})
	}))
	loadedExperience, err := experienceRepo.GetExperience(ctx, ports.GetExperienceRepositoryCommand{ResumeID: resumeID, ExperienceID: experience.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Acme", loadedExperience.Company())
	require.Equal(t, []string{"Built API"}, loadedExperience.Achievements())
	require.Len(t, loadedExperience.Stack(), 1)
	updatedExperience, err := domain.NewExperience(domain.ExperienceParams{
		ID:           experience.ID(),
		ResumeID:     resumeID,
		Company:      "Acme Updated",
		Position:     "Staff Go Developer",
		Location:     stringPtr("Berlin"),
		StartDate:    phase6YearMonth(t, "2023-02"),
		Achievements: []string{"Built API", "Led migration"},
		Stack:        []entities.Technology{phase6Technology(t, "Go"), phase6Technology(t, "PostgreSQL")},
	})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return experienceRepo.UpdateExperience(tx, ports.UpdateExperienceRepositoryCommand{Experience: updatedExperience, ResumeID: resumeID, UserID: owner})
	}))
	loadedExperience, err = experienceRepo.GetExperience(ctx, ports.GetExperienceRepositoryCommand{ResumeID: resumeID, ExperienceID: experience.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Acme Updated", loadedExperience.Company())
	require.Len(t, loadedExperience.Achievements(), 2)
	require.Len(t, loadedExperience.Stack(), 2)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return experienceRepo.DeleteExperience(tx, ports.DeleteExperienceRepositoryCommand{ResumeID: resumeID, ExperienceID: experience.ID(), UserID: owner})
	}))
	_, err = experienceRepo.GetExperience(ctx, ports.GetExperienceRepositoryCommand{ResumeID: resumeID, ExperienceID: experience.ID(), UserID: owner})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	educationRepo := repository.NewEducationRepository(fixture.gorm)
	education, err := domain.NewEducation(domain.EducationParams{
		ID:          valueobjects.NewID(),
		ResumeID:    resumeID,
		Location:    "Paris",
		University:  "University",
		Grade:       enums.EducationLevelMaster,
		ProgramName: "Computer Science",
		Faculty:     "Engineering",
		Course:      intPtr(2),
		EndDate:     ptrYear(phase6Year(t, "2025")),
	})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return educationRepo.CreateEducation(tx, ports.CreateEducationRepositoryCommand{Education: education, ResumeID: resumeID, UserID: owner})
	}))
	loadedEducation, err := educationRepo.GetEducation(ctx, ports.GetEducationRepositoryCommand{ResumeID: resumeID, EducationID: education.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "University", loadedEducation.University())
	updatedEducation, err := domain.NewEducation(domain.EducationParams{
		ID:          education.ID(),
		ResumeID:    resumeID,
		Location:    "Berlin",
		University:  "Updated University",
		Grade:       enums.EducationLevelBachelor,
		ProgramName: "Software Engineering",
		Faculty:     "Computer Science",
		EndDate:     ptrYear(phase6Year(t, "2026")),
	})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return educationRepo.UpdateEducation(tx, ports.UpdateEducationRepositoryCommand{Education: updatedEducation, ResumeID: resumeID, UserID: owner})
	}))
	loadedEducation, err = educationRepo.GetEducation(ctx, ports.GetEducationRepositoryCommand{ResumeID: resumeID, EducationID: education.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Updated University", loadedEducation.University())
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return educationRepo.DeleteEducation(tx, ports.DeleteEducationRepositoryCommand{ResumeID: resumeID, EducationID: education.ID(), UserID: owner})
	}))
	_, err = educationRepo.GetEducation(ctx, ports.GetEducationRepositoryCommand{ResumeID: resumeID, EducationID: education.ID(), UserID: owner})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	courseRepo := repository.NewCourseRepository(fixture.gorm)
	course, err := domain.NewCourse(domain.CourseParams{ID: valueobjects.NewID(), ResumeID: resumeID, Name: "Go Architecture", Organization: "Academy", EndDate: phase6YearMonth(t, "2025-01"), Hours: intPtr(40)})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return courseRepo.CreateCourse(tx, ports.CreateCourseRepositoryCommand{Course: course, ResumeID: resumeID, UserID: owner})
	}))
	loadedCourse, err := courseRepo.GetCourse(ctx, ports.GetCourseRepositoryCommand{ResumeID: resumeID, CourseID: course.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Go Architecture", loadedCourse.Name())
	updatedCourse, err := domain.NewCourse(domain.CourseParams{ID: course.ID(), ResumeID: resumeID, Name: "Advanced Go", Organization: "Academy 2", EndDate: phase6YearMonth(t, "2026-01"), Hours: intPtr(60)})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return courseRepo.UpdateCourse(tx, ports.UpdateCourseRepositoryCommand{Course: updatedCourse, ResumeID: resumeID, UserID: owner})
	}))
	loadedCourse, err = courseRepo.GetCourse(ctx, ports.GetCourseRepositoryCommand{ResumeID: resumeID, CourseID: course.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Advanced Go", loadedCourse.Name())
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return courseRepo.DeleteCourse(tx, ports.DeleteCourseRepositoryCommand{ResumeID: resumeID, CourseID: course.ID(), UserID: owner})
	}))
	_, err = courseRepo.GetCourse(ctx, ports.GetCourseRepositoryCommand{ResumeID: resumeID, CourseID: course.ID(), UserID: owner})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	activityRepo := repository.NewActivityRepository(fixture.gorm)
	activity, err := domain.NewActivity(domain.ActivityParams{ID: valueobjects.NewID(), ResumeID: resumeID, Type: enums.ActivityTypeOpenSource, Title: "Open source", Description: stringPtr("Initial")})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return activityRepo.CreateActivity(tx, ports.CreateActivityRepositoryCommand{Activity: activity, ResumeID: resumeID, UserID: owner})
	}))
	loadedActivity, err := activityRepo.GetActivity(ctx, ports.GetActivityRepositoryCommand{ResumeID: resumeID, ActivityID: activity.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Open source", loadedActivity.Title())
	updatedActivity, err := domain.NewActivity(domain.ActivityParams{ID: activity.ID(), ResumeID: resumeID, Type: enums.ActivityTypeMentoring, Title: "Mentoring", Description: stringPtr("Updated")})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return activityRepo.UpdateActivity(tx, ports.UpdateActivityRepositoryCommand{Activity: updatedActivity, ResumeID: resumeID, UserID: owner})
	}))
	loadedActivity, err = activityRepo.GetActivity(ctx, ports.GetActivityRepositoryCommand{ResumeID: resumeID, ActivityID: activity.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Mentoring", loadedActivity.Title())
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return activityRepo.DeleteActivity(tx, ports.DeleteActivityRepositoryCommand{ResumeID: resumeID, ActivityID: activity.ID(), UserID: owner})
	}))
	_, err = activityRepo.GetActivity(ctx, ports.GetActivityRepositoryCommand{ResumeID: resumeID, ActivityID: activity.ID(), UserID: owner})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	projectRepo := repository.NewProjectRepository(fixture.gorm)
	screenshot, err := valueobjects.NewExternalURL("https://example.com/one.png")
	require.NoError(t, err)
	project, err := domain.NewProject(domain.ProjectParams{ID: valueobjects.NewID(), ResumeID: resumeID, Name: "CV", Role: stringPtr("Developer"), Screenshots: []valueobjects.ExternalURL{screenshot}, Stack: []entities.Technology{phase6Technology(t, "Go")}})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return projectRepo.CreateProject(tx, ports.CreateProjectRepositoryCommand{Project: project, ResumeID: resumeID, UserID: owner})
	}))
	loadedProject, err := projectRepo.GetProject(ctx, ports.GetProjectRepositoryCommand{ResumeID: resumeID, ProjectID: project.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "CV", loadedProject.Name())
	require.Len(t, loadedProject.Screenshots(), 1)
	require.Len(t, loadedProject.Stack(), 1)
	updatedProject, err := domain.NewProject(domain.ProjectParams{ID: project.ID(), ResumeID: resumeID, Name: "CV Updated", Role: stringPtr("Architect"), Stack: []entities.Technology{phase6Technology(t, "PostgreSQL")}})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return projectRepo.UpdateProject(tx, ports.UpdateProjectRepositoryCommand{Project: updatedProject, ResumeID: resumeID, UserID: owner})
	}))
	loadedProject, err = projectRepo.GetProject(ctx, ports.GetProjectRepositoryCommand{ResumeID: resumeID, ProjectID: project.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "CV Updated", loadedProject.Name())
	require.Len(t, loadedProject.Screenshots(), 0)
	require.Len(t, loadedProject.Stack(), 1)
	require.Equal(t, "PostgreSQL", loadedProject.Stack()[0].Name())
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return projectRepo.DeleteProject(tx, ports.DeleteProjectRepositoryCommand{ResumeID: resumeID, ProjectID: project.ID(), UserID: owner})
	}))
	_, err = projectRepo.GetProject(ctx, ports.GetProjectRepositoryCommand{ResumeID: resumeID, ProjectID: project.ID(), UserID: owner})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}

func TestPostgresResumeDeleteCascadesAcrossAggregateTables(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeID := fixture.resume.ID()

	conn, err := pgx.Connect(ctx, fixture.dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	profileID := valueobjects.NewID()
	contactsID := valueobjects.NewID()
	skillsID := valueobjects.NewID()
	languageID := valueobjects.NewID()
	technologyID := valueobjects.NewID()
	experienceID := valueobjects.NewID()
	educationID := valueobjects.NewID()
	courseID := valueobjects.NewID()
	activityID := valueobjects.NewID()
	projectID := valueobjects.NewID()
	contributorID := valueobjects.NewID()
	now := time.Now().UTC()

	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO resume_overviews (resume_id,greeting_first,greeting_second,greeting_highlight,greeting_third) VALUES ($1,'a','b','c','d')`, []any{resumeID.UUID()}},
		{`INSERT INTO resume_overview_sections (resume_id,section_id,position,enabled,collapsed) VALUES ($1,'about',0,true,false)`, []any{resumeID.UUID()}},
		{`INSERT INTO resume_profiles (id,resume_id,name,surname,position,photo_path,date_of_birth,location,relocation) VALUES ($1,$2,'A','B','Developer','https://example.com/photo','1998-01-01','Paris',false)`, []any{profileID.UUID(), resumeID.UUID()}},
		{`INSERT INTO resume_profile_work_formats (profile_id,value,position) VALUES ($1,'remote',0)`, []any{profileID.UUID()}},
		{`INSERT INTO resume_contacts (id,resume_id,email) VALUES ($1,$2,'test@example.com')`, []any{contactsID.UUID(), resumeID.UUID()}},
		{`INSERT INTO resume_contact_hidden (contacts_id,value,position) VALUES ($1,'email',0)`, []any{contactsID.UUID()}},
		{`INSERT INTO resume_contacts_text (contacts_id,title) VALUES ($1,'Contact')`, []any{contactsID.UUID()}},
		{`INSERT INTO resume_skills (id,resume_id) VALUES ($1,$2)`, []any{skillsID.UUID(), resumeID.UUID()}},
		{`INSERT INTO resume_skill_hard (skills_id,value,position) VALUES ($1,'architecture',0)`, []any{skillsID.UUID()}},
		{`INSERT INTO resume_skill_soft (skills_id,value,position) VALUES ($1,'communication',0)`, []any{skillsID.UUID()}},
		{`INSERT INTO resume_languages (id,resume_id,language,position,grade) VALUES ($1,$2,'en',0,'B2')`, []any{languageID.UUID(), resumeID.UUID()}},
		{`INSERT INTO resume_technologies (id,resume_id,name,position,grade) VALUES ($1,$2,'Go',0,'senior')`, []any{technologyID.UUID(), resumeID.UUID()}},
		{`INSERT INTO resume_experiences (id,resume_id,sort_order,company,position,start_date) VALUES ($1,$2,0,'Acme','Developer','2024-01-01')`, []any{experienceID.UUID(), resumeID.UUID()}},
		{`INSERT INTO resume_experience_achievements (experience_id,position,value) VALUES ($1,0,'Built API')`, []any{experienceID.UUID()}},
		{`INSERT INTO resume_experience_stack (experience_id,position,name,grade) VALUES ($1,0,'Go','senior')`, []any{experienceID.UUID()}},
		{`INSERT INTO resume_educations (id,resume_id,sort_order,location,university,grade,program_name,faculty) VALUES ($1,$2,0,'Paris','University',5,'Computer Science','Engineering')`, []any{educationID.UUID(), resumeID.UUID()}},
		{`INSERT INTO resume_courses (id,resume_id,sort_order,name,organization,end_date) VALUES ($1,$2,0,'Go Course','Academy','2025-01-01')`, []any{courseID.UUID(), resumeID.UUID()}},
		{`INSERT INTO resume_activities (id,resume_id,sort_order,type,title) VALUES ($1,$2,0,'openSource','Contributor')`, []any{activityID.UUID(), resumeID.UUID()}},
		{`INSERT INTO resume_projects (id,resume_id,sort_order,name) VALUES ($1,$2,0,'CV')`, []any{projectID.UUID(), resumeID.UUID()}},
		{`INSERT INTO resume_project_screenshots (project_id,position,path) VALUES ($1,0,'https://example.com/shot.png')`, []any{projectID.UUID()}},
		{`INSERT INTO resume_project_technologies (project_id,position,name,grade) VALUES ($1,0,'Go','senior')`, []any{projectID.UUID()}},
		{`INSERT INTO resume_project_contributors (id,project_id,position,name,surname,role,email) VALUES ($1,$2,0,'Jane','Doe','Developer','jane@example.com')`, []any{contributorID.UUID(), projectID.UUID()}},
		{`UPDATE resumes SET updated_at = $1 WHERE id = $2`, []any{now, resumeID.UUID()}},
	}
	for _, statement := range statements {
		_, err = conn.Exec(ctx, statement.query, statement.args...)
		require.NoError(t, err)
	}

	resumeRepo := repository.NewResumeRepository(fixture.gorm)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.DeleteResume(tx, ports.DeleteResumeRepositoryCommand{ResumeID: resumeID, UserID: fixture.owner})
	}))

	tables := []string{
		"resumes", "resume_overviews", "resume_overview_sections", "resume_profiles", "resume_profile_work_formats",
		"resume_contacts", "resume_contact_hidden", "resume_contacts_text", "resume_skills", "resume_skill_hard",
		"resume_skill_soft", "resume_languages", "resume_technologies", "resume_experiences", "resume_experience_achievements",
		"resume_experience_stack", "resume_educations", "resume_courses", "resume_activities", "resume_projects",
		"resume_project_screenshots", "resume_project_technologies", "resume_project_contributors",
	}
	for _, table := range tables {
		var count int
		require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM "+table).Scan(&count))
		require.Zero(t, count, table)
	}
}

func TestPostgresPersistenceEnforcesUniquenessAndDatabaseChecks(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeID := fixture.resume.ID()
	owner := fixture.owner
	gormDB := fixture.gorm
	resumeRepo := repository.NewResumeRepository(gormDB)
	languagesRepo := repository.NewLanguagesRepository(gormDB)

	secondResume, err := domain.NewResume(domain.ResumeParams{ID: valueobjects.NewID(), UserID: owner, Title: "Second resume", Language: enums.ResumeLanguageEN})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: secondResume})
	}))
	duplicateLanguageResume, err := domain.NewResume(domain.ResumeParams{ID: valueobjects.NewID(), UserID: owner, Title: "Duplicate language", Language: enums.ResumeLanguageRU})
	require.NoError(t, err)
	err = fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: duplicateLanguageResume})
	})
	require.ErrorIs(t, err, domain.ErrDuplicateEntity)

	language, err := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{ID: valueobjects.NewID(), Language: enums.SpokenLanguageEN, Position: 0})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return languagesRepo.CreateLanguage(tx, ports.CreateLanguageRepositoryCommand{Language: language, ResumeID: resumeID, UserID: owner})
	}))
	duplicateLanguage, err := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{ID: valueobjects.NewID(), Language: enums.SpokenLanguageEN, Position: 1})
	require.NoError(t, err)
	err = fixture.uow.Do(ctx, func(tx context.Context) error {
		return languagesRepo.CreateLanguage(tx, ports.CreateLanguageRepositoryCommand{Language: duplicateLanguage, ResumeID: resumeID, UserID: owner})
	})
	require.ErrorIs(t, err, domain.ErrDuplicateEntity)

	conn, err := pgx.Connect(ctx, fixture.dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	checks := []struct {
		name  string
		query string
		args  []any
	}{
		{
			name:  "negative experience order",
			query: `INSERT INTO resume_experiences (id,resume_id,sort_order,company,position,start_date) VALUES (gen_random_uuid(),$1,-1,'Acme','Developer','2024-01-01')`,
			args:  []any{resumeID.UUID()},
		},
		{
			name:  "invalid activity enum",
			query: `INSERT INTO resume_activities (id,resume_id,sort_order,type,title) VALUES (gen_random_uuid(),$1,0,'unknown','Activity')`,
			args:  []any{resumeID.UUID()},
		},
		{
			name:  "external URL without host",
			query: `INSERT INTO resume_overviews (resume_id,greeting_first,greeting_second,greeting_highlight,greeting_third,resume_link) VALUES ($1,'a','b','c','d','https:///missing-host')`,
			args:  []any{resumeID.UUID()},
		},
	}
	for _, check := range checks {
		_, err = conn.Exec(ctx, check.query, check.args...)
		require.Error(t, err, check.name)
	}
}

func stringPtr(value string) *string                     { return &value }
func intPtr(value int) *int                              { return &value }
func ptrYear(value valueobjects.Year) *valueobjects.Year { return &value }
