package experiences

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

type aggregateRepositoryStub struct {
	resume *domain.Resume
	saved  bool
}

func (r *aggregateRepositoryStub) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}
func (r *aggregateRepositoryStub) Save(_ context.Context, resume *domain.Resume) error {
	r.resume, r.saved = resume, true
	return nil
}
func (r *aggregateRepositoryStub) Create(context.Context, *domain.Resume) error { return nil }
func (r *aggregateRepositoryStub) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

type uowStub struct{}

func (uowStub) Do(_ context.Context, fn func(context.Context) error) error {
	return fn(context.Background())
}

func TestReplaceExperiencesBindsFullSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := useridentity.NewID()
	resume, err := domain.NewResume(domain.ResumeParams{ID: valueobjects.NewID(), UserID: uid, Title: "Resume", Language: enums.ResumeLanguageEN, Version: 1})
	require.NoError(t, err)
	repo := &aggregateRepositoryStub{resume: resume}
	handler := NewHandler(commands.NewReplaceExperiencesHandler(repo, uowStub{}))

	engine := gin.New()
	engine.PUT("/api/v1/lk/resume/:resumeId/experiences", func(c *gin.Context) {
		c.Set("user_id", uid.UUID())
		handler.Replace(c)
	})
	body := `{"version":1,"items":[{"company":"Acme","position":"Go Developer","startDate":"2024-01","achievements":["Built API"],"stack":[{"name":"Go","grade":"senior"}]}]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+resume.ID().String()+"/experiences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.saved)
	require.Len(t, repo.resume.Experiences(), 1)
	require.Equal(t, "Acme", repo.resume.Experiences()[0].Company())
	require.Len(t, repo.resume.Experiences()[0].Stack(), 1)
}
