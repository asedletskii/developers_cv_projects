package technologies

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
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeTechnologiesHTTPRepository struct{ resume *domain.Resume }

func (r *fakeTechnologiesHTTPRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}
func (r *fakeTechnologiesHTTPRepository) Save(context.Context, *domain.Resume) error {
	return nil
}
func (r *fakeTechnologiesHTTPRepository) Create(context.Context, *domain.Resume) error { return nil }
func (r *fakeTechnologiesHTTPRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

type fakeTechnologiesHTTPUnitOfWork struct{}

func (fakeTechnologiesHTTPUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestReplaceTechnologiesBindsFullSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := sharedvalueobjects.NewID()
	resume, err := domain.NewResume(domain.ResumeParams{
		ID: sharedvalueobjects.NewID(), UserID: uid, Title: "Resume", Language: domain.ResumeLanguageEN,
	})
	require.NoError(t, err)

	handler := NewHandler(commands.NewReplaceTechnologiesHandler(
		&fakeTechnologiesHTTPRepository{resume: resume},
		fakeTechnologiesHTTPUnitOfWork{},
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
