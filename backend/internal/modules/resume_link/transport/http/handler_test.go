package http

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/platform/http/middleware"
)

func newTestHandler(err error) *Handler {
	commandHandler := commands.NewCreateResumeLinkHandler(
		&resumeLinkRepositoryStub{err: err},
		resumeOwnershipReaderStub{},
		unitOfWorkStub{},
	)
	return NewHandler(HandlerDependencies{Create: commandHandler})
}

func TestCreateHandlerRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, middleware.ErrUnauthenticated
	}))
	RegisterRoutes(newTestHandler(nil), engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/link", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateHandlerReturnsCreatedID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	resumeID := uuid.New()
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	RegisterRoutes(newTestHandler(nil), engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/link", bytes.NewBufferString(`{"resumeId":"`+resumeID.String()+`","name":"HH"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Contains(t, w.Body.String(), `"id"`)
}

func TestCreateHandlerMapsBadRequestAndNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	handler := newTestHandler(nil)
	RegisterRoutes(handler, engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil)

	cases := []struct {
		name   string
		body   string
		status int
	}{
		{"malformed json", `{"resumeId":`, http.StatusBadRequest},
		{"invalid uuid", `{"resumeId":"nope","name":"HH"}`, http.StatusBadRequest},
		{"missing name", `{"resumeId":"` + uuid.New().String() + `","name":"   "}`, http.StatusBadRequest},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/link", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			require.Equal(t, tt.status, w.Code)
		})
	}

	missingResumeHandler := NewHandler(HandlerDependencies{Create: commands.NewCreateResumeLinkHandler(&resumeLinkRepositoryStub{}, resumeOwnershipReaderStub{err: resumequeries.ErrResumeNotFound}, unitOfWorkStub{})})
	engine = gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	RegisterRoutes(missingResumeHandler, engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/link", bytes.NewBufferString(`{"resumeId":"`+uuid.New().String()+`","name":"HH"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateHandlerMapsUnexpectedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	RegisterRoutes(newTestHandler(errors.New("db down")), engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/link", bytes.NewBufferString(`{"resumeId":"`+uuid.New().String()+`","name":"HH"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.NotContains(t, w.Body.String(), "db down")
	require.Contains(t, w.Body.String(), `"code":"INTERNAL_ERROR"`)
}

func TestCreateHandlerMapsLimitExceeded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	RegisterRoutes(newTestHandler(domain.ErrResumeLinkLimitExceeded), engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/link", bytes.NewBufferString(`{"resumeId":"`+uuid.New().String()+`","name":"HH"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
	require.Contains(t, w.Body.String(), "Resume link limit exceeded")
}

func TestCreateHandlerMapsDuplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	RegisterRoutes(newTestHandler(domain.ErrDuplicateResumeLink), engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/link", bytes.NewBufferString(`{"resumeId":"`+uuid.New().String()+`","name":"HH"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
}
