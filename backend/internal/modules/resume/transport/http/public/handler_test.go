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
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type readerStub struct {
	view        queries.ResumeView
	err         error
	publicID    valueobjects.ID
	publicCalls int
}

func (s *readerStub) GetPublic(_ context.Context, resumeID valueobjects.ID) (queries.ResumeView, error) {
	s.publicID = resumeID
	s.publicCalls++
	if s.err != nil {
		return queries.ResumeView{}, s.err
	}
	return s.view, nil
}

func TestHandlerUsesPublicReadMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	reader := &readerStub{view: queries.ResumeView{ID: "resume"}}
	query := queries.NewPublicResumeQueryHandler(reader)
	handler := NewHandler(query)

	engine := gin.New()
	engine.GET("/api/v1/cv/:resumeId", handler.Handle)

	id := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+id.String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, id, reader.publicID.UUID())
	require.Equal(t, 1, reader.publicCalls)
}
