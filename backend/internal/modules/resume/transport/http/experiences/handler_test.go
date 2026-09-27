package experiences

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/platform/http/middleware"
)

type createExperienceRepositoryStub struct {
	experience *domain.Experience
}

func (s *createExperienceRepositoryStub) CreateExperience(_ context.Context, command ports.CreateExperienceRepositoryCommand) error {
	s.experience = command.Experience
	return nil
}

type updateExperienceRepositoryStub struct {
	experience *domain.Experience
}

func (s *updateExperienceRepositoryStub) UpdateExperience(_ context.Context, command ports.UpdateExperienceRepositoryCommand) error {
	s.experience = command.Experience
	return nil
}

type unitOfWorkStub struct{}

func (unitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type listExperienceReaderStub struct {
	items []queries.ExperienceView
}

func (s listExperienceReaderStub) ListExperiences(context.Context, queries.ListExperiencesQuery) ([]queries.ExperienceView, error) {
	return s.items, nil
}

func (s listExperienceReaderStub) GetExperience(context.Context, queries.GetExperienceQuery) (queries.ExperienceView, error) {
	if len(s.items) == 0 {
		return queries.ExperienceView{}, nil
	}
	return s.items[0], nil
}

func newExperienceTestEngine(userID, resumeID uuid.UUID, handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	engine.POST("/api/v1/lk/resume/:resumeId/experiences", handler.Create)
	engine.GET("/api/v1/lk/resume/:resumeId/experiences", handler.List)
	engine.GET("/api/v1/lk/resume/:resumeId/experiences/:experienceId", handler.Get)
	engine.PUT("/api/v1/lk/resume/:resumeId/experiences/:experienceId", handler.Update)
	return engine
}

func TestCreateAcceptsStack(t *testing.T) {
	userID := uuid.New()
	resumeID := uuid.New()
	repository := &createExperienceRepositoryStub{}
	createHandler := commands.NewCreateExperienceHandler(repository, unitOfWorkStub{})
	handler := NewHandler(createHandler, nil, nil, nil, nil)
	engine := newExperienceTestEngine(userID, resumeID, handler)

	body := `{
		"company": "Acme",
		"position": "Go Developer",
		"startDate": "2024-01",
		"achievements": [],
		"stack": [
			{"name": "Go", "grade": "senior"}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/resume/"+resumeID.String()+"/experiences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.NotNil(t, repository.experience)
	require.Len(t, repository.experience.Stack(), 1)
	require.Equal(t, "Go", repository.experience.Stack()[0].Name())
}

func TestCreateRejectsLegacyTechnologiesField(t *testing.T) {
	userID := uuid.New()
	resumeID := uuid.New()
	repository := &createExperienceRepositoryStub{}
	createHandler := commands.NewCreateExperienceHandler(repository, unitOfWorkStub{})
	handler := NewHandler(createHandler, nil, nil, nil, nil)
	engine := newExperienceTestEngine(userID, resumeID, handler)

	body := `{
		"company": "Acme",
		"position": "Go Developer",
		"startDate": "2024-01",
		"achievements": [],
		"technologies": [{"name": "Go"}]
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/lk/resume/"+resumeID.String()+"/experiences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Nil(t, repository.experience)
}

func TestUpdateAcceptsStack(t *testing.T) {
	userID := uuid.New()
	resumeID := uuid.New()
	experienceID := uuid.New()
	repository := &updateExperienceRepositoryStub{}
	updateHandler := commands.NewUpdateExperienceHandler(repository, unitOfWorkStub{})
	handler := NewHandler(nil, nil, nil, updateHandler, nil)
	engine := newExperienceTestEngine(userID, resumeID, handler)

	body := `{
		"company": "Acme",
		"position": "Staff Go Developer",
		"startDate": "2024-01",
		"achievements": [],
		"stack": [{"name": "Go", "grade": "senior"}]
	}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/lk/resume/"+resumeID.String()+"/experiences/"+experienceID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.NotNil(t, repository.experience)
	require.Len(t, repository.experience.Stack(), 1)
	require.Equal(t, "Go", repository.experience.Stack()[0].Name())
}

func TestListReturnsStack(t *testing.T) {
	userID := uuid.New()
	resumeID := uuid.New()
	reader := listExperienceReaderStub{
		items: []queries.ExperienceView{
			{
				ID:           uuid.NewString(),
				ResumeID:     resumeID.String(),
				Company:      "Acme",
				Position:     "Go Developer",
				StartDate:    "2024-01",
				Achievements: []string{},
				Stack: []queries.ExperienceStackItemView{
					{Name: "Go", Grade: stringPtr("senior")},
				},
			},
		},
	}
	listHandler := queries.NewListExperiencesQueryHandler(reader)
	handler := NewHandler(nil, listHandler, nil, nil, nil)
	engine := newExperienceTestEngine(userID, resumeID, handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/resume/"+resumeID.String()+"/experiences", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `[
		{
			"id": "`+reader.items[0].ID+`",
			"resumeId": "`+resumeID.String()+`",
			"company": "Acme",
			"position": "Go Developer",
			"startDate": "2024-01",
			"achievements": [],
			"stack": [{"name": "Go", "grade": "senior"}]
		}
	]`, w.Body.String())
	require.NotContains(t, w.Body.String(), `"technologies"`)
}

func TestGetReturnsStack(t *testing.T) {
	userID := uuid.New()
	resumeID := uuid.New()
	experienceID := uuid.New()
	reader := listExperienceReaderStub{
		items: []queries.ExperienceView{
			{
				ID:           experienceID.String(),
				ResumeID:     resumeID.String(),
				Company:      "Acme",
				Position:     "Go Developer",
				StartDate:    "2024-01",
				Achievements: []string{},
				Stack: []queries.ExperienceStackItemView{
					{Name: "Go", Grade: stringPtr("senior")},
				},
			},
		},
	}
	getHandler := queries.NewGetExperienceQueryHandler(reader)
	handler := NewHandler(nil, nil, getHandler, nil, nil)
	engine := newExperienceTestEngine(userID, resumeID, handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/resume/"+resumeID.String()+"/experiences/"+experienceID.String(), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{
		"id": "`+experienceID.String()+`",
		"resumeId": "`+resumeID.String()+`",
		"company": "Acme",
		"position": "Go Developer",
		"startDate": "2024-01",
		"achievements": [],
		"stack": [{"name": "Go", "grade": "senior"}]
	}`, w.Body.String())
	require.NotContains(t, w.Body.String(), `"technologies"`)
}

func stringPtr(value string) *string {
	return &value
}

var _ ports.CreateExperienceRepository = (*createExperienceRepositoryStub)(nil)
var _ ports.UpdateExperienceRepository = (*updateExperienceRepositoryStub)(nil)
var _ commands.UnitOfWork = unitOfWorkStub{}
var _ queries.ExperiencesReader = listExperienceReaderStub{}
