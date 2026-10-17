package http

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/transport/http/handlers"
)

func TestRegisterRoutesPreservesResumeLinkRouteContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := engine.Group("/api/v1")
	RegisterRoutes(NewHandler(HandlerDependencies{}), api, api, api.Group("/lk"), nil, nil)

	registered := make(map[string]struct{}, len(engine.Routes()))
	for _, route := range engine.Routes() {
		registered[route.Method+" "+route.Path] = struct{}{}
	}

	expected := []string{
		http.MethodGet + " /api/v1/cv/:linkId",
		http.MethodPost + " /api/v1/cv/:linkId/contact-views",
		http.MethodGet + " /api/v1/lk/links",
		http.MethodPost + " /api/v1/lk/link",
		http.MethodDelete + " /api/v1/lk/link/:linkId",
		http.MethodPut + " /api/v1/lk/link/:linkId/status",
		http.MethodGet + " /api/v1/lk/link/:linkId/statistics",
		http.MethodGet + " /api/v1/lk/link/:linkId/statistics/daily",
		http.MethodGet + " /api/v1/lk/link/:linkId/views",
		http.MethodGet + " /api/v1/lk/link/:linkId/statistics/sections",
		http.MethodGet + " /api/v1/lk/link/:linkId/statistics/technologies",
		http.MethodGet + " /api/v1/lk/resume/:resumeId/statistics",
	}
	for _, route := range expected {
		_, ok := registered[route]
		require.Truef(t, ok, "expected route %q to be registered", route)
	}

	_, hasClicks := registered[http.MethodGet+" /api/v1/lk/link/:linkId/clicks"]
	require.False(t, hasClicks, "click-map route must remain conditional on its use case")
	_, hasEvents := registered[http.MethodPost+" /api/v1/cv/:linkId/events"]
	require.False(t, hasEvents, "analytics ingestion route must remain conditional on its handler")
}

func TestRegisterRoutesAddsOptionalAnalyticsRoutesWhenConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := engine.Group("/api/v1")
	handler := newClickHTTPHandler(t)
	handler.analyticsEvents = &handlers.ResumeLinkAnalyticsEventsHandler{}
	RegisterRoutes(handler, api, api, api.Group("/lk"), nil, nil)

	registered := make(map[string]struct{}, len(engine.Routes()))
	for _, route := range engine.Routes() {
		registered[route.Method+" "+route.Path] = struct{}{}
	}

	_, hasClicks := registered[http.MethodGet+" /api/v1/lk/link/:linkId/clicks"]
	require.True(t, hasClicks, "click-map route should be registered when its use case is configured")
	_, hasEvents := registered[http.MethodPost+" /api/v1/cv/:linkId/events"]
	require.True(t, hasEvents, "analytics ingestion route should be registered when its handler is configured")
}
