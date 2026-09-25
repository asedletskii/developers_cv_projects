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

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	publichttp "developers-cv/backend/internal/modules/resume/transport/http/public"
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
		resumehttp.NewHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil),
		userhttp.NewHandler(nil),
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
		resumehttp.NewHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil),
		userhttp.NewHandler(nil),
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
		resumehttp.NewHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil),
		userhttp.NewHandler(nil),
		func(context.Context, *http.Request) (uuid.UUID, error) {
			return uuid.Nil, middleware.ErrUnauthenticated
		},
	)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/overview", nil)
	w := httptest.NewRecorder()
	router.Handler().ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Contains(t, w.Body.String(), "UNAUTHORIZED")
}

type routerPublicResumeReader struct{}

func (routerPublicResumeReader) GetPublicResume(
	_ context.Context,
	query resumequeries.PublicResumeQuery,
) (resumequeries.PublicResumeView, error) {
	return resumequeries.PublicResumeView{ID: query.ID.String(), Name: "Public Resume"}, nil
}

func TestRouterAllowsPublicResumeRouteWithoutAuthentication(t *testing.T) {
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
	handler := resumehttp.NewHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, publichttp.NewHandler(resumequeries.NewPublicResumeQueryHandler(routerPublicResumeReader{})))
	router, err := NewRouter(
		cfg,
		log,
		health.NewService(),
		nil,
		handler,
		userhttp.NewHandler(nil),
		func(context.Context, *http.Request) (uuid.UUID, error) {
			return uuid.Nil, middleware.ErrUnauthenticated
		},
	)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/v1/cv/10000000-0000-0000-0000-000000000001", nil)
	w := httptest.NewRecorder()
	router.Handler().ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "Public Resume")
}

type routerGetProfileReader struct{}

func (routerGetProfileReader) GetProfile(_ context.Context, query userqueries.GetProfileQuery) (userqueries.Profile, error) {
	return userqueries.Profile{
		ID:      query.UserID.String(),
		Name:    "Anton",
		Surname: "User",
	}, nil
}

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
		resumehttp.NewHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil),
		userhttp.NewHandler(userqueries.NewGetProfileHandler(routerGetProfileReader{})),
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
		resumehttp.NewHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil),
		userhttp.NewHandler(nil),
		fakeAuthenticate,
	)
	require.NoError(t, err)

	expected := map[string]struct{}{
		"GET /api/v1/lk/resumes":                                        {},
		"POST /api/v1/lk/resume":                                        {},
		"GET /api/v1/lk/resume/:resumeId":                               {},
		"PUT /api/v1/lk/resume/:resumeId":                               {},
		"GET /api/v1/lk/resume/:resumeId/overview":                      {},
		"PUT /api/v1/lk/resume/:resumeId/overview":                      {},
		"DELETE /api/v1/lk/resume/:resumeId":                            {},
		"POST /api/v1/lk/resume/:resumeId/profile":                      {},
		"GET /api/v1/lk/resume/:resumeId/profile":                       {},
		"PUT /api/v1/lk/resume/:resumeId/profile":                       {},
		"POST /api/v1/lk/resume/:resumeId/contacts":                     {},
		"GET /api/v1/lk/resume/:resumeId/contacts":                      {},
		"PUT /api/v1/lk/resume/:resumeId/contacts":                      {},
		"POST /api/v1/lk/resume/:resumeId/skills":                       {},
		"GET /api/v1/lk/resume/:resumeId/skills":                        {},
		"PUT /api/v1/lk/resume/:resumeId/skills":                        {},
		"POST /api/v1/lk/resume/:resumeId/languages":                    {},
		"GET /api/v1/lk/resume/:resumeId/languages":                     {},
		"GET /api/v1/lk/resume/:resumeId/languages/:languageId":         {},
		"PUT /api/v1/lk/resume/:resumeId/languages/:languageId":         {},
		"DELETE /api/v1/lk/resume/:resumeId/languages/:languageId":      {},
		"POST /api/v1/lk/resume/:resumeId/technologies":                 {},
		"GET /api/v1/lk/resume/:resumeId/technologies":                  {},
		"GET /api/v1/lk/resume/:resumeId/technologies/:technologyId":    {},
		"PUT /api/v1/lk/resume/:resumeId/technologies/:technologyId":    {},
		"DELETE /api/v1/lk/resume/:resumeId/technologies/:technologyId": {},
		"POST /api/v1/lk/resume/:resumeId/experiences":                  {},
		"GET /api/v1/lk/resume/:resumeId/experiences":                   {},
		"GET /api/v1/lk/resume/:resumeId/experiences/:experienceId":     {},
		"PUT /api/v1/lk/resume/:resumeId/experiences/:experienceId":     {},
		"DELETE /api/v1/lk/resume/:resumeId/experiences/:experienceId":  {},
		"POST /api/v1/lk/resume/:resumeId/educations":                   {},
		"GET /api/v1/lk/resume/:resumeId/educations":                    {},
		"GET /api/v1/lk/resume/:resumeId/educations/:educationId":       {},
		"PUT /api/v1/lk/resume/:resumeId/educations/:educationId":       {},
		"DELETE /api/v1/lk/resume/:resumeId/educations/:educationId":    {},
		"POST /api/v1/lk/resume/:resumeId/courses":                      {},
		"GET /api/v1/lk/resume/:resumeId/courses":                       {},
		"GET /api/v1/lk/resume/:resumeId/courses/:courseId":             {},
		"PUT /api/v1/lk/resume/:resumeId/courses/:courseId":             {},
		"DELETE /api/v1/lk/resume/:resumeId/courses/:courseId":          {},
		"POST /api/v1/lk/resume/:resumeId/activities":                   {},
		"GET /api/v1/lk/resume/:resumeId/activities":                    {},
		"GET /api/v1/lk/resume/:resumeId/activities/:activityId":        {},
		"PUT /api/v1/lk/resume/:resumeId/activities/:activityId":        {},
		"DELETE /api/v1/lk/resume/:resumeId/activities/:activityId":     {},
		"POST /api/v1/lk/resume/:resumeId/projects":                     {},
		"GET /api/v1/lk/resume/:resumeId/projects":                      {},
		"GET /api/v1/lk/resume/:resumeId/projects/:projectId":           {},
		"PUT /api/v1/lk/resume/:resumeId/projects/:projectId":           {},
		"DELETE /api/v1/lk/resume/:resumeId/projects/:projectId":        {},
		"GET /api/v1/cv/:resumeId":                                      {},
	}

	actual := make(map[string]struct{})
	for _, route := range router.engine.Routes() {
		if strings.HasPrefix(route.Path, "/api/v1/lk/resume") || strings.HasPrefix(route.Path, "/api/v1/cv/") {
			actual[route.Method+" "+route.Path] = struct{}{}
		}
	}

	require.Equal(t, expected, actual)
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
		resumehttp.NewHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil),
		userhttp.NewHandler(nil),
		fakeAuthenticate,
	)
	require.NoError(t, err)

	for _, methodPath := range []struct {
		method string
		path   string
	}{
		// The singular resource path is canonical now, so only methods that
		// are intentionally not registered on the resource itself belong here.
		{method: http.MethodPost, path: "/api/v1/lk/resume/10000000-0000-0000-0000-000000000001"},
		{method: http.MethodGet, path: "/api/v1/lk/resumes/10000000-0000-0000-0000-000000000001"},
		{method: http.MethodPost, path: "/api/v1/lk/resumes/10000000-0000-0000-0000-000000000001"},
		{method: http.MethodPut, path: "/api/v1/lk/resumes/10000000-0000-0000-0000-000000000001"},
		{method: http.MethodDelete, path: "/api/v1/lk/resumes/10000000-0000-0000-0000-000000000001"},
		{method: http.MethodGet, path: "/cv/10000000-0000-0000-0000-000000000001"},
	} {
		req := httptest.NewRequest(methodPath.method, methodPath.path, nil)
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
	paths := []string{
		"/api/v1/lk/resumes",
		"/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/overview",
		"/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/profile",
		"/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/contacts",
		"/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/skills",
		"/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/languages",
		"/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/technologies",
		"/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/experiences",
		"/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/educations",
		"/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/courses",
		"/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/activities",
		"/api/v1/lk/resume/10000000-0000-0000-0000-000000000001/projects",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			router, err := NewRouter(
				cfg,
				log,
				health.NewService(),
				nil,
				resumehttp.NewHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil),
				userhttp.NewHandler(nil),
				func(context.Context, *http.Request) (uuid.UUID, error) {
					return uuid.Nil, middleware.ErrUnauthenticated
				},
			)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.Handler().ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
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
	paths := []string{
		"/api/v1/lk/resume/not-a-uuid/overview",
		"/api/v1/lk/resume/not-a-uuid/profile",
		"/api/v1/lk/resume/not-a-uuid/contacts",
		"/api/v1/lk/resume/not-a-uuid/skills",
		"/api/v1/lk/resume/not-a-uuid/languages",
		"/api/v1/lk/resume/not-a-uuid/technologies",
		"/api/v1/lk/resume/not-a-uuid/experiences",
		"/api/v1/lk/resume/not-a-uuid/educations",
		"/api/v1/lk/resume/not-a-uuid/courses",
		"/api/v1/lk/resume/not-a-uuid/activities",
		"/api/v1/lk/resume/not-a-uuid/projects",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			router, err := NewRouter(
				cfg,
				log,
				health.NewService(),
				nil,
				resumehttp.NewHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil),
				userhttp.NewHandler(nil),
				fakeAuthenticate,
			)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.Handler().ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "BAD_REQUEST")
		})
	}
}
