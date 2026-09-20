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

func TestReplaceCoursesReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resume, err := domain.NewResume(domain.ResumeParams{UserID: sharedvalueobjects.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	h := NewCourseHandler(commands.NewReplaceCoursesHandler(&resumefixture.Repository{Resume: resume}, resumefixture.UnitOfWork{}))
	e := gin.New()
	e.PUT("/resume/:resumeId/courses", func(c *gin.Context) { c.Set("user_id", resume.UserID().UUID()); h.Replace(c) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/courses", strings.NewReader(`{"version":1,"items":[]}`))
	req.Header.Set("Content-Type", "application/json")
	e.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
