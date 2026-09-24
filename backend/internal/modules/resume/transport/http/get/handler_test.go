package get

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/platform/http/middleware"
)

type stubGetResumeReader struct {
	resume queries.Resume
	err    error
	query  queries.GetResumeQuery
}

func (r *stubGetResumeReader) GetResume(_ context.Context, query queries.GetResumeQuery) (queries.Resume, error) {
	r.query = query
	return r.resume, r.err
}

func newTestHandler(reader queries.GetResumeReader) *Handler {
	return NewHandler(queries.NewGetResumeHandler(reader))
}

func TestHandlerPrivateReturnsResumeForCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	userID := uuid.New()
	reader := &stubGetResumeReader{
		resume: queries.Resume{ID: id.String(), Name: "Anton Sedletskii"},
	}
	handler := newTestHandler(reader)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.GET("/api/v1/lk/resume/:resumeId", handler.Private)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/resume/"+id.String()+"?lang=en", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "Anton Sedletskii")
	require.NotNil(t, reader.query.UserID)
	require.Equal(t, userID, reader.query.UserID.UUID())
	require.Equal(t, queries.LanguageEN, reader.query.Lang)
}

func TestHandlerPublicDoesNotRequireAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	reader := &stubGetResumeReader{
		resume: queries.Resume{ID: id.String(), Name: "Anton Sedletskii"},
	}
	handler := newTestHandler(reader)

	engine := gin.New()
	engine.GET("/api/v1/cv/:resumeId", handler.Public)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+id.String()+"?lang=en", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "Anton Sedletskii")
	require.Nil(t, reader.query.UserID)
	require.Equal(t, queries.LanguageEN, reader.query.Lang)
}

func TestHandlerPrivateReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newTestHandler(&stubGetResumeReader{err: queries.ErrResumeNotFound})
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.GET("/api/v1/lk/resume/:resumeId", handler.Private)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/resume/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "Resume not found")
}

func TestHandlerPublicReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newTestHandler(&stubGetResumeReader{err: queries.ErrResumeNotFound})
	engine := gin.New()
	engine.GET("/api/v1/cv/:resumeId", handler.Public)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "Resume not found")
}

func TestHandlerPrivateRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newTestHandler(&stubGetResumeReader{})
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.GET("/api/v1/lk/resume/:resumeId", handler.Private)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/resume/not-a-uuid", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid resume id")
}

func TestHandlerPublicRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newTestHandler(&stubGetResumeReader{})
	engine := gin.New()
	engine.GET("/api/v1/cv/:resumeId", handler.Public)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/not-a-uuid", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid resume id")
}

func TestHandlerPrivateRejectsInvalidLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newTestHandler(&stubGetResumeReader{})
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.GET("/api/v1/lk/resume/:resumeId", handler.Private)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/resume/"+uuid.NewString()+"?lang=de", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid language")
}

func TestHandlerPublicRejectsInvalidLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newTestHandler(&stubGetResumeReader{})
	engine := gin.New()
	engine.GET("/api/v1/cv/:resumeId", handler.Public)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+uuid.NewString()+"?lang=de", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid language")
}
