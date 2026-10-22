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

func TestReplaceLanguagesBindsFullSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := sharedvalueobjects.NewID()
	resume, err := domain.NewResume(domain.ResumeParams{UserID: uid, Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &resumefixture.Repository{Resume: resume}
	handler := NewResumeSpokenLanguagesHandler(commands.NewReplaceLanguagesHandler(repo, resumefixture.UnitOfWork{}))

	engine := gin.New()
	engine.Use(func(c *gin.Context) { c.Set("user_id", uid.UUID()); c.Next() })
	engine.PUT("/resume/:resumeId/languages", handler.Replace)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/languages", strings.NewReader(`{"version":1,"items":[{"language":"en","grade":"C1","position":0},{"language":"ru","position":1}]}`))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Len(t, resume.SpokenLanguages().Entries(), 2)
	require.Equal(t, domain.SpokenLanguageEN, resume.SpokenLanguages().Entries()[0].Language())
	require.Equal(t, domain.SpokenLanguageRU, resume.SpokenLanguages().Entries()[1].Language())
}
