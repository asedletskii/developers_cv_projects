package patch

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

func (r *updateResumeRepositoryStub) UpdateTitle(_ context.Context, command commands.UpdateResumeTitleRepositoryCommand) error {
	if r.err != nil {
		return r.err
	}
	if r.resume != nil {
		if err := r.resume.Rename(command.Title); err != nil {
			return err
		}
	}
	return nil
}

type updateResumeUoWStub struct{}

func (updateResumeUoWStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
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

func TestHandlerChangesTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	owner, err := useridentity.IDFromUUID(userID)
	require.NoError(t, err)

	repo := &updateResumeRepositoryStub{
		resume: newHTTPTestResume(t, resumeID, owner, "Old title"),
	}
	handler := NewHandler(commands.NewPatchResumeHandler(repo, updateResumeUoWStub{}))

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.PATCH("/api/v1/lk/resume/:resumeId", handler.Handle)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/lk/resume/"+resumeID.UUID().String(),
		strings.NewReader(`{"title":"New title"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"id":"`+resumeID.String()+`"`)
	require.Contains(t, w.Body.String(), `"title":"New title"`)
	require.Equal(t, "New title", repo.resume.Title())
	require.Equal(t, "Senior Go Developer", repo.resume.Position())
	require.Nil(t, repo.saved)
}

func TestHandlerRejectsForeignResume(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)

	repo := &updateResumeRepositoryStub{
		err: domain.ErrEntityNotFound,
	}
	handler := NewHandler(commands.NewPatchResumeHandler(repo, updateResumeUoWStub{}))

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.PATCH("/api/v1/lk/resume/:resumeId", handler.Handle)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/lk/resume/"+resumeID.UUID().String(),
		strings.NewReader(`{"title":"New title"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Nil(t, repo.saved)
}

func TestHandlerRejectsUnknownFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	owner, err := useridentity.IDFromUUID(userID)
	require.NoError(t, err)

	repo := &updateResumeRepositoryStub{
		resume: newHTTPTestResume(t, resumeID, owner, "Old title"),
	}
	handler := NewHandler(commands.NewPatchResumeHandler(repo, updateResumeUoWStub{}))

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.PATCH("/api/v1/lk/resume/:resumeId", handler.Handle)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/lk/resume/"+resumeID.UUID().String(),
		strings.NewReader(`{"title":"New title","userId":"00000000-0000-0000-0000-000000000099"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, "Old title", repo.resume.Title())
}
