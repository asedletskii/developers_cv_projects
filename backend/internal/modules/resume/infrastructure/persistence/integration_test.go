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

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/reader"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
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

	repo := repository.NewPersonRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	person := newTestPerson(t)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Save(txCtx, person) }))

	loaded, err := repo.GetByID(ctx, person.ID())
	require.NoError(t, err)
	require.Equal(t, person.ID().String(), loaded.ID().String())
	require.Equal(t, person.Name(), loaded.Name())
	require.Len(t, loaded.Experience(), 1)

	readModel, err := reader.NewPersonReader(pool).GetPerson(ctx, queries.GetPersonQuery{
		ID:   person.ID().UUID(),
		Lang: queries.LanguageEN,
	})
	require.NoError(t, err)
	require.Equal(t, person.ID().String(), readModel.ID)
	require.Equal(t, "Hello", readModel.Greeting.First)
	require.Equal(t, "Developer", readModel.Position)
	require.Len(t, readModel.Experience, 1)
	require.Equal(t, "2024-01", readModel.Experience[0].StartDate)
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
}

func newTestPerson(t *testing.T) *domain.Person {
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
	person, err := domain.NewPerson(
		domain.PersonParams{
			ID:          id,
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
	return person
}
