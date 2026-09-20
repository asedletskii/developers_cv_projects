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

func TestReplaceTechnologiesBindsFullSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := sharedvalueobjects.NewID()
	resume, err := domain.NewResume(domain.ResumeParams{
		UserID: uid, Title: "Resume", Language: domain.ResumeLanguageEN,
	})
	require.NoError(t, err)

	handler := NewResumeTechnologiesHandler(commands.NewReplaceTechnologiesHandler(
		&resumefixture.Repository{Resume: resume},
		resumefixture.UnitOfWork{},
	))

	engine := gin.New()
	engine.Use(func(c *gin.Context) { c.Set("user_id", uid.UUID()); c.Next() })
	engine.PUT("/resume/:resumeId/technologies", handler.Replace)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/technologies",
		strings.NewReader(`{"version":1,"items":[{"name":" Go ","position":0},{"name":"Rust","position":1}]}`))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Len(t, resume.Technologies().Entries(), 2)
	require.Equal(t, "Go", resume.Technologies().Entries()[0].Name())
	require.Equal(t, "Rust", resume.Technologies().Entries()[1].Name())
}
