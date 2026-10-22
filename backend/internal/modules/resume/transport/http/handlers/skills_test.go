package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/resumefixture"
)

func TestReplaceSkillsBindsFullSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := sharedvalueobjects.NewID()
	resume, err := domain.NewResume(domain.ResumeParams{UserID: uid, Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &resumefixture.Repository{Resume: resume}
	handler := NewSkillsHandler(commands.NewReplaceSkillsHandler(repo, resumefixture.UnitOfWork{}))

	e := gin.New()
	e.Use(func(c *gin.Context) { c.Set("user_id", uid.UUID()); c.Next() })
	e.PUT("/resume/:resumeId/skills", handler.Replace)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/skills", bytes.NewBufferString(`{"version":1,"hardSkills":["testing"],"softSkills":["teamwork"]}`))
	req.Header.Set("Content-Type", "application/json")
	e.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.NotNil(t, resume.Skills())
	require.Equal(t, []string{"testing"}, resume.Skills().HardSkills())
	require.Equal(t, []string{"teamwork"}, resume.Skills().SoftSkills())
}
