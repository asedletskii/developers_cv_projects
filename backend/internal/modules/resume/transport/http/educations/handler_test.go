package educations

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
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type repoStub struct{ resume *domain.Resume }

func (s repoStub) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return s.resume, nil
}
func (s repoStub) Save(context.Context, *domain.Resume) error             { return nil }
func (s repoStub) Create(context.Context, *domain.Resume) error           { return nil }
func (s repoStub) Delete(context.Context, ports.DeleteResumeParams) error { return nil }

type uowStub struct{}

func (uowStub) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }
func TestReplaceEducationsHandlerReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, err := domain.NewResume(domain.ResumeParams{ID: identity.NewID(), UserID: identity.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	uid := r.UserID()
	rid := r.ID()
	h := NewHandler(commands.NewReplaceEducationsHandler(repoStub{resume: r}, uowStub{}))
	e := gin.New()
	e.PUT("/resume/:resumeId/educations", func(c *gin.Context) { c.Set("user_id", uid.UUID()); h.Replace(c) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/resume/"+rid.String()+"/educations", strings.NewReader(`{"version":1,"items":[]}`))
	req.Header.Set("Content-Type", "application/json")
	e.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
