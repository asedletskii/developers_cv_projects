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

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type lifecycleReaderStub struct {
	items []queries.ResumeLinkListItem
	err   error
}

func (s lifecycleReaderStub) ListByResumeID(context.Context, queries.ListResumeLinksQuery) ([]queries.ResumeLinkListItem, error) {
	return s.items, s.err
}

func newLifecycleHandler(
	reader queries.ResumeLinkReader,
	ownership interface {
		resumecontracts.ResumeOwnershipChecker
		resumecontracts.OwnedResumeLocker
	},
	repo ports.ResumeLinkRepository,
) *Handler {
	list := queries.NewListResumeLinksQueryHandler(reader, ownership)
	deleteHandler := commands.NewDeleteResumeLinkHandler(repo, ownership, &analyticsCleanupStub{}, unitOfWorkStub{})
	return NewHandler(HandlerDependencies{List: list, Delete: deleteHandler})
}

func restoreHTTPResumeLink(t *testing.T, resumeID sharedvalueobjects.ID, version int64) *domain.ResumeLink {
	t.Helper()
	link, err := domain.RestoreResumeLink(domain.ResumeLinkSnapshot{
		ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "HH",
		Status: domain.ResumeLinkStatusActive, Version: version,
	})
	require.NoError(t, err)
	return link
}

func TestReplaceStatusHandlerDeactivatesOwnedLink(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	link := restoreHTTPResumeLink(t, resumeID, 3)
	repo := &resumeLinkRepositoryStub{link: link}
	ownership := resumeOwnershipReaderStub{}
	h := NewHandler(HandlerDependencies{
		Lifecycle: commands.NewReplaceStatusResumeLinkHandler(repo, ownership, unitOfWorkStub{}),
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
			&resumeLinkRepositoryStub{},
			resumeOwnershipReaderStub{},
			unitOfWorkStub{},
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
			&resumeLinkRepositoryStub{},
			resumeOwnershipReaderStub{},
			unitOfWorkStub{},
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
	h := NewHandler(HandlerDependencies{Lifecycle: commands.NewReplaceStatusResumeLinkHandler(&resumeLinkRepositoryStub{}, resumeOwnershipReaderStub{}, unitOfWorkStub{})})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/link/"+uuid.NewString()+"/status", strings.NewReader(`{"status":"paused","version":1}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReplaceStatusHandlerRejectsStaleVersion(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	link := restoreHTTPResumeLink(t, resumeID, 4)
	repo := &resumeLinkRepositoryStub{link: link}
	h := NewHandler(HandlerDependencies{Lifecycle: commands.NewReplaceStatusResumeLinkHandler(repo, resumeOwnershipReaderStub{}, unitOfWorkStub{})})
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
	RegisterRoutes(handler, engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil)
	return engine
}

func TestListHandlerReturnsLinks(t *testing.T) {
	userID := uuid.New()
	resumeID := sharedvalueobjects.NewID()
	h := newLifecycleHandler(lifecycleReaderStub{
		items: []queries.ResumeLinkListItem{{ID: uuid.NewString(), Name: "HH", Status: "active", Version: 1}},
	}, resumeOwnershipReaderStub{}, &resumeLinkRepositoryStub{})
	engine := newAuthenticatedEngine(h, userID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/links?resumeId="+resumeID.String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"name":"HH"`)
}

func TestListHandlerReturnsNotFoundForForeignResume(t *testing.T) {
	userID := uuid.New()
	h := newLifecycleHandler(lifecycleReaderStub{}, resumeOwnershipReaderStub{err: resumecontracts.ErrResumeNotFound}, &resumeLinkRepositoryStub{})
	engine := newAuthenticatedEngine(h, userID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/links?resumeId="+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteHandlerDeletesOwnedLink(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	link := restoreHTTPResumeLink(t, resumeID, 1)
	repo := &resumeLinkRepositoryStub{link: link}
	h := newLifecycleHandler(lifecycleReaderStub{}, resumeOwnershipReaderStub{}, repo)
	engine := newAuthenticatedEngine(h, userID.UUID())

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/lk/link/"+link.ID().String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.deleted)
}

func TestDeleteHandlerReturnsNotFoundForMissingLink(t *testing.T) {
	h := newLifecycleHandler(lifecycleReaderStub{}, resumeOwnershipReaderStub{}, &resumeLinkRepositoryStub{getErr: domain.ErrEntityNotFound})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/lk/link/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestRenameHandlerRenamesOwnedLink(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	link := restoreHTTPResumeLink(t, resumeID, 3)
	repo := &resumeLinkRepositoryStub{link: link}
	h := NewHandler(HandlerDependencies{
		Rename: commands.NewRenameResumeLinkHandler(repo, resumeOwnershipReaderStub{}, unitOfWorkStub{}),
	})
	engine := newAuthenticatedEngine(h, userID.UUID())

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/link/"+link.ID().String(), strings.NewReader(`{"name":"New Name","version":3}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, repo.renamed)
	require.Equal(t, "New Name", repo.renameName)
	require.Equal(t, int64(3), repo.renameVersion)
}

func TestRenameHandlerUsesStandardErrorForInvalidVersion(t *testing.T) {
	h := NewHandler(HandlerDependencies{
		Rename: commands.NewRenameResumeLinkHandler(
			&resumeLinkRepositoryStub{},
			resumeOwnershipReaderStub{},
			unitOfWorkStub{},
		),
	})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/link/"+uuid.NewString(), strings.NewReader(`{"name":"New","version":0}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"error":{"code":"BAD_REQUEST","message":"Version must be positive","requestId":""}}`, w.Body.String())
}

func TestRenameHandlerRejectsStaleVersion(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	link := restoreHTTPResumeLink(t, resumeID, 4)
	repo := &resumeLinkRepositoryStub{link: link}
	h := NewHandler(HandlerDependencies{
		Rename: commands.NewRenameResumeLinkHandler(repo, resumeOwnershipReaderStub{}, unitOfWorkStub{}),
	})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/link/"+link.ID().String(), strings.NewReader(`{"name":"New","version":3}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusPreconditionFailed, w.Code)
	require.JSONEq(t, `{"error":{"code":"PRECONDITION_FAILED","message":"Resume link was modified","requestId":""}}`, w.Body.String())
	require.False(t, repo.renamed)
}

func TestRenameHandlerRejectsInvalidName(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	link := restoreHTTPResumeLink(t, resumeID, 1)
	repo := &resumeLinkRepositoryStub{link: link}
	h := NewHandler(HandlerDependencies{
		Rename: commands.NewRenameResumeLinkHandler(repo, resumeOwnershipReaderStub{}, unitOfWorkStub{}),
	})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/link/"+link.ID().String(), strings.NewReader(`{"name":"   ","version":1}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"error":{"code":"BAD_REQUEST","message":"Invalid resume link","requestId":""}}`, w.Body.String())
	require.False(t, repo.renamed)
}

func TestReplaceStatusHandlerReturnsPreconditionFailedEnvelope(t *testing.T) {
	resumeID := sharedvalueobjects.NewID()
	link := restoreHTTPResumeLink(t, resumeID, 4)
	repo := &resumeLinkRepositoryStub{link: link}
	h := NewHandler(HandlerDependencies{Lifecycle: commands.NewReplaceStatusResumeLinkHandler(repo, resumeOwnershipReaderStub{}, unitOfWorkStub{})})
	engine := newAuthenticatedEngine(h, uuid.New())

	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/link/"+link.ID().String()+"/status", strings.NewReader(`{"status":"inactive","version":3}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusPreconditionFailed, w.Code)
	require.JSONEq(t, `{"error":{"code":"PRECONDITION_FAILED","message":"Resume link was modified","requestId":""}}`, w.Body.String())
	require.False(t, repo.replaced)
}
