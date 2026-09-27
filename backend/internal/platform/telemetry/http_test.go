package telemetry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"developers-cv/backend/internal/platform/http/middleware"
)

func TestHTTPMiddlewareCreatesServerSpanAndRecordsRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	previousProvider := otel.GetTracerProvider()
	otel.SetTracerProvider(provider)
	t.Cleanup(func() {
		err := provider.Shutdown(context.Background())
		require.NoError(t, err)
		otel.SetTracerProvider(previousProvider)
	})

	engine := gin.New()
	engine.Use(middleware.RequestID(), HTTPMiddleware("developers-cv-api"))
	engine.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("X-Request-ID", "request-123")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	spans := recorder.Ended()
	require.Len(t, spans, 1)
	span := spans[0]
	require.Equal(t, "GET /health", span.Name())
	require.Equal(t, trace.SpanKindServer, span.SpanKind())

	var requestID string
	for _, keyValue := range span.Attributes() {
		if keyValue.Key == attribute.Key("request.id") {
			requestID = keyValue.Value.AsString()
			break
		}
	}

	require.Equal(t, "request-123", requestID)
}

func TestHTTPMiddlewareMarksServerErrorSpan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	previousProvider := otel.GetTracerProvider()
	otel.SetTracerProvider(provider)
	t.Cleanup(func() {
		err := provider.Shutdown(context.Background())
		require.NoError(t, err)
		otel.SetTracerProvider(previousProvider)
	})

	engine := gin.New()
	engine.Use(HTTPMiddleware("developers-cv-api"))
	engine.GET("/boom", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})

	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/boom", nil))

	require.Equal(t, http.StatusInternalServerError, response.Code)
	spans := recorder.Ended()
	require.Len(t, spans, 1)
	require.Equal(t, codes.Error, spans[0].Status().Code)
}
