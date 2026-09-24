//go:build integration

package persistence_test

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/platform/database"
	"developers-cv/backend/internal/platform/postgres"
)

func TestPostgresPersistenceRoundTrip(t *testing.T) {
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

	t.Cleanup(func() {
		require.NoError(t, testcontainers.TerminateContainer(container))
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	applyMigrations(t, ctx, dsn)

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, conn.Close(ctx))
	})

	require.NoError(t, conn.Ping(ctx))

	var rootCount, legacyCount, presentationCount int
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resumes").Scan(&legacyCount))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_roots").Scan(&rootCount))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_presentations").Scan(&presentationCount))
	require.Equal(t, legacyCount, rootCount)
	require.Equal(t, legacyCount, presentationCount)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, database.Close(gormDB))
	})

	pool, err := postgres.NewPool(
		ctx,
		postgres.Config{
			DSN:             dsn,
			MaxConns:        4,
			MinConns:        1,
			MaxConnLifetime: time.Hour,
			MaxConnIdleTime: time.Minute,
		},
	)
	require.NoError(t, err)

	t.Cleanup(pool.Close)

	repo := repository.NewResumeRepository(gormDB)

	userID := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, userID)

	uow := database.NewUnitOfWork(gormDB)

	resume := newTestResume(t, userID)

	// Phase 1 root vertical slice: the new root repository owns only root + presentation tables.
	presentation := domain.ResumePresentation{
		Greeting:     resume.Greeting(),
		Availability: resume.Availability(),
		ResumeLink:   resume.ResumeLink(),
		Theme:        resume.Theme(),
		Sections:     resume.Sections(),
	}
	root, err := domain.NewResumeRoot(domain.ResumeRootParams{
		ID: resume.ID(), UserID: userID, Title: resume.Title(),
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Presentation: presentation,
	})
	require.NoError(t, err)

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.CreateRoot(txCtx, commands.CreateResumeRootRepositoryCommand{Resume: root})
	}))

	loadedRoot, err := repo.GetRoot(ctx, commands.GetResumeRootRepositoryQuery{UserID: userID, ResumeID: root.ID()})
	require.NoError(t, err)
	require.Equal(t, root.ID(), loadedRoot.ID())
	require.Equal(t, root.Title(), loadedRoot.Title())
	require.Equal(t, root.Presentation().Greeting.First().RU(), loadedRoot.Presentation().Greeting.First().RU())
	require.Len(t, loadedRoot.Presentation().Sections, len(root.Presentation().Sections))

	updatedRoot, err := domain.NewResumeRoot(domain.ResumeRootParams{
		ID: root.ID(), UserID: userID, Title: "Updated root",
		CreatedAt: root.CreatedAt(), UpdatedAt: time.Now().UTC(), Presentation: root.Presentation(),
	})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.UpdateRoot(txCtx, commands.UpdateResumeRootRepositoryCommand{Resume: updatedRoot})
	}))

	loadedRoot, err = repo.GetRoot(ctx, commands.GetResumeRootRepositoryQuery{UserID: userID, ResumeID: root.ID()})
	require.NoError(t, err)
	require.Equal(t, "Updated root", loadedRoot.Title())

	startDate, err := valueobjects.ParseYearMonth("2024-01")
	require.NoError(t, err)
	grade := enums.TechnologyGradeSenior
	experienceID := valueobjects.NewID()
	experienceInput := commands.ExperienceInput{
		ID: experienceID, ResumeID: root.ID(), Company: "Acme", Position: "Developer",
		StartDate: startDate, Achievements: []string{"shipped v2"},
		Technologies: []commands.ExperienceTechnologyInput{{Name: "Go", Grade: &grade}},
	}
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.CreateExperience(txCtx, commands.CreateExperienceCommand{ID: experienceID, ResumeID: root.ID(), UserID: userID, Input: experienceInput})
	}))
	loadedExperiences, err := repo.ListExperiences(ctx, commands.ListExperiencesQuery{UserID: userID, ResumeID: root.ID()})
	require.NoError(t, err)
	require.Len(t, loadedExperiences, 1)

	loadedExperience, err := repo.GetExperience(ctx, commands.GetExperienceQuery{UserID: userID, ResumeID: root.ID(), ExperienceID: experienceID})
	require.NoError(t, err)
	require.Equal(t, "Acme", loadedExperience.Company())
	require.Equal(t, []string{"shipped v2"}, loadedExperience.Achievements())
	require.Equal(t, "Go", loadedExperience.Technologies()[0].Name())
	require.Equal(t, grade, *loadedExperience.Technologies()[0].Grade())

	updatedStartDate, err := valueobjects.ParseYearMonth("2024-02")
	require.NoError(t, err)
	experienceInput.StartDate = updatedStartDate
	experienceInput.Company = "Updated Acme"
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.UpdateExperience(txCtx, commands.UpdateExperienceCommand{ResumeID: root.ID(), ExperienceID: experienceID, UserID: userID, Input: experienceInput})
	}))
	loadedExperience, err = repo.GetExperience(ctx, commands.GetExperienceQuery{UserID: userID, ResumeID: root.ID(), ExperienceID: experienceID})
	require.NoError(t, err)
	require.Equal(t, "Updated Acme", loadedExperience.Company())
	require.Equal(t, "2024-02", loadedExperience.StartDate().String())

	foreignExperienceUserID := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, foreignExperienceUserID)
	_, err = repo.GetExperience(ctx, commands.GetExperienceQuery{UserID: foreignExperienceUserID, ResumeID: root.ID(), ExperienceID: experienceID})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	missingExperienceID := valueobjects.NewID()
	err = uow.Do(ctx, func(txCtx context.Context) error {
		return repo.UpdateExperience(txCtx, commands.UpdateExperienceCommand{ResumeID: root.ID(), ExperienceID: missingExperienceID, UserID: userID, Input: experienceInput})
	})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.DeleteExperience(txCtx, commands.DeleteExperienceCommand{ResumeID: root.ID(), ExperienceID: experienceID, UserID: userID})
	}))
	_, err = repo.GetExperience(ctx, commands.GetExperienceQuery{UserID: userID, ResumeID: root.ID(), ExperienceID: experienceID})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	educationID := valueobjects.NewID()
	educationEnd, err := valueobjects.ParseYear("2023")
	require.NoError(t, err)
	educationInput := commands.EducationInput{ID: educationID, ResumeID: root.ID(), Location: "Paris", University: "University", Grade: enums.EducationLevelMaster, ProgramName: "Computer Science", Faculty: "Engineering", EndDate: &educationEnd}
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.CreateEducation(txCtx, commands.CreateEducationCommand{ID: educationID, ResumeID: root.ID(), UserID: userID, Input: educationInput})
	}))
	loadedEducation, err := repo.GetEducation(ctx, commands.GetEducationQuery{UserID: userID, ResumeID: root.ID(), EducationID: educationID})
	require.NoError(t, err)
	require.Equal(t, "University", loadedEducation.University())
	require.Equal(t, "2023", loadedEducation.EndDate().String())

	updatedEducation := educationInput
	updatedEducation.University = "Updated University"
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.UpdateEducation(txCtx, commands.UpdateEducationCommand{ResumeID: root.ID(), EducationID: educationID, UserID: userID, Input: updatedEducation})
	}))
	loadedEducation, err = repo.GetEducation(ctx, commands.GetEducationQuery{UserID: userID, ResumeID: root.ID(), EducationID: educationID})
	require.NoError(t, err)
	require.Equal(t, "Updated University", loadedEducation.University())

	courseID := valueobjects.NewID()
	courseEnd, err := valueobjects.ParseYearMonth("2024-06")
	require.NoError(t, err)
	hours := 40
	certificate, err := valueobjects.NewExternalURL("https://example.com/certificate")
	require.NoError(t, err)
	courseInput := commands.CourseInput{ID: courseID, ResumeID: root.ID(), Name: "Go Course", Organization: "Acme Academy", EndDate: courseEnd, Hours: &hours, CertificateLink: &certificate}
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.CreateCourse(txCtx, commands.CreateCourseCommand{ID: courseID, ResumeID: root.ID(), UserID: userID, Input: courseInput})
	}))
	loadedCourse, err := repo.GetCourse(ctx, commands.GetCourseQuery{UserID: userID, ResumeID: root.ID(), CourseID: courseID})
	require.NoError(t, err)
	require.Equal(t, "Go Course", loadedCourse.Name())
	require.Equal(t, "2024-06", loadedCourse.EndDate().String())
	require.Equal(t, certificate.String(), loadedCourse.CertificateLink().String())

	activityID := valueobjects.NewID()
	activityDate, err := valueobjects.ParseYearMonth("2025-01")
	require.NoError(t, err)
	activityInput := commands.ActivityInput{ID: activityID, ResumeID: root.ID(), Type: enums.ActivityTypeOpenSource, Title: "Open Source", Date: &activityDate, Link: &certificate}
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.CreateActivity(txCtx, commands.CreateActivityCommand{ID: activityID, ResumeID: root.ID(), UserID: userID, Input: activityInput})
	}))
	loadedActivity, err := repo.GetActivity(ctx, commands.GetActivityQuery{UserID: userID, ResumeID: root.ID(), ActivityID: activityID})
	require.NoError(t, err)
	require.Equal(t, enums.ActivityTypeOpenSource, loadedActivity.Type())
	require.Equal(t, "2025-01", loadedActivity.Date().String())

	foreignActivityUserID := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, foreignActivityUserID)
	_, err = repo.GetActivity(ctx, commands.GetActivityQuery{UserID: foreignActivityUserID, ResumeID: root.ID(), ActivityID: activityID})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	missingCourseID := valueobjects.NewID()
	err = uow.Do(ctx, func(txCtx context.Context) error {
		return repo.UpdateCourse(txCtx, commands.UpdateCourseCommand{ResumeID: root.ID(), CourseID: missingCourseID, UserID: userID, Input: courseInput})
	})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.DeleteEducation(txCtx, commands.DeleteEducationCommand{ResumeID: root.ID(), EducationID: educationID, UserID: userID})
	}))
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.DeleteCourse(txCtx, commands.DeleteCourseCommand{ResumeID: root.ID(), CourseID: courseID, UserID: userID})
	}))
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.DeleteActivity(txCtx, commands.DeleteActivityCommand{ResumeID: root.ID(), ActivityID: activityID, UserID: userID})
	}))

	foreignUserID := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, foreignUserID)
	_, err = repo.GetRoot(ctx, commands.GetResumeRootRepositoryQuery{UserID: foreignUserID, ResumeID: root.ID()})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	missingRoot, err := domain.NewResumeRoot(domain.ResumeRootParams{
		ID: valueobjects.NewID(), UserID: userID, Title: "Missing",
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Presentation: presentation,
	})
	require.NoError(t, err)
	err = uow.Do(ctx, func(txCtx context.Context) error {
		return repo.UpdateRoot(txCtx, commands.UpdateResumeRootRepositoryCommand{Resume: missingRoot})
	})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.DeleteRoot(txCtx, commands.GetResumeRootRepositoryQuery{UserID: userID, ResumeID: root.ID()})
	}))
	_, err = repo.GetRoot(ctx, commands.GetResumeRootRepositoryQuery{UserID: userID, ResumeID: root.ID()})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	require.NoError(
		t,
		uow.Do(
			ctx,
			func(txCtx context.Context) error {
				return repo.Create(txCtx, commands.CreateResumeRepositoryCommand{Resume: resume})
			},
		),
	)

	loaded, err := repo.Get(ctx, commands.GetResumeRepositoryQuery{UserID: userID, ResumeID: resume.ID()})
	require.NoError(t, err)

	require.Equal(t, resume.ID().String(), loaded.ID().String())
	require.Equal(t, resume.Name(), loaded.Name())
	require.Equal(t, resume.Title(), loaded.Title())
	require.Equal(t, userID.UUID(), loaded.UserID().UUID())
	require.Len(t, loaded.Experience(), 1)

	otherUserID := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, otherUserID)

	otherResume := newTestResume(t, otherUserID)

	require.NoError(
		t,
		uow.Do(
			ctx,
			func(txCtx context.Context) error {
				return repo.Create(txCtx, commands.CreateResumeRepositoryCommand{Resume: otherResume})
			},
		),
	)

	_, err = repo.Get(ctx, commands.GetResumeRepositoryQuery{UserID: otherUserID, ResumeID: resume.ID()})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	listItems, err := readers.NewListResumesReader(pool).ListResumes(
		ctx,
		queries.ListResumesQuery{
			UserID: userID,
		},
	)
	require.NoError(t, err)

	require.Len(t, listItems, 1)
	require.Equal(t, resume.ID().String(), listItems[0].ID)
	require.Equal(t, resume.Title(), listItems[0].Title)

	otherListItems, err := readers.NewListResumesReader(pool).ListResumes(
		ctx,
		queries.ListResumesQuery{
			UserID: otherUserID,
		},
	)
	require.NoError(t, err)

	require.Len(t, otherListItems, 1)
	require.Equal(t, otherResume.ID().String(), otherListItems[0].ID)

	readModel, err := readers.NewGetResumeReader(pool).GetResume(
		ctx,
		queries.GetResumeQuery{
			UserID: &userID,
			ID:     resume.ID().UUID(),
			Lang:   queries.LanguageEN,
		},
	)
	require.NoError(t, err)

	require.Equal(t, resume.ID().String(), readModel.ID)
	require.Equal(t, "Hello", readModel.Greeting.First)
	require.Equal(t, "Developer", readModel.Position)
	require.Equal(t, "Developer", readModel.Title)
	require.Equal(t, []string{"remote", "office"}, readModel.WorkFormat)
	require.Len(t, readModel.Experience, 1)
	require.Equal(t, "2024-01", readModel.Experience[0].StartDate)

	_, err = readers.NewGetResumeReader(pool).GetResume(
		ctx,
		queries.GetResumeQuery{
			UserID: &otherUserID,
			ID:     resume.ID().UUID(),
			Lang:   queries.LanguageEN,
		},
	)
	require.ErrorIs(t, err, queries.ErrResumeNotFound)

	publicReadModel, err := readers.NewGetResumeReader(pool).GetResume(
		ctx,
		queries.GetResumeQuery{
			ID:   resume.ID().UUID(),
			Lang: queries.LanguageEN,
		},
	)
	require.NoError(t, err)

	require.Equal(t, resume.ID().String(), publicReadModel.ID)
	require.Equal(t, "Hello", publicReadModel.Greeting.First)

	patchHandler := commands.NewPatchResumeHandler(repo, uow)

	_, err = patchHandler.Handle(
		ctx,
		commands.PatchResumeCommand{
			ResumeID: resume.ID(),
			UserID:   userID,
			Title:    "Go Backend Engineer",
		},
	)
	require.NoError(t, err)

	loadedAfterPatch, err := repo.Get(ctx, commands.GetResumeRepositoryQuery{UserID: userID, ResumeID: resume.ID()})
	require.NoError(t, err)

	require.Equal(t, "Go Backend Engineer", loadedAfterPatch.Title())
	require.Equal(t, resume.Position(), loadedAfterPatch.Position())

	require.Len(t, loadedAfterPatch.Projects(), 1)
	require.Equal(
		t,
		resume.Projects()[0].ID(),
		loadedAfterPatch.Projects()[0].ID(),
	)
	require.Equal(
		t,
		resume.Projects()[0].Name(),
		loadedAfterPatch.Projects()[0].Name(),
	)

	require.Len(t, loadedAfterPatch.Projects()[0].Screenshots(), 1)
	require.Equal(
		t,
		resume.Projects()[0].Screenshots()[0].String(),
		loadedAfterPatch.Projects()[0].Screenshots()[0].String(),
	)

	require.Len(t, loadedAfterPatch.Projects()[0].Stack(), 1)
	require.Equal(
		t,
		resume.Projects()[0].Stack()[0].Name(),
		loadedAfterPatch.Projects()[0].Stack()[0].Name(),
	)

	require.Len(t, loadedAfterPatch.Projects()[0].Contributors(), 1)
	require.Equal(
		t,
		resume.Projects()[0].Contributors()[0].Name(),
		loadedAfterPatch.Projects()[0].Contributors()[0].Name(),
	)

	patched, err := readers.NewGetResumeReader(pool).GetResume(
		ctx,
		queries.GetResumeQuery{
			ID:   resume.ID().UUID(),
			Lang: queries.LanguageEN,
		},
	)
	require.NoError(t, err)

	require.Equal(t, "Go Backend Engineer", patched.Title)
	require.Equal(t, "Developer", patched.Position)
	require.Equal(t, []string{"remote", "office"}, patched.WorkFormat)

	projectID := valueobjects.NewID()
	projectLink, err := valueobjects.NewExternalURL("https://example.com/project")
	require.NoError(t, err)
	projectScreenshot, err := valueobjects.NewAssetPath("/assets/project.png")
	require.NoError(t, err)
	projectGrade := enums.TechnologyGradeSenior
	projectTechnology, err := entities.NewTechnology(entities.TechnologyParams{Name: "Go", Grade: &projectGrade})
	require.NoError(t, err)
	contributorEmail, err := valueobjects.NewEmail("contributor@example.com")
	require.NoError(t, err)
	contributorContacts, err := entities.NewContacts(entities.ContactsParams{Email: &contributorEmail})
	require.NoError(t, err)
	contributor, err := entities.NewContributor(entities.ContributorParams{Name: "Jane", Surname: "Doe", Role: "Developer", Contacts: contributorContacts})
	require.NoError(t, err)
	projectInput := commands.ProjectInput{
		ID: projectID, ResumeID: root.ID(), Name: "Project", Description: stringPtr("Project description"), Role: stringPtr("Lead"),
		Screenshots: []valueobjects.AssetPath{projectScreenshot}, Link: &projectLink,
		Stack: []entities.Technology{projectTechnology}, Contributors: []entities.Contributor{contributor},
	}
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.CreateProject(tx, commands.CreateProjectCommand{ID: projectID, ResumeID: root.ID(), UserID: userID, Input: projectInput})
	}))
	loadedProject, err := repo.GetProject(ctx, commands.GetProjectQuery{UserID: userID, ResumeID: root.ID(), ProjectID: projectID})
	require.NoError(t, err)
	require.Equal(t, "Project", loadedProject.Name())
	require.Len(t, loadedProject.Screenshots(), 1)
	require.Equal(t, "Go", loadedProject.Stack()[0].Name())
	require.Equal(t, "Jane", loadedProject.Contributors()[0].Name())
	require.Equal(t, "contributor@example.com", loadedProject.Contributors()[0].Contacts().Email().String())

	projectInput.Name = "Updated Project"
	projectInput.Screenshots = nil
	projectInput.Contributors = nil
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.UpdateProject(tx, commands.UpdateProjectCommand{ResumeID: root.ID(), ProjectID: projectID, UserID: userID, Input: projectInput})
	}))
	loadedProject, err = repo.GetProject(ctx, commands.GetProjectQuery{UserID: userID, ResumeID: root.ID(), ProjectID: projectID})
	require.NoError(t, err)
	require.Equal(t, "Updated Project", loadedProject.Name())
	require.Empty(t, loadedProject.Screenshots())
	require.Empty(t, loadedProject.Contributors())

	_, err = repo.GetProject(ctx, commands.GetProjectQuery{UserID: foreignExperienceUserID, ResumeID: root.ID(), ProjectID: projectID})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
	missingProjectID := valueobjects.NewID()
	err = uow.Do(ctx, func(tx context.Context) error {
		return repo.UpdateProject(tx, commands.UpdateProjectCommand{ResumeID: root.ID(), ProjectID: missingProjectID, UserID: userID, Input: projectInput})
	})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.DeleteProject(tx, commands.DeleteProjectCommand{ResumeID: root.ID(), ProjectID: projectID, UserID: userID})
	}))
	_, err = repo.GetProject(ctx, commands.GetProjectQuery{UserID: userID, ResumeID: root.ID(), ProjectID: projectID})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

}

func applyMigrations(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()

	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)

	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, conn.Close(ctx))
	})

	migrationsDir := filepath.Join(
		"..",
		"..",
		"..",
		"..",
		"..",
		"db",
		"migrations",
	)

	entries, err := os.ReadDir(migrationsDir)
	require.NoError(t, err)

	migrationFiles := make([]string, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if filepath.Ext(name) != ".sql" {
			continue
		}

		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		migrationFiles = append(migrationFiles, name)
	}

	sort.Strings(migrationFiles)

	require.NotEmpty(t, migrationFiles)

	for _, migrationFile := range migrationFiles {
		content, err := os.ReadFile(
			filepath.Join(migrationsDir, migrationFile),
		)
		require.NoError(t, err)

		_, err = conn.Exec(ctx, string(content))
		require.NoErrorf(
			t,
			err,
			"failed to apply migration %s",
			migrationFile,
		)
	}
}

func insertResumeTestUser(t *testing.T, ctx context.Context, dsn string, userID useridentity.ID) {
	t.Helper()

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, conn.Close(ctx))
	})

	_, err = conn.Exec(ctx, `INSERT INTO users (id) VALUES ($1)`, userID.UUID())
	require.NoError(t, err)
}

func newTestResume(t *testing.T, userID useridentity.ID) *domain.Resume {
	t.Helper()

	id, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)

	dob, err := valueobjects.ParseDate("1998-04-12")
	require.NoError(t, err)

	photo, err := valueobjects.NewAssetPath("/assets/photo.svg")
	require.NoError(t, err)

	first, err := valueobjects.NewLocalizedText(
		valueobjects.LocalizedTextParams{
			RU: "Привет",
			EN: "Hello",
		},
	)
	require.NoError(t, err)

	second, err := valueobjects.NewLocalizedText(
		valueobjects.LocalizedTextParams{
			RU: "Я",
			EN: "I am",
		},
	)
	require.NoError(t, err)

	highlight, err := valueobjects.NewLocalizedText(
		valueobjects.LocalizedTextParams{
			RU: "Тест",
			EN: "Test",
		},
	)
	require.NoError(t, err)

	third, err := valueobjects.NewLocalizedText(
		valueobjects.LocalizedTextParams{
			RU: "Разработчик",
			EN: "Developer",
		},
	)
	require.NoError(t, err)

	greeting := valueobjects.NewGreeting(
		valueobjects.GreetingParams{
			First:           first,
			Second:          second,
			SecondHighlight: highlight,
			Third:           third,
		},
	)

	phone := "+79999999999"

	contacts, err := entities.NewContacts(
		entities.ContactsParams{
			Phone: phone,
		},
	)
	require.NoError(t, err)

	start, err := valueobjects.ParseYearMonth("2024-01")
	require.NoError(t, err)

	experienceID, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)

	experience, err := entities.NewExperience(
		entities.ExperienceParams{
			ID:           experienceID,
			Company:      "Example",
			Position:     "Developer",
			StartDate:    start,
			Achievements: []string{"Shipped feature"},
		},
	)
	require.NoError(t, err)

	projectID, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)

	screenshot, err := valueobjects.NewAssetPath("/assets/project.png")
	require.NoError(t, err)

	projectLink, err := valueobjects.NewExternalURL(
		"https://example.com/project",
	)
	require.NoError(t, err)

	projectTechnology, err := entities.NewTechnology(
		entities.TechnologyParams{
			Name: "Go",
		},
	)
	require.NoError(t, err)

	contributor, err := entities.NewContributor(
		entities.ContributorParams{
			Name:     "John",
			Surname:  "Doe",
			Role:     "Developer",
			Contacts: entities.EmptyContacts(),
		},
	)
	require.NoError(t, err)

	project, err := entities.NewProject(
		entities.ProjectParams{
			ID:           projectID,
			Name:         "CV Platform",
			Description:  stringPtr("Platform project"),
			Role:         stringPtr("Backend"),
			Screenshots:  []valueobjects.AssetPath{screenshot},
			Link:         &projectLink,
			Stack:        []entities.Technology{projectTechnology},
			Contributors: []entities.Contributor{contributor},
		},
	)
	require.NoError(t, err)

	resume, err := domain.NewResume(
		domain.ResumeParams{
			ID:       id,
			UserID:   userID,
			Name:     "Test",
			Surname:  "User",
			Position: "Developer",
			WorkFormats: []enums.WorkFormat{
				enums.WorkFormatRemote,
				enums.WorkFormatOffice,
			},
			Greeting:    greeting,
			PhotoPath:   photo,
			DateOfBirth: dob,
			Location:    "Remote",
			Contacts:    contacts,
			Experience:  []entities.Experience{experience},
			Projects:    []entities.Project{project},
			Sections:    []valueobjects.SectionSettings{},
		},
	)
	require.NoError(t, err)

	return resume
}

func stringPtr(value string) *string {
	return &value
}

func TestPostgresContactsVerticalSlice(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, testcontainers.TerminateContainer(container)) })
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyMigrations(t, ctx, dsn)
	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })
	pool, err := postgres.NewPool(ctx, postgres.Config{DSN: dsn, MaxConns: 4, MinConns: 1, MaxConnLifetime: time.Hour, MaxConnIdleTime: time.Minute})
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	owner := useridentity.NewID()
	foreign := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	insertResumeTestUser(t, ctx, dsn, foreign)
	legacyResume := newTestResume(t, owner)
	repo := repository.NewResumeRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	presentation := domain.ResumePresentation{Greeting: legacyResume.Greeting(), Availability: legacyResume.Availability(), ResumeLink: legacyResume.ResumeLink(), Theme: legacyResume.Theme(), Sections: legacyResume.Sections()}
	root, err := domain.NewResumeRoot(domain.ResumeRootParams{ID: legacyResume.ID(), UserID: owner, Title: legacyResume.Title(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Presentation: presentation})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.CreateRoot(tx, commands.CreateResumeRootRepositoryCommand{Resume: root})
	}))

	email, err := valueobjects.NewEmail("contacts@example.com")
	require.NoError(t, err)
	contactsEntity, err := entities.NewContacts(entities.ContactsParams{Phone: "+79999999999", Email: &email})
	require.NoError(t, err)
	contacts, err := domain.NewContacts(domain.ContactsParams{ID: valueobjects.NewID(), ResumeID: root.ID(), Contacts: contactsEntity, HiddenContacts: []enums.ContactKey{enums.ContactKeyPhone}})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.CreateContacts(tx, commands.CreateResumeContactsRepositoryCommand{Contacts: contacts, UserID: owner})
	}))

	loaded, err := repo.GetContacts(ctx, commands.GetResumeContactsQuery{ResumeID: root.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "+79999999999", *loaded.Contacts().Phone())
	require.Equal(t, "contacts@example.com", loaded.Contacts().Email().String())
	require.Equal(t, []enums.ContactKey{enums.ContactKeyPhone}, loaded.HiddenContacts())

	_, err = repo.GetContacts(ctx, commands.GetResumeContactsQuery{ResumeID: root.ID(), UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	beforeRoot, err := repo.GetRoot(ctx, commands.GetResumeRootRepositoryQuery{ResumeID: root.ID(), UserID: owner})
	require.NoError(t, err)
	updatedEmail, err := valueobjects.NewEmail("updated@example.com")
	require.NoError(t, err)
	updatedEntity, err := entities.NewContacts(entities.ContactsParams{Email: &updatedEmail})
	require.NoError(t, err)
	updatedContacts, err := domain.NewContacts(domain.ContactsParams{ID: contacts.ID(), ResumeID: root.ID(), Contacts: updatedEntity, HiddenContacts: []enums.ContactKey{enums.ContactKeyEmail}})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.UpdateContacts(tx, commands.UpdateResumeContactsRepositoryCommand{Contacts: updatedContacts, UserID: owner})
	}))
	afterRoot, err := repo.GetRoot(ctx, commands.GetResumeRootRepositoryQuery{ResumeID: root.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, beforeRoot.Presentation().Greeting.First().RU(), afterRoot.Presentation().Greeting.First().RU())
	require.Equal(t, beforeRoot.Presentation().Sections, afterRoot.Presentation().Sections)
}

func TestPostgresSkillsLanguagesTechnologiesVerticalSlice(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, testcontainers.TerminateContainer(container)) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyMigrations(t, ctx, dsn)

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	var hardCount, softCount, languageCount, technologyCount int
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_hard_v2").Scan(&hardCount))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_soft_v2").Scan(&softCount))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_languages_v2").Scan(&languageCount))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_technologies_v2").Scan(&technologyCount))
	require.Greater(t, hardCount, 0)
	require.Greater(t, softCount, 0)
	require.Greater(t, languageCount, 0)
	require.Greater(t, technologyCount, 0)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	owner := useridentity.NewID()
	foreign := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	insertResumeTestUser(t, ctx, dsn, foreign)

	resume := newTestResume(t, owner)
	repo := repository.NewResumeRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	presentation := domain.ResumePresentation{
		Greeting:     resume.Greeting(),
		Availability: resume.Availability(),
		ResumeLink:   resume.ResumeLink(),
		Theme:        resume.Theme(),
		Sections:     resume.Sections(),
	}
	root, err := domain.NewResumeRoot(domain.ResumeRootParams{
		ID: resume.ID(), UserID: owner, Title: resume.Title(),
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Presentation: presentation,
	})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.CreateRoot(tx, commands.CreateResumeRootRepositoryCommand{Resume: root})
	}))

	hard := []enums.HardSkill{enums.HardSkillArchitecture, enums.HardSkillTesting}
	soft := []enums.SoftSkill{enums.SoftSkillCommunication, enums.SoftSkillTeamwork}
	skills, err := domain.NewSkills(domain.SkillsParams{ID: valueobjects.NewID(), ResumeID: root.ID(), Hard: hard, Soft: soft})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.CreateSkills(tx, skills, owner)
	}))
	loadedSkills, err := repo.GetSkills(ctx, commands.GetSkillsQuery{ResumeID: root.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, hard, loadedSkills.HardSkills())
	require.Equal(t, soft, loadedSkills.SoftSkills())
	_, err = repo.GetSkills(ctx, commands.GetSkillsQuery{ResumeID: root.ID(), UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	languageID := valueobjects.NewID()
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.CreateLanguage(tx, commands.LanguageInput{ID: languageID, Language: enums.LanguageEN, Grade: languageGradePtr(enums.LanguageGradeB2), Position: 0}, root.ID(), owner)
	}))
	languages, err := repo.ListLanguages(ctx, commands.ListLanguagesQuery{ResumeID: root.ID(), UserID: owner})
	require.NoError(t, err)
	require.Len(t, languages, 1)
	require.Equal(t, languageID, languages[0].ID())
	_, err = repo.GetLanguage(ctx, commands.GetLanguageQuery{ResumeID: root.ID(), LanguageID: languageID, UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	technologyID := valueobjects.NewID()
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.CreateTechnology(tx, commands.TechnologyInput{ID: technologyID, Name: "Go", Grade: technologyGradePtr(enums.TechnologyGradeSenior), Position: 0}, root.ID(), owner)
	}))
	technology, err := repo.GetTechnology(ctx, commands.GetTechnologyQuery{ResumeID: root.ID(), TechnologyID: technologyID, UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Go", technology.Name())
	_, err = repo.GetTechnology(ctx, commands.GetTechnologyQuery{ResumeID: root.ID(), TechnologyID: technologyID, UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}

func languageGradePtr(value enums.LanguageGrade) *enums.LanguageGrade {
	return &value
}

func technologyGradePtr(value enums.TechnologyGrade) *enums.TechnologyGrade {
	return &value
}
