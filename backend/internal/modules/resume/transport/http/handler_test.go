package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/platform/http/middleware"
)

type stubResumeReader struct {
	resume queries.Resume
	err    error
	query  queries.GetResumeQuery
}

func (r *stubResumeReader) GetResume(_ context.Context, query queries.GetResumeQuery) (queries.Resume, error) {
	r.query = query
	return r.resume, r.err
}

func TestGetResume(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.New()
	userID := uuid.New()
	reader := &stubResumeReader{
		resume: queries.Resume{ID: id.String(), Name: "Anton Sedletskii"},
	}
	handler := NewHandler(queries.NewGetResumeHandler(reader))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.GET("/cv/:cvId", handler.GetResume)

	req := httptest.NewRequest("GET", "/cv/"+id.String()+"?lang=en", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), "Anton Sedletskii")
	require.Equal(t, userID, reader.query.UserID.UUID())
}

func TestGetResumeNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(&stubResumeReader{err: queries.ErrResumeNotFound}))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.GET("/cv/:cvId", handler.GetResume)

	req := httptest.NewRequest("GET", "/cv/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 404, w.Code)
	require.Contains(t, w.Body.String(), "CV not found")
}

func TestGetResumeRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(&stubResumeReader{}))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.GET("/cv/:cvId", handler.GetResume)

	req := httptest.NewRequest("GET", "/cv/not-a-uuid", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
	require.Contains(t, w.Body.String(), "Invalid CV id")
}

func TestGetResumeRejectsInvalidLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(&stubResumeReader{}))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.GET("/cv/:cvId", handler.GetResume)

	req := httptest.NewRequest("GET", "/cv/"+uuid.NewString()+"?lang=de", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
	require.Contains(t, w.Body.String(), "Invalid language")
}

var _ queries.ResumeReader = &stubResumeReader{}

type fakeCreateResumeRepository struct{ saved *domain.Resume }

func (f *fakeCreateResumeRepository) GetByID(
	context.Context,
	useridentity.ID,
	valueobjects.ID,
) (*domain.Resume, error) {
	return nil, nil
}
func (f *fakeCreateResumeRepository) Save(_ context.Context, resume *domain.Resume) error {
	f.saved = resume
	return nil
}

type fakeCreateResumeUoW struct{ calls int }

func (f *fakeCreateResumeUoW) Do(ctx context.Context, fn func(context.Context) error) error {
	f.calls++
	return fn(ctx)
}

func TestCreateResume(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &fakeCreateResumeRepository{}
	uow := &fakeCreateResumeUoW{}
	createHandler := commands.NewCreateResumeHandler(repo, uow)
	getHandler := queries.NewGetResumeHandler(&stubResumeReader{})
	handler := NewHandler(getHandler, createHandler)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.POST("/cv", handler.CreateResume)

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
	req := httptest.NewRequest("POST", "/cv", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, 1, uow.calls)
	require.NotNil(t, repo.saved)
	require.False(t, repo.saved.ID().IsZero())
	require.NotEqual(t, uuid.Nil, repo.saved.UserID().UUID())
}

func TestCreateResumeRejectsInvalidDomainInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &fakeCreateResumeRepository{}
	uow := &fakeCreateResumeUoW{}
	handler := NewHandler(
		queries.NewGetResumeHandler(&stubResumeReader{}),
		commands.NewCreateResumeHandler(repo, uow),
	)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.POST("/cv", handler.CreateResume)

	body := `{"name":"","surname":"User","position":"Developer","greeting":{"first":{"ru":"Привет","en":"Hello"},"second":{"ru":"Я","en":"I am"},"secondHighlight":{"ru":"Тест","en":"Test"},"third":{"ru":"Разработчик","en":"Developer"}},"photoPath":"/assets/photo.svg","dateOfBirth":"1998-04-12","location":"Remote","contacts":{"phone":"+79999999999"}}`
	req := httptest.NewRequest("POST", "/cv", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Nil(t, repo.saved)
}
