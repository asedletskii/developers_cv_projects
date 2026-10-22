//go:build integration

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	analyticscommands "developers-cv/backend/internal/modules/analytics/application/commands"
	analyticsqueries "developers-cv/backend/internal/modules/analytics/application/queries"
	analyticsadapters "developers-cv/backend/internal/modules/analytics/infrastructure/analytics"
	analyticsreaders "developers-cv/backend/internal/modules/analytics/infrastructure/persistence/readers"
	analyticsrepository "developers-cv/backend/internal/modules/analytics/infrastructure/persistence/repository"
	analyticshttp "developers-cv/backend/internal/modules/analytics/transport/http"
	resumecommands "developers-cv/backend/internal/modules/resume/application/commands"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumereaders "developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	resumerepository "developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/analyticsaccess"
	resumelinkreaders "developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/readers"
	resumelinkrepository "developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/repository"
	"developers-cv/backend/internal/modules/resume_link/transport/http/handlers"
	"developers-cv/backend/internal/platform/database"
	"developers-cv/backend/internal/platform/http/middleware"
	"developers-cv/backend/internal/testutil/pgtest"
)

type integrationTokenVerifierStub struct {
	err error
}

func (s *integrationTokenVerifierStub) Verify(context.Context, string, string) error {
	return s.err
}

func newAnalyticsPublicLinkQuery(t *testing.T, linkReader *resumelinkreaders.ResumeLinkReader, resumeReader *resumereaders.ResumeReader, analyticsRepository *analyticsrepository.AnalyticsRepository, uow ports.UnitOfWork) *queries.PublicResumeLinkQueryHandler {
	t.Helper()
	visitorIdentifier, err := analyticsadapters.NewVisitorIdentifier("integration-analytics-secret-32-bytes-minimum", false)
	require.NoError(t, err)
	subjects := analyticsaccess.NewReader(linkReader, resumeReader, resumeReader, nil, nil)
	visit := analyticscommands.NewCreatePublicVisitHandler(
		analyticsRepository,
		uow,
		visitorIdentifier,
		analyticsadapters.NewTrafficClassifier(5, time.Minute),
		analyticsadapters.NewLayoutSnapshotBuilder(),
		subjects,
	)
	return queries.NewPublicResumeLinkQueryHandlerWithVisit(
		queries.NewActivePublicResumeLinkResolver(linkReader),
		resumequeries.NewPublicResumeQueryHandler(resumeReader),
		resumeReader,
		visit,
		func(failure queries.PublicAnalyticsFailure) {
			t.Logf("public analytics stage %s failed: %v", failure.Stage, failure.Err)
		},
	)
}

func TestCreateResumeLinkIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn := pgtest.NewDatabase(t)
	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	userID := uuid.New()
	resumeID := uuid.New()
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	_, err = conn.Exec(ctx, `INSERT INTO users.users (id) VALUES ($1)`, userID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())`, resumeID, userID, "Test Resume", "en", "draft", 1)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, $2, $3, $4, $5)`, resumeID, "Hello", "from", "", "test")
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	linkRepository := resumelinkrepository.NewResumeLinkRepository(gormDB)
	create := commands.NewCreateResumeLinkHandler(linkRepository, resumeReader, database.NewUnitOfWork(gormDB))
	handler := NewHandler(HandlerDependencies{Create: create})

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	RegisterRoutes(handler, engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/link", bytes.NewBufferString(`{"resumeId":"`+resumeID.String()+`","name":"HH"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var count int
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_link.resume_links WHERE resume_id = $1 AND name = 'HH'`, resumeID).Scan(&count))
	require.Equal(t, 1, count)
}

func TestPublicResumeThroughLinkIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn := pgtest.NewDatabase(t)
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

	_, err = conn.Exec(ctx, `INSERT INTO users.users (id) VALUES ($1)`, userID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())`, resumeID, userID, "Public Resume", "en", "opened", 1)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())`, draftResumeID, userID, "Draft Resume", "en", "draft", 1)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, $2, $3, $4, $5)`, resumeID, "Hello", "from", "", "public")
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_link.resume_links (id, resume_id, name, status, version, created_at, updated_at) VALUES ($1, $2, 'HH', 'active', 1, now(), now()), ($3, $2, 'Old', 'inactive', 1, now(), now()), ($4, $5, 'Draft', 'active', 1, now(), now())`, activeLinkID, resumeID, inactiveLinkID, draftLinkID, draftResumeID)
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	linkReader := resumelinkreaders.NewResumeLinkReader(gormDB)
	publicResume := resumequeries.NewPublicResumeQueryHandler(resumeReader)
	publicLink := queries.NewPublicResumeLinkQueryHandler(queries.NewActivePublicResumeLinkResolver(linkReader), publicResume)
	handler := handlers.NewPublicResumeLinkHandler(publicLink)

	engine := gin.New()
	api := engine.Group("/api/v1")
	RegisterRoutes(NewHandler(HandlerDependencies{Public: handler}), api, api, api.Group("/lk"), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+activeLinkID.String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"title":"Public Resume"`)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+inactiveLinkID.String(), nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusGone, w.Code)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+resumeID.String(), nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+draftLinkID.String(), nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

}

func TestContactViewThroughLinkIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn := pgtest.NewDatabase(t)
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

	_, err = conn.Exec(ctx, `INSERT INTO users.users (id) VALUES ($1)`, userID)
	require.NoError(t, err)
	for _, item := range []struct {
		id, status string
	}{
		{resumeID.String(), "opened"},
		{draftResumeID.String(), "draft"},
		{closedResumeID.String(), "closed"},
	} {
		_, err = conn.Exec(ctx, `INSERT INTO resume.resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())`, uuid.MustParse(item.id), userID, item.status+" Resume", "en", item.status, 1)
		require.NoError(t, err)
	}
	_, err = conn.Exec(ctx, `INSERT INTO resume.resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, $2, $3, $4, $5)`, resumeID, "Hello", "from", "", "public")
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resume_contacts (id, resume_id, phone, email) VALUES ($1, $2, $3, $4)`, contactsID, resumeID, "+33123456789", "private@example.com")
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resume_contact_hidden (contacts_id, value, position) VALUES ($1, 'email', 0)`, contactsID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_link.resume_links (id, resume_id, name, status, version, created_at, updated_at) VALUES ($1, $2, 'Active', 'active', 1, now(), now()), ($3, $2, 'Inactive', 'inactive', 1, now(), now()), ($4, $5, 'Draft', 'active', 1, now(), now()), ($6, $7, 'Closed', 'active', 1, now(), now())`, activeLinkID, resumeID, inactiveLinkID, draftLinkID, draftResumeID, closedLinkID, closedResumeID)
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	linkReader := resumelinkreaders.NewResumeLinkReader(gormDB)
	analyticsRepository := analyticsrepository.NewAnalyticsRepository(gormDB)
	publicLink := newAnalyticsPublicLinkQuery(t, linkReader, resumeReader, analyticsRepository, database.NewUnitOfWork(gormDB))
	contactView := resumequeries.NewContactViewQueryHandler(resumeReader)
	downloadContacts := resumequeries.NewDownloadContactsQueryHandler(resumeReader)
	contactEventRecorder := analyticscommands.NewRecordPublicInteractionEventHandler(analyticsRepository, analyticsRepository, database.NewUnitOfWork(gormDB))
	contactThroughLink := queries.NewContactViewThroughLinkQueryHandlerWithAnalytics(
		queries.NewActivePublicResumeLinkResolver(linkReader),
		contactView,
		downloadContacts,
		contactEventRecorder,
		func(err error) { t.Errorf("record contact view analytics: %v", err) },
	)
	verifier := &integrationTokenVerifierStub{}
	contactHandler := handlers.NewPublicContactViewHandler(contactThroughLink, verifier)

	engine := gin.New()
	api := engine.Group("/api/v1")
	linkHandler := NewHandler(HandlerDependencies{Public: handlers.NewPublicResumeLinkHandler(publicLink), ContactView: contactHandler})
	RegisterRoutes(linkHandler, api, api, api.Group("/lk"), nil)

	getPublic := func(linkID uuid.UUID) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+linkID.String(), nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}
	publicResponse := getPublic(activeLinkID)
	require.Equal(t, http.StatusOK, publicResponse.Code)
	viewID := publicResponse.Header().Get("X-View-Id")
	require.NotEmpty(t, viewID)

	postContact := func(linkID uuid.UUID, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+linkID.String()+"/contact-views", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}

	verifier.err = ports.ErrTokenVerificationFailed
	response := postContact(activeLinkID, `{"contact":"email","turnstileToken":"rejected-token"}`)
	require.Equal(t, http.StatusForbidden, response.Code)
	require.NotContains(t, response.Body.String(), "private@example.com")
	require.NotContains(t, response.Body.String(), "rejected-token")

	verifier.err = nil
	response = postContact(activeLinkID, `{"contact":"email","turnstileToken":"test-token","viewId":"`+viewID+`"}`)
	require.Equal(t, http.StatusOK, response.Code)
	require.Contains(t, response.Body.String(), `"value":"private@example.com"`)
	require.NotContains(t, response.Body.String(), "test-token")

	response = postContact(activeLinkID, `{"initiator":"download","turnstileToken":"test-token","viewId":"`+viewID+`"}`)
	require.Equal(t, http.StatusOK, response.Code)
	require.JSONEq(t, `{"contacts":[{"contact":"email","value":"private@example.com"},{"contact":"phone","value":"+33123456789"}]}`, response.Body.String())

	response = postContact(activeLinkID, `{"contact":"phone","turnstileToken":"test-token"}`)
	require.Equal(t, http.StatusConflict, response.Code)

	response = postContact(inactiveLinkID, `{"contact":"email","turnstileToken":"test-token"}`)
	require.Equal(t, http.StatusGone, response.Code)

	response = postContact(draftLinkID, `{"contact":"email","turnstileToken":"test-token"}`)
	require.Equal(t, http.StatusNotFound, response.Code)

	response = postContact(closedLinkID, `{"contact":"email","turnstileToken":"test-token"}`)
	require.Equal(t, http.StatusNotFound, response.Code)

	response = postContact(uuid.New(), `{"contact":"email","turnstileToken":"test-token"}`)
	require.Equal(t, http.StatusNotFound, response.Code)

	response = postContact(resumeID, `{"contact":"email","turnstileToken":"test-token"}`)
	require.Equal(t, http.StatusNotFound, response.Code)

	var contactEvents int
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM analytics.analytics_events WHERE link_id = $1`, activeLinkID).Scan(&contactEvents))
	require.Equal(t, 1, contactEvents)
}

func TestLinkStatisticsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

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

	_, err = conn.Exec(ctx, `INSERT INTO users.users (id) VALUES ($1)`, userID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())`, resumeID, userID, "Statistics Resume", "en", "opened", 1)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, $2, $3, $4, $5)`, resumeID, "Hello", "from", "", "statistics")
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resume_contacts (id, resume_id, phone, email) VALUES ($1, $2, $3, $4)`, contactsID, resumeID, "+33123456789", "private@example.com")
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resume_contact_hidden (contacts_id, value, position) VALUES ($1, 'email', 0)`, contactsID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_link.resume_links (id, resume_id, name, status, version, created_at, updated_at) VALUES ($1, $2, 'HH', 'active', 1, now(), now())`, linkID, resumeID)
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	linkReader := resumelinkreaders.NewResumeLinkReader(gormDB)
	access := analyticsaccess.NewReader(linkReader, resumeReader, resumeReader, resumequeries.NewOwnedResumeSectionsQueryHandler(resumeReader), nil)
	analyticsReader := analyticsreaders.NewAnalyticsReader(gormDB, access)
	analyticsRepository := analyticsrepository.NewAnalyticsRepository(gormDB)
	publicLink := newAnalyticsPublicLinkQuery(t, linkReader, resumeReader, analyticsRepository, database.NewUnitOfWork(gormDB))
	contactView := resumequeries.NewContactViewQueryHandler(resumeReader)
	downloadContacts := resumequeries.NewDownloadContactsQueryHandler(resumeReader)
	contactEventRecorder := analyticscommands.NewRecordPublicInteractionEventHandler(analyticsRepository, analyticsRepository, database.NewUnitOfWork(gormDB))
	contactThroughLink := queries.NewContactViewThroughLinkQueryHandlerWithAnalytics(
		queries.NewActivePublicResumeLinkResolver(linkReader),
		contactView,
		downloadContacts,
		contactEventRecorder,
		func(err error) { t.Errorf("record contact view analytics: %v", err) },
	)
	statistics := analyticsqueries.NewLinkStatisticsHandler(access, analyticsReader)

	publicHandler := handlers.NewPublicResumeLinkHandler(publicLink)
	verifier := &integrationTokenVerifierStub{}
	contactHandler := handlers.NewPublicContactViewHandler(contactThroughLink, verifier)
	analyticsHandler := analyticshttp.NewHandler(analyticshttp.HandlerDependencies{
		Ingest:     analyticscommands.NewIngestAnalyticsBatchHandler(analyticsRepository, database.NewUnitOfWork(gormDB), access, nil),
		Statistics: statistics,
	})
	linkHandler := NewHandler(HandlerDependencies{
		Public:      publicHandler,
		ContactView: contactHandler,
	})

	engine := gin.New()
	api := engine.Group("/api/v1")
	lk := api.Group("/lk")
	lk.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	RegisterRoutes(linkHandler, api, api, lk, nil)
	analyticshttp.RegisterRoutes(analyticsHandler, api, lk, nil)

	getPublic := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+linkID.String(), nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/140.0 Safari/537.36")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}
	firstPublic := getPublic()
	require.Equal(t, http.StatusOK, firstPublic.Code)
	viewID := firstPublic.Header().Get("X-View-Id")
	require.NotEmpty(t, viewID)
	postContact := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+linkID.String()+"/contact-views", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}

	require.Equal(t, http.StatusOK, getPublic().Code)
	require.Equal(t, http.StatusOK, postContact(`{"contact":"email","turnstileToken":"test-token","viewId":"`+viewID+`"}`).Code)
	require.Equal(t, http.StatusOK, postContact(`{"initiator":"download","turnstileToken":"test-token","viewId":"`+viewID+`"}`).Code)

	eventBody := `{"viewId":"` + viewID + `","seq":1,"final":false,"state":{"activeMs":1000},"events":[{"id":"` + uuid.NewString() + `","type":"download","at":1000}],"clicks":[]}`
	eventRequest := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+linkID.String()+"/events", strings.NewReader(eventBody))
	eventRequest.Header.Set("Content-Type", "application/json")
	eventResponse := httptest.NewRecorder()
	engine.ServeHTTP(eventResponse, eventRequest)
	require.Equal(t, http.StatusNoContent, eventResponse.Code)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+linkID.String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var stats struct {
		Views        int `json:"views"`
		ContactViews int `json:"contactViews"`
		Downloads    int `json:"downloads"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &stats))
	require.Equal(t, 2, stats.Views)
	require.Equal(t, 1, stats.ContactViews)
	require.Equal(t, 1, stats.Downloads)
}

func TestResumeLinkLifecycleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn := pgtest.NewDatabase(t)
	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	userID := uuid.New()
	resumeID := uuid.New()
	linkID := uuid.New()
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	_, err = conn.Exec(ctx, `INSERT INTO users.users (id) VALUES ($1)`, userID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resumes (id, user_id, title, language, status, version, created_at, updated_at) VALUES ($1, $2, 'Lifecycle Resume', 'en', 'opened', 1, now(), now())`, resumeID, userID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, 'Hello', 'from', '', 'lifecycle')`, resumeID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume_link.resume_links (id, resume_id, name, status, version, created_at, updated_at) VALUES ($1, $2, 'HH', 'active', 1, now(), now())`, linkID, resumeID)
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	linkReader := resumelinkreaders.NewResumeLinkReader(gormDB)
	linkRepository := resumelinkrepository.NewResumeLinkRepository(gormDB)
	analyticsRepository := analyticsrepository.NewAnalyticsRepository(gormDB)
	publicLink := newAnalyticsPublicLinkQuery(t, linkReader, resumeReader, analyticsRepository, database.NewUnitOfWork(gormDB))
	lifecycle := commands.NewReplaceStatusResumeLinkHandler(linkRepository, resumeReader, database.NewUnitOfWork(gormDB))

	linkHandler := NewHandler(HandlerDependencies{
		Public:    handlers.NewPublicResumeLinkHandler(publicLink),
		Lifecycle: lifecycle,
	})
	engine := gin.New()
	api := engine.Group("/api/v1")
	lk := api.Group("/lk")
	lk.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	RegisterRoutes(linkHandler, api, api, lk, nil)

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
	require.Equal(t, http.StatusGone, getPublic().Code)

	require.Equal(t, http.StatusPreconditionFailed, changeStatus("active", 1).Code)
	require.Equal(t, http.StatusNoContent, changeStatus("active", 2).Code)
	require.Equal(t, http.StatusOK, getPublic().Code)

	var eventCount int
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM analytics.analytics_events WHERE link_id = $1`, linkID).Scan(&eventCount))
	require.Equal(t, 0, eventCount)
}

func TestResumeLinksFullIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

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

	_, err = conn.Exec(ctx, `INSERT INTO users.users (id) VALUES ($1), ($2)`, userA, userB)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `
		INSERT INTO resume.resumes (id, user_id, title, language, status, version, created_at, updated_at)
		VALUES
			($1, $4, 'Resume A', 'en', 'opened', 1, now(), now()),
			($2, $4, 'Resume B', 'en', 'opened', 1, now(), now()),
			($3, $5, 'Resume C', 'en', 'opened', 1, now(), now())
	`, resumeA, resumeB, resumeC, userA, userB)
	require.NoError(t, err)
	for _, resumeID := range []uuid.UUID{resumeA, resumeB, resumeC} {
		_, err = conn.Exec(ctx, `INSERT INTO resume.resume_overviews (resume_id, greeting_first, greeting_second, greeting_highlight, greeting_third) VALUES ($1, 'Hello', 'from', '', 'integration')`, resumeID)
		require.NoError(t, err)
	}

	_, err = conn.Exec(ctx, `INSERT INTO resume.resume_contacts (id, resume_id, phone, email) VALUES ($1, $3, '+33111111111', 'a@example.com'), ($2, $4, '+33222222222', 'b@example.com'), ($5, $6, '+33333333333', 'github@example.com')`, contactsA, contactsC, resumeA, resumeC, contactsB, resumeB)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO resume.resume_contact_hidden (contacts_id, value, position) VALUES ($1, 'email', 0), ($2, 'email', 0), ($3, 'email', 0)`, contactsA, contactsC, contactsB)
	require.NoError(t, err)

	resumeReader := resumereaders.NewResumeReader(gormDB)
	resumeRepository := resumerepository.NewResumeRepository(gormDB)
	linkReader := resumelinkreaders.NewResumeLinkReader(gormDB)
	access := analyticsaccess.NewReader(linkReader, resumeReader, resumeReader, resumequeries.NewOwnedResumeSectionsQueryHandler(resumeReader), nil)
	analyticsReader := analyticsreaders.NewAnalyticsReader(gormDB, access)
	linkRepository := resumelinkrepository.NewResumeLinkRepository(gormDB)
	analyticsRepository := analyticsrepository.NewAnalyticsRepository(gormDB)
	publicLink := newAnalyticsPublicLinkQuery(t, linkReader, resumeReader, analyticsRepository, database.NewUnitOfWork(gormDB))
	contactView := resumequeries.NewContactViewQueryHandler(resumeReader)
	downloadContacts := resumequeries.NewDownloadContactsQueryHandler(resumeReader)
	contactEventRecorder := analyticscommands.NewRecordPublicInteractionEventHandler(analyticsRepository, analyticsRepository, database.NewUnitOfWork(gormDB))
	contactThroughLink := queries.NewContactViewThroughLinkQueryHandlerWithAnalytics(
		queries.NewActivePublicResumeLinkResolver(linkReader),
		contactView,
		downloadContacts,
		contactEventRecorder,
		func(err error) { t.Errorf("record contact view analytics: %v", err) },
	)
	statistics := analyticsqueries.NewLinkStatisticsHandler(access, analyticsReader)
	createLink := commands.NewCreateResumeLinkHandler(linkRepository, resumeReader, database.NewUnitOfWork(gormDB))
	listLinks := queries.NewListResumeLinksQueryHandler(linkReader, resumeReader)
	deleteLink := commands.NewDeleteResumeLinkHandler(linkRepository, resumeReader, analyticsRepository, database.NewUnitOfWork(gormDB))
	lifecycle := commands.NewReplaceStatusResumeLinkHandler(linkRepository, resumeReader, database.NewUnitOfWork(gormDB))
	resumeStatus := resumecommands.NewChangeResumeStatusHandler(resumeRepository, database.NewUnitOfWork(gormDB))
	resumeLinkCleanup := commands.NewDeleteResumeLinksHandler(linkRepository, analyticsRepository)
	deleteResume := resumecommands.NewDeleteResumeHandler(resumeRepository, resumeLinkCleanup, database.NewUnitOfWork(gormDB))

	linkHandler := NewHandler(HandlerDependencies{
		Create:      createLink,
		List:        listLinks,
		Delete:      deleteLink,
		Public:      handlers.NewPublicResumeLinkHandler(publicLink),
		ContactView: handlers.NewPublicContactViewHandler(contactThroughLink, &integrationTokenVerifierStub{}),
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
	resumehttp.RegisterRoutes(resumeHandler, lk)
	RegisterRoutes(linkHandler, api, api, lk, nil)
	analyticshttp.RegisterRoutes(analyticshttp.NewHandler(analyticshttp.HandlerDependencies{Statistics: statistics}), api, lk, nil)

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
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+linkID.String(), nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/140.0 Safari/537.36")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}
	postContact := func(linkID uuid.UUID, viewID string) *httptest.ResponseRecorder {
		return do(http.MethodPost, "/api/v1/cv/"+linkID.String()+"/contact-views", `{"contact":"email","turnstileToken":"test-token","viewId":"`+viewID+`"}`)
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
	require.Equal(t, http.StatusOK, postContact(linkAHH, publicAHH.Header().Get("X-View-Id")).Code)
	require.Equal(t, http.StatusOK, postContact(linkAGitHub, publicAGitHub.Header().Get("X-View-Id")).Code)
	for _, item := range []struct {
		id           uuid.UUID
		wantViews    int
		wantContacts int
	}{
		{linkAHH, 3, 1},
		{linkATelegram, 2, 0},
		{linkAGitHub, 1, 1},
	} {
		response := getStatistics(item.id)
		require.Equal(t, http.StatusOK, response.Code)
		var stats struct {
			Views        int `json:"views"`
			ContactViews int `json:"contactViews"`
		}
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &stats))
		require.Equal(t, item.wantViews, stats.Views)
		require.Equal(t, item.wantContacts, stats.ContactViews)
	}

	currentUser = userB
	response := getStatistics(linkBHH)
	require.Equal(t, http.StatusOK, response.Code)
	var linkBStats struct {
		Views        int `json:"views"`
		ContactViews int `json:"contactViews"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &linkBStats))
	require.Equal(t, 1, linkBStats.Views)
	require.Equal(t, 0, linkBStats.ContactViews)
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
	require.Equal(t, http.StatusGone, getPublic(linkAHH).Code)
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

	response = do(http.MethodDelete, "/api/v1/lk/resume/"+resumeA.String(), `{"version":3}`)
	require.Equal(t, http.StatusNoContent, response.Code)

	var linksAfterDelete, eventsAfterDelete int
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM resume_link.resume_links WHERE resume_id = $1`, resumeA).Scan(&linksAfterDelete))
	require.NoError(t, conn.QueryRow(ctx, `SELECT count(*) FROM analytics.analytics_events WHERE link_id IN ($1, $2)`, linkAHH, linkATelegram).Scan(&eventsAfterDelete))
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
