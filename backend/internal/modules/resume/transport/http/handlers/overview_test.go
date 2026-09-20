package handlers

import (
	"bytes"
	"context"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/platform/http/middleware"
	"developers-cv/backend/internal/testutil/resumefixture"
)

func TestReplaceOverviewReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resume, repo := newOverviewTestResume(t)
	handler := NewResumeOverviewHandler(commands.NewReplaceOverviewHandler(repo, resumefixture.UnitOfWork{}))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return resume.UserID().UUID(), nil }))
	engine.PUT("/resume/:resumeId/overview", handler.Replace)

	body := `{"version":1,"greeting":{"first":"Updated","second":"","secondHighlight":"","third":""},"sections":[]}`
	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/overview", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.Saved)
}

func TestReplaceOverviewRejectsMissingVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resume, repo := newOverviewTestResume(t)
	handler := NewResumeOverviewHandler(commands.NewReplaceOverviewHandler(repo, resumefixture.UnitOfWork{}))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return resume.UserID().UUID(), nil }))
	engine.PUT("/resume/:resumeId/overview", handler.Replace)

	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/overview", bytes.NewBufferString(`{"greeting":{"first":"Updated","second":"","secondHighlight":"","third":""},"sections":[]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.False(t, repo.Saved)
}

func TestReplaceOverviewMapsConcurrencyConflictToPreconditionFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resume, repo := newOverviewTestResume(t)
	handler := NewResumeOverviewHandler(commands.NewReplaceOverviewHandler(repo, resumefixture.UnitOfWork{}))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return resume.UserID().UUID(), nil }))
	engine.PUT("/resume/:resumeId/overview", handler.Replace)

	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/overview", bytes.NewBufferString(`{"version":2,"greeting":{"first":"Updated","second":"","secondHighlight":"","third":""},"sections":[]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusPreconditionFailed, w.Code)
	require.False(t, repo.Saved)
}
