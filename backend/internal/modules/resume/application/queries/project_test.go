package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type projectsReaderStub struct {
	list func(context.Context, ListProjectsQuery) ([]ProjectView, error)
	get  func(context.Context, GetProjectQuery) (ProjectView, error)
}

func (s projectsReaderStub) ListProjects(ctx context.Context, q ListProjectsQuery) ([]ProjectView, error) {
	return s.list(ctx, q)
}

func (s projectsReaderStub) GetProject(ctx context.Context, q GetProjectQuery) (ProjectView, error) {
	return s.get(ctx, q)
}

func TestProjectQueryHandlersDelegateToReader(t *testing.T) {
	resumeID := valueobjects.NewID()
	projectID := valueobjects.NewID()
	userID := useridentity.NewID()
	want := ProjectView{ID: projectID.String(), ResumeID: resumeID.String(), Name: "Project"}

	reader := projectsReaderStub{
		list: func(context.Context, ListProjectsQuery) ([]ProjectView, error) { return []ProjectView{want}, nil },
		get:  func(context.Context, GetProjectQuery) (ProjectView, error) { return want, nil },
	}

	items, err := NewListProjectsQueryHandler(reader).Handle(context.Background(), ListProjectsQuery{ResumeID: resumeID, UserID: userID})
	require.NoError(t, err)
	require.Equal(t, []ProjectView{want}, items)

	item, err := NewGetProjectQueryHandler(reader).Handle(context.Background(), GetProjectQuery{ResumeID: resumeID, ProjectID: projectID, UserID: userID})
	require.NoError(t, err)
	require.Equal(t, want, item)
}

func TestProjectQueryHandlersRejectInvalidQueries(t *testing.T) {
	reader := projectsReaderStub{}

	_, err := NewListProjectsQueryHandler(reader).Handle(context.Background(), ListProjectsQuery{})
	require.Error(t, err)

	_, err = NewGetProjectQueryHandler(reader).Handle(context.Background(), GetProjectQuery{})
	require.Error(t, err)
}
