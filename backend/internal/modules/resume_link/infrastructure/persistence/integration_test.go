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

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	resumereaders "developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/readers"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func applyMigrations(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	migrationsDir := filepath.Join("..", "..", "..", "..", "..", "db", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	require.NoError(t, err)
	var migrationFiles []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		migrationFiles = append(migrationFiles, entry.Name())
	}
	sort.Strings(migrationFiles)
	for _, name := range migrationFiles {
		content, err := os.ReadFile(filepath.Join(migrationsDir, name))
		require.NoError(t, err)
		_, err = conn.Exec(ctx, string(content))
		require.NoErrorf(t, err, "failed to apply migration %s", name)
	}
}

func newDatabase(t *testing.T) (context.Context, string, *gorm.DB, *repository.ResumeLinkRepository, *readers.ResumeLinkReader) {
	t.Helper()
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
	return ctx, dsn, gormDB, repository.NewResumeLinkRepository(gormDB), readers.NewResumeLinkReader(gormDB)
}

func insertUserAndResume(t *testing.T, ctx context.Context, dsn string, userID sharedvalueobjects.ID, resumeID sharedvalueobjects.ID) {
	t.Helper()
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	_, err = conn.Exec(ctx, `INSERT INTO users (id) VALUES ($1) ON CONFLICT (id) DO NOTHING`, userID.UUID())
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())`, resumeID.UUID(), userID.UUID(), "Test Resume", "ru", "draft", 1)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, 'Test', 'Resume', '', '')`, resumeID.UUID())
	require.NoError(t, err)
}

func TestResumeLinkPersistence(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "  HH  "})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return repo.Create(tx, link) }))

	loaded, err := repo.GetByID(ctx, link.ID())
	require.NoError(t, err)
	require.Equal(t, link.ID(), loaded.ID())
	require.Equal(t, resumeID, loaded.ResumeID())
	require.Equal(t, "HH", loaded.Name())
	require.True(t, loaded.IsActive())
	require.Equal(t, int64(1), loaded.Version())

	var status string
	var createdAt, updatedAt time.Time
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	require.NoError(t, conn.QueryRow(ctx, `SELECT status, created_at, updated_at FROM resume_links WHERE id = $1`, link.ID().UUID()).Scan(&status, &createdAt, &updatedAt))
	require.Equal(t, "active", status)
	require.False(t, createdAt.IsZero())
	require.False(t, updatedAt.IsZero())

	links, err := reader.ListByResumeID(ctx, queries.ListResumeLinksQuery{ResumeID: resumeID})
	require.NoError(t, err)
	require.Len(t, links, 1)
	require.Equal(t, link.ID().String(), links[0].ID)
	require.Equal(t, "HH", links[0].Name)
	require.Equal(t, "active", links[0].Status)
	require.Equal(t, int64(1), links[0].Version)
	require.False(t, links[0].CreatedAt.IsZero())
	require.False(t, links[0].UpdatedAt.IsZero())

	linkView, err := reader.GetByID(ctx, queries.PublicResumeLinkQuery{LinkID: link.ID()})
	require.NoError(t, err)
	require.Equal(t, link.ID().String(), linkView.ID)
	require.Equal(t, resumeID.String(), linkView.ResumeID)

	second, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "GitHub",
	})
	require.NoError(t, err)
	foreignResumeID := sharedvalueobjects.NewID()
	foreignOwnerID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, foreignOwnerID, foreignResumeID)
	foreign, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ID: sharedvalueobjects.NewID(), ResumeID: foreignResumeID, Name: "Telegram",
	})
	require.NoError(t, err)

	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		if err := repo.Create(tx, second); err != nil {
			return err
		}
		return repo.Create(tx, foreign)
	}))

	links, err = reader.ListByResumeID(ctx, queries.ListResumeLinksQuery{UserID: owner, ResumeID: resumeID})
	require.NoError(t, err)
	require.Len(t, links, 2)
	require.Equal(t, "HH", links[0].Name)
	require.Equal(t, "GitHub", links[1].Name)
	require.Equal(t, link.ID().String(), links[0].ID)
	require.Equal(t, second.ID().String(), links[1].ID)

	empty, err := reader.ListByResumeID(ctx, queries.ListResumeLinksQuery{UserID: owner, ResumeID: sharedvalueobjects.NewID()})
	require.NoError(t, err)
	require.Empty(t, empty)
}

func TestResumeLinkMigrationCascade(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "GitHub"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return repo.Create(tx, link) }))

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	_, err = conn.Exec(ctx, `DELETE FROM resumes WHERE id = $1`, resumeID.UUID())
	require.NoError(t, err)
	var links int
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_links WHERE id = $1`, link.ID().UUID()).Scan(&links))
	require.Zero(t, links)
}

func TestReplaceStatusResumeLinkIntegration(t *testing.T) {
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

	userID := uuid.New()
	resumeID := uuid.New()
	linkID := uuid.New()
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	_, err = conn.Exec(ctx, `INSERT INTO users (id) VALUES ($1)`, userID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, 'Lifecycle Resume', 'en', 'opened', 1, now(), now())`, resumeID, userID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, 'Hello', 'from', '', 'lifecycle')`, resumeID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_links (id, resume_id, name, status, version, created_at, updated_at) VALUES ($1, $2, 'HH', 'active', 1, now(), now())`, linkID, resumeID)
	require.NoError(t, err)

	repo := repository.NewResumeLinkRepository(gormDB)
	resumeReader := resumereaders.NewResumeReader(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	handler := commands.NewReplaceStatusResumeLinkHandler(repo, resumeReader, uow)

	linkValue, err := sharedvalueobjects.IDFromUUID(linkID)
	require.NoError(t, err)
	userValue, err := sharedvalueobjects.IDFromUUID(userID)
	require.NoError(t, err)
	require.NoError(t, handler.Handle(ctx, commands.ReplaceStatusResumeLinkCommand{UserID: userValue, LinkID: linkValue, Active: false, Version: 1}))

	var status string
	var version int64
	require.NoError(t, conn.QueryRow(ctx, `SELECT status, version FROM resume_links WHERE id = $1`, linkID).Scan(&status, &version))
	require.Equal(t, "inactive", status)
	require.Equal(t, int64(2), version)

	require.ErrorIs(t, handler.Handle(ctx, commands.ReplaceStatusResumeLinkCommand{UserID: userValue, LinkID: linkValue, Active: true, Version: 1}), domain.ErrConcurrencyConflict)
	require.NoError(t, handler.Handle(ctx, commands.ReplaceStatusResumeLinkCommand{UserID: userValue, LinkID: linkValue, Active: true, Version: 2}))
	require.NoError(t, conn.QueryRow(ctx, `SELECT status, version FROM resume_links WHERE id = $1`, linkID).Scan(&status, &version))
	require.Equal(t, "active", status)
	require.Equal(t, int64(3), version)

}

func TestResumeLinkPersistenceEnforcesDatabaseNameLength(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, _, _, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	validName := strings.Repeat("ж", domain.MaxResumeLinkNameLength)
	_, err = conn.Exec(ctx, `INSERT INTO resume_links (id, resume_id, name) VALUES ($1, $2, $3)`, uuid.New(), resumeID.UUID(), validName)
	require.NoError(t, err)

	tooLongName := strings.Repeat("ж", domain.MaxResumeLinkNameLength+1)
	_, err = conn.Exec(ctx, `INSERT INTO resume_links (id, resume_id, name) VALUES ($1, $2, $3)`, uuid.New(), resumeID.UUID(), tooLongName)
	require.Error(t, err, "database must reject names longer than the domain limit")
}

func TestResumeLinkPersistenceEnforcesUniqueNameWithinResume(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)

	first, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "HH"})
	require.NoError(t, err)
	second, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "HH"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return repo.Create(tx, first) }))
	err = uow.Do(ctx, func(tx context.Context) error { return repo.Create(tx, second) })
	require.ErrorIs(t, err, domain.ErrDuplicateResumeLink)

	otherResumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, otherResumeID)
	otherResumeLink, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: otherResumeID, Name: "HH"})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return repo.Create(tx, otherResumeLink) }))
}

func TestReplaceStatusResumeLinkNoOpPreservesVersionAndUpdatedAt(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "HH", Version: 3})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return repo.Create(tx, link) }))

	var before time.Time
	require.NoError(t, gormDB.Raw(`SELECT updated_at FROM resume_links WHERE id = ?`, link.ID().UUID()).Scan(&before).Error)

	time.Sleep(10 * time.Millisecond)
	handler := commands.NewReplaceStatusResumeLinkHandler(repo, resumereaders.NewResumeReader(gormDB), uow)
	ownerValue := owner
	require.NoError(t, handler.Handle(ctx, commands.ReplaceStatusResumeLinkCommand{
		UserID:  ownerValue,
		LinkID:  link.ID(),
		Active:  true,
		Version: 3,
	}))

	var row struct {
		Status    string
		Version   int64
		UpdatedAt time.Time
	}
	require.NoError(t, gormDB.Raw(`SELECT status, version, updated_at FROM resume_links WHERE id = ?`, link.ID().UUID()).Scan(&row).Error)
	status, version, after := row.Status, row.Version, row.UpdatedAt
	require.Equal(t, string(domain.ResumeLinkStatusActive), status)
	require.Equal(t, int64(3), version)
	require.Equal(t, before, after)
}

func TestResumeLinkPersistenceEnforcesMaximumLinksConcurrently(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)

	const attempts = 55
	errs := make(chan error, attempts)
	var wg sync.WaitGroup
	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: fmt.Sprintf("link-%d", i)})
			if err != nil {
				errs <- err
				return
			}
			uow := database.NewUnitOfWork(gormDB)
			errs <- uow.Do(ctx, func(tx context.Context) error { return repo.Create(tx, link) })
		}(i)
	}
	wg.Wait()
	close(errs)

	var success, limited int
	for err := range errs {
		if err == nil {
			success++
			continue
		}
		require.ErrorIs(t, err, domain.ErrResumeLinkLimitExceeded)
		limited++
	}
	require.Equal(t, domain.MaxResumeLinksPerResume, success)
	require.Equal(t, attempts-domain.MaxResumeLinksPerResume, limited)

	var count int64
	require.NoError(t, gormDB.Model(&gormmodel.ResumeLink{}).Where("resume_id = ?", resumeID.UUID()).Count(&count).Error)
	require.Equal(t, int64(domain.MaxResumeLinksPerResume), count)
}
