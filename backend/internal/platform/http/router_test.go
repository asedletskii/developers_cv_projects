package http

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumemock "developers-cv/backend/internal/modules/resume/infrastructure/mock"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/health"
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
		resumehttp.NewHandler(resumequeries.NewGetPersonHandler(resumemock.NewPersonReader())),
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
		resumehttp.NewHandler(resumequeries.NewGetPersonHandler(resumemock.NewPersonReader())),
	)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/v1/unknown", nil)
	w := httptest.NewRecorder()
	router.Handler().ServeHTTP(w, req)

	require.Equal(t, 404, w.Code)
	require.Contains(t, w.Body.String(), "NOT_FOUND")
}
