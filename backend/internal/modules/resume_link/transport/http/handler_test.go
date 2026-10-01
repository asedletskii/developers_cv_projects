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
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/platform/http/middleware"
	"developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeResumeOwnershipReader struct{ err error }

func (r fakeResumeOwnershipReader) GetOwned(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	if r.err != nil {
		return resumequeries.ResumeView{}, r.err
	}
	return resumequeries.ResumeView{ID: "resume"}, nil
}

type fakeRepository struct{ err error }

func (r *fakeRepository) Create(context.Context, *domain.ResumeLink) error { return r.err }

func (r *fakeRepository) GetByID(context.Context, valueobjects.ID) (*domain.ResumeLink, error) {
	return nil, r.err
}
func (r *fakeRepository) ReplaceStatus(context.Context, valueobjects.ID, bool, int64) error {
	return r.err
}

func (r *fakeRepository) Delete(context.Context, valueobjects.ID) error              { return r.err }
func (r *fakeRepository) RecordEvent(context.Context, *domain.ResumeLinkEvent) error { return r.err }

var _ ports.ResumeLinkRepository = (*fakeRepository)(nil)

func newTestHandler(err error) *Handler {
	commandHandler := commands.NewCreateResumeLinkHandler(
		&fakeRepository{err: err},
		fakeResumeOwnershipReader{},
		fakeUnitOfWork{},
	)
	return NewHandler(HandlerDependencies{Create: commandHandler})
}

type fakeUnitOfWork struct{}

func (fakeUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

func TestCreateHandlerRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, middleware.ErrUnauthenticated
	}))
	RegisterRoutes(newTestHandler(nil), engine.Group("/api/v1"), engine.Group("/api/v1/lk"))

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
	RegisterRoutes(newTestHandler(nil), engine.Group("/api/v1"), engine.Group("/api/v1/lk"))

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
	RegisterRoutes(handler, engine.Group("/api/v1"), engine.Group("/api/v1/lk"))

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

	missingResumeHandler := NewHandler(HandlerDependencies{Create: commands.NewCreateResumeLinkHandler(&fakeRepository{}, fakeResumeOwnershipReader{err: resumequeries.ErrResumeNotFound}, fakeUnitOfWork{})})
	engine = gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	RegisterRoutes(missingResumeHandler, engine.Group("/api/v1"), engine.Group("/api/v1/lk"))
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
	RegisterRoutes(newTestHandler(errors.New("db down")), engine.Group("/api/v1"), engine.Group("/api/v1/lk"))

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
	RegisterRoutes(newTestHandler(domain.ErrResumeLinkLimitExceeded), engine.Group("/api/v1"), engine.Group("/api/v1/lk"))

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
	RegisterRoutes(newTestHandler(domain.ErrDuplicateResumeLink), engine.Group("/api/v1"), engine.Group("/api/v1/lk"))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/link", bytes.NewBufferString(`{"resumeId":"`+uuid.New().String()+`","name":"HH"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
}
