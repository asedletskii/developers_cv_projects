package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/resumefixture"
)

func TestReplaceExperiencesMapsConcurrencyConflictToPreconditionFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := sharedvalueobjects.NewID()
	resume, err := domain.NewResume(domain.ResumeParams{
		UserID: uid, Title: "Resume", Language: domain.ResumeLanguageEN,
	})
	require.NoError(t, err)
	repo := &resumefixture.Repository{Resume: resume}
	handler := NewExperienceHandler(commands.NewReplaceExperiencesHandler(repo, resumefixture.UnitOfWork{}))

	engine := gin.New()
	engine.PUT("/api/v1/lk/resume/:resumeId/experiences", func(c *gin.Context) {
		c.Set("user_id", uid.UUID())
		handler.Replace(c)
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+resume.ID().String()+"/experiences", strings.NewReader(`{"version":2,"items":[]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusPreconditionFailed, w.Code)
	require.False(t, repo.Saved)
}

func TestReplaceExperiencesBindsFullSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := sharedvalueobjects.NewID()
	resume, err := domain.NewResume(domain.ResumeParams{UserID: uid, Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &resumefixture.Repository{Resume: resume}
	handler := NewExperienceHandler(commands.NewReplaceExperiencesHandler(repo, resumefixture.UnitOfWork{}))

	engine := gin.New()
	engine.PUT("/api/v1/lk/resume/:resumeId/experiences", func(c *gin.Context) {
		c.Set("user_id", uid.UUID())
		handler.Replace(c)
	})
	body := `{"version":1,"items":[{"company":"Acme","position":"Go Developer","startDate":"2024-01","achievements":["Built API"],"tasks":["Implement backend","Review PRs"],"stack":[{"name":"Go"}]}]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+resume.ID().String()+"/experiences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.Saved)
	require.Len(t, repo.Resume.Experiences(), 1)
	require.Equal(t, "Acme", repo.Resume.Experiences()[0].Company())
	require.Equal(t, []string{"Implement backend", "Review PRs"}, repo.Resume.Experiences()[0].Tasks())
	require.Len(t, repo.Resume.Experiences()[0].Stack(), 1)
}
