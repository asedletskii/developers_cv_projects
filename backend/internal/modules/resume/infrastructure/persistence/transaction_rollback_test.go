//go:build integration

package persistence_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/platform/database"
	"developers-cv/backend/internal/platform/postgres"

	"github.com/jackc/pgx/v5"
	testcontainers "github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresMultiTableWritesRollbackAsOneUnit(t *testing.T) {
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
	pool, err := postgres.NewPool(ctx, postgres.Config{DSN: dsn, MaxConns: 4, MinConns: 1})
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	owner := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	resumeRepo := repository.NewResumeRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	resume := newTestResume(t, owner)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: root})
	}))

	photo, err := valueobjects.NewAssetPath("/assets/profile.svg")
	require.NoError(t, err)
	date, err := valueobjects.NewDate(time.Date(1998, time.April, 12, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	profile, err := domain.NewProfile(domain.ProfileParams{ID: valueobjects.NewID(), ResumeID: resume.ID(), Name: "John", Surname: "Doe", Position: "Developer", PhotoPath: photo, DateOfBirth: date, Location: "Paris", WorkFormats: []enums.WorkFormat{enums.WorkFormatRemote}})
	require.NoError(t, err)

	email, err := valueobjects.NewEmail("rollback@example.com")
	require.NoError(t, err)
	contactEntity, err := entities.NewContacts(entities.ContactsParams{Email: &email})
	require.NoError(t, err)
	contacts, err := domain.NewContacts(domain.ContactsParams{ID: valueobjects.NewID(), ResumeID: resume.ID(), Contacts: contactEntity, HiddenContacts: []enums.ContactKey{enums.ContactKeyEmail}})
	require.NoError(t, err)

	skills, err := domain.NewSkills(domain.SkillsParams{ID: valueobjects.NewID(), ResumeID: resume.ID(), Hard: []enums.HardSkill{enums.HardSkillArchitecture, enums.HardSkillTesting}, Soft: []enums.SoftSkill{enums.SoftSkillCommunication}})
	require.NoError(t, err)

	start, err := valueobjects.NewYearMonth(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	experience, err := domain.NewExperience(domain.ExperienceParams{ID: valueobjects.NewID(), ResumeID: resume.ID(), Company: "Example", Position: "Developer", StartDate: start, Achievements: []string{"Built things"}, Technologies: []domain.ExperienceTechnologyParams{{Name: "Go"}}})
	require.NoError(t, err)

	project, err := domain.NewProject(domain.ProjectParams{ID: valueobjects.NewID(), ResumeID: resume.ID(), Name: "Rollback project", Screenshots: []valueobjects.AssetPath{photo}})
	require.NoError(t, err)

	profileRepo := repository.NewProfileRepository(gormDB)
	contactsRepo := repository.NewContactsRepository(gormDB)
	skillsRepo := repository.NewSkillsRepository(gormDB)
	experienceRepo := repository.NewExperiencesRepository(gormDB)
	projectRepo := repository.NewProjectRepository(gormDB)

	rollbackErr := errors.New("force rollback after aggregate writes")
	err = uow.Do(ctx, func(tx context.Context) error {
		if err := profileRepo.CreateProfile(tx, ports.CreateResumeProfileRepositoryCommand{Profile: profile, UserID: owner}); err != nil {
			return err
		}
		if err := contactsRepo.CreateContacts(tx, ports.CreateResumeContactsRepositoryCommand{Contacts: contacts, UserID: owner}); err != nil {
			return err
		}
		if err := skillsRepo.CreateSkills(tx, ports.CreateSkillsRepositoryCommand{Skills: skills, UserID: owner}); err != nil {
			return err
		}
		if err := experienceRepo.CreateExperience(tx, ports.CreateExperienceRepositoryCommand{Experience: experience, ResumeID: resume.ID(), UserID: owner}); err != nil {
			return err
		}
		if err := projectRepo.CreateProject(tx, ports.CreateProjectRepositoryCommand{Project: project, ResumeID: resume.ID(), UserID: owner}); err != nil {
			return err
		}
		return rollbackErr
	})
	require.ErrorIs(t, err, rollbackErr)

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	for _, check := range []struct {
		table string
		id    valueobjects.ID
	}{
		{"resume_profiles", profile.ID()},
		{"resume_contacts", contacts.ID()},
		{"resume_skills", skills.ID()},
		{"resume_experiences", experience.ID()},
		{"resume_projects", project.ID()},
	} {
		var count int
		require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM "+check.table+" WHERE id = $1", check.id.UUID()).Scan(&count))
		require.Zero(t, count, check.table+" must rollback with the surrounding transaction")
	}

	var childCounts struct{ profile, hidden, hard, soft, achievements, technologies, screenshots int }
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_profile_work_formats").Scan(&childCounts.profile))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_contact_hidden").Scan(&childCounts.hidden))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_hard").Scan(&childCounts.hard))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_soft").Scan(&childCounts.soft))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_experience_achievements").Scan(&childCounts.achievements))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_experience_technologies").Scan(&childCounts.technologies))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_project_screenshots").Scan(&childCounts.screenshots))
	require.Zero(t, childCounts.profile)
	require.Zero(t, childCounts.hidden)
	require.Zero(t, childCounts.hard)
	require.Zero(t, childCounts.soft)
	require.Zero(t, childCounts.achievements)
	require.Zero(t, childCounts.technologies)
	require.Zero(t, childCounts.screenshots)
}
