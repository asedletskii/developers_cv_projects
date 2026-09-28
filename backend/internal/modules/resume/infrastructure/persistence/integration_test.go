//go:build integration

package persistence_test

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func languageGradePtr(value enums.LanguageGrade) *enums.LanguageGrade {
	return &value
}

func newTestResume(t *testing.T, userID useridentity.ID) *domain.Resume {
	t.Helper()
	root, err := domain.NewResume(domain.ResumeParams{
		ID: valueobjects.NewID(), UserID: userID, Title: "Test Resume", Language: enums.ResumeLanguageRU,
	})
	require.NoError(t, err)
	return root
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

func TestPostgresSkillsLanguagesAggregateVerticalSlice(t *testing.T) {
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

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	owner := useridentity.NewID()
	foreign := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	insertResumeTestUser(t, ctx, dsn, foreign)

	resumeRepo := repository.NewResumeRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	resume := newTestResume(t, owner)
	require.NoError(t, resume.ReplaceOverview(domain.ResumeOverviewParams{
		Greeting: valueobjects.NewGreeting(valueobjects.GreetingParams{First: "Hello", Second: "I am", SecondHighlight: "Developer", Third: "Welcome"}),
	}))
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Create(tx, resume)
	}))

	hard := []enums.HardSkill{enums.HardSkillArchitecture, enums.HardSkillTesting}
	soft := []enums.SoftSkill{enums.SoftSkillCommunication, enums.SoftSkillTeamwork}
	require.NoError(t, resume.ReplaceSkills(domain.SkillsParams{Hard: hard, Soft: soft}))
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, resume)
	}))
	loadedResume, err := resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.NotNil(t, loadedResume.Skills())
	require.Equal(t, hard, loadedResume.Skills().HardSkills())
	require.Equal(t, soft, loadedResume.Skills().SoftSkills())
	var hardCount, softCount int
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_hard").Scan(&hardCount))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_soft").Scan(&softCount))
	require.Equal(t, 2, hardCount)
	require.Equal(t, 2, softCount)

	require.NoError(t, loadedResume.ReplaceLanguages(domain.ResumeLanguagesParams{Entries: []domain.ResumeLanguageEntryParams{
		{Language: enums.SpokenLanguageEN, Grade: languageGradePtr(enums.LanguageGradeB2), Position: 0},
		{Language: enums.SpokenLanguageFR, Position: 1},
	}}))
	firstLanguageID := loadedResume.Languages().Entries()[0].ID()
	secondLanguageID := loadedResume.Languages().Entries()[1].ID()
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, loadedResume)
	}))
	loadedResume, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.Len(t, loadedResume.Languages().Entries(), 2)
	require.Equal(t, firstLanguageID, loadedResume.Languages().Entries()[0].ID())
	require.Equal(t, secondLanguageID, loadedResume.Languages().Entries()[1].ID())

	require.NoError(t, loadedResume.ReplaceLanguages(domain.ResumeLanguagesParams{Entries: []domain.ResumeLanguageEntryParams{
		{ID: secondLanguageID, Language: enums.SpokenLanguageDE, Position: 0},
		{ID: firstLanguageID, Language: enums.SpokenLanguageEN, Grade: languageGradePtr(enums.LanguageGradeC1), Position: 1},
	}}))
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, loadedResume)
	}))
	loadedResume, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, []enums.SpokenLanguage{enums.SpokenLanguageDE, enums.SpokenLanguageEN}, []enums.SpokenLanguage{loadedResume.Languages().Entries()[0].Language(), loadedResume.Languages().Entries()[1].Language()})

}
