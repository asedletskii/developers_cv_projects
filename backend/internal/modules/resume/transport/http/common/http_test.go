package common

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
)

type bindTestRequest struct {
	Name string `json:"name"`
}

func TestBindJSONRejectsDuplicateFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"first","name":"second"}`))

	var request bindTestRequest
	err := BindJSON(ctx, &request)

	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate JSON field")
}

func TestBindJSONRejectsNestedDuplicateFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"ok","nested":{"value":1,"value":2}}`))

	var request map[string]any
	err := BindJSON(ctx, &request)

	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate JSON field")
}

func TestBindJSONRejectsMultipleJSONValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"first"}{"name":"second"}`))

	var request bindTestRequest
	err := BindJSON(ctx, &request)

	require.Error(t, err)
}

func TestBindJSONReturnsBodyTooLarge(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"long"}`))
	request.Body = http.MaxBytesReader(recorder, request.Body, 4)
	ctx.Request = request

	var body bindTestRequest
	err := BindJSON(ctx, &body)

	require.ErrorIs(t, err, ErrRequestBodyTooLarge)
}

func TestBindJSONRejectsNonObjectBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(`null`))

	var request bindTestRequest
	err := BindJSON(ctx, &request)

	require.Error(t, err)
	require.Contains(t, err.Error(), "JSON object")
}

func TestRequirePositiveVersionRejectsNonPositiveVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/", nil)

	require.False(t, RequirePositiveVersion(c, 0))
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), `"code":"BAD_REQUEST"`)
}

func TestWriteCommandErrorMapsConcurrencyConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/", nil)

	WriteCommandError(c, domain.ErrConcurrencyConflict)

	require.Equal(t, http.StatusPreconditionFailed, w.Code)
	require.Contains(t, w.Body.String(), `"code":"PRECONDITION_FAILED"`)
}

func TestWriteCommandErrorMapsStatusTransitionConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/", nil)

	WriteCommandError(c, domain.ErrInvalidResumeStatusTransition)

	require.Equal(t, http.StatusConflict, w.Code)
	require.Contains(t, w.Body.String(), `"code":"CONFLICT"`)
}
