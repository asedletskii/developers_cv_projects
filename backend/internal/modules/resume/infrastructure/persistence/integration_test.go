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

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/platform/database"
	"developers-cv/backend/internal/platform/postgres"
)

func newTestRoot(t *testing.T, userID useridentity.ID) *domain.ResumeRoot {
	t.Helper()
	first, err := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: "Привет", EN: "Hello"})
	require.NoError(t, err)
	second, err := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: "Я", EN: "I am"})
	require.NoError(t, err)
	highlight, err := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: "Тест", EN: "Test"})
	require.NoError(t, err)
	third, err := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: "Разработчик", EN: "Developer"})
	require.NoError(t, err)
	greeting := valueobjects.NewGreeting(valueobjects.GreetingParams{First: first, Second: second, SecondHighlight: highlight, Third: third})
	presentation := domain.ResumePresentation{Greeting: greeting, Sections: []valueobjects.SectionSettings{}}
	root, err := domain.NewResumeRoot(domain.ResumeRootParams{
		ID: valueobjects.NewID(), UserID: userID, Title: "Test Resume",
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Presentation: presentation,
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
	rootRepo := repository.NewResumeRootRepository(gormDB)
	contactsRepo := repository.NewContactsRepository(gormDB)
	skillsRepo := repository.NewSkillsRepository(gormDB)
	languagesRepo := repository.NewLanguagesRepository(gormDB)
	technologiesRepo := repository.NewTechnologiesRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	root := newTestRoot(t, owner)
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return rootRepo.CreateRoot(tx, commands.CreateResumeRootRepositoryCommand{Resume: root})
	}))

	email, err := valueobjects.NewEmail("contacts@example.com")
	require.NoError(t, err)
	contactsEntity, err := entities.NewContacts(entities.ContactsParams{Phone: "+79999999999", Email: &email})
	require.NoError(t, err)
	contacts, err := domain.NewContacts(domain.ContactsParams{ID: valueobjects.NewID(), ResumeID: root.ID(), Contacts: contactsEntity, HiddenContacts: []enums.ContactKey{enums.ContactKeyPhone}})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return contactsRepo.CreateContacts(tx, commands.CreateResumeContactsRepositoryCommand{Contacts: contacts, UserID: owner})
	}))

	loaded, err := contactsRepo.GetContacts(ctx, commands.GetResumeContactsQuery{ResumeID: root.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "+79999999999", *loaded.Contacts().Phone())
	require.Equal(t, "contacts@example.com", loaded.Contacts().Email().String())
	require.Equal(t, []enums.ContactKey{enums.ContactKeyPhone}, loaded.HiddenContacts())

	_, err = contactsRepo.GetContacts(ctx, commands.GetResumeContactsQuery{ResumeID: root.ID(), UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	beforeRoot, err := rootRepo.GetRoot(ctx, commands.GetResumeRootRepositoryQuery{ResumeID: root.ID(), UserID: owner})
	require.NoError(t, err)
	updatedEmail, err := valueobjects.NewEmail("updated@example.com")
	require.NoError(t, err)
	updatedEntity, err := entities.NewContacts(entities.ContactsParams{Email: &updatedEmail})
	require.NoError(t, err)
	updatedContacts, err := domain.NewContacts(domain.ContactsParams{ID: contacts.ID(), ResumeID: root.ID(), Contacts: updatedEntity, HiddenContacts: []enums.ContactKey{enums.ContactKeyEmail}})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return contactsRepo.UpdateContacts(tx, commands.UpdateResumeContactsRepositoryCommand{Contacts: updatedContacts, UserID: owner})
	}))
	afterRoot, err := rootRepo.GetRoot(ctx, commands.GetResumeRootRepositoryQuery{ResumeID: root.ID(), UserID: owner})
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
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_hard").Scan(&hardCount))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_soft").Scan(&softCount))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_languages").Scan(&languageCount))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_technologies").Scan(&technologyCount))
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

	rootRepo := repository.NewResumeRootRepository(gormDB)
	contactsRepo := repository.NewContactsRepository(gormDB)
	skillsRepo := repository.NewSkillsRepository(gormDB)
	languagesRepo := repository.NewLanguagesRepository(gormDB)
	technologiesRepo := repository.NewTechnologiesRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	root := newTestRoot(t, owner)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return rootRepo.CreateRoot(tx, commands.CreateResumeRootRepositoryCommand{Resume: root})
	}))

	hard := []enums.HardSkill{enums.HardSkillArchitecture, enums.HardSkillTesting}
	soft := []enums.SoftSkill{enums.SoftSkillCommunication, enums.SoftSkillTeamwork}
	skills, err := domain.NewSkills(domain.SkillsParams{ID: valueobjects.NewID(), ResumeID: root.ID(), Hard: hard, Soft: soft})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return skillsRepo.CreateSkills(tx, skills, owner)
	}))
	loadedSkills, err := skillsRepo.GetSkills(ctx, commands.GetSkillsQuery{ResumeID: root.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, hard, loadedSkills.HardSkills())
	require.Equal(t, soft, loadedSkills.SoftSkills())
	_, err = skillsRepo.GetSkills(ctx, commands.GetSkillsQuery{ResumeID: root.ID(), UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	languageID := valueobjects.NewID()
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return languagesRepo.CreateLanguage(tx, commands.LanguageInput{ID: languageID, Language: enums.LanguageEN, Grade: languageGradePtr(enums.LanguageGradeB2), Position: 0}, root.ID(), owner)
	}))
	languages, err := languagesRepo.ListLanguages(ctx, commands.ListLanguagesQuery{ResumeID: root.ID(), UserID: owner})
	require.NoError(t, err)
	require.Len(t, languages, 1)
	require.Equal(t, languageID, languages[0].ID())
	_, err = languagesRepo.GetLanguage(ctx, commands.GetLanguageQuery{ResumeID: root.ID(), LanguageID: languageID, UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	technologyID := valueobjects.NewID()
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return technologiesRepo.CreateTechnology(tx, commands.TechnologyInput{ID: technologyID, Name: "Go", Grade: technologyGradePtr(enums.TechnologyGradeSenior), Position: 0}, root.ID(), owner)
	}))
	technology, err := technologiesRepo.GetTechnology(ctx, commands.GetTechnologyQuery{ResumeID: root.ID(), TechnologyID: technologyID, UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Go", technology.Name())
	_, err = technologiesRepo.GetTechnology(ctx, commands.GetTechnologyQuery{ResumeID: root.ID(), TechnologyID: technologyID, UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}

func languageGradePtr(value enums.LanguageGrade) *enums.LanguageGrade {
	return &value
}

func technologyGradePtr(value enums.TechnologyGrade) *enums.TechnologyGrade {
	return &value
}
