//go:build integration

package persistence_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
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
	"developers-cv/backend/internal/platform/database"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
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
		return resumeRepo.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: resume})
	}))

	var createdAt, updatedAt time.Time
	require.NoError(t, conn.QueryRow(ctx, "SELECT created_at, updated_at FROM resumes WHERE id = $1", resume.ID().UUID()).Scan(&createdAt, &updatedAt))
	require.False(t, createdAt.IsZero())
	require.False(t, updatedAt.IsZero())
	require.WithinDuration(t, createdAt, updatedAt, time.Second)

	beforeReplace, err := resumeRepo.GetResume(ctx, ports.GetResumeRepositoryCommand{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	beforeReplaceUpdatedAt := updatedAt
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.ReplaceResumeTitle(tx, ports.ReplaceResumeTitleRepositoryCommand{
			ResumeID: resume.ID(),
			UserID:   owner,
			Title:    "Updated Resume Title",
		})
	}))
	afterReplace, err := resumeRepo.GetResume(ctx, ports.GetResumeRepositoryCommand{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Updated Resume Title", afterReplace.Title())
	require.Equal(t, beforeReplace.Language(), afterReplace.Language())
	require.NoError(t, conn.QueryRow(ctx, "SELECT updated_at FROM resumes WHERE id = $1", resume.ID().UUID()).Scan(&updatedAt))
	require.True(t, updatedAt.After(beforeReplaceUpdatedAt))

	_, err = resumeRepo.GetResume(ctx, ports.GetResumeRepositoryCommand{ResumeID: resume.ID(), UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
	err = uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.ReplaceResumeTitle(tx, ports.ReplaceResumeTitleRepositoryCommand{
			ResumeID: resume.ID(),
			UserID:   foreign,
			Title:    "Foreign Update",
		})
	})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
	unchanged, err := resumeRepo.GetResume(ctx, ports.GetResumeRepositoryCommand{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, "Updated Resume Title", unchanged.Title())

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

	secondResume := newTestResume(t, foreign)
	secondResumeID := secondResume.ID()
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: secondResume})
	}))
	_, err = conn.Exec(ctx, `INSERT INTO resume_contacts (id,resume_id,phone,telegram,email) VALUES ($1,$2,'   ',NULL,NULL)`, valueobjects.NewID().UUID(), secondResumeID.UUID())
	require.Error(t, err)
	contactsID := valueobjects.NewID()
	_, err = conn.Exec(ctx, `INSERT INTO resume_contacts (id,resume_id,email) VALUES ($1,$2,'second@example.com')`, contactsID.UUID(), secondResumeID.UUID())
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_contacts_text (contacts_id,title,subtitle) VALUES ($1,'   ',NULL)`, contactsID.UUID())
	require.Error(t, err)
}

func TestPostgresOrderedCreatesSerializeOnResumeLock(t *testing.T) {
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

	owner := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	resume := newTestResume(t, owner)
	uow := database.NewUnitOfWork(gormDB)
	resumeRepository := repository.NewResumeRepository(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepository.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: resume})
	}))

	activityRepository := repository.NewActivityRepository(gormDB)
	activities := make([]*domain.Activity, 2)
	for i := range activities {
		activity, createErr := domain.NewActivity(domain.ActivityParams{
			ID:       valueobjects.NewID(),
			ResumeID: resume.ID(),
			Type:     enums.ActivityTypeOpenSource,
			Title:    fmt.Sprintf("Activity %d", i),
		})
		require.NoError(t, createErr)
		activities[i] = activity
	}

	start := make(chan struct{})
	errs := make(chan error, len(activities))
	var wg sync.WaitGroup
	for _, activity := range activities {
		activity := activity
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			errs <- uow.Do(ctx, func(tx context.Context) error {
				return activityRepository.CreateActivity(tx, ports.CreateActivityRepositoryCommand{Activity: activity, ResumeID: resume.ID(), UserID: owner})
			})
		}()
	}
	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}

	var orders []int
	require.NoError(t, gormDB.Raw("SELECT sort_order FROM resume_activities WHERE resume_id = ? ORDER BY sort_order", resume.ID().UUID()).Scan(&orders).Error)
	require.Equal(t, []int{0, 1}, orders)
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

	owner := useridentity.NewID()
	foreign := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	insertResumeTestUser(t, ctx, dsn, foreign)

	resumeRepo := repository.NewResumeRepository(gormDB)
	languagesReader := readers.NewLanguagesReader(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	skillsRepo := repository.NewSkillsRepository(gormDB)
	languagesRepo := repository.NewLanguagesRepository(gormDB)
	technologiesRepo := repository.NewTechnologiesRepository(gormDB)
	resume := newTestResume(t, owner)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: resume})
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

func TestPostgresOrderedCollectionCreateAndReorder(t *testing.T) {
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

	owner := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	resumeRepo := repository.NewResumeRepository(gormDB)
	experiencesRepo := repository.NewExperiencesRepository(gormDB)
	languagesRepo := repository.NewLanguagesRepository(gormDB)
	technologiesRepo := repository.NewTechnologiesRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	resume := newTestResume(t, owner)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: resume})
	}))

	const experienceCount = 8
	errs := make(chan error, experienceCount)
	for i := 0; i < experienceCount; i++ {
		go func(index int) {
			itemID := valueobjects.NewID()
			startDate, parseErr := valueobjects.ParseYearMonth(fmt.Sprintf("202%d-01", index+1))
			if parseErr != nil {
				errs <- parseErr
				return
			}
			item, createErr := domain.NewExperience(domain.ExperienceParams{
				ID: itemID, ResumeID: resume.ID(), Company: "Acme", Position: "Developer", StartDate: startDate,
			})
			if createErr != nil {
				errs <- createErr
				return
			}
			errs <- uow.Do(ctx, func(tx context.Context) error {
				return experiencesRepo.CreateExperience(tx, ports.CreateExperienceRepositoryCommand{Experience: item, ResumeID: resume.ID(), UserID: owner})
			})
		}(i)
	}
	for i := 0; i < experienceCount; i++ {
		require.NoError(t, <-errs)
	}

	var orders []int
	require.NoError(t, gormDB.Raw("SELECT sort_order FROM resume_experiences WHERE resume_id = ? ORDER BY sort_order", resume.ID().UUID()).Scan(&orders).Error)
	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7}, orders)

	languageIDs := make([]valueobjects.ID, 3)
	for i, language := range []enums.SpokenLanguage{enums.SpokenLanguageEN, enums.SpokenLanguageFR, enums.SpokenLanguageDE} {
		languageIDs[i] = valueobjects.NewID()
		item, createErr := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{ID: languageIDs[i], Language: language, Position: 0})
		require.NoError(t, createErr)
		require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
			return languagesRepo.CreateLanguage(tx, ports.CreateLanguageRepositoryCommand{Language: item, ResumeID: resume.ID(), UserID: owner})
		}))
	}

	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return languagesRepo.ReorderLanguages(tx, ports.ReorderLanguagesRepositoryCommand{ResumeID: resume.ID(), UserID: owner, LanguageIDs: []valueobjects.ID{languageIDs[2], languageIDs[0], languageIDs[1]}})
	}))

	var languagePositions []int
	require.NoError(t, gormDB.Raw("SELECT position FROM resume_languages WHERE resume_id = ? ORDER BY position", resume.ID().UUID()).Scan(&languagePositions).Error)
	require.Equal(t, []int{0, 1, 2}, languagePositions)
	var orderedLanguageIDs []string
	require.NoError(t, gormDB.Raw("SELECT id::text FROM resume_languages WHERE resume_id = ? ORDER BY position", resume.ID().UUID()).Scan(&orderedLanguageIDs).Error)
	require.Equal(t, []string{languageIDs[2].String(), languageIDs[0].String(), languageIDs[1].String()}, orderedLanguageIDs)

	technologyIDs := make([]valueobjects.ID, 3)
	for i, name := range []string{"Go", "PostgreSQL", "Docker"} {
		technologyIDs[i] = valueobjects.NewID()
		item, createErr := domain.NewResumeTechnologyEntry(domain.ResumeTechnologyEntryParams{ID: technologyIDs[i], Name: name, Position: 0})
		require.NoError(t, createErr)
		require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
			return technologiesRepo.CreateTechnology(tx, ports.CreateTechnologyRepositoryCommand{Technology: item, ResumeID: resume.ID(), UserID: owner})
		}))
	}

	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return technologiesRepo.ReorderTechnologies(tx, ports.ReorderTechnologiesRepositoryCommand{ResumeID: resume.ID(), UserID: owner, TechnologyIDs: []valueobjects.ID{technologyIDs[2], technologyIDs[0], technologyIDs[1]}})
	}))
	var orderedTechnologyIDs []string
	require.NoError(t, gormDB.Raw("SELECT id::text FROM resume_technologies WHERE resume_id = ? ORDER BY position", resume.ID().UUID()).Scan(&orderedTechnologyIDs).Error)
	require.Equal(t, []string{technologyIDs[2].String(), technologyIDs[0].String(), technologyIDs[1].String()}, orderedTechnologyIDs)
}

func languageGradePtr(value enums.LanguageGrade) *enums.LanguageGrade {
	return &value
}

func technologyGradePtr(value enums.TechnologyGrade) *enums.TechnologyGrade {
	return &value
}
