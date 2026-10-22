//go:build integration

package persistence_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	analyticsrepository "developers-cv/backend/internal/modules/analytics/infrastructure/persistence/repository"
	authcommands "developers-cv/backend/internal/modules/auth/application/commands"
	authrepository "developers-cv/backend/internal/modules/auth/infrastructure/persistence/repository"
	resumecommands "developers-cv/backend/internal/modules/resume/application/commands"
	resumerepository "developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	resumelinkcommands "developers-cv/backend/internal/modules/resume_link/application/commands"
	resumelinkrepository "developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/repository"
	usercommands "developers-cv/backend/internal/modules/user/application/commands"
	userrepository "developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/pgtest"
)

type deleteUserFixture struct {
	userA, userB          uuid.UUID
	resumeA1, resumeA2    uuid.UUID
	resumeB               uuid.UUID
	linkA1, linkA2, linkB uuid.UUID
	viewA1, viewA2, viewB uuid.UUID
}

func TestDeleteUserRemovesAllModuleOwnedData(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)
	db, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(db)) })

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	fixture := seedDeleteUserFixture(t, ctx, conn)
	userRepo := userrepository.NewUserRepository(db)
	analyticsRepo := analyticsrepository.NewAnalyticsRepository(db)
	resumeLinkCleanup := resumelinkcommands.NewDeleteResumeLinksHandler(resumelinkrepository.NewResumeLinkRepository(db), analyticsRepo)
	resumeCleanup := resumecommands.NewDeleteUserResumesHandler(resumerepository.NewResumeRepository(db), resumeLinkCleanup)
	authCleanup := authcommands.NewDeleteUserDataHandler(authrepository.NewUserDataCleanupRepository(db))
	handler := usercommands.NewDeleteUserHandler(userRepo, authCleanup, resumeCleanup, database.NewUnitOfWork(db))

	userAID, err := sharedvalueobjects.IDFromUUID(fixture.userA)
	require.NoError(t, err)
	require.NoError(t, handler.Handle(ctx, usercommands.DeleteUserCommand{UserID: userAID}))
	require.NoError(t, handler.Handle(ctx, usercommands.DeleteUserCommand{UserID: userAID}), "repeated deletion must be a no-op")

	var count int64
	require.NoError(t, db.Table("users.users").Where("id = ?", fixture.userA).Count(&count).Error)
	require.Zero(t, count, "users.users")
	for _, table := range []string{
		"users.user_profiles", "users.user_profile_localizations",
		"auth.user_identities", "auth.auth_sessions",
		"resume.resumes",
	} {
		require.NoError(t, db.Table(table).Where("user_id = ?", fixture.userA).Count(&count).Error)
		require.Zero(t, count, table)
	}
	require.NoError(t, db.Table("resume.resume_overviews").Where("resume_id IN ?", []uuid.UUID{fixture.resumeA1, fixture.resumeA2}).Count(&count).Error)
	require.Zero(t, count, "resume.resume_overviews")
	require.NoError(t, db.Table("resume_link.resume_links").Where("id IN ?", []uuid.UUID{fixture.linkA1, fixture.linkA2}).Count(&count).Error)
	require.Zero(t, count, "resume_link.resume_links")

	for _, table := range []string{"analytics.analytics_views", "analytics.analytics_events", "analytics.analytics_clicks", "analytics.analytics_visitors"} {
		require.NoError(t, db.Table(table).Where("link_id IN ?", []uuid.UUID{fixture.linkA1, fixture.linkA2}).Count(&count).Error)
		require.Zero(t, count, table)
	}

	require.NoError(t, db.Table("users.users").Where("id = ?", fixture.userB).Count(&count).Error)
	require.Equal(t, int64(1), count)
	require.NoError(t, db.Table("resume.resumes").Where("user_id = ?", fixture.userB).Count(&count).Error)
	require.Equal(t, int64(1), count)
	require.NoError(t, db.Table("resume_link.resume_links").Where("resume_id = ?", fixture.resumeB).Count(&count).Error)
	require.Equal(t, int64(1), count)
	for _, table := range []string{"analytics.analytics_views", "analytics.analytics_events", "analytics.analytics_clicks", "analytics.analytics_visitors"} {
		require.NoError(t, db.Table(table).Where("link_id = ?", fixture.linkB).Count(&count).Error)
		require.Equal(t, int64(1), count, table)
	}

	pgtest.RequireNoOrphans(t, ctx, conn)
}

func TestDeleteUserRollsBackWhenAnalyticsCleanupFails(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)
	db, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(db)) })

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	fixture := seedDeleteUserFixture(t, ctx, conn)
	userRepo := userrepository.NewUserRepository(db)
	authCleanup := authcommands.NewDeleteUserDataHandler(authrepository.NewUserDataCleanupRepository(db))
	resumeRepo := resumerepository.NewResumeRepository(db)
	linkRepo := resumelinkrepository.NewResumeLinkRepository(db)
	uow := database.NewUnitOfWork(db)
	analyticsErr := errors.New("analytics cleanup failed")
	failingLinks := resumelinkcommands.NewDeleteResumeLinksHandler(linkRepo, failingAnalyticsCleanup{err: analyticsErr})
	failingResumes := resumecommands.NewDeleteUserResumesHandler(resumeRepo, failingLinks)
	handler := usercommands.NewDeleteUserHandler(userRepo, authCleanup, failingResumes, uow)

	userAID, err := sharedvalueobjects.IDFromUUID(fixture.userA)
	require.NoError(t, err)
	err = handler.Handle(ctx, usercommands.DeleteUserCommand{UserID: userAID})
	require.ErrorIs(t, err, analyticsErr)

	var count int64
	require.NoError(t, db.Table("users.users").Where("id = ?", fixture.userA).Count(&count).Error)
	require.Equal(t, int64(1), count, "users.users must survive rollback")
	for _, table := range []string{"users.user_profiles", "auth.user_identities", "auth.auth_sessions"} {
		require.NoError(t, db.Table(table).Where("user_id = ?", fixture.userA).Count(&count).Error)
		require.Equal(t, int64(1), count, table+" must survive rollback")
	}
	require.NoError(t, db.Table("resume.resumes").Where("user_id = ?", fixture.userA).Count(&count).Error)
	require.Equal(t, int64(2), count, "resumes must survive rollback")
	require.NoError(t, db.Table("resume_link.resume_links").Where("id IN ?", []uuid.UUID{fixture.linkA1, fixture.linkA2}).Count(&count).Error)
	require.Equal(t, int64(2), count, "links must survive rollback")
	require.NoError(t, db.Table("analytics.analytics_views").Where("link_id IN ?", []uuid.UUID{fixture.linkA1, fixture.linkA2}).Count(&count).Error)
	require.Equal(t, int64(2), count, "analytics must survive rollback")
	pgtest.RequireNoOrphans(t, ctx, conn)

	cleanLinks := resumelinkcommands.NewDeleteResumeLinksHandler(linkRepo, analyticsrepository.NewAnalyticsRepository(db))
	cleanResumes := resumecommands.NewDeleteUserResumesHandler(resumeRepo, cleanLinks)
	cleanHandler := usercommands.NewDeleteUserHandler(userRepo, authCleanup, cleanResumes, uow)
	require.NoError(t, cleanHandler.Handle(ctx, usercommands.DeleteUserCommand{UserID: userAID}))
	pgtest.RequireNoOrphans(t, ctx, conn)
}

type failingAnalyticsCleanup struct {
	err error
}

func (f failingAnalyticsCleanup) DeleteByLinkIDs(context.Context, []sharedvalueobjects.ID) error {
	return f.err
}

func seedDeleteUserFixture(t *testing.T, ctx context.Context, conn *pgx.Conn) deleteUserFixture {
	t.Helper()
	fixture := deleteUserFixture{
		userA: uuid.New(), userB: uuid.New(),
		resumeA1: uuid.New(), resumeA2: uuid.New(), resumeB: uuid.New(),
		linkA1: uuid.New(), linkA2: uuid.New(), linkB: uuid.New(),
		viewA1: uuid.New(), viewA2: uuid.New(), viewB: uuid.New(),
	}
	_, err := conn.Exec(ctx, `INSERT INTO users.users (id) VALUES ($1), ($2)`, fixture.userA, fixture.userB)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO users.user_profiles (user_id) VALUES ($1), ($2)`, fixture.userA, fixture.userB)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO users.user_profile_localizations (user_id, language, name) VALUES ($1, 'en', 'A'), ($2, 'en', 'B')`, fixture.userA, fixture.userB)
	require.NoError(t, err)

	_, err = conn.Exec(ctx, `
		INSERT INTO auth.user_identities (id, user_id, provider, provider_subject, created_at, updated_at, last_login_at)
		VALUES ($1, $2, 'telegram', 'subject-a', now(), now(), now()),
		       ($3, $4, 'telegram', 'subject-b', now(), now(), now())
	`, uuid.New(), fixture.userA, uuid.New(), fixture.userB)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `
		INSERT INTO auth.auth_sessions (id, user_id, token_hash, created_at, last_seen_at, expires_at)
		VALUES ($1, $2, decode(repeat('00', 32), 'hex'), now(), now(), now() + interval '1 hour'),
		       ($3, $4, decode(repeat('11', 32), 'hex'), now(), now(), now() + interval '1 hour')
	`, uuid.New(), fixture.userA, uuid.New(), fixture.userB)
	require.NoError(t, err)

	_, err = conn.Exec(ctx, `
		INSERT INTO resume.resumes (id, user_id, title, language, status, version, created_at, updated_at)
		VALUES ($1, $4, 'A1', 'en', 'opened', 1, now(), now()),
		       ($2, $4, 'A2', 'en', 'draft', 1, now(), now()),
		       ($3, $5, 'B', 'en', 'opened', 1, now(), now())
	`, fixture.resumeA1, fixture.resumeA2, fixture.resumeB, fixture.userA, fixture.userB)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `
		INSERT INTO resume.resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third)
		VALUES ($1, 'Hello', 'from', '', 'A1'), ($2, 'Hello', 'from', '', 'A2'), ($3, 'Hello', 'from', '', 'B')
	`, fixture.resumeA1, fixture.resumeA2, fixture.resumeB)
	require.NoError(t, err)

	_, err = conn.Exec(ctx, `
		INSERT INTO resume_link.resume_links (id, resume_id, name, status, version, created_at, updated_at)
		VALUES ($1, $4, 'A1-link', 'active', 1, now(), now()),
		       ($2, $5, 'A2-link', 'active', 1, now(), now()),
		       ($3, $6, 'B-link', 'active', 1, now(), now())
	`, fixture.linkA1, fixture.linkA2, fixture.linkB, fixture.resumeA1, fixture.resumeA2, fixture.resumeB)
	require.NoError(t, err)

	_, err = conn.Exec(ctx, `
		INSERT INTO analytics.analytics_views
			(id, link_id, resume_id, visitor_id, ip_hash, device, browser, os, layout, started_at)
		VALUES
			($1, $4, $7, $10, decode(repeat('22', 32), 'hex'), 'desktop', 'chrome', 'macos', '{}', now()),
			($2, $5, $8, $11, decode(repeat('33', 32), 'hex'), 'desktop', 'chrome', 'macos', '{}', now()),
			($3, $6, $9, $12, decode(repeat('44', 32), 'hex'), 'desktop', 'chrome', 'macos', '{}', now())
	`,
		fixture.viewA1, fixture.viewA2, fixture.viewB,
		fixture.linkA1, fixture.linkA2, fixture.linkB,
		fixture.resumeA1, fixture.resumeA2, fixture.resumeB,
		uuid.New(), uuid.New(), uuid.New(),
	)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `
		INSERT INTO analytics.analytics_events (id, view_id, link_id, event_type, payload, occurred_at)
		VALUES ($1, $4, $7, 'download', '{}', now()),
		       ($2, $5, $8, 'download', '{}', now()),
		       ($3, $6, $9, 'download', '{}', now())
	`, uuid.New(), uuid.New(), uuid.New(), fixture.viewA1, fixture.viewA2, fixture.viewB, fixture.linkA1, fixture.linkA2, fixture.linkB)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `
		INSERT INTO analytics.analytics_clicks
			(view_id, link_id, section, x_pct, y_pct, viewport_w, layout_hash, occurred_at)
		VALUES ($1, $4, 'about', 10, 20, 1024, '0123456789abcdef', now()),
		       ($2, $5, 'about', 10, 20, 1024, '0123456789abcdef', now()),
		       ($3, $6, 'about', 10, 20, 1024, '0123456789abcdef', now())
	`, fixture.viewA1, fixture.viewA2, fixture.viewB, fixture.linkA1, fixture.linkA2, fixture.linkB)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `
		INSERT INTO analytics.analytics_visitors (link_id, visitor_id, first_seen_at, first_view_id, ordinal)
		VALUES ($1, $4, now(), $7, 1),
		       ($2, $5, now(), $8, 1),
		       ($3, $6, now(), $9, 1)
	`, fixture.linkA1, fixture.linkA2, fixture.linkB, uuid.New(), uuid.New(), uuid.New(), fixture.viewA1, fixture.viewA2, fixture.viewB)
	require.NoError(t, err)
	return fixture
}
