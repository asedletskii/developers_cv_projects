//go:build integration

package persistence_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/reader"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	userdomain "developers-cv/backend/internal/modules/user/domain"
	userrepository "developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	"developers-cv/backend/internal/platform/postgres"
)

func TestPostgresPersistenceRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine",
		testpostgres.WithDatabase("developers_cv"),
		testpostgres.WithUsername("developers_cv"),
		testpostgres.WithPassword("test-password"),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, container.Terminate(ctx)) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyMigrations(t, ctx, dsn)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	pool, err := postgres.NewPool(ctx,
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
	uow := database.NewUnitOfWork(gormDB)
	resume := newTestResume(t, userID)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Save(txCtx, resume) }))

	loaded, err := repo.GetByID(ctx, userID, resume.ID())
	require.NoError(t, err)
	require.Equal(t, resume.ID().String(), loaded.ID().String())
	require.Equal(t, resume.Name(), loaded.Name())
	require.Equal(t, resume.Title(), loaded.Title())
	require.Equal(t, userID.UUID(), loaded.UserID().UUID())
	require.Len(t, loaded.Experience(), 1)

	otherUserID := useridentity.NewID()
	otherUser := newTestUser(t, otherUserID)
	require.NoError(t, userRepository.Save(ctx, otherUser))
	otherResume := newTestResume(t, otherUserID)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Save(txCtx, otherResume) }))
	_, err = repo.GetByID(ctx, otherUserID, resume.ID())
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	listItems, err := reader.NewListResumesReader(pool).ListResumes(ctx, queries.GetListResumesQuery{UserID: userID})
	require.NoError(t, err)
	require.Len(t, listItems, 1)
	require.Equal(t, resume.ID().String(), listItems[0].ID)
	require.Equal(t, resume.Title(), listItems[0].Title)
	otherListItems, err := reader.NewListResumesReader(pool).ListResumes(ctx, queries.GetListResumesQuery{UserID: otherUserID})
	require.NoError(t, err)
	require.Len(t, otherListItems, 1)
	require.Equal(t, otherResume.ID().String(), otherListItems[0].ID)

	readModel, err := reader.NewGetResumeReader(pool).GetResume(ctx, queries.GetResumeQuery{
		UserID: &userID,
		ID:     resume.ID().UUID(),
		Lang:   queries.LanguageEN,
	})
	require.NoError(t, err)
	require.Equal(t, resume.ID().String(), readModel.ID)
	require.Equal(t, "Hello", readModel.Greeting.First)
	require.Equal(t, "Developer", readModel.Position)
	require.Equal(t, "Developer", readModel.Title)
	require.Len(t, readModel.Experience, 1)
	require.Equal(t, "2024-01", readModel.Experience[0].StartDate)

	_, err = reader.NewGetResumeReader(pool).GetResume(ctx, queries.GetResumeQuery{
		UserID: &otherUserID,
		ID:     resume.ID().UUID(),
		Lang:   queries.LanguageEN,
	})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
	publicReadModel, err := reader.NewGetResumeReader(pool).GetResume(ctx, queries.GetResumeQuery{
		ID:   resume.ID().UUID(),
		Lang: queries.LanguageEN,
	})
	require.NoError(t, err)
	require.Equal(t, resume.ID().String(), publicReadModel.ID)
	require.Equal(t, "Hello", publicReadModel.Greeting.First)

	require.NoError(t, resume.Rename("Go Backend Engineer"))
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Save(txCtx, resume) }))
	patched, err := reader.NewGetResumeReader(pool).GetResume(ctx, queries.GetResumeQuery{
		ID:   resume.ID().UUID(),
		Lang: queries.LanguageEN,
	})
	require.NoError(t, err)
	require.Equal(t, "Go Backend Engineer", patched.Title)
	require.Equal(t, "Developer", patched.Position)
}

func applyMigrations(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	content, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", "000001_resume_core.up.sql"))
	require.NoError(t, err)
	_, err = conn.Exec(ctx, string(content))
	require.NoError(t, err)

	content, err = os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", "000002_resume_settings.up.sql"))
	require.NoError(t, err)
	_, err = conn.Exec(ctx, string(content))
	require.NoError(t, err)

	content, err = os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", "000003_resume_content.up.sql"))
	require.NoError(t, err)
	_, err = conn.Exec(ctx, string(content))
	require.NoError(t, err)

	content, err = os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", "000004_resume_projects.up.sql"))
	require.NoError(t, err)
	_, err = conn.Exec(ctx, string(content))
	require.NoError(t, err)

	content, err = os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", "000005_resume_constraints.up.sql"))
	require.NoError(t, err)
	_, err = conn.Exec(ctx, string(content))
	require.NoError(t, err)

	content, err = os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", "000006_seed_demo_people.up.sql"))
	require.NoError(t, err)
	_, err = conn.Exec(ctx, string(content))
	require.NoError(t, err)

	content, err = os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", "000007_resume_terminology.up.sql"))
	require.NoError(t, err)
	_, err = conn.Exec(ctx, string(content))
	require.NoError(t, err)

	for _, migration := range []string{"000008_users.up.sql", "000009_resume_users.up.sql", "000010_resume_title.up.sql"} {
		content, err = os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", migration))
		require.NoError(t, err)
		_, err = conn.Exec(ctx, string(content))
		require.NoError(t, err)
	}
}

func newTestUser(t *testing.T, userID useridentity.ID) *userdomain.User {
	t.Helper()
	user, err := userdomain.NewUser(userdomain.UserParams{ID: userID})
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
		valueobjects.LocalizedTextParams{RU: "Привет", EN: "Hello"},
	)
	require.NoError(t, err)
	second, err := valueobjects.NewLocalizedText(
		valueobjects.LocalizedTextParams{RU: "Я", EN: "I am"},
	)
	require.NoError(t, err)
	highlight, err := valueobjects.NewLocalizedText(
		valueobjects.LocalizedTextParams{RU: "Тест", EN: "Test"},
	)
	require.NoError(t, err)
	third, err := valueobjects.NewLocalizedText(
		valueobjects.LocalizedTextParams{RU: "Разработчик", EN: "Developer"},
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
	contacts, err := entities.NewContacts(entities.ContactsParams{Phone: phone})
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
	resume, err := domain.NewResume(
		domain.ResumeParams{
			ID:          id,
			UserID:      userID,
			Name:        "Test",
			Surname:     "User",
			Position:    "Developer",
			Greeting:    greeting,
			PhotoPath:   photo,
			DateOfBirth: dob,
			Location:    "Remote",
			Contacts:    contacts,
			Experience:  []entities.Experience{experience},
			Sections:    []valueobjects.SectionSettings{},
		},
	)
	require.NoError(t, err)
	return resume
}
