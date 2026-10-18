//go:build integration

package persistence_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	analyticsreaders "developers-cv/backend/internal/modules/analytics/infrastructure/persistence/readers"
	analyticsrepository "developers-cv/backend/internal/modules/analytics/infrastructure/persistence/repository"
	resumelinkdomain "developers-cv/backend/internal/modules/resume_link/domain"
	resumelinkrepository "developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/pgtest"
)

type analyticsTestRepository struct {
	*resumelinkrepository.ResumeLinkRepository
	*analyticsrepository.AnalyticsRepository
}

type analyticsTestReader struct {
	*analyticsreaders.AnalyticsReader
}

func newDatabase(t *testing.T) (context.Context, string, *gorm.DB, *analyticsTestRepository, *analyticsTestReader) {
	t.Helper()
	ctx, dsn := pgtest.NewDatabase(t)
	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })
	return ctx, dsn, gormDB, &analyticsTestRepository{
		ResumeLinkRepository: resumelinkrepository.NewResumeLinkRepository(gormDB),
		AnalyticsRepository:  analyticsrepository.NewAnalyticsRepository(gormDB),
	}, &analyticsTestReader{
		AnalyticsReader: analyticsreaders.NewAnalyticsReader(gormDB),
	}
}

func insertUserAndResume(t *testing.T, ctx context.Context, dsn string, userID, resumeID sharedvalueobjects.ID) {
	t.Helper()
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	_, err = conn.Exec(ctx, `INSERT INTO users (id) VALUES ($1) ON CONFLICT (id) DO NOTHING`, userID.UUID())
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resumes (id, user_id, title, language, status, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, now(), now())`,
		resumeID.UUID(), userID.UUID(), "Test Resume", "ru", "draft", 1)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_overviews
		(resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third)
		VALUES ($1, 'Test', 'Resume', '', '')`, resumeID.UUID())
	require.NoError(t, err)
}

func newTestResumeLink(t *testing.T, resumeID sharedvalueobjects.ID, name string) *resumelinkdomain.ResumeLink {
	t.Helper()
	link, err := resumelinkdomain.NewResumeLink(resumelinkdomain.ResumeLinkParams{
		ResumeID: resumeID, Name: name,
	})
	require.NoError(t, err)
	return link
}
