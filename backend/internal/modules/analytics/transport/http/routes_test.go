package http

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/analytics/application/commands"
	"developers-cv/backend/internal/modules/analytics/application/queries"
)

func TestRegisterRoutesPreservesAnalyticsRouteContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := engine.Group("/api/v1")
	handler := NewHandler(HandlerDependencies{
		Ingest: commands.NewIngestAnalyticsBatchHandler(nil, nil, nil, nil),
		Clicks: queries.NewAnalyticsClickQueryHandler(nil, nil, nil),
	})
	RegisterRoutes(handler, api, api.Group("/lk"), nil)

	registered := make(map[string]struct{}, len(engine.Routes()))
	for _, route := range engine.Routes() {
		registered[route.Method+" "+route.Path] = struct{}{}
	}

	expected := []string{
		http.MethodPost + " /api/v1/cv/:linkId/events",
		http.MethodGet + " /api/v1/lk/link/:linkId/statistics",
		http.MethodGet + " /api/v1/lk/link/:linkId/statistics/daily",
		http.MethodGet + " /api/v1/lk/link/:linkId/views",
		http.MethodGet + " /api/v1/lk/link/:linkId/statistics/sections",
		http.MethodGet + " /api/v1/lk/link/:linkId/statistics/technologies",
		http.MethodGet + " /api/v1/lk/link/:linkId/clicks",
		http.MethodGet + " /api/v1/lk/resume/:resumeId/statistics",
	}
	for _, route := range expected {
		_, ok := registered[route]
		require.Truef(t, ok, "expected route %q to be registered", route)
	}
}

func TestRegisterRoutesKeepsIngestionAndClickRoutesConditional(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := engine.Group("/api/v1")
	RegisterRoutes(NewHandler(HandlerDependencies{}), api, api.Group("/lk"), nil)

	registered := make(map[string]struct{}, len(engine.Routes()))
	for _, route := range engine.Routes() {
		registered[route.Method+" "+route.Path] = struct{}{}
	}
	_, hasEvents := registered[http.MethodPost+" /api/v1/cv/:linkId/events"]
	require.False(t, hasEvents)
	_, hasClicks := registered[http.MethodGet+" /api/v1/lk/link/:linkId/clicks"]
	require.False(t, hasClicks)
}
