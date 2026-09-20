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

func TestReplaceEducationsHandlerReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, err := domain.NewResume(domain.ResumeParams{UserID: sharedvalueobjects.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	uid := r.UserID()
	rid := r.ID()
	h := NewEducationHandler(commands.NewReplaceEducationsHandler(&resumefixture.Repository{Resume: r}, resumefixture.UnitOfWork{}))
	e := gin.New()
	e.PUT("/resume/:resumeId/educations", func(c *gin.Context) { c.Set("user_id", uid.UUID()); h.Replace(c) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/resume/"+rid.String()+"/educations", strings.NewReader(`{"version":1,"items":[]}`))
	req.Header.Set("Content-Type", "application/json")
	e.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
