package skills

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
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

type fakeSkillsHTTPRepository struct {
	resume *domain.Resume
}

func (r *fakeSkillsHTTPRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}
func (r *fakeSkillsHTTPRepository) Save(context.Context, *domain.Resume) error   { return nil }
func (r *fakeSkillsHTTPRepository) Create(context.Context, *domain.Resume) error { return nil }
func (r *fakeSkillsHTTPRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

type fakeSkillsHTTPUnitOfWork struct{}

func (fakeSkillsHTTPUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestReplaceSkillsBindsFullSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := useridentity.NewID()
	resume, err := domain.NewResume(domain.ResumeParams{ID: valueobjects.NewID(), UserID: uid, Title: "Resume", Language: enums.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &fakeSkillsHTTPRepository{resume: resume}
	handler := NewHandler(commands.NewReplaceSkillsHandler(repo, fakeSkillsHTTPUnitOfWork{}))

	e := gin.New()
	e.Use(func(c *gin.Context) { c.Set("user_id", uid.UUID()); c.Next() })
	e.PUT("/resume/:resumeId/skills", handler.Replace)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/skills", bytes.NewBufferString(`{"version":1,"hardSkills":["testing"],"softSkills":["teamwork"]}`))
	req.Header.Set("Content-Type", "application/json")
	e.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.NotNil(t, resume.Skills())
	require.Equal(t, []enums.HardSkill{enums.HardSkillTesting}, resume.Skills().HardSkills())
	require.Equal(t, []enums.SoftSkill{enums.SoftSkillTeamwork}, resume.Skills().SoftSkills())
}
