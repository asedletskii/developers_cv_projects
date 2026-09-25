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

	"developers-cv/backend/internal/modules/resume/application/ports"
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

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	pool, err := postgres.NewPool(ctx, postgres.Config{DSN: dsn, MaxConns: 4, MinConns: 1, MaxConnLifetime: time.Hour, MaxConnIdleTime: time.Minute})
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	owner := useridentity.NewID()
	foreign := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	insertResumeTestUser(t, ctx, dsn, foreign)
	resumeRepo := repository.NewResumeRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	contactsRepo := repository.NewContactsRepository(gormDB)
	resume := newTestResume(t, owner)
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: root})
	}))

	var createdAt, updatedAt time.Time
	require.NoError(t, conn.QueryRow(ctx, "SELECT created_at, updated_at FROM resumes WHERE id = $1", resume.ID().UUID()).Scan(&createdAt, &updatedAt))
	require.False(t, createdAt.IsZero())
	require.False(t, updatedAt.IsZero())
	require.WithinDuration(t, createdAt, updatedAt, time.Second)

	email, err := valueobjects.NewEmail("contacts@example.com")
	require.NoError(t, err)
	contactsEntity, err := entities.NewContacts(entities.ContactsParams{Phone: "+79999999999", Email: &email})
	require.NoError(t, err)
	contacts, err := domain.NewContacts(domain.ContactsParams{ID: valueobjects.NewID(), ResumeID: resume.ID(), Contacts: contactsEntity, HiddenContacts: []enums.ContactKey{enums.ContactKeyPhone}})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return contactsRepo.CreateContacts(tx, ports.CreateResumeContactsRepositoryCommand{Contacts: contacts, UserID: owner})
	}))

	loaded, err := contactsRepo.GetContacts(ctx, ports.GetResumeContactsRepositoryCommand{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "+79999999999", *loaded.Contacts().Phone())
	require.Equal(t, "contacts@example.com", loaded.Contacts().Email().String())
	require.Equal(t, []enums.ContactKey{enums.ContactKeyPhone}, loaded.HiddenContacts())

	_, err = contactsRepo.GetContacts(ctx, ports.GetResumeContactsRepositoryCommand{ResumeID: resume.ID(), UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	beforeResume, err := resumeRepo.GetResume(ctx, ports.GetResumeRepositoryCommand{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	updatedEmail, err := valueobjects.NewEmail("updated@example.com")
	require.NoError(t, err)
	updatedEntity, err := entities.NewContacts(entities.ContactsParams{Email: &updatedEmail})
	require.NoError(t, err)
	updatedContacts, err := domain.NewContacts(domain.ContactsParams{ID: contacts.ID(), ResumeID: resume.ID(), Contacts: updatedEntity, HiddenContacts: []enums.ContactKey{enums.ContactKeyEmail}})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return contactsRepo.UpdateContacts(tx, ports.UpdateResumeContactsRepositoryCommand{Contacts: updatedContacts, UserID: owner})
	}))
	afterResume, err := resumeRepo.GetResume(ctx, ports.GetResumeRepositoryCommand{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	loadedUpdated, err := contactsRepo.GetContacts(ctx, ports.GetResumeContactsRepositoryCommand{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "updated@example.com", loadedUpdated.Contacts().Email().String())
	require.Equal(t, []enums.ContactKey{enums.ContactKeyEmail}, loadedUpdated.HiddenContacts())
	require.Equal(t, beforeResume.Title(), afterResume.Title())
	require.Equal(t, beforeResume.Language(), afterResume.Language())
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

	resumeRepo := repository.NewResumeRepository(gormDB)
	languagesReader := readers.NewLanguagesReader(pool)
	uow := database.NewUnitOfWork(gormDB)
	skillsRepo := repository.NewSkillsRepository(gormDB)
	languagesRepo := repository.NewLanguagesRepository(gormDB)
	technologiesRepo := repository.NewTechnologiesRepository(gormDB)
	resume := newTestResume(t, owner)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: root})
	}))

	hard := []enums.HardSkill{enums.HardSkillArchitecture, enums.HardSkillTesting}
	soft := []enums.SoftSkill{enums.SoftSkillCommunication, enums.SoftSkillTeamwork}
	skills, err := domain.NewSkills(domain.SkillsParams{ID: valueobjects.NewID(), ResumeID: resume.ID(), Hard: hard, Soft: soft})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return skillsRepo.CreateSkills(tx, ports.CreateSkillsRepositoryCommand{Skills: skills, UserID: owner})
	}))
	loadedSkills, err := skillsRepo.GetSkills(ctx, ports.GetSkillsRepositoryCommand{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, hard, loadedSkills.HardSkills())
	require.Equal(t, soft, loadedSkills.SoftSkills())
	var hardCount, softCount int
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_hard").Scan(&hardCount))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_soft").Scan(&softCount))
	require.Equal(t, 2, hardCount)
	require.Equal(t, 2, softCount)
	_, err = skillsRepo.GetSkills(ctx, ports.GetSkillsRepositoryCommand{ResumeID: resume.ID(), UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	languageID := valueobjects.NewID()
	language, err := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{ID: languageID, Language: enums.SpokenLanguageEN, Grade: languageGradePtr(enums.LanguageGradeB2), Position: 0})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return languagesRepo.CreateLanguage(tx, ports.CreateLanguageRepositoryCommand{Language: language, ResumeID: resume.ID(), UserID: owner})
	}))
	languages, err := languagesReader.ListLanguages(ctx, queries.ListLanguagesQuery{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.Len(t, languages, 1)
	require.Equal(t, languageID.String(), languages[0].ID)
	var languageCount int
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_languages").Scan(&languageCount))
	require.Equal(t, 1, languageCount)
	_, err = languagesRepo.GetLanguage(ctx, ports.GetLanguageRepositoryCommand{ResumeID: resume.ID(), LanguageID: languageID, UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	technologyID := valueobjects.NewID()
	technology, err := domain.NewResumeTechnologyEntry(domain.ResumeTechnologyEntryParams{ID: technologyID, Name: "Go", Grade: technologyGradePtr(enums.TechnologyGradeSenior), Position: 0})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return technologiesRepo.CreateTechnology(tx, ports.CreateTechnologyRepositoryCommand{Technology: technology, ResumeID: resume.ID(), UserID: owner})
	}))
	loadedTechnology, err := technologiesRepo.GetTechnology(ctx, ports.GetTechnologyRepositoryCommand{ResumeID: resume.ID(), TechnologyID: technologyID, UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Go", loadedTechnology.Name())
	var technologyCount int
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_technologies").Scan(&technologyCount))
	require.Equal(t, 1, technologyCount)
	_, err = technologiesRepo.GetTechnology(ctx, ports.GetTechnologyRepositoryCommand{ResumeID: resume.ID(), TechnologyID: technologyID, UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}

func languageGradePtr(value enums.LanguageGrade) *enums.LanguageGrade {
	return &value
}

func technologyGradePtr(value enums.TechnologyGrade) *enums.TechnologyGrade {
	return &value
}
