package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
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
		resumehttp.NewHandler(nil, nil, nil, nil, nil),
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
		resumehttp.NewHandler(nil, nil, nil, nil, nil),
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
		resumehttp.NewHandler(nil, nil, nil, nil, nil),
		func(context.Context, *http.Request) (uuid.UUID, error) {
			return uuid.Nil, middleware.ErrUnauthenticated
		},
	)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/lk/resume/10000000-0000-0000-0000-000000000001", nil)
	w := httptest.NewRecorder()
	router.Handler().ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Contains(t, w.Body.String(), "UNAUTHORIZED")
}

type routerGetResumeReader struct{}

func (routerGetResumeReader) GetResume(
	_ context.Context,
	query resumequeries.GetResumeQuery,
) (resumequeries.Resume, error) {
	return resumequeries.Resume{ID: query.ID.String(), Name: "Public Resume"}, nil
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
	handler := resumehttp.NewHandler(resumequeries.NewGetResumeHandler(routerGetResumeReader{}), nil, nil, nil, nil)
	router, err := NewRouter(
		cfg,
		log,
		health.NewService(),
		nil,
		handler,
		func(context.Context, *http.Request) (uuid.UUID, error) {
			return uuid.Nil, middleware.ErrUnauthenticated
		},
	)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/cv/10000000-0000-0000-0000-000000000001", nil)
	w := httptest.NewRecorder()
	router.Handler().ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "Public Resume")
}
