//go:build integration

package persistence_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/pgtest"
)

func languageGradePtr(value domain.LanguageGrade) *domain.LanguageGrade {
	return &value
}

func newTestResume(t *testing.T, userID sharedvalueobjects.ID) *domain.Resume {
	t.Helper()
	root, err := domain.NewResume(domain.ResumeParams{
		UserID: userID, Title: "Test Resume", Language: domain.ResumeLanguageRU,
	})
	require.NoError(t, err)
	return root
}

func insertResumeTestUser(t *testing.T, ctx context.Context, dsn string, userID sharedvalueobjects.ID) {
	t.Helper()

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, conn.Close(ctx))
	})

	_, err = conn.Exec(ctx, `INSERT INTO users.users (id) VALUES ($1)`, userID.UUID())
	require.NoError(t, err)
}

func TestPostgresSkillsLanguagesAggregateVerticalSlice(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	owner := sharedvalueobjects.NewID()
	foreign := sharedvalueobjects.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	insertResumeTestUser(t, ctx, dsn, foreign)

	resumeRepo := repository.NewResumeRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	resume := newTestResume(t, owner)
	require.NoError(t, resume.ReplaceOverview(domain.ResumeOverviewParams{
		Greeting: mustGreeting(t, domain.GreetingParams{First: "Hello", Second: "I am", SecondHighlight: "Developer", Third: "Welcome"}),
	}))
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Create(tx, resume)
	}))

	hard := []string{"architecture", "testing"}
	soft := []string{"communication", "teamwork"}
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
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume.resume_skill_hard").Scan(&hardCount))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM resume.resume_skill_soft").Scan(&softCount))
	require.Equal(t, 2, hardCount)
	require.Equal(t, 2, softCount)

	require.NoError(t, loadedResume.ReplaceSpokenLanguages(domain.ResumeSpokenLanguagesParams{Entries: []domain.SpokenLanguageEntryParams{
		{Language: domain.SpokenLanguageEN, Grade: languageGradePtr(domain.LanguageGradeB2), Position: 0},
		{Language: domain.SpokenLanguageRU, Position: 1},
	}}))
	firstLanguageID := loadedResume.SpokenLanguages().Entries()[0].ID()
	secondLanguageID := loadedResume.SpokenLanguages().Entries()[1].ID()
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, loadedResume)
	}))
	loadedResume, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.Len(t, loadedResume.SpokenLanguages().Entries(), 2)
	require.Equal(t, firstLanguageID, loadedResume.SpokenLanguages().Entries()[0].ID())
	require.Equal(t, secondLanguageID, loadedResume.SpokenLanguages().Entries()[1].ID())

	require.NoError(t, loadedResume.ReplaceSpokenLanguages(domain.ResumeSpokenLanguagesParams{Entries: []domain.SpokenLanguageEntryParams{
		{ID: secondLanguageID, Language: domain.SpokenLanguageRU, Position: 0},
		{ID: firstLanguageID, Language: domain.SpokenLanguageEN, Grade: languageGradePtr(domain.LanguageGradeC1), Position: 1},
	}}))
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, loadedResume)
	}))
	loadedResume, err = resumeRepo.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.Equal(t, []domain.SpokenLanguage{domain.SpokenLanguageRU, domain.SpokenLanguageEN}, []domain.SpokenLanguage{loadedResume.SpokenLanguages().Entries()[0].Language(), loadedResume.SpokenLanguages().Entries()[1].Language()})

}

func mustGreeting(t *testing.T, params domain.GreetingParams) domain.Greeting {
	t.Helper()
	value, err := domain.NewGreeting(params)
	require.NoError(t, err)
	return value
}
