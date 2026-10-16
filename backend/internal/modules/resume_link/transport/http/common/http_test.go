package common

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
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

func TestParseResumeIDReadsPathParameter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/lk/resume/90e9ee51-7a2d-4dcb-b0ab-086347d5e5c3/statistics", nil)
	c.Params = gin.Params{{Key: "resumeId", Value: "90e9ee51-7a2d-4dcb-b0ab-086347d5e5c3"}}

	id, ok := ParseResumeID(c)

	require.True(t, ok)
	require.Equal(t, "90e9ee51-7a2d-4dcb-b0ab-086347d5e5c3", id.String())
	require.Equal(t, http.StatusOK, w.Code)
}

func TestParseResumeIDReadsQueryParameter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/lk/links?resumeId=90e9ee51-7a2d-4dcb-b0ab-086347d5e5c3", nil)

	id, ok := ParseResumeIDQuery(c)

	require.True(t, ok)
	require.Equal(t, "90e9ee51-7a2d-4dcb-b0ab-086347d5e5c3", id.String())
	require.Equal(t, http.StatusOK, w.Code)
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

func TestWriteCommandErrorMapsResumeOwnershipNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)

	WriteCommandError(c, resumequeries.ErrResumeNotFound)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestWriteStatisticsErrorMapsMissingResumeLink(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	WriteStatisticsError(c, domain.ErrEntityNotFound)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestWriteCommandErrorMapsMissingResume(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)

	WriteCommandError(c, resumedomain.ErrEntityNotFound)

	require.Equal(t, http.StatusNotFound, w.Code)
}
