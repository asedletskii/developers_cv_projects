package list

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/platform/http/middleware"
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

	handler := NewHandler(queries.NewListResumesHandler(&stubListReader{}))
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, middleware.ErrUnauthenticated
	}))
	engine.GET("/api/v1/lk/resumes", handler.Handle)

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
	handler := NewHandler(queries.NewListResumesHandler(reader))

	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.GET("/api/v1/lk/resumes", handler.Handle)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/resumes", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"resume-a"`)
	require.Contains(t, w.Body.String(), `"Senior Go Developer"`)
	require.Equal(t, userID, reader.userID.UUID())
}
