package profile

import (
	"bytes"
	"context"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/application/commands"
	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/platform/http/middleware"
	"developers-cv/backend/internal/shared/domain/valueobjects"
)

type stubReplaceProfileHandler struct {
	err     error
	called  bool
	command commands.ReplaceProfileCommand
}

func (s *stubReplaceProfileHandler) Handle(_ context.Context, command commands.ReplaceProfileCommand) error {
	s.called = true
	s.command = command
	return s.err
}

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
		ID: userID.String(), RU: queries.LocalizedProfile{Name: stringPtr("Антон"), Surname: stringPtr("Седлецкий"), Patronymic: stringPtr("Александрович"), Position: stringPtr("Frontend"), Location: stringPtr("Новосибирск")}, PhotoPath: stringPtr("/mock/avatar.svg"),
	}}
	engine := newTestEngine(userID, nil, NewHandler(queries.NewGetProfileHandler(reader), nil))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/profile", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"userId":"`+userID.String()+`"`)
	require.Contains(t, w.Body.String(), `"ru":{"name":"Антон"`)
	require.Contains(t, w.Body.String(), `"surname":"Седлецкий"`)
	require.Contains(t, w.Body.String(), `"position":"Frontend"`)
	require.Equal(t, identity, reader.userID)
}

func TestGetRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := newTestEngine(uuid.Nil, middleware.ErrUnauthenticated, NewHandler(queries.NewGetProfileHandler(&stubGetProfileReader{}), nil))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/profile", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Contains(t, w.Body.String(), "UNAUTHORIZED")
}

func TestGetReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := newTestEngine(uuid.New(), nil, NewHandler(queries.NewGetProfileHandler(&stubGetProfileReader{err: queries.ErrProfileNotFound}), nil))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/profile", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "Profile not found")
}

func TestGetIgnoresUserIDQueryParameter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	reader := &stubGetProfileReader{profile: queries.Profile{ID: userID.String(), RU: queries.LocalizedProfile{Name: stringPtr("Anton"), Surname: stringPtr("User")}}}
	engine := newTestEngine(userID, nil, NewHandler(queries.NewGetProfileHandler(reader), nil))

	foreignID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/profile?userId="+foreignID.String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, userID.String(), reader.userID.String())
}

func stringPtr(value string) *string { return &value }

func TestReplaceReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	replace := &stubReplaceProfileHandler{}
	engine := newTestEngineWithProfile(userID, nil, NewHandler(queries.NewGetProfileHandler(&stubGetProfileReader{}), replace), replace)
	body := []byte(`{"ru":{"name":"Дмитрий","surname":"Костерин","patronymic":"Александрович","position":"Backend Developer","location":"Томск"},"en":{"name":"Dmitry","surname":"Kosterin","patronymic":"Alexandrovich","position":"Backend Developer","location":"Tomsk"},"photoPath":"https://example.com/profiles/dmitry.svg"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
	require.Empty(t, w.Body.String())
	require.True(t, replace.called)
	require.Equal(t, "Дмитрий", *replace.command.RU.Name)
	require.Equal(t, "Костерин", *replace.command.RU.Surname)
	require.Equal(t, "Александрович", *replace.command.RU.Patronymic)
	require.Equal(t, "Backend Developer", *replace.command.RU.Position)
	require.Equal(t, "Томск", *replace.command.RU.Location)
	require.Equal(t, "Dmitry", *replace.command.EN.Name)
	require.Equal(t, "Kosterin", *replace.command.EN.Surname)
	require.Equal(t, "Alexandrovich", *replace.command.EN.Patronymic)
	require.Equal(t, "Backend Developer", *replace.command.EN.Position)
	require.Equal(t, "Tomsk", *replace.command.EN.Location)
	require.Equal(t, "https://example.com/profiles/dmitry.svg", *replace.command.PhotoPath)
}

func newTestEngineWithProfile(userID uuid.UUID, authErr error, handler *Handler, replace *stubReplaceProfileHandler) *gin.Engine {
	engine := newTestEngine(userID, authErr, handler)
	engine.PUT("/api/v1/lk/profile", handler.Replace)
	return engine
}

func TestReplaceRejectsEmptyNameBeforeCommand(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	replace := &stubReplaceProfileHandler{}
	engine := newTestEngineWithProfile(userID, nil, NewHandler(queries.NewGetProfileHandler(&stubGetProfileReader{}), replace), replace)
	body := []byte(`{"ru":{"name":"   ","surname":"Костерин"}}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "BAD_REQUEST")
	require.False(t, replace.called)
}

func TestReplaceRejectsEmptySurnameBeforeCommand(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	replace := &stubReplaceProfileHandler{}
	engine := newTestEngineWithProfile(userID, nil, NewHandler(queries.NewGetProfileHandler(&stubGetProfileReader{}), replace), replace)
	body := []byte(`{"en":{"name":"Dmitry","surname":"\t"}}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "BAD_REQUEST")
	require.False(t, replace.called)
}

func TestReplaceRejectsInvalidPhotoPathBeforeCommand(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	replace := &stubReplaceProfileHandler{}
	engine := newTestEngineWithProfile(userID, nil, NewHandler(queries.NewGetProfileHandler(&stubGetProfileReader{}), replace), replace)
	body := []byte(`{"ru":{"name":"Дмитрий","surname":"Костерин"},"photoPath":"/profiles/dmitry.svg"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "BAD_REQUEST")
	require.False(t, replace.called)
}
