package http

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
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/platform/http/middleware"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type stubListReader struct {
	items  []queries.ResumeListItem
	userID useridentity.ID
}

func (r *stubListReader) ListResumes(_ context.Context, query queries.ListResumesQuery) ([]queries.ResumeListItem, error) {
	r.userID = query.UserID
	return r.items, nil
}

func TestHandlerRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newResumeHandler(nil, nil, nil, nil, queries.NewListResumesQueryHandler(&stubListReader{}))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, middleware.ErrUnauthenticated
	}))
	engine.GET("/api/v1/lk/resumes", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/resumes", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandlerReturnsCurrentUsersResumes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	reader := &stubListReader{
		items: []queries.ResumeListItem{
			{ID: "resume-a", Title: "Senior Go Developer"},
		},
	}
	handler := newResumeHandler(nil, nil, nil, nil, queries.NewListResumesQueryHandler(reader))

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.GET("/api/v1/lk/resumes", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/resumes", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"resume-a"`)
	require.Contains(t, w.Body.String(), `"Senior Go Developer"`)
	require.Equal(t, userID, reader.userID.UUID())
}

type stubReplaceResumeTitleRepository struct {
	cmd ports.ReplaceResumeTitleRepositoryCommand
}

func (r *stubReplaceResumeTitleRepository) ReplaceResumeTitle(_ context.Context, cmd ports.ReplaceResumeTitleRepositoryCommand) error {
	r.cmd = cmd
	return nil
}

type stubResumeUnitOfWork struct{}

func (stubResumeUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestHandlerReplacesResumeTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID := uuid.New()
	repository := &stubReplaceResumeTitleRepository{}
	replaceHandler := commands.NewReplaceResumeTitleHandler(repository, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, replaceHandler, nil, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.ReplaceTitle)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+resumeID.String(), bytes.NewBufferString(`{"title":"Senior Go Developer"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Equal(t, resumeID, repository.cmd.ResumeID.UUID())
	require.Equal(t, userID, repository.cmd.UserID.UUID())
	require.Equal(t, "Senior Go Developer", repository.cmd.Title)
}

func TestHandlerRejectsLanguageInResumeTitleReplacement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	replaceHandler := commands.NewReplaceResumeTitleHandler(&stubReplaceResumeTitleRepository{}, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, replaceHandler, nil, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.ReplaceTitle)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+uuid.New().String(), bytes.NewBufferString(`{"title":"Senior Go Developer","language":"en"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerReturnsBadRequestForInvalidResumeTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	replaceHandler := commands.NewReplaceResumeTitleHandler(&stubReplaceResumeTitleRepository{}, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, replaceHandler, nil, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.ReplaceTitle)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+uuid.New().String(), bytes.NewBufferString(`{"title":"   "}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
