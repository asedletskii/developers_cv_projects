package profile

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/platform/http/middleware"
)

type stubGetProfileReader struct {
	profile queries.Profile
	err     error
	userID  valueobjects.ID
}

func (r *stubGetProfileReader) GetProfile(_ context.Context, query queries.GetProfileQuery) (queries.Profile, error) {
	r.userID = query.UserID
	return r.profile, r.err
}

func newTestEngine(userID uuid.UUID, authErr error, handler *Handler) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		if authErr != nil {
			return uuid.Nil, authErr
		}
		return userID, nil
	}))
	engine.GET("/api/v1/lk/profile", handler.Get)
	return engine
}

func TestGetReturnsAuthenticatedUserProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	identity, err := valueobjects.IDFromUUID(userID)
	require.NoError(t, err)

	reader := &stubGetProfileReader{profile: queries.Profile{
		ID:         userID.String(),
		Name:       "Антон",
		Surname:    "Седлецкий",
		Patronymic: stringPtr("Александрович"),
		Position:   stringPtr("Frontend"),
		Location:   stringPtr("Новосибирск"),
		PhotoPath:  stringPtr("/mock/avatar.svg"),
	}}
	engine := newTestEngine(userID, nil, NewHandler(queries.NewGetProfileHandler(reader)))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/profile", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"userId":"`+userID.String()+`"`)
	require.Contains(t, w.Body.String(), `"name":"Антон"`)
	require.Contains(t, w.Body.String(), `"surname":"Седлецкий"`)
	require.Contains(t, w.Body.String(), `"position":"Frontend"`)
	require.Equal(t, identity, reader.userID)
}

func TestGetRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := newTestEngine(uuid.Nil, middleware.ErrUnauthenticated, NewHandler(queries.NewGetProfileHandler(&stubGetProfileReader{})))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/profile", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Contains(t, w.Body.String(), "UNAUTHORIZED")
}

func TestGetReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := newTestEngine(uuid.New(), nil, NewHandler(queries.NewGetProfileHandler(&stubGetProfileReader{err: queries.ErrProfileNotFound})))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/profile", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "Profile not found")
}

func TestGetIgnoresUserIDQueryParameter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	reader := &stubGetProfileReader{profile: queries.Profile{ID: userID.String(), Name: "Anton", Surname: "User"}}
	engine := newTestEngine(userID, nil, NewHandler(queries.NewGetProfileHandler(reader)))

	foreignID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/profile?userId="+foreignID.String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, userID.String(), reader.userID.String())
}

func stringPtr(value string) *string { return &value }
