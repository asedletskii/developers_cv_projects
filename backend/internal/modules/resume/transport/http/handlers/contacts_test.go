package handlers

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
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/resumefixture"
)

func TestHandlerReplaceContacts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rid := uuid.New()
	uid := uuid.New()
	owner, err := sharedvalueobjects.IDFromUUID(uid)
	require.NoError(t, err)
	resume, err := domain.NewResume(domain.ResumeParams{UserID: owner, Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)

	repo := &resumefixture.Repository{Resume: resume}
	h := NewContactsHandler(commands.NewReplaceContactsHandler(repo, resumefixture.UnitOfWork{}))
	e := gin.New()
	e.Use(func(c *gin.Context) { c.Set("user_id", uid); c.Next() })
	e.PUT("/resume/:resumeId/contacts", h.Replace)

	body := `{"version":1,"email":"john@example.com","hiddenContacts":["email"],"contactsText":{"title":"Contacts"}}`
	req := httptest.NewRequest(http.MethodPut, "/resume/"+rid.String()+"/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.Saved)
	require.NotNil(t, resume.Contacts())
}

func TestHandlerReplaceContactsDropsEmptyText(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rid := uuid.New()
	uid := uuid.New()
	owner, err := sharedvalueobjects.IDFromUUID(uid)
	require.NoError(t, err)
	resume, err := domain.NewResume(domain.ResumeParams{UserID: owner, Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)

	repo := &resumefixture.Repository{Resume: resume}
	h := NewContactsHandler(commands.NewReplaceContactsHandler(repo, resumefixture.UnitOfWork{}))
	e := gin.New()
	e.Use(func(c *gin.Context) { c.Set("user_id", uid); c.Next() })
	e.PUT("/resume/:resumeId/contacts", h.Replace)

	body := `{"version":1,"email":"john@example.com","hiddenContacts":[],"contactsText":{"title":"","subtitle":"   "}}`
	req := httptest.NewRequest(http.MethodPut, "/resume/"+rid.String()+"/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.Saved)
	require.NotNil(t, resume.Contacts())
	require.Nil(t, resume.Contacts().Text())
}

func TestHandlerReplaceContactsRejectsMissingVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewContactsHandler(nil)
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
	h := NewContactsHandler(nil)
	e := gin.New()
	e.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, middleware.ErrUnauthenticated
	}))
	e.PUT("/resume/:resumeId/contacts", h.Replace)

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/resume/"+uuid.New().String()+"/contacts", nil))
	require.Equal(t, http.StatusUnauthorized, w.Code)
}
