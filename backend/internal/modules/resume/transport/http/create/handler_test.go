package create

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/platform/http/middleware"
)

type fakeCreateResumeRepository struct {
	saved *domain.Resume
}

func (f *fakeCreateResumeRepository) Create(_ context.Context, command commands.CreateResumeRepositoryCommand) error {
	f.saved = command.Resume
	return nil
}

type fakeCreateResumeUoW struct {
	calls int
}

func (f *fakeCreateResumeUoW) Do(ctx context.Context, fn func(context.Context) error) error {
	f.calls++
	return fn(ctx)
}

type stubGetResumeReader struct {
	resume queries.Resume
	err    error
}

func (r *stubGetResumeReader) GetResume(_ context.Context, _ queries.GetResumeQuery) (queries.Resume, error) {
	return r.resume, r.err
}

func TestHandlerCreatesResume(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeCreateResumeRepository{}
	uow := &fakeCreateResumeUoW{}
	getHandler := queries.NewGetResumeHandler(&stubGetResumeReader{
		resume: queries.Resume{ID: "created", Name: "Test"},
	})
	handler := NewHandler(
		commands.NewCreateResumeHandler(repo, uow),
		getHandler,
	)

	engine := gin.New()
	userID := uuid.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.POST("/api/v1/lk/resume", handler.Handle)

	body := `{
		"name":"Test",
		"surname":"User",
		"position":"Developer",
		"greeting":{
			"first":{"ru":"Привет","en":"Hello"},
			"second":{"ru":"Я","en":"I am"},
			"secondHighlight":{"ru":"Тест","en":"Test"},
			"third":{"ru":"Разработчик","en":"Developer"}
		},
		"photoPath":"/assets/photo.svg",
		"dateOfBirth":"1998-04-12",
		"location":"Remote",
		"contacts":{"phone":"+79999999999"}
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/resume", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Contains(t, w.Body.String(), `"id":"created"`)
	require.Equal(t, 1, uow.calls)
	require.NotNil(t, repo.saved)
	require.False(t, repo.saved.ID().IsZero())
	require.Equal(t, userID, repo.saved.UserID().UUID())
	require.Equal(t, "Developer", repo.saved.Title())
}

func TestHandlerRejectsInvalidDomainInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeCreateResumeRepository{}
	uow := &fakeCreateResumeUoW{}
	handler := NewHandler(
		commands.NewCreateResumeHandler(repo, uow),
		queries.NewGetResumeHandler(&stubGetResumeReader{}),
	)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.POST("/api/v1/lk/resume", handler.Handle)

	body := `{
		"name":"",
		"surname":"User",
		"position":"Developer",
		"greeting":{
			"first":{"ru":"Привет","en":"Hello"},
			"second":{"ru":"Я","en":"I am"},
			"secondHighlight":{"ru":"Тест","en":"Test"},
			"third":{"ru":"Разработчик","en":"Developer"}
		},
		"photoPath":"/assets/photo.svg",
		"dateOfBirth":"1998-04-12",
		"location":"Remote",
		"contacts":{"phone":"+79999999999"}
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/resume", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Nil(t, repo.saved)
}

func TestHandlerRejectsUnknownFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeCreateResumeRepository{}
	uow := &fakeCreateResumeUoW{}
	handler := NewHandler(
		commands.NewCreateResumeHandler(repo, uow),
		queries.NewGetResumeHandler(&stubGetResumeReader{}),
	)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.POST("/api/v1/lk/resume", handler.Handle)

	body := `{"name":"Test","surname":"User","position":"Developer","unknown":"field"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/resume", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Nil(t, repo.saved)
}
