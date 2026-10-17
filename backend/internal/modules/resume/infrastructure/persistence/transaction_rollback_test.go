//go:build integration

package persistence_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/pgtest"
)

func TestPostgresMultiTableWritesRollbackAsOneUnit(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	owner := sharedvalueobjects.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	resumeRepo := repository.NewResumeRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)

	resume := newTestResume(t, owner)
	require.NoError(t, resume.ReplaceOverview(domain.ResumeOverviewParams{
		Greeting: mustGreeting(t, domain.GreetingParams{First: "Hello", Second: "I am", SecondHighlight: "Developer", Third: "Welcome"}),
	}))
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Create(tx, resume)
	}))

	photo, err := domain.NewExternalURL("https://example.com/assets/profile.svg")
	require.NoError(t, err)
	date, err := domain.NewDate(time.Date(1998, time.April, 12, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	start, err := domain.NewYearMonth(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	email, err := domain.NewEmail("rollback@example.com")
	require.NoError(t, err)
	contactEntity, err := domain.NewContactMethods(domain.ContactMethodsParams{Email: &email})
	require.NoError(t, err)

	rollbackErr := errors.New("force rollback after aggregate writes")
	err = uow.Do(ctx, func(tx context.Context) error {
		if err := resume.ReplaceProfile(domain.ProfileParams{
			Name:        "John",
			Surname:     "Doe",
			Position:    "Developer",
			PhotoURL:    photo,
			DateOfBirth: date,
			Location:    "Paris",
			WorkFormats: []domain.WorkFormat{domain.WorkFormatRemote},
		}); err != nil {
			return err
		}
		if err := resume.ReplaceSkills(domain.SkillsParams{
			Hard: []string{"architecture", "testing"},
			Soft: []string{"communication"},
		}); err != nil {
			return err
		}
		if err := resume.ReplaceExperiences([]domain.ExperienceParams{{
			Company:      "Example",
			Position:     "Developer",
			StartDate:    start,
			Achievements: []string{"Built things"},
			Stack:        []domain.ExperienceStackParams{{Name: "Go"}},
		}}); err != nil {
			return err
		}
		if err := resume.ReplaceProjects([]domain.ProjectParams{{
			Name:        "Rollback project",
			Screenshots: []domain.ExternalURL{photo},
		}}); err != nil {
			return err
		}
		if err := resume.ReplaceContacts(domain.ContactsParams{
			Contacts:       contactEntity,
			HiddenContacts: []domain.ContactKey{domain.ContactKeyEmail},
		}); err != nil {
			return err
		}

		if err := resumeRepo.Save(tx, resume); err != nil {
			return err
		}

		reloaded, err := resumeRepo.Get(tx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: owner})
		if err != nil {
			return err
		}
		if reloaded.Profile() == nil || reloaded.Skills() == nil || len(reloaded.Experiences()) != 1 || len(reloaded.Projects()) != 1 || reloaded.Contacts() == nil {
			return errors.New("aggregate children were not persisted inside transaction")
		}

		updatedProfile := reloaded.Profile()
		if err := reloaded.ReplaceProfile(domain.ProfileParams{
			ID:          updatedProfile.ID(),
			Name:        "Jane",
			Surname:     updatedProfile.Surname(),
			Patronymic:  updatedProfile.Patronymic(),
			Position:    "Senior Developer",
			PhotoURL:    updatedProfile.PhotoURL(),
			DateOfBirth: updatedProfile.DateOfBirth(),
			Location:    updatedProfile.Location(),
			Relocation:  updatedProfile.Relocation(),
			WorkFormats: []domain.WorkFormat{domain.WorkFormatOffice, domain.WorkFormatHybrid},
		}); err != nil {
			return err
		}
		if err := reloaded.ReplaceContacts(domain.ContactsParams{
			ID:             reloaded.Contacts().ID(),
			Contacts:       contactEntity,
			HiddenContacts: []domain.ContactKey{domain.ContactKeyEmail},
		}); err != nil {
			return err
		}
		if err := reloaded.ReplaceExperiences([]domain.ExperienceParams{{
			ID:           reloaded.Experiences()[0].ID(),
			Company:      "Example",
			Position:     "Senior Developer",
			StartDate:    reloaded.Experiences()[0].StartDate(),
			Achievements: reloaded.Experiences()[0].Achievements(),
			Stack:        experienceStackParams(reloaded.Experiences()[0]),
		}}); err != nil {
			return err
		}
		if err := reloaded.ReplaceProjects([]domain.ProjectParams{{
			ID:          reloaded.Projects()[0].ID(),
			Name:        reloaded.Projects()[0].Name(),
			Screenshots: reloaded.Projects()[0].Screenshots(),
		}}); err != nil {
			return err
		}
		if err := resumeRepo.Save(tx, reloaded); err != nil {
			return err
		}

		return rollbackErr
	})
	require.ErrorIs(t, err, rollbackErr)

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	var version int64
	require.NoError(t, conn.QueryRow(ctx, "SELECT version FROM resumes WHERE id = $1", resume.ID().UUID()).Scan(&version))
	require.Equal(t, int64(1), version, "rollback must restore the aggregate version")

	for _, table := range []string{
		"resume_profiles",
		"resume_skills",
		"resume_experiences",
		"resume_projects",
		"resume_contacts",
	} {
		var count int
		require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM "+table+" WHERE resume_id = $1", resume.ID().UUID()).Scan(&count))
		require.Zero(t, count, table+" must rollback with the surrounding transaction")
	}

	var childCounts struct{ profile, contacts, contactHidden, contactText, hard, soft, achievements, tasks, stack, screenshots int }
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_profile_work_formats").Scan(&childCounts.profile))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_contacts").Scan(&childCounts.contacts))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_contact_hidden").Scan(&childCounts.contactHidden))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_contacts_text").Scan(&childCounts.contactText))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_hard").Scan(&childCounts.hard))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_skill_soft").Scan(&childCounts.soft))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_experience_achievements").Scan(&childCounts.achievements))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_experience_tasks").Scan(&childCounts.tasks))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_experience_stack").Scan(&childCounts.stack))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume_project_screenshots").Scan(&childCounts.screenshots))
	require.Zero(t, childCounts.profile)
	require.Zero(t, childCounts.contacts)
	require.Zero(t, childCounts.contactHidden)
	require.Zero(t, childCounts.contactText)
	require.Zero(t, childCounts.hard)
	require.Zero(t, childCounts.soft)
	require.Zero(t, childCounts.achievements)
	require.Zero(t, childCounts.tasks)
	require.Zero(t, childCounts.stack)
	require.Zero(t, childCounts.screenshots)
}
