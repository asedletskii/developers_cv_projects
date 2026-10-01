package languages

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

type fakeLanguagesHTTPRepository struct{ resume *domain.Resume }

func (r *fakeLanguagesHTTPRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}
func (r *fakeLanguagesHTTPRepository) Save(context.Context, *domain.Resume) error   { return nil }
func (r *fakeLanguagesHTTPRepository) Create(context.Context, *domain.Resume) error { return nil }
func (r *fakeLanguagesHTTPRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

type fakeLanguagesHTTPUnitOfWork struct{}

func (fakeLanguagesHTTPUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestReplaceLanguagesBindsFullSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := sharedvalueobjects.NewID()
	resume, err := domain.NewResume(domain.ResumeParams{ID: sharedvalueobjects.NewID(), UserID: uid, Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &fakeLanguagesHTTPRepository{resume: resume}
	handler := NewHandler(commands.NewReplaceLanguagesHandler(repo, fakeLanguagesHTTPUnitOfWork{}))

	engine := gin.New()
	engine.Use(func(c *gin.Context) { c.Set("user_id", uid.UUID()); c.Next() })
	engine.PUT("/resume/:resumeId/languages", handler.Replace)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/languages", strings.NewReader(`{"version":1,"items":[{"language":"en","grade":"C1","position":0},{"language":"fr","position":1}]}`))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Len(t, resume.Languages().Entries(), 2)
	require.Equal(t, domain.SpokenLanguageEN, resume.Languages().Entries()[0].Language())
	require.Equal(t, domain.SpokenLanguageFR, resume.Languages().Entries()[1].Language())
}
