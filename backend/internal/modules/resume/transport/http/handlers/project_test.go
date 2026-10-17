package handlers

import (
	"context"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type projectHTTPRepoStub struct {
	resume *domain.Resume
	saved  bool
}

func (s *projectHTTPRepoStub) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return s.resume, nil
}
func (s *projectHTTPRepoStub) Save(context.Context, *domain.Resume) error {
	s.saved = true
	return nil
}
func (r *projectHTTPRepoStub) SaveSection(ctx context.Context, resume *domain.Resume, _ ports.ResumeMutationSection) error {
	return r.Save(ctx, resume)
}
func (s *projectHTTPRepoStub) Create(context.Context, *domain.Resume) error { return nil }
func (s *projectHTTPRepoStub) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

type projectHTTPUOWStub struct{}

func (projectHTTPUOWStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestReplaceProjectsReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resume, err := domain.NewResume(domain.ResumeParams{
		UserID: sharedvalueobjects.NewID(),
		Title:  "Resume", Language: domain.ResumeLanguageEN,
	})
	require.NoError(t, err)

	h := NewProjectHandler(commands.NewReplaceProjectsHandler(&projectHTTPRepoStub{resume: resume}, projectHTTPUOWStub{}))
	engine := gin.New()
	engine.PUT("/resume/:resumeId/projects", func(c *gin.Context) {
		c.Set("user_id", resume.UserID().UUID())
		h.Replace(c)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/projects", strings.NewReader(`{"version":1,"items":[{"name":"Project","screenshots":[],"stack":[],"contributors":[]}]}`))
	request.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestReplaceProjectsRejectsNonPositiveVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resume, err := domain.NewResume(domain.ResumeParams{
		UserID: sharedvalueobjects.NewID(),
		Title:  "Resume", Language: domain.ResumeLanguageEN,
	})
	require.NoError(t, err)

	h := NewProjectHandler(commands.NewReplaceProjectsHandler(&projectHTTPRepoStub{resume: resume}, projectHTTPUOWStub{}))
	engine := gin.New()
	engine.PUT("/resume/:resumeId/projects", func(c *gin.Context) {
		c.Set("user_id", resume.UserID().UUID())
		h.Replace(c)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/projects", strings.NewReader(`{"version":0,"items":[]}`))
	request.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), `"code":"BAD_REQUEST"`)
}
