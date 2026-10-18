package http

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRegisterRoutesPreservesResumeLinkRouteContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := engine.Group("/api/v1")
	RegisterRoutes(NewHandler(HandlerDependencies{}), api, api, api.Group("/lk"), nil)

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
	}
	for _, route := range expected {
		_, ok := registered[route]
		require.Truef(t, ok, "expected route %q to be registered", route)
	}

	_, hasEvents := registered[http.MethodPost+" /api/v1/cv/:linkId/events"]
	require.False(t, hasEvents, "analytics ingestion route is registered by the Analytics transport")
}
