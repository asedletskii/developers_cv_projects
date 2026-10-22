package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationStoresUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()

	engine := gin.New()
	engine.Use(RequestID(), Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.GET("/test", func(c *gin.Context) {
		actual, ok := GetUserID(c)
		require.True(t, ok)
		require.Equal(t, userID, actual)
		c.Status(204)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()
	engine.ServeHTTP(res, req)

	require.Equal(t, 204, res.Code)
}

func TestAuthenticationReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.Use(RequestID(), Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, ErrUnauthenticated
	}))
	engine.GET("/test", func(c *gin.Context) {
		c.Status(204)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()
	engine.ServeHTTP(res, req)

	require.Equal(t, 401, res.Code)
	require.Contains(t, res.Body.String(), "UNAUTHORIZED")
}

func TestAuthenticationReturnsInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.Use(RequestID(), Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, errors.New("database unavailable")
	}))
	engine.GET("/test", func(c *gin.Context) {
		c.Status(204)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()
	engine.ServeHTTP(res, req)

	require.Equal(t, 500, res.Code)
	require.Contains(t, res.Body.String(), "INTERNAL_ERROR")
}

func TestOptionalAuthenticationAllowsUnauthenticatedRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(RequestID(), OptionalAuthentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, ErrUnauthenticated
	}))
	engine.GET("/test", func(c *gin.Context) {
		_, ok := GetUserID(c)
		require.False(t, ok)
		c.Status(http.StatusNoContent)
	})

	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/test", nil))

	require.Equal(t, http.StatusNoContent, response.Code)
}

func TestOptionalAuthenticationReturnsInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(RequestID(), OptionalAuthentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, errors.New("session store unavailable")
	}))
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/test", nil))

	require.Equal(t, http.StatusInternalServerError, response.Code)
	require.Contains(t, response.Body.String(), "INTERNAL_ERROR")
	require.NotContains(t, response.Body.String(), "session store unavailable")
}

func TestOptionalAuthenticationRejectsEmptyUserIDWithoutError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(RequestID(), OptionalAuthentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, nil
	}))
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/test", nil))

	require.Equal(t, http.StatusInternalServerError, response.Code)
	require.Contains(t, response.Body.String(), "INTERNAL_ERROR")
}
