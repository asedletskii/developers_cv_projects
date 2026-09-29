package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumevalueobjects "developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/domain/valueobjects"
	"developers-cv/backend/internal/platform/http/middleware"
)

type statisticsHTTPLinkReader struct {
	view queries.ResumeLinkPublicView
	err  error
}

func (s statisticsHTTPLinkReader) GetByID(context.Context, queries.PublicResumeLinkQuery) (queries.ResumeLinkPublicView, error) {
	return s.view, s.err
}

type statisticsHTTPOwnershipReader struct {
	err error
}

func (s statisticsHTTPOwnershipReader) GetOwned(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	if s.err != nil {
		return resumequeries.ResumeView{}, s.err
	}
	return resumequeries.ResumeView{ID: "resume"}, nil
}

type statisticsHTTPReader struct {
	result queries.ResumeLinkStatistics
	err    error
}

func (s statisticsHTTPReader) GetStatistics(context.Context, queries.GetResumeLinkStatisticsQuery) (queries.ResumeLinkStatistics, error) {
	return s.result, s.err
}

func newStatisticsHTTPHandler(
	linkReader queries.ResumeLinkPublicReader,
	ownership ports.ResumeOwnershipReader,
	reader queries.ResumeLinkStatisticsReader,
) *Handler {
	statistics := queries.NewGetResumeLinkStatisticsQueryHandler(linkReader, ownership, reader)
	return NewHandler(HandlerDependencies{Statistics: statistics})
}

func newStatisticsHTTPServer(handler *Handler, userID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	RegisterRoutes(handler, engine.Group("/api/v1"), engine.Group("/api/v1/lk"))
	return engine
}

func TestStatisticsHandlerReturnsStatistics(t *testing.T) {
	userID := uuid.New()
	linkID := valueobjects.NewID()
	resumeID := resumevalueobjects.NewID()
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{view: queries.ResumeLinkPublicView{ResumeID: resumeID.String()}},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{result: queries.ResumeLinkStatistics{Views: 125, ContactViews: 17}},
	), userID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+linkID.String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"views":125,"contactViews":17}`, w.Body.String())
}

func TestStatisticsHandlerReturnsZeroState(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{view: queries.ResumeLinkPublicView{ResumeID: resumevalueobjects.NewID().String()}},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{result: queries.ResumeLinkStatistics{}},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+valueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"views":0,"contactViews":0}`, w.Body.String())
}

func TestStatisticsHandlerRejectsForeignLink(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{view: queries.ResumeLinkPublicView{ResumeID: resumevalueobjects.NewID().String()}},
		statisticsHTTPOwnershipReader{err: resumequeries.ErrResumeNotFound},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+valueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestStatisticsHandlerReturnsNotFoundForMissingLink(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{err: domain.ErrEntityNotFound},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+valueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestStatisticsHandlerReturnsBadRequestForInvalidUUID(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/not-a-uuid/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStatisticsHandlerReturnsInternalErrorForReaderFailure(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{view: queries.ResumeLinkPublicView{ResumeID: resumevalueobjects.NewID().String()}},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{err: errors.New("db down")},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+valueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}
