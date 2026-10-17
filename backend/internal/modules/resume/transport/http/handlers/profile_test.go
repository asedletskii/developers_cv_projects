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

func TestHandlerReplaceProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := uuid.New()
	owner, err := sharedvalueobjects.IDFromUUID(uid)
	require.NoError(t, err)
	resume, err := domain.NewResume(domain.ResumeParams{UserID: owner, Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &resumefixture.Repository{Resume: resume}
	h := NewProfileHandler(commands.NewReplaceProfileHandler(repo, resumefixture.UnitOfWork{}))
	e := gin.New()
	e.Use(func(c *gin.Context) { c.Set("user_id", uid); c.Next() })
	e.PUT("/resume/:resumeId/profile", h.Replace)
	body := `{"version":1,"name":"John","surname":"Doe","position":"Developer","photoPath":"https://example.com/p.jpg","dateOfBirth":"1998-04-12","location":"Paris","workFormat":["remote"]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/resume/"+resume.ID().String()+"/profile", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	e.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.Saved)
}
func TestHandlerReplaceProfileRejectsMissingVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewProfileHandler(nil)
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
	h := NewProfileHandler(nil)
	e := gin.New()
	e.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, middleware.ErrUnauthenticated
	}))
	e.PUT("/resume/:resumeId/profile", h.Replace)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/resume/"+uuid.New().String()+"/profile", nil))
	require.Equal(t, http.StatusUnauthorized, w.Code)
}
