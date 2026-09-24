package public

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type readerStub struct {
	view  queries.PublicResumeView
	err   error
	query queries.PublicResumeQuery
}

func (s *readerStub) GetPublicResume(_ context.Context, q queries.PublicResumeQuery) (queries.PublicResumeView, error) {
	s.query = q
	return s.view, s.err
}

func TestHandlerDoesNotRequireLanguageQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	reader := &readerStub{view: queries.PublicResumeView{ID: "resume"}}
	query := queries.NewPublicResumeQueryHandler(reader)
	handler := NewHandler(query)

	engine := gin.New()
	engine.GET("/api/v1/cv/:resumeId", handler.Handle)

	id := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+id.String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, id, reader.query.ID)
	require.NotContains(t, req.URL.RawQuery, "lang")
}
