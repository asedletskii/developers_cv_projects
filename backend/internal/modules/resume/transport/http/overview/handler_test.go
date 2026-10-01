package overview

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
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/platform/http/middleware"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type handlerAggregateRepository struct {
	resume *domain.Resume
	saved  bool
}

func (r *handlerAggregateRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}
func (r *handlerAggregateRepository) Save(context.Context, *domain.Resume) error {
	r.saved = true
	return nil
}
func (r *handlerAggregateRepository) Create(context.Context, *domain.Resume) error { return nil }
func (r *handlerAggregateRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

type handlerUOW struct{}

func (handlerUOW) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

func newOverviewTestResume(t *testing.T) *domain.Resume {
	t.Helper()
	resume, err := domain.NewResume(domain.ResumeParams{UserID: identity.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	overview := domain.ResumeOverviewParams{Greeting: domain.NewGreeting(domain.GreetingParams{}), Sections: []domain.SectionSettings{}}
	require.NoError(t, resume.ReplaceOverview(overview))
	return resume
}

func TestReplaceOverviewReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resume := newOverviewTestResume(t)
	repo := &handlerAggregateRepository{resume: resume}
	handler := NewHandler(commands.NewReplaceOverviewHandler(repo, handlerUOW{}))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return resume.UserID().UUID(), nil }))
	engine.PUT("/resume/:resumeId/overview", handler.Replace)

	body := `{"version":1,"greeting":{"first":"Updated","second":"","secondHighlight":"","third":""},"sections":[]}`
	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/overview", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.saved)
}

func TestReplaceOverviewRejectsMissingVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resume := newOverviewTestResume(t)
	repo := &handlerAggregateRepository{resume: resume}
	handler := NewHandler(commands.NewReplaceOverviewHandler(repo, handlerUOW{}))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return resume.UserID().UUID(), nil }))
	engine.PUT("/resume/:resumeId/overview", handler.Replace)

	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/overview", bytes.NewBufferString(`{"greeting":{"first":"Updated","second":"","secondHighlight":"","third":""},"sections":[]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.False(t, repo.saved)
}

func TestReplaceOverviewMapsConcurrencyConflictToPreconditionFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resume := newOverviewTestResume(t)
	repo := &handlerAggregateRepository{resume: resume}
	handler := NewHandler(commands.NewReplaceOverviewHandler(repo, handlerUOW{}))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return resume.UserID().UUID(), nil }))
	engine.PUT("/resume/:resumeId/overview", handler.Replace)

	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/overview", bytes.NewBufferString(`{"version":2,"greeting":{"first":"Updated","second":"","secondHighlight":"","third":""},"sections":[]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusPreconditionFailed, w.Code)
	require.False(t, repo.saved)
}
