package handlers

import (
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/testutil/resumefixture"
)

func TestReplaceProjectsReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resume, repo := newHandlerTestResume(t)

	h := NewProjectHandler(commands.NewReplaceProjectsHandler(repo, resumefixture.UnitOfWork{}))
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
	resume, repo := newHandlerTestResume(t)

	h := NewProjectHandler(commands.NewReplaceProjectsHandler(repo, resumefixture.UnitOfWork{}))
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
