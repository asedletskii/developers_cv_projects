package contacts

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/platform/http/middleware"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
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

func TestHandlerReplaceContacts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rid := uuid.New()
	uid := uuid.New()
	resumeID, err := valueobjects.IDFromUUID(rid)
	require.NoError(t, err)
	owner, err := useridentity.IDFromUUID(uid)
	require.NoError(t, err)
	resume, err := domain.NewResume(domain.ResumeParams{ID: resumeID, UserID: owner, Title: "Resume", Language: enums.ResumeLanguageEN, Version: 1})
	require.NoError(t, err)

	repo := &fakeRepo{resume: resume}
	h := NewHandler(commands.NewReplaceContactsHandler(repo, fakeUOW{}))
	e := gin.New()
	e.Use(func(c *gin.Context) { c.Set("user_id", uid); c.Next() })
	e.PUT("/resume/:resumeId/contacts", h.Replace)

	body := `{"version":1,"email":"john@example.com","hiddenContacts":["email"],"contactsText":{"title":"Contacts"}}`
	req := httptest.NewRequest(http.MethodPut, "/resume/"+rid.String()+"/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.saved)
	require.NotNil(t, resume.Contacts())
}

func TestHandlerReplaceContactsRejectsMissingVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(nil)
	e := gin.New()
	e.Use(func(c *gin.Context) { c.Set("user_id", uuid.New()); c.Next() })
	e.PUT("/resume/:resumeId/contacts", h.Replace)

	req := httptest.NewRequest(http.MethodPut, "/resume/"+uuid.New().String()+"/contacts", bytes.NewBufferString(`{"email":"john@example.com"}`))
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerReplaceContactsRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(nil)
	e := gin.New()
	e.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, middleware.ErrUnauthenticated
	}))
	e.PUT("/resume/:resumeId/contacts", h.Replace)

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/resume/"+uuid.New().String()+"/contacts", nil))
	require.Equal(t, http.StatusUnauthorized, w.Code)
}
