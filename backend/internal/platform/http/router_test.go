package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	resumelinkhttp "developers-cv/backend/internal/modules/resume_link/transport/http"
	resumelinkpublichandler "developers-cv/backend/internal/modules/resume_link/transport/http/public"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
	userhttp "developers-cv/backend/internal/modules/user/transport/http"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/health"
	"developers-cv/backend/internal/platform/http/middleware"
)

func TestRouterHealthAndRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		App: config.AppConfig{
			Env: "test",
		},
		CORS: config.CORSConfig{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowCredentials: true,
		},
		HTTP: config.HTTPConfig{
			TLS: config.TLSConfig{Enabled: true},
		},
	}
	log := zaptest.NewLogger(t)
	router, err := NewRouter(
		cfg,
		log,
		health.NewService(),
		nil,
		resumehttp.NewHandler(resumehttp.HandlerDependencies{}),
		userhttp.NewHandler(nil, nil),
		testResumeLinkHandler(),
		fakeAuthenticate,
	)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.Handler().ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Header().Get("X-Request-ID"))
	require.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	require.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
}

func TestRouterReturnsNotFoundForUnknownRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		App:  config.AppConfig{Env: "test"},
		HTTP: config.HTTPConfig{},
		CORS: config.CORSConfig{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowCredentials: true,
		},
	}
	log := zaptest.NewLogger(t)
	router, err := NewRouter(
		cfg,
		log,
		health.NewService(),
		nil,
		resumehttp.NewHandler(resumehttp.HandlerDependencies{}),
		userhttp.NewHandler(nil, nil),
		testResumeLinkHandler(),
		fakeAuthenticate,
	)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/v1/unknown", nil)
	w := httptest.NewRecorder()
	router.Handler().ServeHTTP(w, req)

	require.Equal(t, 404, w.Code)
	require.Contains(t, w.Body.String(), "NOT_FOUND")
}

func fakeAuthenticate(context.Context, *http.Request) (uuid.UUID, error) {
	return uuid.New(), nil
}

func testResumeLinkHandler() *resumelinkhttp.Handler {
	return resumelinkhttp.NewHandler(resumelinkhttp.HandlerDependencies{
		Public: resumelinkpublichandler.NewHandler(nil),
	})
}

func TestRouterRequiresAuthenticationForPrivateResumeRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		App: config.AppConfig{Env: "test"},
		CORS: config.CORSConfig{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowCredentials: true,
		},
		HTTP: config.HTTPConfig{},
	}
	log := zaptest.NewLogger(t)
	router, err := NewRouter(
		cfg,
		log,
		health.NewService(),
		nil,
		resumehttp.NewHandler(resumehttp.HandlerDependencies{}),
		userhttp.NewHandler(nil, nil),
		testResumeLinkHandler(),
		func(context.Context, *http.Request) (uuid.UUID, error) {
			return uuid.Nil, middleware.ErrUnauthenticated
		},
	)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/overview", nil)
	w := httptest.NewRecorder()
	router.Handler().ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Contains(t, w.Body.String(), "UNAUTHORIZED")
}

type routerGetProfileReader struct{}

func (routerGetProfileReader) GetProfile(_ context.Context, query userqueries.GetProfileQuery) (userqueries.Profile, error) {
	return userqueries.Profile{
		ID: query.UserID.String(),
		RU: userqueries.LocalizedProfile{Name: stringPtr("Anton"), Surname: stringPtr("User")},
	}, nil
}

func stringPtr(value string) *string { return &value }

func TestRouterRegistersProfileRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		App: config.AppConfig{Env: "test"},
		CORS: config.CORSConfig{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowCredentials: true,
		},
		HTTP: config.HTTPConfig{},
	}
	log := zaptest.NewLogger(t)
	userID := uuid.New()
	router, err := NewRouter(
		cfg,
		log,
		health.NewService(),
		nil,
		resumehttp.NewHandler(resumehttp.HandlerDependencies{}),
		userhttp.NewHandler(userqueries.NewGetProfileHandler(routerGetProfileReader{}), nil),
		testResumeLinkHandler(),
		func(context.Context, *http.Request) (uuid.UUID, error) {
			return userID, nil
		},
	)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/profile", nil)
	w := httptest.NewRecorder()
	router.Handler().ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"userId":"`+userID.String()+`"`)
}

func TestRouterRegistersCanonicalResumeRoutesOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		App: config.AppConfig{Env: "test"},
		CORS: config.CORSConfig{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowCredentials: true,
		},
		HTTP: config.HTTPConfig{},
	}
	log := zaptest.NewLogger(t)
	router, err := NewRouter(
		cfg,
		log,
		health.NewService(),
		nil,
		resumehttp.NewHandler(resumehttp.HandlerDependencies{}),
		userhttp.NewHandler(nil, nil),
		testResumeLinkHandler(),
		fakeAuthenticate,
	)
	require.NoError(t, err)

	expected := map[string]struct{}{
		"GET /api/v1/lk/resumes":                       {},
		"POST /api/v1/lk/resume":                       {},
		"GET /api/v1/lk/resume/:resumeId":              {},
		"PUT /api/v1/lk/resume/:resumeId":              {},
		"PUT /api/v1/lk/resume/:resumeId/status":       {},
		"PUT /api/v1/lk/resume/:resumeId/overview":     {},
		"DELETE /api/v1/lk/resume/:resumeId":           {},
		"PUT /api/v1/lk/resume/:resumeId/profile":      {},
		"PUT /api/v1/lk/resume/:resumeId/contacts":     {},
		"PUT /api/v1/lk/resume/:resumeId/skills":       {},
		"PUT /api/v1/lk/resume/:resumeId/languages":    {},
		"PUT /api/v1/lk/resume/:resumeId/technologies": {},
		"PUT /api/v1/lk/resume/:resumeId/experiences":  {},
		"PUT /api/v1/lk/resume/:resumeId/educations":   {},
		"PUT /api/v1/lk/resume/:resumeId/courses":      {},
		"PUT /api/v1/lk/resume/:resumeId/activities":   {},
		"PUT /api/v1/lk/resume/:resumeId/projects":     {},
		"GET /api/v1/cv/:linkId":                       {},
		"POST /api/v1/cv/:linkId/contact-views":        {},
		"GET /api/v1/lk/links":                         {},
		"POST /api/v1/lk/link":                         {},
		"DELETE /api/v1/lk/link/:linkId":               {},
		"PUT /api/v1/lk/link/:linkId/status":           {},
		"GET /api/v1/lk/link/:linkId/statistics":       {},
	}

	actual := make(map[string]struct{})
	for _, route := range router.engine.Routes() {
		if strings.HasPrefix(route.Path, "/api/v1/lk/resume") || strings.HasPrefix(route.Path, "/api/v1/lk/link") || strings.HasPrefix(route.Path, "/api/v1/cv/") {
			actual[route.Method+" "+route.Path] = struct{}{}
		}
	}

	require.Equal(t, expected, actual)
	require.NotContains(t, actual, "PATCH /api/v1/lk/resume/:resumeId")
}

func TestRouterDoesNotExposeLegacyResumeAliases(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		App: config.AppConfig{Env: "test"},
		CORS: config.CORSConfig{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowCredentials: true,
		},
		HTTP: config.HTTPConfig{},
	}
	log := zaptest.NewLogger(t)
	router, err := NewRouter(
		cfg,
		log,
		health.NewService(),
		nil,
		resumehttp.NewHandler(resumehttp.HandlerDependencies{}),
		userhttp.NewHandler(nil, nil),
		testResumeLinkHandler(),
		fakeAuthenticate,
	)
	require.NoError(t, err)

	for _, methodPath := range []struct {
		method string
		path   string
	}{
		{method: http.MethodPost, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001"},
		{method: http.MethodGet, path: "/api/v1/lk/resumes/10000000-0000-0000-0000-000000000001"},
		{method: http.MethodPost, path: "/api/v1/lk/resumes/10000000-0000-0000-0000-000000000001"},
		{method: http.MethodPut, path: "/api/v1/lk/resumes/10000000-0000-0000-0000-000000000001"},
		{method: http.MethodDelete, path: "/api/v1/lk/resumes/10000000-0000-0000-0000-000000000001"},
		{method: http.MethodGet, path: "/cv/10000000-0000-0000-0000-000000000001"},
		{method: http.MethodPost, path: "/api/v1/cv/10000000-0000-0000-0000-000000000001"},
	} {
		req := httptest.NewRequest(methodPath.method, methodPath.path, strings.NewReader(""))
		w := httptest.NewRecorder()
		router.Handler().ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code, methodPath.method+" "+methodPath.path)
	}
}

func TestPrivateResumeRoutesRejectUnauthenticatedRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		App:  config.AppConfig{Env: "test"},
		CORS: config.CORSConfig{AllowedOrigins: []string{"http://localhost:5173"}, AllowCredentials: true},
		HTTP: config.HTTPConfig{},
	}
	log := zaptest.NewLogger(t)
	paths := []struct {
		method string
		path   string
	}{
		{method: http.MethodPut, path: "/api/v1/lk/profile"},
		{method: http.MethodGet, path: "/api/v1/lk/resumes"},
		{method: http.MethodPost, path: "/api/v1/lk/link"},
		{method: http.MethodGet, path: "/api/v1/lk/link/10000000-0000-0000-0000-000000000001/statistics"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/overview"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/profile"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/contacts"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/skills"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/languages"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/technologies"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/experiences"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/educations"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/courses"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/activities"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/projects"},
	}

	for _, route := range paths {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			router, err := NewRouter(
				cfg,
				log,
				health.NewService(),
				nil,
				resumehttp.NewHandler(resumehttp.HandlerDependencies{}),
				userhttp.NewHandler(nil, nil),
				testResumeLinkHandler(),
				func(context.Context, *http.Request) (uuid.UUID, error) {
					return uuid.Nil, middleware.ErrUnauthenticated
				},
			)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.Handler().ServeHTTP(w, httptest.NewRequest(route.method, route.path, nil))
			require.Equal(t, http.StatusUnauthorized, w.Code)
			require.Contains(t, w.Body.String(), "UNAUTHORIZED")
		})
	}
}
func TestPrivateResumeRoutesRejectMalformedResumeID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		App:  config.AppConfig{Env: "test"},
		CORS: config.CORSConfig{AllowedOrigins: []string{"http://localhost:5173"}, AllowCredentials: true},
		HTTP: config.HTTPConfig{},
	}
	log := zaptest.NewLogger(t)
	paths := []struct {
		method string
		path   string
	}{
		{method: http.MethodPut, path: "/api/v1/lk/resume/not-a-uuid/overview"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/not-a-uuid/profile"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/not-a-uuid/contacts"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/not-a-uuid/skills"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/not-a-uuid/languages"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/not-a-uuid/technologies"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/not-a-uuid/experiences"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/not-a-uuid/educations"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/not-a-uuid/courses"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/not-a-uuid/activities"},
		{method: http.MethodPut, path: "/api/v1/lk/resume/not-a-uuid/projects"},
	}

	for _, route := range paths {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			router, err := NewRouter(
				cfg,
				log,
				health.NewService(),
				nil,
				resumehttp.NewHandler(resumehttp.HandlerDependencies{}),
				userhttp.NewHandler(nil, nil),
				testResumeLinkHandler(),
				fakeAuthenticate,
			)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.Handler().ServeHTTP(w, httptest.NewRequest(route.method, route.path, nil))
			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "BAD_REQUEST")
		})
	}
}

type routerHealthChecker struct {
	name string
	err  error
}

func (c routerHealthChecker) Name() string {
	return c.name
}

func (c routerHealthChecker) Check(context.Context) error {
	return c.err
}

func TestRouterReadinessReturnsUnavailableWhenDependencyFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		App:  config.AppConfig{Env: "test"},
		HTTP: config.HTTPConfig{},
		CORS: config.CORSConfig{AllowedOrigins: []string{"http://localhost:5173"}, AllowCredentials: true},
	}
	log := zaptest.NewLogger(t)
	router, err := NewRouter(
		cfg,
		log,
		health.NewService(routerHealthChecker{name: "postgres", err: context.DeadlineExceeded}),
		nil,
		resumehttp.NewHandler(resumehttp.HandlerDependencies{}),
		userhttp.NewHandler(nil, nil),
		testResumeLinkHandler(),
		fakeAuthenticate,
	)
	require.NoError(t, err)

	response := httptest.NewRecorder()
	router.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/ready", nil))

	require.Equal(t, http.StatusServiceUnavailable, response.Code)
	require.Contains(t, response.Body.String(), `"status":"not_ready"`)
	require.Contains(t, response.Body.String(), `"postgres":"failed"`)
}

func TestRouterReadinessReturnsReadyWhenDependenciesAreHealthy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		App:  config.AppConfig{Env: "test"},
		HTTP: config.HTTPConfig{},
		CORS: config.CORSConfig{AllowedOrigins: []string{"http://localhost:5173"}, AllowCredentials: true},
	}
	log := zaptest.NewLogger(t)
	router, err := NewRouter(
		cfg,
		log,
		health.NewService(routerHealthChecker{name: "postgres"}),
		nil,
		resumehttp.NewHandler(resumehttp.HandlerDependencies{}),
		userhttp.NewHandler(nil, nil),
		testResumeLinkHandler(),
		fakeAuthenticate,
	)
	require.NoError(t, err)

	response := httptest.NewRecorder()
	router.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/ready", nil))

	require.Equal(t, http.StatusOK, response.Code)
	require.Contains(t, response.Body.String(), `"status":"ready"`)
	require.Contains(t, response.Body.String(), `"postgres":"ok"`)
}

func TestRouterRateLimitCoversPublicResumeLinkRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		App:  config.AppConfig{Env: "test"},
		CORS: config.CORSConfig{AllowedOrigins: []string{"http://localhost:5173"}, AllowCredentials: true},
		HTTP: config.HTTPConfig{
			RateLimit: config.RateLimitConfig{Enabled: true, Rate: 1, Burst: 1},
		},
	}
	log := zaptest.NewLogger(t)
	router, err := NewRouter(
		cfg,
		log,
		health.NewService(),
		nil,
		resumehttp.NewHandler(resumehttp.HandlerDependencies{}),
		userhttp.NewHandler(nil, nil),
		testResumeLinkHandler(),
		fakeAuthenticate,
	)
	require.NoError(t, err)

	path := "/api/v1/cv/10000000-0000-0000-0000-000000000001"
	first := httptest.NewRecorder()
	router.Handler().ServeHTTP(first, httptest.NewRequest(http.MethodGet, path, nil))
	second := httptest.NewRecorder()
	router.Handler().ServeHTTP(second, httptest.NewRequest(http.MethodGet, path, nil))

	require.NotEqual(t, http.StatusTooManyRequests, first.Code)
	require.Equal(t, http.StatusTooManyRequests, second.Code)
	require.Contains(t, second.Body.String(), "TOO_MANY_REQUESTS")
}
