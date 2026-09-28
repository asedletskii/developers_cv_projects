package activities

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type repoStub struct {
	resume *domain.Resume
	saved  bool
}

func (s *repoStub) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return s.resume, nil
}
func (s *repoStub) Save(context.Context, *domain.Resume) error             { s.saved = true; return nil }
func (s *repoStub) Create(context.Context, *domain.Resume) error           { return nil }
func (s *repoStub) Delete(context.Context, ports.DeleteResumeParams) error { return nil }

type uowStub struct{}

func (uowStub) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

func TestReplaceActivitiesReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resume, err := domain.NewResume(domain.ResumeParams{ID: valueobjects.NewID(), UserID: useridentity.NewID(), Title: "Resume", Language: enums.ResumeLanguageEN})
	require.NoError(t, err)
	h := NewHandler(commands.NewReplaceActivitiesHandler(&repoStub{resume: resume}, uowStub{}))
	e := gin.New()
	e.PUT("/resume/:resumeId/activities", func(c *gin.Context) { c.Set("user_id", resume.UserID().UUID()); h.Replace(c) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/activities", strings.NewReader(`{"version":1,"items":[]}`))
	req.Header.Set("Content-Type", "application/json")
	e.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
