package http

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type stubResumeReader struct {
	resume queries.Resume
	err    error
}

func (r stubResumeReader) GetResume(context.Context, queries.GetResumeQuery) (queries.Resume, error) {
	return r.resume, r.err
}

func TestGetResume(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.New()
	handler := NewHandler(queries.NewGetResumeHandler(stubResumeReader{
		resume: queries.Resume{ID: id.String(), Name: "Anton Sedletskii"},
	}))
	engine := gin.New()
	engine.GET("/cv/:cvId", handler.GetResume)

	req := httptest.NewRequest("GET", "/cv/"+id.String()+"?lang=en", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), "Anton Sedletskii")
}

func TestGetResumeNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(stubResumeReader{err: queries.ErrResumeNotFound}))
	engine := gin.New()
	engine.GET("/cv/:cvId", handler.GetResume)

	req := httptest.NewRequest("GET", "/cv/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 404, w.Code)
	require.Contains(t, w.Body.String(), "CV not found")
}

func TestGetResumeRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(stubResumeReader{}))
	engine := gin.New()
	engine.GET("/cv/:cvId", handler.GetResume)

	req := httptest.NewRequest("GET", "/cv/not-a-uuid", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
	require.Contains(t, w.Body.String(), "Invalid CV id")
}

func TestGetResumeRejectsInvalidLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(stubResumeReader{}))
	engine := gin.New()
	engine.GET("/cv/:cvId", handler.GetResume)

	req := httptest.NewRequest("GET", "/cv/"+uuid.NewString()+"?lang=de", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
	require.Contains(t, w.Body.String(), "Invalid language")
}

var _ queries.ResumeReader = stubResumeReader{}
