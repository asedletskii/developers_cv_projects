//go:build integration

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap/zaptest"

	resumecommands "developers-cv/backend/internal/modules/resume/application/commands"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumereaders "developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	resumerepository "developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/analytics"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/readers"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/repository"
	"developers-cv/backend/internal/modules/resume_link/transport/http/public"
	"developers-cv/backend/internal/platform/database"
	"developers-cv/backend/internal/platform/http/middleware"
)

func applyAllMigrations(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	migrationsDir := filepath.Join("..", "..", "..", "..", "..", "db", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	require.NoError(t, err)
	var files []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)
	for _, name := range files {
		content, err := os.ReadFile(filepath.Join(migrationsDir, name))
		require.NoError(t, err)
		_, err = conn.Exec(ctx, string(content))
		require.NoErrorf(t, err, "migration %s", name)
	}
}

func TestCreateResumeLinkIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, testcontainers.TerminateContainer(container)) })
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyAllMigrations(t, ctx, dsn)
	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	userID := uuid.New()
	resumeID := uuid.New()
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	_, err = conn.Exec(ctx, `INSERT INTO users (id) VALUES ($1)`, userID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())`, resumeID, userID, "Test Resume", "en", "draft", 1)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, $2, $3, $4, $5)`, resumeID, "Hello", "from", "", "test")
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	linkRepository := repository.NewResumeLinkRepository(gormDB)
	create := commands.NewCreateResumeLinkHandler(linkRepository, resumeReader, database.NewUnitOfWork(gormDB))
	handler := NewHandler(HandlerDependencies{Create: create})

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	RegisterRoutes(handler, engine.Group("/api/v1"), engine.Group("/api/v1/lk"))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/link", bytes.NewBufferString(`{"resumeId":"`+resumeID.String()+`","name":"HH"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var count int
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_links WHERE resume_id = $1 AND name = 'HH'`, resumeID).Scan(&count))
	require.Equal(t, 1, count)
}

func TestPublicResumeThroughLinkIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, testcontainers.TerminateContainer(container)) })
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyAllMigrations(t, ctx, dsn)
	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	userID := uuid.New()
	resumeID := uuid.New()
	activeLinkID := uuid.New()
	inactiveLinkID := uuid.New()
	draftResumeID := uuid.New()
	draftLinkID := uuid.New()
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	_, err = conn.Exec(ctx, `INSERT INTO users (id) VALUES ($1)`, userID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())`, resumeID, userID, "Public Resume", "en", "opened", 1)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())`, draftResumeID, userID, "Draft Resume", "en", "draft", 1)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, $2, $3, $4, $5)`, resumeID, "Hello", "from", "", "public")
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_links (id, resume_id, name, status, version, created_at, updated_at) VALUES ($1, $2, 'HH', 'active', 1, now(), now()), ($3, $2, 'Old', 'inactive', 1, now(), now()), ($4, $5, 'Draft', 'active', 1, now(), now())`, activeLinkID, resumeID, inactiveLinkID, draftLinkID, draftResumeID)
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	linkReader := readers.NewResumeLinkReader(gormDB)
	publicResume := resumequeries.NewPublicResumeQueryHandler(resumeReader)
	eventRecorder := analytics.NewRecorder(repository.NewResumeLinkRepository(gormDB), zaptest.NewLogger(t))
	publicLink := queries.NewPublicResumeLinkQueryHandler(linkReader, publicResume, eventRecorder)
	handler := public.NewHandler(publicLink)

	engine := gin.New()
	api := engine.Group("/api/v1")
	RegisterRoutes(NewHandler(HandlerDependencies{Public: handler}), api, api.Group("/lk"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+activeLinkID.String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"title":"Public Resume"`)

	var viewEvents int
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_link_events WHERE link_id = $1 AND event_type = 'view'`, activeLinkID).Scan(&viewEvents))
	require.Equal(t, 1, viewEvents)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+inactiveLinkID.String(), nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+resumeID.String(), nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+draftLinkID.String(), nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_link_events WHERE link_id = $1`, activeLinkID).Scan(&viewEvents))
	require.Equal(t, 1, viewEvents)
}

func TestContactViewThroughLinkIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, testcontainers.TerminateContainer(container)) })
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyAllMigrations(t, ctx, dsn)
	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	userID := uuid.New()
	resumeID := uuid.New()
	inactiveLinkID := uuid.New()
	activeLinkID := uuid.New()
	draftResumeID := uuid.New()
	draftLinkID := uuid.New()
	closedResumeID := uuid.New()
	closedLinkID := uuid.New()
	contactsID := uuid.New()

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	_, err = conn.Exec(ctx, `INSERT INTO users (id) VALUES ($1)`, userID)
	require.NoError(t, err)
	for _, item := range []struct {
		id, status string
	}{
		{resumeID.String(), "opened"},
		{draftResumeID.String(), "draft"},
		{closedResumeID.String(), "closed"},
	} {
		_, err = conn.Exec(ctx, `INSERT INTO resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())`, uuid.MustParse(item.id), userID, item.status+" Resume", "en", item.status, 1)
		require.NoError(t, err)
	}
	_, err = conn.Exec(ctx, `INSERT INTO resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, $2, $3, $4, $5)`, resumeID, "Hello", "from", "", "public")
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_contacts (id, resume_id, phone, email) VALUES ($1, $2, $3, $4)`, contactsID, resumeID, "+33123456789", "private@example.com")
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_contact_hidden (contacts_id, value, position) VALUES ($1, 'email', 0)`, contactsID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_links (id, resume_id, name, status, version, created_at, updated_at) VALUES ($1, $2, 'Active', 'active', 1, now(), now()), ($3, $2, 'Inactive', 'inactive', 1, now(), now()), ($4, $5, 'Draft', 'active', 1, now(), now()), ($6, $7, 'Closed', 'active', 1, now(), now())`, activeLinkID, resumeID, inactiveLinkID, draftLinkID, draftResumeID, closedLinkID, closedResumeID)
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	linkReader := readers.NewResumeLinkReader(gormDB)
	contactView := resumequeries.NewContactViewQueryHandler(resumeReader)
	eventRecorder := analytics.NewRecorder(repository.NewResumeLinkRepository(gormDB), zaptest.NewLogger(t))
	contactThroughLink := queries.NewContactViewThroughLinkQueryHandler(linkReader, contactView, eventRecorder)
	contactHandler := public.NewContactViewHandler(contactThroughLink)

	engine := gin.New()
	api := engine.Group("/api/v1")
	RegisterRoutes(NewHandler(HandlerDependencies{ContactView: contactHandler}), api, api.Group("/lk"))

	postContact := func(linkID uuid.UUID, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+linkID.String()+"/contact-views", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}

	response := postContact(activeLinkID, `{"contact":"email"}`)
	require.Equal(t, http.StatusOK, response.Code)
	require.Contains(t, response.Body.String(), `"value":"private@example.com"`)

	var contactEvents int
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_link_events WHERE link_id = $1 AND event_type = 'contact_view'`, activeLinkID).Scan(&contactEvents))
	require.Equal(t, 1, contactEvents)

	response = postContact(activeLinkID, `{"contact":"phone"}`)
	require.Equal(t, http.StatusConflict, response.Code)

	response = postContact(inactiveLinkID, `{"contact":"email"}`)
	require.Equal(t, http.StatusNotFound, response.Code)

	response = postContact(draftLinkID, `{"contact":"email"}`)
	require.Equal(t, http.StatusNotFound, response.Code)

	response = postContact(closedLinkID, `{"contact":"email"}`)
	require.Equal(t, http.StatusNotFound, response.Code)

	response = postContact(uuid.New(), `{"contact":"email"}`)
	require.Equal(t, http.StatusNotFound, response.Code)

	response = postContact(resumeID, `{"contact":"email"}`)
	require.Equal(t, http.StatusNotFound, response.Code)

	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_link_events WHERE link_id = $1`, activeLinkID).Scan(&contactEvents))
	require.Equal(t, 1, contactEvents)
}

func TestResumeLinkStatisticsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, testcontainers.TerminateContainer(container)) })
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyAllMigrations(t, ctx, dsn)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	userID := uuid.New()
	resumeID := uuid.New()
	linkID := uuid.New()
	contactsID := uuid.New()

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	_, err = conn.Exec(ctx, `INSERT INTO users (id) VALUES ($1)`, userID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())`, resumeID, userID, "Statistics Resume", "en", "opened", 1)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, $2, $3, $4, $5)`, resumeID, "Hello", "from", "", "statistics")
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_contacts (id, resume_id, phone, email) VALUES ($1, $2, $3, $4)`, contactsID, resumeID, "+33123456789", "private@example.com")
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_contact_hidden (contacts_id, value, position) VALUES ($1, 'email', 0)`, contactsID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_links (id, resume_id, name, status, version, created_at, updated_at) VALUES ($1, $2, 'HH', 'active', 1, now(), now())`, linkID, resumeID)
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	linkReader := readers.NewResumeLinkReader(gormDB)
	linkRepository := repository.NewResumeLinkRepository(gormDB)
	eventRecorder := analytics.NewRecorder(linkRepository, zaptest.NewLogger(t))
	publicResume := resumequeries.NewPublicResumeQueryHandler(resumeReader)
	publicLink := queries.NewPublicResumeLinkQueryHandler(linkReader, publicResume, eventRecorder)
	contactView := resumequeries.NewContactViewQueryHandler(resumeReader)
	contactThroughLink := queries.NewContactViewThroughLinkQueryHandler(linkReader, contactView, eventRecorder)
	statistics := queries.NewGetResumeLinkStatisticsQueryHandler(linkReader, resumeReader, linkReader)

	publicHandler := public.NewHandler(publicLink)
	contactHandler := public.NewContactViewHandler(contactThroughLink)
	linkHandler := NewHandler(HandlerDependencies{
		Public:      publicHandler,
		ContactView: contactHandler,
		Statistics:  statistics,
	})

	engine := gin.New()
	api := engine.Group("/api/v1")
	lk := api.Group("/lk")
	lk.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	RegisterRoutes(linkHandler, api, lk)

	getPublic := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+linkID.String(), nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}
	postContact := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+linkID.String()+"/contact-views", bytes.NewBufferString(`{"contact":"email"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}

	require.Equal(t, http.StatusOK, getPublic().Code)
	require.Equal(t, http.StatusOK, getPublic().Code)
	require.Equal(t, http.StatusOK, postContact().Code)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+linkID.String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"views":2,"contactViews":1}`, w.Body.String())
}

func TestResumeLinkLifecycleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, testcontainers.TerminateContainer(container)) })
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyAllMigrations(t, ctx, dsn)
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
	_, err = conn.Exec(ctx, `INSERT INTO resume_link_events (id, link_id, event_type, occurred_at) VALUES ($1, $2, 'view', now())`, uuid.New(), linkID)
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	linkReader := readers.NewResumeLinkReader(gormDB)
	linkRepository := repository.NewResumeLinkRepository(gormDB)
	eventRecorder := analytics.NewRecorder(linkRepository, zaptest.NewLogger(t))
	publicResume := resumequeries.NewPublicResumeQueryHandler(resumeReader)
	publicLink := queries.NewPublicResumeLinkQueryHandler(linkReader, publicResume, eventRecorder)
	lifecycle := commands.NewReplaceStatusResumeLinkHandler(linkRepository, resumeReader, database.NewUnitOfWork(gormDB))

	linkHandler := NewHandler(HandlerDependencies{
		Public:    public.NewHandler(publicLink),
		Lifecycle: lifecycle,
	})
	engine := gin.New()
	api := engine.Group("/api/v1")
	lk := api.Group("/lk")
	lk.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	RegisterRoutes(linkHandler, api, lk)

	getPublic := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+linkID.String(), nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}
	changeStatus := func(status string, version int) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/link/"+linkID.String()+"/status", bytes.NewBufferString(`{"status":"`+status+`","version":`+strconv.Itoa(version)+`}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}

	require.Equal(t, http.StatusOK, getPublic().Code)
	require.Equal(t, http.StatusNoContent, changeStatus("inactive", 1).Code)
	require.Equal(t, http.StatusNotFound, getPublic().Code)

	var eventCount int
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_link_events WHERE link_id = $1`, linkID).Scan(&eventCount))
	require.Equal(t, 2, eventCount)

	require.Equal(t, http.StatusPreconditionFailed, changeStatus("active", 1).Code)
	require.Equal(t, http.StatusNoContent, changeStatus("active", 2).Code)
	require.Equal(t, http.StatusOK, getPublic().Code)

	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_link_events WHERE link_id = $1`, linkID).Scan(&eventCount))
	require.Equal(t, 3, eventCount)
}

func TestResumeLinksFullIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, testcontainers.TerminateContainer(container)) })
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyAllMigrations(t, ctx, dsn)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	userA := uuid.New()
	userB := uuid.New()
	resumeA := uuid.New()
	resumeB := uuid.New()
	resumeC := uuid.New()
	contactsA := uuid.New()
	contactsB := uuid.New()
	contactsC := uuid.New()

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	_, err = conn.Exec(ctx, `INSERT INTO users (id) VALUES ($1), ($2)`, userA, userB)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `
		INSERT INTO resumes (id, user_id, title, language, status, version, created_at, updated_at)
		VALUES
			($1, $4, 'Resume A', 'en', 'opened', 1, now(), now()),
			($2, $4, 'Resume B', 'en', 'opened', 1, now(), now()),
			($3, $5, 'Resume C', 'en', 'opened', 1, now(), now())
	`, resumeA, resumeB, resumeC, userA, userB)
	require.NoError(t, err)
	for _, resumeID := range []uuid.UUID{resumeA, resumeB, resumeC} {
		_, err = conn.Exec(ctx, `INSERT INTO resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, 'Hello', 'from', '', 'integration')`, resumeID)
		require.NoError(t, err)
	}

	_, err = conn.Exec(ctx, `INSERT INTO resume_contacts (id, resume_id, phone, email) VALUES ($1, $3, '+33111111111', 'a@example.com'), ($2, $4, '+33222222222', 'b@example.com'), ($5, $6, '+33333333333', 'github@example.com')`, contactsA, contactsC, resumeA, resumeC, contactsB, resumeB)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_contact_hidden (contacts_id, value, position) VALUES ($1, 'email', 0), ($2, 'email', 0), ($3, 'email', 0)`, contactsA, contactsC, contactsB)
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	resumeRepository := resumerepository.NewResumeRepository(gormDB)
	linkReader := readers.NewResumeLinkReader(gormDB)
	linkRepository := repository.NewResumeLinkRepository(gormDB)
	eventRecorder := analytics.NewRecorder(linkRepository, zaptest.NewLogger(t))
	publicResume := resumequeries.NewPublicResumeQueryHandler(resumeReader)
	publicLink := queries.NewPublicResumeLinkQueryHandler(linkReader, publicResume, eventRecorder)
	contactView := resumequeries.NewContactViewQueryHandler(resumeReader)
	contactThroughLink := queries.NewContactViewThroughLinkQueryHandler(linkReader, contactView, eventRecorder)
	statistics := queries.NewGetResumeLinkStatisticsQueryHandler(linkReader, resumeReader, linkReader)
	createLink := commands.NewCreateResumeLinkHandler(linkRepository, resumeReader, database.NewUnitOfWork(gormDB))
	listLinks := queries.NewListResumeLinksQueryHandler(linkReader, resumeReader)
	deleteLink := commands.NewDeleteResumeLinkHandler(linkRepository, resumeReader, database.NewUnitOfWork(gormDB))
	lifecycle := commands.NewReplaceStatusResumeLinkHandler(linkRepository, resumeReader, database.NewUnitOfWork(gormDB))
	resumeStatus := resumecommands.NewChangeResumeStatusHandler(resumeRepository, database.NewUnitOfWork(gormDB))
	deleteResume := resumecommands.NewDeleteResumeHandler(resumeRepository, database.NewUnitOfWork(gormDB))

	linkHandler := NewHandler(HandlerDependencies{
		Create:      createLink,
		List:        listLinks,
		Delete:      deleteLink,
		Public:      public.NewHandler(publicLink),
		ContactView: public.NewContactViewHandler(contactThroughLink),
		Statistics:  statistics,
		Lifecycle:   lifecycle,
	})
	resumeHandler := resumehttp.NewHandler(resumehttp.HandlerDependencies{
		DeleteResume:       deleteResume,
		ChangeResumeStatus: resumeStatus,
	})

	var currentUser uuid.UUID = userA
	engine := gin.New()
	api := engine.Group("/api/v1")
	lk := api.Group("/lk")
	lk.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return currentUser, nil
	}))
	resumehttp.RegisterRoutes(resumeHandler, api, lk)
	RegisterRoutes(linkHandler, api, lk)

	do := func(method, path, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}
	createLinkHTTP := func(resumeID uuid.UUID, name string) uuid.UUID {
		response := do(http.MethodPost, "/api/v1/lk/link", `{"resumeId":"`+resumeID.String()+`","name":"`+name+`"}`)
		require.Equal(t, http.StatusCreated, response.Code)
		var payload struct {
			ID string `json:"id"`
		}
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &payload))
		linkID, err := uuid.Parse(payload.ID)
		require.NoError(t, err)
		return linkID
	}
	getPublic := func(linkID uuid.UUID) *httptest.ResponseRecorder {
		return do(http.MethodGet, "/api/v1/cv/"+linkID.String(), "")
	}
	postContact := func(linkID uuid.UUID) *httptest.ResponseRecorder {
		return do(http.MethodPost, "/api/v1/cv/"+linkID.String()+"/contact-views", `{"contact":"email"}`)
	}
	getStatistics := func(linkID uuid.UUID) *httptest.ResponseRecorder {
		return do(http.MethodGet, "/api/v1/lk/link/"+linkID.String()+"/statistics", "")
	}

	linkAHH := createLinkHTTP(resumeA, "HH")
	linkATelegram := createLinkHTTP(resumeA, "Telegram")
	linkAGitHub := createLinkHTTP(resumeB, "GitHub")

	currentUser = userB
	linkBHH := createLinkHTTP(resumeC, "HH")
	currentUser = userA

	publicAHH := getPublic(linkAHH)
	publicATelegram := getPublic(linkATelegram)
	publicAGitHub := getPublic(linkAGitHub)
	publicBHH := getPublic(linkBHH)
	require.Contains(t, publicAHH.Body.String(), `"title":"Resume A"`)
	require.Contains(t, publicATelegram.Body.String(), `"title":"Resume A"`)
	require.Contains(t, publicAGitHub.Body.String(), `"title":"Resume B"`)
	require.Contains(t, publicBHH.Body.String(), `"title":"Resume C"`)

	require.Equal(t, http.StatusOK, getPublic(linkAHH).Code) // HH: 3 views
	require.Equal(t, http.StatusOK, getPublic(linkAHH).Code)
	require.Equal(t, http.StatusOK, getPublic(linkATelegram).Code) // Telegram: 2 views
	require.Equal(t, http.StatusOK, postContact(linkAHH).Code)
	require.Equal(t, http.StatusOK, postContact(linkAGitHub).Code)
	for _, item := range []struct {
		id   uuid.UUID
		body string
	}{
		{linkAHH, `{"views":3,"contactViews":1}`},
		{linkATelegram, `{"views":2,"contactViews":0}`},
		{linkAGitHub, `{"views":1,"contactViews":1}`},
	} {
		response := getStatistics(item.id)
		require.Equal(t, http.StatusOK, response.Code)
		require.JSONEq(t, item.body, response.Body.String())
	}

	currentUser = userB
	response := getStatistics(linkBHH)
	require.Equal(t, http.StatusOK, response.Code)
	require.JSONEq(t, `{"views":1,"contactViews":0}`, response.Body.String())
	currentUser = userA

	currentUser = userB
	require.Equal(t, http.StatusNotFound, do(http.MethodGet, "/api/v1/lk/links?resumeId="+resumeA.String(), "").Code)
	require.Equal(t, http.StatusNotFound, do(http.MethodDelete, "/api/v1/lk/link/"+linkAHH.String(), "").Code)
	require.Equal(t, http.StatusNotFound, do(http.MethodPut, "/api/v1/lk/link/"+linkAHH.String()+"/status", `{"status":"inactive","version":1}`).Code)
	require.Equal(t, http.StatusNotFound, getStatistics(linkAHH).Code)
	currentUser = userA

	response = do(http.MethodPut, "/api/v1/lk/resume/"+resumeA.String()+"/status", `{"status":"closed","version":1}`)
	require.Equal(t, http.StatusNoContent, response.Code)
	require.Equal(t, http.StatusNotFound, getPublic(linkAHH).Code)
	require.Equal(t, http.StatusNotFound, getPublic(linkATelegram).Code)
	require.Equal(t, http.StatusOK, getPublic(linkAGitHub).Code)
	require.Equal(t, http.StatusOK, getPublic(linkBHH).Code)

	response = do(http.MethodPut, "/api/v1/lk/resume/"+resumeA.String()+"/status", `{"status":"opened","version":2}`)
	require.Equal(t, http.StatusNoContent, response.Code)
	require.Equal(t, http.StatusOK, getPublic(linkAHH).Code)
	require.Equal(t, http.StatusOK, getPublic(linkATelegram).Code)

	response = do(http.MethodPut, "/api/v1/lk/link/"+linkAHH.String()+"/status", `{"status":"inactive","version":1}`)
	require.Equal(t, http.StatusNoContent, response.Code)
	require.Equal(t, http.StatusNotFound, getPublic(linkAHH).Code)
	require.Equal(t, http.StatusOK, getPublic(linkATelegram).Code)
	response = getStatistics(linkAHH)
	require.Equal(t, http.StatusOK, response.Code)
	var inactiveStats struct {
		Views        int `json:"views"`
		ContactViews int `json:"contactViews"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &inactiveStats))
	require.Equal(t, 4, inactiveStats.Views)
	require.Equal(t, 1, inactiveStats.ContactViews)

	response = do(http.MethodPut, "/api/v1/lk/link/"+linkAHH.String()+"/status", `{"status":"active","version":2}`)
	require.Equal(t, http.StatusNoContent, response.Code)
	require.Equal(t, http.StatusOK, getPublic(linkAHH).Code)

	response = do(http.MethodDelete, "/api/v1/lk/resume/"+resumeA.String(), "")
	require.Equal(t, http.StatusNoContent, response.Code)

	var linksAfterDelete, eventsAfterDelete int
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_links WHERE resume_id = $1`, resumeA).Scan(&linksAfterDelete))
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_link_events WHERE link_id IN ($1, $2)`, linkAHH, linkATelegram).Scan(&eventsAfterDelete))
	require.Equal(t, 0, linksAfterDelete)
	require.Equal(t, 0, eventsAfterDelete)

	currentUser = userB
	require.Equal(t, http.StatusOK, getPublic(linkBHH).Code)
	require.Equal(t, http.StatusOK, getStatistics(linkBHH).Code)
	currentUser = userA
	require.Equal(t, http.StatusNotFound, getPublic(linkAHH).Code)
	require.Equal(t, http.StatusNotFound, getPublic(linkATelegram).Code)

	response = getPublic(linkAGitHub)
	require.Equal(t, http.StatusOK, response.Code)
	require.Contains(t, response.Body.String(), `"title":"Resume B"`)

	require.Equal(t, http.StatusNotFound, do(http.MethodGet, "/api/v1/cv/"+resumeB.String(), "").Code)
}
