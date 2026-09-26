package errors

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/platform/http/middleware"

	"github.com/stretchr/testify/require"
)

func TestFromPreservesAppError(t *testing.T) {
	expected := New(CodeNotFound, http.StatusNotFound, "Resource not found")

	actual := From(expected)

	require.Same(t, expected, actual)
}

func TestFromConvertsUnknownError(t *testing.T) {
	actual := From(errors.New("database password leaked"))

	require.Equal(t, CodeInternal, actual.Code)
	require.Equal(t, http.StatusInternalServerError, actual.Status)
	require.Equal(t, "Internal server error", actual.Message)
	require.NotContains(t, actual.Message, "database")
}

func TestWriteErrorUsesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.RequestIDKey, "request-789")

	WriteError(c, New(CodeBadRequest, http.StatusBadRequest, "Invalid request"))

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), `"code":"BAD_REQUEST"`)
	require.Contains(t, w.Body.String(), `"requestId":"request-789"`)
}

func TestUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Unauthorized(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Contains(t, w.Body.String(), `"code":"UNAUTHORIZED"`)
}
