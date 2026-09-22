package http

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/mock"
)

func TestGetPerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mock.NewPersonReader()
	handler := NewHandler(queries.NewGetPersonHandler(repo))
	engine := gin.New()
	engine.GET("/cv/:cvId", handler.GetPerson)

	req := httptest.NewRequest("GET", "/cv/"+mock.PersonOneID+"?lang=en", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), "Anton Sedletskii")
	require.Contains(t, w.Body.String(), "Looking for a Senior Frontend role")
}

func TestGetPersonNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mock.NewPersonReader()
	handler := NewHandler(queries.NewGetPersonHandler(repo))
	engine := gin.New()
	engine.GET("/cv/:cvId", handler.GetPerson)

	req := httptest.NewRequest("GET", "/cv/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 404, w.Code)
	require.Contains(t, w.Body.String(), "CV not found")
}

func TestGetPersonRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mock.NewPersonReader()
	handler := NewHandler(queries.NewGetPersonHandler(repo))
	engine := gin.New()
	engine.GET("/cv/:cvId", handler.GetPerson)

	req := httptest.NewRequest("GET", "/cv/not-a-uuid", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
	require.Contains(t, w.Body.String(), "Invalid CV id")
}

func TestGetPersonRejectsInvalidLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mock.NewPersonReader()
	handler := NewHandler(queries.NewGetPersonHandler(repo))
	engine := gin.New()
	engine.GET("/cv/:cvId", handler.GetPerson)

	req := httptest.NewRequest("GET", "/cv/"+mock.PersonOneID+"?lang=de", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
	require.Contains(t, w.Body.String(), "Invalid language")
}
