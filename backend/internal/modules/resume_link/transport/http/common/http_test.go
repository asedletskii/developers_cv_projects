package common

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/domain"
)

type bindTestRequest struct {
	Name string `json:"name"`
}

func TestBindJSONRejectsDuplicateFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"first","name":"second"}`))
	var body bindTestRequest
	err := BindJSON(ctx, &body)
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate JSON field")
}

func TestBindJSONRejectsMultipleValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"first"}{"name":"second"}`))
	var body bindTestRequest
	require.Error(t, BindJSON(ctx, &body))
}

func TestRequirePositiveVersionRejectsNonPositive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/", nil)
	require.False(t, RequirePositiveVersion(c, 0))
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWriteCommandErrorMapsConcurrencyConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/", nil)
	WriteCommandError(c, domain.ErrConcurrencyConflict)
	require.Equal(t, http.StatusPreconditionFailed, w.Code)
}

func TestWriteCommandErrorMapsInvalidName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	WriteCommandError(c, domain.ErrInvalidResumeLinkName)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
