package http

import (
	"context"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/platform/http/middleware"
	valueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type lifecycleReaderStub struct {
	items []queries.ResumeLinkListItem
	err   error
}

func (s lifecycleReaderStub) ListByResumeID(context.Context, queries.ListResumeLinksQuery) ([]queries.ResumeLinkListItem, error) {
	return s.items, s.err
}

type lifecycleOwnershipStub struct {
	err error
}

func (s lifecycleOwnershipStub) GetOwned(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	if s.err != nil {
		return resumequeries.ResumeView{}, s.err
	}
	return resumequeries.ResumeView{ID: "resume"}, nil
}

type lifecycleRepositoryStub struct {
	link           *domain.ResumeLink
	getErr         error
	delErr         error
	deleted        bool
	replaced       bool
	replaceErr     error
	replaceActive  bool
	replaceVersion int64
}

func (s *lifecycleRepositoryStub) Create(context.Context, *domain.ResumeLink) error { return nil }
func (s *lifecycleRepositoryStub) GetByID(context.Context, valueobjects.ID) (*domain.ResumeLink, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.link, nil
}
func (s *lifecycleRepositoryStub) ReplaceStatus(_ context.Context, _ valueobjects.ID, active bool, version int64) error {
	s.replaced = true
	s.replaceActive = active
	s.replaceVersion = version
	return s.replaceErr
}

func (s *lifecycleRepositoryStub) Delete(context.Context, valueobjects.ID) error {
	s.deleted = true
	return s.delErr
}
func (s *lifecycleRepositoryStub) RecordEvent(context.Context, *domain.ResumeLinkEvent) error {
	return nil
}

var _ ports.ResumeLinkRepository = (*lifecycleRepositoryStub)(nil)

func newLifecycleHandler(reader queries.ResumeLinkReader, ownership queries.ResumeOwnershipReader, repo ports.ResumeLinkRepository) *Handler {
	list := queries.NewListResumeLinksQueryHandler(reader, ownership)
	deleteHandler := commands.NewDeleteResumeLinkHandler(repo, ownership, lifecycleUnitOfWork{})
	return NewHandler(HandlerDependencies{List: list, Delete: deleteHandler})
}

type lifecycleUnitOfWork struct{}

func (lifecycleUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestReplaceStatusHandlerDeactivatesOwnedLink(t *testing.T) {
	userID := valueobjects.NewID()
	resumeID := valueobjects.NewID()
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: valueobjects.NewID(), ResumeID: resumeID, Name: "HH", Version: 3})
	require.NoError(t, err)
	repo := &lifecycleRepositoryStub{link: link}
	ownership := lifecycleOwnershipStub{}
	h := NewHandler(HandlerDependencies{
		Lifecycle: commands.NewReplaceStatusResumeLinkHandler(repo, ownership, lifecycleUnitOfWork{}),
	})
	engine := newAuthenticatedEngine(h, userID.UUID())

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/link/"+link.ID().String()+"/status", strings.NewReader(`{"status":"inactive","version":3}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.replaced)
	require.False(t, repo.replaceActive)
	require.Equal(t, int64(3), repo.replaceVersion)
}

func TestReplaceStatusHandlerUsesStandardErrorForInvalidVersion(t *testing.T) {
	h := NewHandler(HandlerDependencies{
		Lifecycle: commands.NewReplaceStatusResumeLinkHandler(
			&lifecycleRepositoryStub{},
			lifecycleOwnershipStub{},
			lifecycleUnitOfWork{},
		),
	})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/lk/link/"+uuid.NewString()+"/status",
		strings.NewReader(`{"status":"inactive","version":0}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"error":{"code":"BAD_REQUEST","message":"Version must be positive","requestId":""}}`, w.Body.String())
}

func TestReplaceStatusHandlerUsesStandardErrorForInvalidStatus(t *testing.T) {
	h := NewHandler(HandlerDependencies{
		Lifecycle: commands.NewReplaceStatusResumeLinkHandler(
			&lifecycleRepositoryStub{},
			lifecycleOwnershipStub{},
			lifecycleUnitOfWork{},
		),
	})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/lk/link/"+uuid.NewString()+"/status",
		strings.NewReader(`{"status":"paused","version":1}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"error":{"code":"BAD_REQUEST","message":"Invalid status","requestId":""}}`, w.Body.String())
}

func TestReplaceStatusHandlerRejectsInvalidBody(t *testing.T) {
	h := NewHandler(HandlerDependencies{Lifecycle: commands.NewReplaceStatusResumeLinkHandler(&lifecycleRepositoryStub{}, lifecycleOwnershipStub{}, lifecycleUnitOfWork{})})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/link/"+uuid.NewString()+"/status", strings.NewReader(`{"status":"paused","version":1}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReplaceStatusHandlerRejectsStaleVersion(t *testing.T) {
	resumeID := valueobjects.NewID()
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: valueobjects.NewID(), ResumeID: resumeID, Name: "HH", Version: 4})
	require.NoError(t, err)
	repo := &lifecycleRepositoryStub{link: link}
	h := NewHandler(HandlerDependencies{Lifecycle: commands.NewReplaceStatusResumeLinkHandler(repo, lifecycleOwnershipStub{}, lifecycleUnitOfWork{})})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/link/"+link.ID().String()+"/status", strings.NewReader(`{"status":"inactive","version":3}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusPreconditionFailed, w.Code)
	require.False(t, repo.replaced)
}

func newAuthenticatedEngine(handler *Handler, userID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	RegisterRoutes(handler, engine.Group("/api/v1"), engine.Group("/api/v1/lk"))
	return engine
}

func TestListHandlerReturnsLinks(t *testing.T) {
	userID := uuid.New()
	resumeID := valueobjects.NewID()
	h := newLifecycleHandler(lifecycleReaderStub{
		items: []queries.ResumeLinkListItem{{ID: uuid.NewString(), Name: "HH", Status: "active", Version: 1}},
	}, lifecycleOwnershipStub{}, &lifecycleRepositoryStub{})
	engine := newAuthenticatedEngine(h, userID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/links?resumeId="+resumeID.String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"name":"HH"`)
}

func TestListHandlerReturnsNotFoundForForeignResume(t *testing.T) {
	userID := uuid.New()
	h := newLifecycleHandler(lifecycleReaderStub{}, lifecycleOwnershipStub{err: resumequeries.ErrResumeNotFound}, &lifecycleRepositoryStub{})
	engine := newAuthenticatedEngine(h, userID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/links?resumeId="+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteHandlerDeletesOwnedLink(t *testing.T) {
	userID := valueobjects.NewID()
	resumeID := valueobjects.NewID()
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: valueobjects.NewID(), ResumeID: resumeID, Name: "HH"})
	require.NoError(t, err)
	repo := &lifecycleRepositoryStub{link: link}
	h := newLifecycleHandler(lifecycleReaderStub{}, lifecycleOwnershipStub{}, repo)
	engine := newAuthenticatedEngine(h, userID.UUID())

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/lk/link/"+link.ID().String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.deleted)
}

func TestDeleteHandlerReturnsNotFoundForMissingLink(t *testing.T) {
	h := newLifecycleHandler(lifecycleReaderStub{}, lifecycleOwnershipStub{}, &lifecycleRepositoryStub{getErr: domain.ErrEntityNotFound})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/lk/link/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestReplaceStatusHandlerReturnsPreconditionFailedEnvelope(t *testing.T) {
	resumeID := valueobjects.NewID()
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: valueobjects.NewID(), ResumeID: resumeID, Name: "HH", Version: 4})
	require.NoError(t, err)
	repo := &lifecycleRepositoryStub{link: link}
	h := NewHandler(HandlerDependencies{Lifecycle: commands.NewReplaceStatusResumeLinkHandler(repo, lifecycleOwnershipStub{}, lifecycleUnitOfWork{})})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/link/"+link.ID().String()+"/status", strings.NewReader(`{"status":"inactive","version":3}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusPreconditionFailed, w.Code)
	require.JSONEq(t, `{"error":{"code":"PRECONDITION_FAILED","message":"Resume link was modified","requestId":""}}`, w.Body.String())
	require.False(t, repo.replaced)
}
