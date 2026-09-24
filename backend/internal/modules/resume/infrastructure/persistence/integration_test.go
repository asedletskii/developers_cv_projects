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
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/reader"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
	userdomain "developers-cv/backend/internal/modules/user/domain"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	userreader "developers-cv/backend/internal/modules/user/infrastructure/persistence/reader"
	userrepository "developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"
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
	userRepository := userrepository.NewUserRepository(gormDB)

	userID := useridentity.NewID()
	user := newTestUser(t, userID)

	require.NoError(t, userRepository.Save(ctx, user))

	profileConn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, profileConn.Close(ctx)) })
	_, err = profileConn.Exec(ctx, `
		INSERT INTO user_profiles (user_id, name, surname, position, location, photo_path)
		VALUES ($1, 'Profile Name', 'Profile Surname', 'Profile Position', 'Profile Location', '/profiles/profile.svg')
	`, userID.UUID())
	require.NoError(t, err)

	uow := database.NewUnitOfWork(gormDB)

	resume := newTestResume(t, userID)

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

	profile, err := userreader.NewGetProfileReader(pool).GetProfile(
		ctx,
		userqueries.GetProfileQuery{UserID: userID},
	)
	require.NoError(t, err)
	require.Equal(t, "Profile Name", profile.Name)
	require.Equal(t, "Profile Surname", profile.Surname)
	require.Equal(t, "Profile Position", *profile.Position)
	require.NotEqual(t, loaded.Name(), profile.Name)
	require.NotEqual(t, loaded.Surname(), profile.Surname)

	otherUserID := useridentity.NewID()
	otherUser := newTestUser(t, otherUserID)

	require.NoError(t, userRepository.Save(ctx, otherUser))

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

	listItems, err := reader.NewListResumesReader(pool).ListResumes(
		ctx,
		queries.ListResumesQuery{
			UserID: userID,
		},
	)
	require.NoError(t, err)

	require.Len(t, listItems, 1)
	require.Equal(t, resume.ID().String(), listItems[0].ID)
	require.Equal(t, resume.Title(), listItems[0].Title)

	otherListItems, err := reader.NewListResumesReader(pool).ListResumes(
		ctx,
		queries.ListResumesQuery{
			UserID: otherUserID,
		},
	)
	require.NoError(t, err)

	require.Len(t, otherListItems, 1)
	require.Equal(t, otherResume.ID().String(), otherListItems[0].ID)

	readModel, err := reader.NewGetResumeReader(pool).GetResume(
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

	_, err = reader.NewGetResumeReader(pool).GetResume(
		ctx,
		queries.GetResumeQuery{
			UserID: &otherUserID,
			ID:     resume.ID().UUID(),
			Lang:   queries.LanguageEN,
		},
	)
	require.ErrorIs(t, err, queries.ErrResumeNotFound)

	publicReadModel, err := reader.NewGetResumeReader(pool).GetResume(
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

	patched, err := reader.NewGetResumeReader(pool).GetResume(
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

func newTestUser(t *testing.T, userID useridentity.ID) *userdomain.User {
	t.Helper()

	user, err := userdomain.NewUser(
		userdomain.UserParams{
			ID: userID,
		},
	)
	require.NoError(t, err)

	return user
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
