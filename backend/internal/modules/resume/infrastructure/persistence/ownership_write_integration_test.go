//go:build integration

package persistence_test

import (
	"context"
	"testing"

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

func TestPostgresNestedWritesEnforceResumeOwnership(t *testing.T) {
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

	owner := useridentity.NewID()
	foreign := useridentity.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	insertResumeTestUser(t, ctx, dsn, foreign)

	resume := newTestResume(t, owner)
	resumeRepository := repository.NewResumeRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepository.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: resume})
	}))

	startDate, err := valueobjects.ParseYearMonth("2024-01")
	require.NoError(t, err)
	year, err := valueobjects.ParseYear("2025")
	require.NoError(t, err)

	experience, err := domain.NewExperience(domain.ExperienceParams{
		ID:        valueobjects.NewID(),
		ResumeID:  resume.ID(),
		Company:   "Acme",
		Position:  "Go Developer",
		StartDate: startDate,
	})
	require.NoError(t, err)
	experienceRepository := repository.NewExperiencesRepository(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return experienceRepository.CreateExperience(tx, ports.CreateExperienceRepositoryCommand{Experience: experience, ResumeID: resume.ID(), UserID: owner})
	}))
	assertForeignNestedWriteRejected(t, uow, func(tx context.Context) error {
		return experienceRepository.UpdateExperience(tx, ports.UpdateExperienceRepositoryCommand{Experience: experience, ResumeID: resume.ID(), UserID: foreign})
	})
	assertForeignNestedWriteRejected(t, uow, func(tx context.Context) error {
		return experienceRepository.DeleteExperience(tx, ports.DeleteExperienceRepositoryCommand{ResumeID: resume.ID(), ExperienceID: experience.ID(), UserID: foreign})
	})

	education, err := domain.NewEducation(domain.EducationParams{
		ID:          valueobjects.NewID(),
		ResumeID:    resume.ID(),
		Location:    "Berlin",
		University:  "Example University",
		Grade:       enums.EducationLevelMaster,
		ProgramName: "Computer Science",
		Faculty:     "Engineering",
		EndDate:     &year,
	})
	require.NoError(t, err)
	educationRepository := repository.NewEducationRepository(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return educationRepository.CreateEducation(tx, ports.CreateEducationRepositoryCommand{Education: education, ResumeID: resume.ID(), UserID: owner})
	}))
	assertForeignNestedWriteRejected(t, uow, func(tx context.Context) error {
		return educationRepository.UpdateEducation(tx, ports.UpdateEducationRepositoryCommand{Education: education, ResumeID: resume.ID(), UserID: foreign})
	})
	assertForeignNestedWriteRejected(t, uow, func(tx context.Context) error {
		return educationRepository.DeleteEducation(tx, ports.DeleteEducationRepositoryCommand{ResumeID: resume.ID(), EducationID: education.ID(), UserID: foreign})
	})

	course, err := domain.NewCourse(domain.CourseParams{
		ID:           valueobjects.NewID(),
		ResumeID:     resume.ID(),
		Name:         "Go Architecture",
		Organization: "Example Academy",
		EndDate:      startDate,
	})
	require.NoError(t, err)
	courseRepository := repository.NewCourseRepository(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return courseRepository.CreateCourse(tx, ports.CreateCourseRepositoryCommand{Course: course, ResumeID: resume.ID(), UserID: owner})
	}))
	assertForeignNestedWriteRejected(t, uow, func(tx context.Context) error {
		return courseRepository.UpdateCourse(tx, ports.UpdateCourseRepositoryCommand{Course: course, ResumeID: resume.ID(), UserID: foreign})
	})
	assertForeignNestedWriteRejected(t, uow, func(tx context.Context) error {
		return courseRepository.DeleteCourse(tx, ports.DeleteCourseRepositoryCommand{ResumeID: resume.ID(), CourseID: course.ID(), UserID: foreign})
	})

	activity, err := domain.NewActivity(domain.ActivityParams{
		ID:       valueobjects.NewID(),
		ResumeID: resume.ID(),
		Type:     enums.ActivityTypeOpenSource,
		Title:    "Open source",
	})
	require.NoError(t, err)
	activityRepository := repository.NewActivityRepository(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return activityRepository.CreateActivity(tx, ports.CreateActivityRepositoryCommand{Activity: activity, ResumeID: resume.ID(), UserID: owner})
	}))
	assertForeignNestedWriteRejected(t, uow, func(tx context.Context) error {
		return activityRepository.UpdateActivity(tx, ports.UpdateActivityRepositoryCommand{Activity: activity, ResumeID: resume.ID(), UserID: foreign})
	})
	assertForeignNestedWriteRejected(t, uow, func(tx context.Context) error {
		return activityRepository.DeleteActivity(tx, ports.DeleteActivityRepositoryCommand{ResumeID: resume.ID(), ActivityID: activity.ID(), UserID: foreign})
	})

	project, err := domain.NewProject(domain.ProjectParams{
		ID:       valueobjects.NewID(),
		ResumeID: resume.ID(),
		Name:     "Resume Platform",
	})
	require.NoError(t, err)
	projectRepository := repository.NewProjectRepository(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return projectRepository.CreateProject(tx, ports.CreateProjectRepositoryCommand{Project: project, ResumeID: resume.ID(), UserID: owner})
	}))
	assertForeignNestedWriteRejected(t, uow, func(tx context.Context) error {
		return projectRepository.UpdateProject(tx, ports.UpdateProjectRepositoryCommand{Project: project, ResumeID: resume.ID(), UserID: foreign})
	})
	assertForeignNestedWriteRejected(t, uow, func(tx context.Context) error {
		return projectRepository.DeleteProject(tx, ports.DeleteProjectRepositoryCommand{ResumeID: resume.ID(), ProjectID: project.ID(), UserID: foreign})
	})
}

func assertForeignNestedWriteRejected(t *testing.T, uow *database.UnitOfWork, operation func(context.Context) error) {
	t.Helper()
	err := uow.Do(context.Background(), operation)
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}
