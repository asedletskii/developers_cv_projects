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
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type phase6PostgresFixture struct {
	ctx    context.Context
	dsn    string
	gorm   *gorm.DB
	owner  sharedvalueobjects.ID
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

	owner := sharedvalueobjects.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	resume := newTestResume(t, owner)
	overview := domain.ResumeOverviewParams{Greeting: domain.NewGreeting(domain.GreetingParams{}), Sections: []domain.SectionSettings{}}
	require.NoError(t, resume.ReplaceOverview(overview))
	uow := database.NewUnitOfWork(gormDB)
	resumeRepo := repository.NewResumeRepository(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Create(tx, resume)
	}))

	return &phase6PostgresFixture{ctx: ctx, dsn: dsn, gorm: gormDB, owner: owner, resume: resume, uow: uow}
}

func phase6YearMonth(t *testing.T, value string) domain.YearMonth {
	t.Helper()
	result, err := domain.ParseYearMonth(value)
	require.NoError(t, err)
	return result
}

func phase6Technology(t *testing.T, name string) domain.Technology {
	t.Helper()
	technology, err := domain.NewTechnology(domain.TechnologyParams{Name: name})
	require.NoError(t, err)
	return technology
}

func TestPostgresResumeAggregateOverviewReplacement(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resume := fixture.resume

	primary, err := domain.NewColor("#112233")
	require.NoError(t, err)
	secondary, err := domain.NewColor("#abc")
	require.NoError(t, err)
	theme, err := domain.NewThemeSettings(domain.ThemeSettingsParams{Primary: &primary, Secondary: &secondary})
	require.NoError(t, err)
	availability, err := domain.NewAvailability(domain.AvailabilityParams{Status: domain.AvailabilityStatusOpen})
	require.NoError(t, err)
	section, err := domain.NewSectionSettings(domain.SectionSettingsParams{ID: domain.SectionAbout, Enabled: true, Collapsed: false})
	require.NoError(t, err)
	link, err := domain.NewExternalURL("https://example.com/cv.pdf")
	require.NoError(t, err)
	overview := domain.ResumeOverviewParams{
		Theme:        &theme,
		Greeting:     domain.NewGreeting(domain.GreetingParams{First: "Hello", Second: "I am", SecondHighlight: "Developer", Third: "Welcome"}),
		Availability: &availability,
		ResumeLink:   &link,
		Sections:     []domain.SectionSettings{section},
	}
	require.NoError(t, resume.ReplaceOverview(overview))

	resumeRepo := repository.NewResumeRepository(fixture.gorm)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, resume)
	}))

	loaded, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
	require.NoError(t, err)
	require.NotNil(t, loaded.Overview())
	require.Equal(t, "Hello", loaded.Overview().Greeting().First())
	require.Equal(t, "#112233", loaded.Overview().Theme().Primary.String())
	require.Len(t, loaded.Overview().Sections(), 1)
	require.Equal(t, domain.SectionAbout, loaded.Overview().Sections()[0].ID())
	require.Equal(t, int64(2), loaded.Version())

	stale, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
	require.NoError(t, err)
	require.Equal(t, int64(2), stale.Version())

	updatedSection, err := domain.NewSectionSettings(domain.SectionSettingsParams{ID: domain.SectionProjects, Enabled: true, Collapsed: true})
	require.NoError(t, err)
	updatedAvailability, err := domain.NewAvailability(domain.AvailabilityParams{Status: domain.AvailabilityStatusBusy, Note: stringPtr("Not available")})
	require.NoError(t, err)
	updatedOverview := domain.ResumeOverviewParams{
		Greeting:     domain.NewGreeting(domain.GreetingParams{First: "Updated", Second: "I build", SecondHighlight: "systems", Third: ""}),
		Availability: &updatedAvailability,
		Sections:     []domain.SectionSettings{updatedSection},
	}
	require.NoError(t, loaded.ReplaceOverview(updatedOverview))
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, loaded)
	}))

	loaded, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
	require.NoError(t, err)
	require.Equal(t, int64(3), loaded.Version())
	require.Equal(t, "Updated", loaded.Overview().Greeting().First())
	require.Nil(t, loaded.Overview().Theme())
	require.Nil(t, loaded.Overview().ResumeLink())
	require.Equal(t, domain.AvailabilityStatusBusy, loaded.Overview().Availability().Status())
	require.Len(t, loaded.Overview().Sections(), 1)
	require.Equal(t, domain.SectionProjects, loaded.Overview().Sections()[0].ID())
	require.True(t, loaded.Overview().Sections()[0].Collapsed())

	require.ErrorIs(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, stale)
	}), domain.ErrConcurrencyConflict)
}

func TestPostgresResumeAggregateContactsReplacement(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeRepo := repository.NewResumeRepository(fixture.gorm)

	email, err := domain.NewEmail("first@example.com")
	require.NoError(t, err)
	contactEntity, err := domain.NewContactMethods(domain.ContactMethodsParams{Phone: "+33123456789", Email: &email})
	require.NoError(t, err)
	title := "Contacts"
	contactsText, err := domain.NewContactsText(domain.ContactsTextParams{Title: &title})
	require.NoError(t, err)
	require.NoError(t, fixture.resume.ReplaceContacts(domain.ContactsParams{
		Contacts:       contactEntity,
		HiddenContacts: []domain.ContactKey{domain.ContactKeyEmail},
		Text:           &contactsText,
	}))
	firstID := fixture.resume.Contacts().ID()

	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, fixture.resume)
	}))

	loaded, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: fixture.resume.ID(), UserID: fixture.owner})
	require.NoError(t, err)
	require.Equal(t, int64(2), loaded.Version())
	require.NotNil(t, loaded.Contacts())
	require.Equal(t, firstID, loaded.Contacts().ID())
	require.Equal(t, "+33123456789", *loaded.Contacts().Contacts().Phone())
	require.Equal(t, []domain.ContactKey{domain.ContactKeyEmail}, loaded.Contacts().HiddenContacts())
	require.Equal(t, "Contacts", *loaded.Contacts().Text().Title())

	stale, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: fixture.resume.ID(), UserID: fixture.owner})
	require.NoError(t, err)
	require.Equal(t, int64(2), stale.Version())

	updatedEmail, err := domain.NewEmail("second@example.com")
	require.NoError(t, err)
	updatedEntity, err := domain.NewContactMethods(domain.ContactMethodsParams{Email: &updatedEmail})
	require.NoError(t, err)
	require.NoError(t, loaded.ReplaceContacts(domain.ContactsParams{Contacts: updatedEntity}))
	require.Equal(t, firstID, loaded.Contacts().ID())
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, loaded)
	}))

	loaded, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: fixture.resume.ID(), UserID: fixture.owner})
	require.NoError(t, err)
	require.Equal(t, int64(3), loaded.Version())
	require.Equal(t, firstID, loaded.Contacts().ID())
	require.Nil(t, loaded.Contacts().Contacts().Phone())
	require.Equal(t, "second@example.com", loaded.Contacts().Contacts().Email().String())
	require.Empty(t, loaded.Contacts().HiddenContacts())
	require.Nil(t, loaded.Contacts().Text())

	require.ErrorIs(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, stale)
	}), domain.ErrConcurrencyConflict)
}

func TestPostgresResumeDeleteCascadesAcrossAggregateTables(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeID := fixture.resume.ID()

	conn, err := pgx.Connect(ctx, fixture.dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	profileID := sharedvalueobjects.NewID()
	contactsID := sharedvalueobjects.NewID()
	skillsID := sharedvalueobjects.NewID()
	languageID := sharedvalueobjects.NewID()
	technologyID := sharedvalueobjects.NewID()
	experienceID := sharedvalueobjects.NewID()
	educationID := sharedvalueobjects.NewID()
	courseID := sharedvalueobjects.NewID()
	activityID := sharedvalueobjects.NewID()
	projectID := sharedvalueobjects.NewID()
	contributorID := sharedvalueobjects.NewID()
	now := time.Now().UTC()

	statements := []struct {
		query string
		args  []any
	}{
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
		{`INSERT INTO resume_experience_tasks (experience_id,position,value) VALUES ($1,0,'Implement backend')`, []any{experienceID.UUID()}},
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
		return resumeRepo.Delete(tx, ports.DeleteResumeParams{ResumeID: resumeID, UserID: fixture.owner, Version: fixture.resume.Version()})
	}))

	tables := []string{
		"resumes", "resume_overviews", "resume_overview_sections", "resume_profiles", "resume_profile_work_formats",
		"resume_contacts", "resume_contact_hidden", "resume_contacts_text", "resume_skills", "resume_skill_hard",
		"resume_skill_soft", "resume_languages", "resume_technologies", "resume_experiences", "resume_experience_achievements",
		"resume_experience_tasks", "resume_experience_stack", "resume_educations", "resume_courses", "resume_activities", "resume_projects",
		"resume_project_screenshots", "resume_project_technologies", "resume_project_contributors",
	}
	for _, table := range tables {
		var count int
		require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM "+table).Scan(&count))
		require.Zero(t, count, table)
	}
}

func TestPostgresResumeDeleteRejectsStaleVersion(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeRepo := repository.NewResumeRepository(fixture.gorm)

	err := fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Delete(tx, ports.DeleteResumeParams{
			ResumeID: fixture.resume.ID(),
			UserID:   fixture.owner,
			Version:  fixture.resume.Version() - 1,
		})
	})

	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)

	var count int64
	require.NoError(t, fixture.gorm.WithContext(ctx).Table("resumes").Where("id = ?", fixture.resume.ID().UUID()).Count(&count).Error)
	require.Equal(t, int64(1), count)
}

func TestPostgresPersistenceEnforcesUniquenessAndDatabaseChecks(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeID := fixture.resume.ID()
	owner := fixture.owner
	gormDB := fixture.gorm
	resumeRepo := repository.NewResumeRepository(gormDB)

	secondResume, err := domain.NewResume(domain.ResumeParams{ID: sharedvalueobjects.NewID(), UserID: owner, Title: "Second resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Create(tx, secondResume)
	}))
	duplicateLanguageResume, err := domain.NewResume(domain.ResumeParams{ID: sharedvalueobjects.NewID(), UserID: owner, Title: "Duplicate language", Language: domain.ResumeLanguageRU})
	require.NoError(t, err)
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Create(tx, duplicateLanguageResume)
	}))

	conn, err := pgx.Connect(ctx, fixture.dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	_, err = conn.Exec(ctx, `INSERT INTO resume_languages (id,resume_id,language,position,grade) VALUES (gen_random_uuid(),$1,'en',0,'B2')`, resumeID.UUID())
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_languages (id,resume_id,language,position,grade) VALUES (gen_random_uuid(),$1,'en',1,'C1')`, resumeID.UUID())
	require.Error(t, err, "duplicate resume language must violate UNIQUE(resume_id, language)")

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
			query: `UPDATE resume_overviews SET resume_link = 'https:///missing-host' WHERE resume_id = $1`,
			args:  []any{resumeID.UUID()},
		},
	}
	for _, check := range checks {
		_, err = conn.Exec(ctx, check.query, check.args...)
		require.Error(t, err, check.name)
	}
}

func stringPtr(value string) *string { return &value }

func experienceStackParams(item *domain.Experience) []domain.ExperienceStackParams {
	stack := item.Stack()
	out := make([]domain.ExperienceStackParams, 0, len(stack))
	for _, technology := range stack {
		out = append(out, domain.ExperienceStackParams{Name: technology.Name(), Grade: technology.Grade()})
	}
	return out
}

func TestPostgresExperiencesAggregateReplacement(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resume := fixture.resume

	require.NoError(t, resume.ReplaceExperiences([]domain.ExperienceParams{
		{Company: "Acme", Position: "Go Developer", StartDate: phase6YearMonth(t, "2024-01"), Achievements: []string{"Built API"}, Tasks: []string{"Implement backend"}, Stack: []domain.ExperienceStackParams{{Name: "Go"}}},
		{Company: "Beta", Position: "Staff Engineer", StartDate: phase6YearMonth(t, "2023-01"), Achievements: []string{"Led migration"}, Stack: []domain.ExperienceStackParams{{Name: "PostgreSQL"}}},
	}))
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return repository.NewResumeRepository(fixture.gorm).Save(tx, resume)
	}))

	loaded, err := repository.NewResumeRepository(fixture.gorm).Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
	require.NoError(t, err)
	require.Len(t, loaded.Experiences(), 2)
	firstID := loaded.Experiences()[0].ID()
	secondID := loaded.Experiences()[1].ID()
	require.NotEqual(t, firstID, secondID)
	require.Len(t, loaded.Experiences()[0].Achievements(), 1)
	require.Equal(t, []string{"Implement backend"}, loaded.Experiences()[0].Tasks())
	require.Len(t, loaded.Experiences()[0].Stack(), 1)

	updated := loaded.Experiences()[1]
	require.NoError(t, loaded.ReplaceExperiences([]domain.ExperienceParams{{
		ID: updated.ID(), ResumeID: loaded.ID(), Company: "Beta Updated", Position: "Principal Engineer",
		StartDate: updated.StartDate(), Achievements: []string{"Led migration", "Reduced latency"},
		Tasks: []string{"Lead migration", "Review PRs"},
		Stack: []domain.ExperienceStackParams{{Name: "Go"}, {Name: "PostgreSQL"}},
	}}))
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return repository.NewResumeRepository(fixture.gorm).Save(tx, loaded)
	}))

	loaded, err = repository.NewResumeRepository(fixture.gorm).Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
	require.NoError(t, err)
	require.Len(t, loaded.Experiences(), 1)
	require.Equal(t, secondID, loaded.Experiences()[0].ID())
	require.Equal(t, "Beta Updated", loaded.Experiences()[0].Company())
	require.Len(t, loaded.Experiences()[0].Achievements(), 2)
	require.Equal(t, []string{"Lead migration", "Review PRs"}, loaded.Experiences()[0].Tasks())
	require.Len(t, loaded.Experiences()[0].Stack(), 2)
}

func TestPostgresOrderedAggregateReplacementPreservesIDsOnReorder(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeRepo := repository.NewResumeRepository(fixture.gorm)

	t.Run("experiences", func(t *testing.T) {
		resume := fixture.resume
		require.NoError(t, resume.ReplaceExperiences([]domain.ExperienceParams{
			{Company: "First", Position: "Engineer", StartDate: phase6YearMonth(t, "2024-01")},
			{Company: "Second", Position: "Lead", StartDate: phase6YearMonth(t, "2023-01")},
		}))
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, resume) }))
		loaded, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		firstID, secondID := loaded.Experiences()[0].ID(), loaded.Experiences()[1].ID()

		require.NoError(t, loaded.ReplaceExperiences([]domain.ExperienceParams{
			{ID: secondID, Company: "Second", Position: "Lead", StartDate: phase6YearMonth(t, "2023-01")},
			{ID: firstID, Company: "First", Position: "Engineer", StartDate: phase6YearMonth(t, "2024-01")},
		}))
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, loaded) }))
		loaded, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		require.Equal(t, []sharedvalueobjects.ID{secondID, firstID}, []sharedvalueobjects.ID{loaded.Experiences()[0].ID(), loaded.Experiences()[1].ID()})
	})

	t.Run("educations", func(t *testing.T) {
		resume, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: fixture.resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		resumeVersion := resume.Version()
		require.NoError(t, resume.ReplaceEducations([]domain.EducationParams{
			{Location: "Paris", University: "First University", Grade: domain.EducationLevelBachelor, ProgramName: "Computer Science", Faculty: "Engineering"},
			{Location: "Lyon", University: "Second University", Grade: domain.EducationLevelMaster, ProgramName: "Software Engineering", Faculty: "Engineering"},
		}))
		require.Equal(t, resumeVersion, resume.Version())
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, resume) }))
		loaded, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		firstID, secondID := loaded.Educations()[0].ID(), loaded.Educations()[1].ID()

		require.NoError(t, loaded.ReplaceEducations([]domain.EducationParams{
			{ID: secondID, Location: "Lyon", University: "Second University", Grade: domain.EducationLevelMaster, ProgramName: "Software Engineering", Faculty: "Engineering"},
			{ID: firstID, Location: "Paris", University: "First University", Grade: domain.EducationLevelBachelor, ProgramName: "Computer Science", Faculty: "Engineering"},
		}))
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, loaded) }))
		loaded, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		require.Equal(t, []sharedvalueobjects.ID{secondID, firstID}, []sharedvalueobjects.ID{loaded.Educations()[0].ID(), loaded.Educations()[1].ID()})
	})

	t.Run("courses", func(t *testing.T) {
		resume, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: fixture.resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		require.NoError(t, resume.ReplaceCourses([]domain.CourseParams{
			{Name: "First Course", Organization: "Academy", EndDate: phase6YearMonth(t, "2024-01")},
			{Name: "Second Course", Organization: "Academy", EndDate: phase6YearMonth(t, "2025-01")},
		}))
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, resume) }))
		loaded, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		firstID, secondID := loaded.Courses()[0].ID(), loaded.Courses()[1].ID()

		require.NoError(t, loaded.ReplaceCourses([]domain.CourseParams{
			{ID: secondID, Name: "Second Course", Organization: "Academy", EndDate: phase6YearMonth(t, "2025-01")},
			{ID: firstID, Name: "First Course", Organization: "Academy", EndDate: phase6YearMonth(t, "2024-01")},
		}))
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, loaded) }))
		loaded, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		require.Equal(t, []sharedvalueobjects.ID{secondID, firstID}, []sharedvalueobjects.ID{loaded.Courses()[0].ID(), loaded.Courses()[1].ID()})
	})

	t.Run("activities", func(t *testing.T) {
		resume, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: fixture.resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		require.NoError(t, resume.ReplaceActivities([]domain.ActivityParams{
			{Type: domain.ActivityTypeOpenSource, Title: "First Activity"},
			{Type: domain.ActivityTypeSpeaker, Title: "Second Activity"},
		}))
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, resume) }))
		loaded, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		firstID, secondID := loaded.Activities()[0].ID(), loaded.Activities()[1].ID()

		require.NoError(t, loaded.ReplaceActivities([]domain.ActivityParams{
			{ID: secondID, Type: domain.ActivityTypeSpeaker, Title: "Second Activity"},
			{ID: firstID, Type: domain.ActivityTypeOpenSource, Title: "First Activity"},
		}))
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, loaded) }))
		loaded, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		require.Equal(t, []sharedvalueobjects.ID{secondID, firstID}, []sharedvalueobjects.ID{loaded.Activities()[0].ID(), loaded.Activities()[1].ID()})
	})

	t.Run("technologies", func(t *testing.T) {
		resume, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: fixture.resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		require.NoError(t, resume.ReplaceTechnologies(domain.ResumeTechnologiesParams{Entries: []domain.ResumeTechnologyEntryParams{
			{Name: "Go", Position: 0},
			{Name: "PostgreSQL", Position: 1},
		}}))
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, resume) }))
		loaded, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		entries := loaded.Technologies().Entries()
		firstID, secondID := entries[0].ID(), entries[1].ID()
		collectionID := loaded.Technologies().ID()

		require.NoError(t, loaded.ReplaceTechnologies(domain.ResumeTechnologiesParams{ID: collectionID, Entries: []domain.ResumeTechnologyEntryParams{
			{ID: secondID, Name: "PostgreSQL", Position: 0},
			{ID: firstID, Name: "Go", Position: 1},
		}}))
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, loaded) }))
		loaded, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		entries = loaded.Technologies().Entries()
		require.Equal(t, []sharedvalueobjects.ID{secondID, firstID}, []sharedvalueobjects.ID{entries[0].ID(), entries[1].ID()})
	})

	t.Run("projects", func(t *testing.T) {
		resume, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: fixture.resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		firstScreenshot, err := domain.NewExternalURL("https://example.com/first.png")
		require.NoError(t, err)
		secondScreenshot, err := domain.NewExternalURL("https://example.com/second.png")
		require.NoError(t, err)
		require.NoError(t, resume.ReplaceProjects([]domain.ProjectParams{
			{Name: "First Project", Screenshots: []domain.ExternalURL{firstScreenshot}, Stack: []domain.Technology{phase6Technology(t, "Go")}},
			{Name: "Second Project", Screenshots: []domain.ExternalURL{secondScreenshot}, Stack: []domain.Technology{phase6Technology(t, "PostgreSQL")}},
		}))
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, resume) }))
		loaded, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		firstID, secondID := loaded.Projects()[0].ID(), loaded.Projects()[1].ID()

		require.NoError(t, loaded.ReplaceProjects([]domain.ProjectParams{
			{ID: secondID, Name: "Second Project", Screenshots: loaded.Projects()[1].Screenshots(), Stack: loaded.Projects()[1].Stack()},
			{ID: firstID, Name: "First Project", Screenshots: loaded.Projects()[0].Screenshots(), Stack: loaded.Projects()[0].Stack()},
		}))
		require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error { return resumeRepo.Save(tx, loaded) }))
		loaded, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: fixture.owner})
		require.NoError(t, err)
		require.Equal(t, []sharedvalueobjects.ID{secondID, firstID}, []sharedvalueobjects.ID{loaded.Projects()[0].ID(), loaded.Projects()[1].ID()})
		require.Len(t, loaded.Projects()[0].Screenshots(), 1)
		require.Len(t, loaded.Projects()[0].Stack(), 1)
		require.Equal(t, "PostgreSQL", loaded.Projects()[0].Stack()[0].Name())
	})
}
