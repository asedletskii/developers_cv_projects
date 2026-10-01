package http

import (
	"bytes"
	"context"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type stubListReader struct {
	items  []queries.ResumeListItem
	userID sharedvalueobjects.ID
}

func (r *stubListReader) ListResumes(_ context.Context, query queries.ListResumesQuery) ([]queries.ResumeListItem, error) {
	r.userID = query.UserID
	return r.items, nil
}

type stubOwnedResumeReader struct {
	view queries.ResumeView
}

func (r stubOwnedResumeReader) GetOwned(context.Context, queries.GetResumeQuery) (queries.ResumeView, error) {
	return r.view, nil
}

func TestResumeMutationHandlersRejectNonPositiveVersionBeforeApplication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID := uuid.New()
	handler := NewHandler(HandlerDependencies{})

	routes := []struct {
		method string
		path   string
	}{
		{http.MethodPut, "/resume/" + resumeID.String()},
		{http.MethodPut, "/resume/" + resumeID.String() + "/status"},
		{http.MethodPut, "/resume/" + resumeID.String() + "/overview"},
		{http.MethodPut, "/resume/" + resumeID.String() + "/profile"},
		{http.MethodPut, "/resume/" + resumeID.String() + "/contacts"},
		{http.MethodPut, "/resume/" + resumeID.String() + "/skills"},
		{http.MethodPut, "/resume/" + resumeID.String() + "/languages"},
		{http.MethodPut, "/resume/" + resumeID.String() + "/technologies"},
		{http.MethodPut, "/resume/" + resumeID.String() + "/experiences"},
		{http.MethodPut, "/resume/" + resumeID.String() + "/educations"},
		{http.MethodPut, "/resume/" + resumeID.String() + "/courses"},
		{http.MethodPut, "/resume/" + resumeID.String() + "/activities"},
		{http.MethodPut, "/resume/" + resumeID.String() + "/projects"},
		{http.MethodDelete, "/resume/" + resumeID.String()},
	}

	for _, route := range routes {
		t.Run(route.path, func(t *testing.T) {
			engine := gin.New()
			engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
				return userID, nil
			}))
			RegisterRoutes(handler, engine.Group("/"))

			req := httptest.NewRequest(route.method, route.path, bytes.NewBufferString(`{"version":0}`))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), `"code":"BAD_REQUEST"`)
		})
	}
}

func TestHandlerDeletesResumeWithVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID := uuid.New()
	resume := newStubResume(t, resumeID, userID, 4)
	repository := &stubResumeRepository{resume: resume}
	deleteHandler := commands.NewDeleteResumeHandler(repository, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, nil, deleteHandler, nil, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.DELETE("/api/v1/lk/resume/:resumeId", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/lk/resume/"+resumeID.String(), bytes.NewBufferString(`{"version":4}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestHandlerRejectsDeleteWithNonPositiveVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID := uuid.New()
	deleteHandler := commands.NewDeleteResumeHandler(&stubResumeRepository{}, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, nil, deleteHandler, nil, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.DELETE("/api/v1/lk/resume/:resumeId", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/lk/resume/"+resumeID.String(), bytes.NewBufferString(`{"version":0}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), `"Version must be positive"`)
}

func TestHandlerMapsDeleteConcurrencyConflictToPreconditionFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID := uuid.New()
	repository := &stubResumeRepository{deleteErr: domain.ErrConcurrencyConflict}
	deleteHandler := commands.NewDeleteResumeHandler(repository, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, nil, deleteHandler, nil, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.DELETE("/api/v1/lk/resume/:resumeId", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/lk/resume/"+resumeID.String(), bytes.NewBufferString(`{"version":4}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusPreconditionFailed, w.Code)
	require.Contains(t, w.Body.String(), `"Resume was modified"`)
}

func TestHandlerReturnsResumeVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID := uuid.New()
	getHandler := queries.NewGetResumeQueryHandler(stubOwnedResumeReader{
		view: queries.ResumeView{ID: resumeID.String(), Title: "CV", Version: int64Ptr(42)},
	})
	handler := newResumeHandler(nil, getHandler, nil, nil, nil, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.GET("/api/v1/lk/resume/:resumeId", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/resume/"+resumeID.String(), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"version":42`)
}

func TestHandlerRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newResumeHandler(nil, nil, nil, nil, nil, queries.NewListResumesQueryHandler(&stubListReader{}))
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
			{ID: "resume-a", Title: "Senior Go Developer", Version: 3},
		},
	}
	handler := newResumeHandler(nil, nil, nil, nil, nil, queries.NewListResumesQueryHandler(reader))

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
	require.Contains(t, w.Body.String(), `"version":3`)
	require.Equal(t, userID, reader.userID.UUID())
}

type stubResumeRepository struct {
	resume    *domain.Resume
	saved     *domain.Resume
	version   int64
	deleteErr error
}

func (r *stubResumeRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}

func (r *stubResumeRepository) Save(_ context.Context, resume *domain.Resume) error {
	r.saved = resume
	r.version = resume.Version()
	return nil
}

func (r *stubResumeRepository) Create(context.Context, *domain.Resume) error { return nil }
func (r *stubResumeRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return r.deleteErr
}

type stubResumeUnitOfWork struct{}

func (stubResumeUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func newStubResume(t *testing.T, resumeID uuid.UUID, userID uuid.UUID, version int64) *domain.Resume {
	t.Helper()
	id, err := sharedvalueobjects.IDFromUUID(resumeID)
	require.NoError(t, err)
	uid, err := sharedvalueobjects.IDFromUUID(userID)
	require.NoError(t, err)
	resume, err := domain.NewResume(domain.ResumeParams{
		ID: id, UserID: uid, Title: "Old title", Language: domain.ResumeLanguageEN, Version: version,
	})
	require.NoError(t, err)
	return resume
}

func TestHandlerReplacesResumeTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID := uuid.New()
	resume := newStubResume(t, resumeID, userID, 4)
	repository := &stubResumeRepository{resume: resume}
	replaceHandler := commands.NewReplaceResumeTitleHandler(repository, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, replaceHandler, nil, nil, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.ReplaceTitle)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+resumeID.String(), bytes.NewBufferString(`{"version":4,"title":"Senior Go Developer"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Equal(t, "Senior Go Developer", repository.saved.Title())
	require.Equal(t, int64(4), repository.version)
}

func TestHandlerRejectsLanguageInResumeTitleReplacement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	replaceHandler := commands.NewReplaceResumeTitleHandler(&stubResumeRepository{}, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, replaceHandler, nil, nil, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.ReplaceTitle)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+uuid.New().String(), bytes.NewBufferString(`{"version":1,"title":"Senior Go Developer","language":"en"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerMapsResumeTitleConcurrencyConflictToPreconditionFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID := uuid.New()
	resume := newStubResume(t, resumeID, userID, 2)
	repository := &stubResumeRepository{resume: resume}
	replaceHandler := commands.NewReplaceResumeTitleHandler(repository, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, replaceHandler, nil, nil, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.ReplaceTitle)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+resumeID.String(), bytes.NewBufferString(`{"version":1,"title":"New title"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusPreconditionFailed, w.Code)
	require.Contains(t, w.Body.String(), `"code":"PRECONDITION_FAILED"`)
	require.Nil(t, repository.saved)
}

func TestHandlerReturnsBadRequestForInvalidResumeTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	resumeID := uuid.New()
	resume := newStubResume(t, resumeID, userID, 1)
	repository := &stubResumeRepository{resume: resume}
	replaceHandler := commands.NewReplaceResumeTitleHandler(repository, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, replaceHandler, nil, nil, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.ReplaceTitle)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+resumeID.String(), bytes.NewBufferString(`{"version":1,"title":"   "}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func int64Ptr(value int64) *int64 {
	return &value
}

func TestHandlerChangesResumeStatusAndReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	resumeID := uuid.New()
	resume := newStubResume(t, resumeID, userID, 4)
	repository := &stubResumeRepository{resume: resume}
	statusHandler := commands.NewChangeResumeStatusHandler(repository, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, nil, nil, statusHandler, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	engine.PUT("/api/v1/lk/resume/:resumeId/status", handler.ChangeStatus)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+resumeID.String()+"/status", bytes.NewBufferString(`{"status":"opened","version":4}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Empty(t, w.Body.String())
	require.Equal(t, domain.ResumeStatusOpened, repository.saved.Status())
}

func TestHandlerRejectsInvalidResumeStatusTransition(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	resumeID := uuid.New()
	resume := newStubResume(t, resumeID, userID, 4)
	statusHandler := commands.NewChangeResumeStatusHandler(&stubResumeRepository{resume: resume}, stubResumeUnitOfWork{})
	handler := newResumeHandler(nil, nil, nil, nil, statusHandler, nil)

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	engine.PUT("/api/v1/lk/resume/:resumeId/status", handler.ChangeStatus)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+resumeID.String()+"/status", bytes.NewBufferString(`{"status":"closed","version":4}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
}
