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

	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/platform/http/middleware"
)

type stubGetResumeReader struct {
	resume queries.Resume
	err    error
	query  queries.GetResumeQuery
}

func (r *stubGetResumeReader) GetResume(_ context.Context, query queries.GetResumeQuery) (queries.Resume, error) {
	r.query = query
	return r.resume, r.err
}

func TestGetResume(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.New()
	userID := uuid.New()
	reader := &stubGetResumeReader{
		resume: queries.Resume{ID: id.String(), Name: "Anton Sedletskii"},
	}
	handler := NewHandler(queries.NewGetResumeHandler(reader), nil, nil, nil, nil)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.GET("/api/v1/lk/resume/:resumeId", handler.GetResume)

	req := httptest.NewRequest("GET", "/api/v1/lk/resume/"+id.String()+"?lang=en", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "Anton Sedletskii")
	require.NotNil(t, reader.query.UserID)
	require.Equal(t, userID, reader.query.UserID.UUID())
}

func TestGetPublicResumeDoesNotRequireAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.New()
	reader := &stubGetResumeReader{
		resume: queries.Resume{ID: id.String(), Name: "Anton Sedletskii"},
	}
	handler := NewHandler(queries.NewGetResumeHandler(reader), nil, nil, nil, nil)
	engine := gin.New()
	engine.GET("/api/v1/cv/:resumeId", handler.GetPublicResume)

	req := httptest.NewRequest("GET", "/api/v1/cv/"+id.String()+"?lang=en", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "Anton Sedletskii")
	require.Nil(t, reader.query.UserID)
}

func TestGetResumeNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(&stubGetResumeReader{err: queries.ErrResumeNotFound}), nil, nil, nil, nil)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.GET("/api/v1/lk/resume/:resumeId", handler.GetResume)

	req := httptest.NewRequest("GET", "/api/v1/lk/resume/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "Resume not found")
}

func TestGetPublicResumeNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(&stubGetResumeReader{err: queries.ErrResumeNotFound}), nil, nil, nil, nil)
	engine := gin.New()
	engine.GET("/api/v1/cv/:resumeId", handler.GetPublicResume)

	req := httptest.NewRequest("GET", "/api/v1/cv/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "Resume not found")
}

func TestGetResumeRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(&stubGetResumeReader{}), nil, nil, nil, nil)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.GET("/api/v1/lk/resume/:resumeId", handler.GetResume)

	req := httptest.NewRequest("GET", "/api/v1/lk/resume/not-a-uuid", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid resume id")
}

func TestGetPublicResumeRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(&stubGetResumeReader{}), nil, nil, nil, nil)
	engine := gin.New()
	engine.GET("/api/v1/cv/:resumeId", handler.GetPublicResume)

	req := httptest.NewRequest("GET", "/api/v1/cv/not-a-uuid", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid resume id")
}

func TestGetResumeRejectsInvalidLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(&stubGetResumeReader{}), nil, nil, nil, nil)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.GET("/api/v1/lk/resume/:resumeId", handler.GetResume)

	req := httptest.NewRequest("GET", "/api/v1/lk/resume/"+uuid.NewString()+"?lang=de", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid language")
}

func TestGetPublicResumeRejectsInvalidLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(queries.NewGetResumeHandler(&stubGetResumeReader{}), nil, nil, nil, nil)
	engine := gin.New()
	engine.GET("/api/v1/cv/:resumeId", handler.GetPublicResume)

	req := httptest.NewRequest("GET", "/api/v1/cv/"+uuid.NewString()+"?lang=de", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid language")
}

var _ queries.GetResumeReader = &stubGetResumeReader{}

type fakeCreateResumeRepository struct{ saved *domain.Resume }

func (f *fakeCreateResumeRepository) GetByID(
	context.Context,
	useridentity.ID,
	valueobjects.ID,
) (*domain.Resume, error) {
	return nil, nil
}
func (f *fakeCreateResumeRepository) Save(_ context.Context, resume *domain.Resume) error {
	f.saved = resume
	return nil
}
func (f *fakeCreateResumeRepository) UpdateTitle(context.Context, useridentity.ID, valueobjects.ID, string) error {
	return nil
}

type fakeCreateResumeUoW struct{ calls int }

func (f *fakeCreateResumeUoW) Do(ctx context.Context, fn func(context.Context) error) error {
	f.calls++
	return fn(ctx)
}

func TestCreateResume(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &fakeCreateResumeRepository{}
	uow := &fakeCreateResumeUoW{}
	createHandler := commands.NewCreateResumeHandler(repo, uow)
	getHandler := queries.NewGetResumeHandler(&stubGetResumeReader{resume: queries.Resume{ID: "created", Name: "Test"}})
	handler := NewHandler(getHandler, nil, createHandler, nil, nil)
	engine := gin.New()
	userID := uuid.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.POST("/api/v1/lk/resume", handler.CreateResume)

	body := `{
		"userId":"00000000-0000-0000-0000-000000000099",
		"title":"Client supplied title",
		"name":"Test",
		"surname":"User",
		"position":"Developer",
		"greeting":{
			"first":{"ru":"Привет","en":"Hello"},
			"second":{"ru":"Я","en":"I am"},
			"secondHighlight":{"ru":"Тест","en":"Test"},
			"third":{"ru":"Разработчик","en":"Developer"}
		},
		"photoPath":"/assets/photo.svg",
		"dateOfBirth":"1998-04-12",
		"location":"Remote",
		"contacts":{"phone":"+79999999999"}
	}`
	req := httptest.NewRequest("POST", "/api/v1/lk/resume", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"id":"created"`)
	require.Equal(t, 1, uow.calls)
	require.NotNil(t, repo.saved)
	require.False(t, repo.saved.ID().IsZero())
	require.Equal(t, userID, repo.saved.UserID().UUID())
	require.Equal(t, "Developer", repo.saved.Title())
}

func TestCreateResumeRejectsInvalidDomainInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &fakeCreateResumeRepository{}
	uow := &fakeCreateResumeUoW{}
	handler := NewHandler(
		queries.NewGetResumeHandler(&stubGetResumeReader{}),
		nil,
		commands.NewCreateResumeHandler(repo, uow),
		nil,
		nil,
	)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.New(), nil
	}))
	engine.POST("/api/v1/lk/resume", handler.CreateResume)

	body := `{"name":"","surname":"User","position":"Developer","greeting":{"first":{"ru":"Привет","en":"Hello"},"second":{"ru":"Я","en":"I am"},"secondHighlight":{"ru":"Тест","en":"Test"},"third":{"ru":"Разработчик","en":"Developer"}},"photoPath":"/assets/photo.svg","dateOfBirth":"1998-04-12","location":"Remote","contacts":{"phone":"+79999999999"}}`
	req := httptest.NewRequest("POST", "/api/v1/lk/resume", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Nil(t, repo.saved)
}

type stubListReader struct {
	items  []queries.ResumeListItem
	userID useridentity.ID
}

func (r *stubListReader) ListResumes(_ context.Context, query queries.GetListResumesQuery) ([]queries.ResumeListItem, error) {
	r.userID = query.UserID
	return r.items, nil
}

func TestListResumesRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	listReader := &stubListReader{}
	handler := NewHandler(queries.NewGetResumeHandler(&stubGetResumeReader{}), queries.NewListResumesHandler(listReader), nil, nil, nil)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return uuid.Nil, middleware.ErrUnauthenticated
	}))
	engine.GET("/api/v1/lk/resumes", handler.ListResumes)

	req := httptest.NewRequest("GET", "/api/v1/lk/resumes", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestListResumesReturnsOnlyCurrentUsersData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	listReader := &stubListReader{items: []queries.ResumeListItem{{ID: "resume-a", Title: "Senior Go Developer"}}}
	handler := NewHandler(queries.NewGetResumeHandler(&stubGetResumeReader{}), queries.NewListResumesHandler(listReader), nil, nil, nil)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.GET("/api/v1/lk/resumes", handler.ListResumes)

	req := httptest.NewRequest("GET", "/api/v1/lk/resumes", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"resume-a"`)
	require.Contains(t, w.Body.String(), `"Senior Go Developer"`)
	require.Equal(t, userID, listReader.userID.UUID())
}

type updateResumeRepositoryStub struct {
	resume *domain.Resume
	saved  *domain.Resume
	err    error
}

func (r *updateResumeRepositoryStub) GetByID(_ context.Context, _ useridentity.ID, _ valueobjects.ID) (*domain.Resume, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.resume, nil
}

func (r *updateResumeRepositoryStub) Save(_ context.Context, resume *domain.Resume) error {
	r.saved = resume
	return nil
}

func (r *updateResumeRepositoryStub) UpdateTitle(_ context.Context, _ useridentity.ID, _ valueobjects.ID, title string) error {
	if r.err != nil {
		return r.err
	}
	if r.resume != nil {
		if err := r.resume.Rename(title); err != nil {
			return err
		}
	}
	return nil
}

type updateResumeUoWStub struct{}

func (updateResumeUoWStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestPatchResumeHTTPChangesTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	resumeID, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	owner, err := useridentity.IDFromUUID(userID)
	require.NoError(t, err)
	resume := newHTTPTestResume(t, resumeID, owner, "Old title")
	repo := &updateResumeRepositoryStub{resume: resume}
	handler := NewHandler(
		queries.NewGetResumeHandler(&stubGetResumeReader{}),
		nil,
		nil,
		nil,
		commands.NewPatchResumeHandler(repo, updateResumeUoWStub{}),
	)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	engine.PATCH("/api/v1/lk/resume/:resumeId", handler.PatchResume)

	req := httptest.NewRequest("PATCH", "/api/v1/lk/resume/"+userIDFromID(resumeID).String(), strings.NewReader(`{"title":"New title"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"id":"`+resumeID.String()+`"`)
	require.Contains(t, w.Body.String(), `"title":"New title"`)
	require.NotNil(t, repo.resume)
	require.Equal(t, "New title", repo.resume.Title())
	require.Equal(t, "Senior Go Developer", repo.resume.Position())
	require.Nil(t, repo.saved)
}

func TestPatchResumeHTTPRejectsForeignResume(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	resumeID, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	repo := &updateResumeRepositoryStub{resume: nil, err: domain.ErrEntityNotFound}
	handler := NewHandler(
		queries.NewGetResumeHandler(&stubGetResumeReader{}), nil, nil, nil,
		commands.NewPatchResumeHandler(repo, updateResumeUoWStub{}),
	)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	engine.PATCH("/api/v1/lk/resume/:resumeId", handler.PatchResume)

	req := httptest.NewRequest("PATCH", "/api/v1/lk/resume/"+userIDFromID(resumeID).String(), strings.NewReader(`{"title":"New title"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Nil(t, repo.saved)
}

func userIDFromID(id valueobjects.ID) uuid.UUID { return id.UUID() }

func newHTTPTestResume(t *testing.T, id valueobjects.ID, userID useridentity.ID, title string) *domain.Resume {
	t.Helper()
	date, err := valueobjects.ParseDate("1998-04-12")
	require.NoError(t, err)
	photo, err := valueobjects.NewAssetPath("/assets/photo.svg")
	require.NoError(t, err)
	resume, err := domain.NewResume(domain.ResumeParams{
		ID: id, UserID: userID, Name: "User", Surname: "Test", Position: "Senior Go Developer", Title: title,
		PhotoPath: photo, DateOfBirth: date, Location: "Remote",
	})
	require.NoError(t, err)
	return resume
}

func TestUpdateResumeHTTPPreservesTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	resumeID, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	owner, err := useridentity.IDFromUUID(userID)
	require.NoError(t, err)
	resume := newHTTPTestResume(t, resumeID, owner, "Manual title")
	repo := &updateResumeRepositoryStub{resume: resume}
	handler := NewHandler(
		queries.NewGetResumeHandler(&stubGetResumeReader{resume: queries.Resume{ID: resumeID.String(), Name: "Updated"}}), nil, nil,
		commands.NewUpdateResumeHandler(repo, updateResumeUoWStub{}), nil,
	)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.UpdateResume)

	body := `{
		"name":"Updated",
		"surname":"User",
		"position":"Staff Go Developer",
		"greeting":{"first":{"ru":"Привет","en":"Hello"},"second":{"ru":"Я","en":"I am"},"secondHighlight":{"ru":"Тест","en":"Test"},"third":{"ru":"Разработчик","en":"Developer"}},
		"photoPath":"/assets/photo.svg",
		"dateOfBirth":"1998-04-12",
		"location":"Remote",
		"contacts":{"phone":"+79999999999"}
	}`
	req := httptest.NewRequest("PUT", "/api/v1/lk/resume/"+resumeID.UUID().String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"id":"`+resumeID.String()+`"`)
	require.Equal(t, "Updated", repo.saved.Name())
	require.Equal(t, "Staff Go Developer", repo.saved.Position())
	require.Equal(t, "Manual title", repo.saved.Title())
}

func TestUpdateResumeHTTPRejectsMalformedID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(
		queries.NewGetResumeHandler(&stubGetResumeReader{}), nil, nil, nil, nil,
	)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return uuid.New(), nil }))
	engine.PUT("/api/v1/lk/resume/:resumeId", handler.UpdateResume)

	req := httptest.NewRequest("PUT", "/api/v1/lk/resume/not-a-uuid", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
