package common

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
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
