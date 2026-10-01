package profile

import (
	"bytes"
	"context"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/platform/http/middleware"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeRepo struct {
	resume *domain.Resume
	saved  bool
}

func (r *fakeRepo) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}
func (r *fakeRepo) Save(context.Context, *domain.Resume) error             { r.saved = true; return nil }
func (r *fakeRepo) Create(context.Context, *domain.Resume) error           { return nil }
func (r *fakeRepo) Delete(context.Context, ports.DeleteResumeParams) error { return nil }

type fakeUOW struct{}

func (fakeUOW) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }
func TestHandlerReplaceProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := uuid.New()
	rid := uuid.New()
	id, err := identity.IDFromUUID(rid)
	require.NoError(t, err)
	owner, err := identity.IDFromUUID(uid)
	require.NoError(t, err)
	resume, err := domain.NewResume(domain.ResumeParams{ID: id, UserID: owner, Title: "Resume", Language: domain.ResumeLanguageEN, Version: 1})
	require.NoError(t, err)
	repo := &fakeRepo{resume: resume}
	h := NewHandler(commands.NewReplaceProfileHandler(repo, fakeUOW{}))
	e := gin.New()
	e.Use(func(c *gin.Context) { c.Set("user_id", uid); c.Next() })
	e.PUT("/resume/:resumeId/profile", h.Replace)
	body := `{"version":1,"name":"John","surname":"Doe","position":"Developer","photoPath":"https://example.com/p.jpg","dateOfBirth":"1998-04-12","location":"Paris","workFormat":["remote"]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/resume/"+rid.String()+"/profile", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	e.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.saved)
}
func TestHandlerReplaceProfileRejectsMissingVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(nil)
	e := gin.New()
	e.Use(func(c *gin.Context) { c.Set("user_id", uuid.New()); c.Next() })
	e.PUT("/resume/:resumeId/profile", h.Replace)
	req := httptest.NewRequest(http.MethodPut, "/resume/"+uuid.New().String()+"/profile", bytes.NewBufferString(`{"name":"John"}`))
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
func TestHandlerRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(nil)
	e := gin.New()
	e.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, middleware.ErrUnauthenticated
	}))
	e.PUT("/resume/:resumeId/profile", h.Replace)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/resume/"+uuid.New().String()+"/profile", nil))
	require.Equal(t, http.StatusUnauthorized, w.Code)
}
