package update

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
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/platform/http/middleware"
)

type updateResumeRepositoryStub struct {
	resume *domain.Resume
	saved  *domain.Resume
	err    error
}

func (r *updateResumeRepositoryStub) Get(_ context.Context, _ commands.GetResumeRepositoryQuery) (*domain.Resume, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.resume, nil
}

func (r *updateResumeRepositoryStub) Update(_ context.Context, command commands.UpdateResumeRepositoryCommand) error {
	r.saved = command.Resume
	return nil
}

type updateResumeUoWStub struct{}

func (updateResumeUoWStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type stubGetResumeReader struct {
	resume queries.Resume
	err    error
}

func (r *stubGetResumeReader) GetResume(_ context.Context, _ queries.GetResumeQuery) (queries.Resume, error) {
	return r.resume, r.err
}

func newHTTPTestResume(t *testing.T, id valueobjects.ID, userID useridentity.ID, title string) *domain.Resume {
	t.Helper()

	date, err := valueobjects.ParseDate("1998-04-12")
	require.NoError(t, err)
	photo, err := valueobjects.NewAssetPath("/assets/photo.svg")
	require.NoError(t, err)

	resume, err := domain.NewResume(domain.ResumeParams{
		ID:          id,
		UserID:      userID,
		Name:        "User",
		Surname:     "Test",
		Position:    "Senior Go Developer",
		Title:       title,
		PhotoPath:   photo,
		DateOfBirth: date,
		Location:    "Remote",
	})
	require.NoError(t, err)

	return resume
}

func TestHandlerUpdatesResumeAndPreservesTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	owner, err := useridentity.IDFromUUID(userID)
	require.NoError(t, err)

	repo := &updateResumeRepositoryStub{
		resume: newHTTPTestResume(t, resumeID, owner, "Manual title"),
	}
	getReader := &stubGetResumeReader{
		resume: queries.Resume{
			ID:   resumeID.String(),
			Name: "Updated",
		},
	}
	handler := NewHandler(
		commands.NewUpdateResumeHandler(repo, repo, updateResumeUoWStub{}),
		queries.NewGetResumeHandler(getReader),
	)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.Handle)

	body := `{
		"name":"Updated",
		"surname":"User",
		"position":"Staff Go Developer",
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
	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/lk/resume/"+resumeID.UUID().String(),
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"id":"`+resumeID.String()+`"`)
	require.NotNil(t, repo.saved)
	require.Equal(t, "Updated", repo.saved.Name())
	require.Equal(t, "Staff Go Developer", repo.saved.Position())
	require.Equal(t, "Manual title", repo.saved.Title())
}

func TestHandlerRejectsMalformedID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(
		commands.NewUpdateResumeHandler(
			&updateResumeRepositoryStub{},
			&updateResumeRepositoryStub{},
			updateResumeUoWStub{},
		),
		queries.NewGetResumeHandler(&stubGetResumeReader{}),
	)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.Handle)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/lk/resume/not-a-uuid",
		strings.NewReader(`{}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerRejectsUnknownFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &updateResumeRepositoryStub{}
	handler := NewHandler(
		commands.NewUpdateResumeHandler(repo, repo, updateResumeUoWStub{}),
		queries.NewGetResumeHandler(&stubGetResumeReader{}),
	)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))

	resumeID := uuid.NewString()
	engine = gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.Handle)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/lk/resume/"+resumeID,
		strings.NewReader(`{"name":"Test","unknown":"field"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Nil(t, repo.saved)
}
