package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/platform/http/requestid"
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

func TestWriteErrorUsesRequestIDFromContextAndHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(requestid.Key, "request-789")

	WriteError(c, New(CodeBadRequest, http.StatusBadRequest, "Invalid request"))

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, "request-789", w.Header().Get(requestid.Header))
	require.JSONEq(t, `{"error":{"code":"BAD_REQUEST","message":"Invalid request","requestId":"request-789"}}`, w.Body.String())
}

func TestWriteErrorHandlesContextWithoutRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	WriteError(c, New(CodeBadRequest, http.StatusBadRequest, "Invalid request"))

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Empty(t, w.Header().Get(requestid.Header))
	require.JSONEq(t, `{"error":{"code":"BAD_REQUEST","message":"Invalid request","requestId":""}}`, w.Body.String())
}

func TestWriteErrorUsesRequestIDFromHeaderWhenContextIsEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set(requestid.Header, "request-header")

	WriteError(c, New(CodeBadRequest, http.StatusBadRequest, "Invalid request"))

	require.Equal(t, "request-header", w.Header().Get(requestid.Header))
	require.JSONEq(t, `{"error":{"code":"BAD_REQUEST","message":"Invalid request","requestId":"request-header"}}`, w.Body.String())
}

func TestWriteErrorPrefersContextRequestIDOverHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set(requestid.Header, "request-header")
	c.Set(requestid.Key, "request-context")

	WriteError(c, New(CodeBadRequest, http.StatusBadRequest, "Invalid request"))

	require.Equal(t, "request-context", w.Header().Get(requestid.Header))
	require.JSONEq(t, `{"error":{"code":"BAD_REQUEST","message":"Invalid request","requestId":"request-context"}}`, w.Body.String())
}

func TestWriteErrorContract(t *testing.T) {
	tests := []struct {
		name    string
		code    Code
		status  int
		message string
	}{
		{name: "bad request", code: CodeBadRequest, status: http.StatusBadRequest, message: "Invalid request"},
		{name: "unauthorized", code: CodeUnauthorized, status: http.StatusUnauthorized, message: "Authentication required"},
		{name: "not found", code: CodeNotFound, status: http.StatusNotFound, message: "Resource not found"},
		{name: "gone", code: CodeGone, status: http.StatusGone, message: "Resource is gone"},
		{name: "conflict", code: CodeConflict, status: http.StatusConflict, message: "Resource already exists"},
		{name: "precondition failed", code: CodePreconditionFailed, status: http.StatusPreconditionFailed, message: "Precondition failed"},
		{name: "request too large", code: CodeRequestTooLarge, status: http.StatusRequestEntityTooLarge, message: "Request body too large"},
		{name: "too many requests", code: CodeTooMany, status: http.StatusTooManyRequests, message: "Too many requests"},
		{name: "internal", code: CodeInternal, status: http.StatusInternalServerError, message: "Internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set(requestid.Key, "request-contract")

			appErr := New(tt.code, tt.status, tt.message)
			if tt.code == CodeInternal {
				appErr = Internal(errors.New("database connection failed"))
			}
			WriteError(c, appErr)

			require.Equal(t, tt.status, w.Code)
			var body struct {
				Error struct {
					Code      string `json:"code"`
					Message   string `json:"message"`
					RequestID string `json:"requestId"`
				} `json:"error"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
			require.Equal(t, string(tt.code), body.Error.Code)
			require.Equal(t, tt.message, body.Error.Message)
			require.Equal(t, "request-contract", body.Error.RequestID)
			require.NotContains(t, w.Body.String(), "database connection failed")
		})
	}
}

func TestUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Unauthorized(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Contains(t, w.Body.String(), `"code":"UNAUTHORIZED"`)
}
