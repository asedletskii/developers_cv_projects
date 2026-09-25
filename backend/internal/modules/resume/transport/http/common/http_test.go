package common

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
)

func TestWriteCommandErrorContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		err  error
		code string
		want int
	}{
		{name: "not found", err: domain.ErrEntityNotFound, code: string(apperrors.CodeNotFound), want: http.StatusNotFound},
		{name: "duplicate", err: domain.ErrDuplicateEntity, code: string(apperrors.CodeConflict), want: http.StatusConflict},
		{name: "validation", err: domain.ErrInvalidEntity, code: string(apperrors.CodeBadRequest), want: http.StatusBadRequest},
		{name: "internal", err: errors.New("database failure"), code: string(apperrors.CodeInternal), want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set(middleware.RequestIDKey, "request-123")

			WriteCommandError(c, tt.err)

			require.Equal(t, tt.want, w.Code)
			require.Contains(t, w.Body.String(), tt.code)
			require.Contains(t, w.Body.String(), "request-123")
		})
	}
}
