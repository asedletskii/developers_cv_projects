package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ListProjectsQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type GetProjectQuery struct {
	ResumeID  valueobjects.ID
	ProjectID valueobjects.ID
	UserID    useridentity.ID
}

type ProjectsReader interface {
	ListProjects(context.Context, ListProjectsQuery) ([]ProjectView, error)
	GetProject(context.Context, GetProjectQuery) (ProjectView, error)
}

type ListProjectsQueryHandler struct{ reader ProjectsReader }

func NewListProjectsQueryHandler(reader ProjectsReader) *ListProjectsQueryHandler {
	return &ListProjectsQueryHandler{reader: reader}
}

func (h *ListProjectsQueryHandler) Handle(ctx context.Context, q ListProjectsQuery) ([]ProjectView, error) {
	if h == nil || h.reader == nil {
		return nil, fmt.Errorf("list projects query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("invalid list projects query")
	}
	return h.reader.ListProjects(ctx, q)
}

type GetProjectQueryHandler struct{ reader ProjectsReader }

func NewGetProjectQueryHandler(reader ProjectsReader) *GetProjectQueryHandler {
	return &GetProjectQueryHandler{reader: reader}
}

func (h *GetProjectQueryHandler) Handle(ctx context.Context, q GetProjectQuery) (ProjectView, error) {
	if h == nil || h.reader == nil {
		return ProjectView{}, fmt.Errorf("get project query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.ProjectID.IsZero() || q.UserID.IsZero() {
		return ProjectView{}, fmt.Errorf("invalid get project query")
	}
	return h.reader.GetProject(ctx, q)
}

type ProjectView struct {
	ID           string                   `json:"id"`
	ResumeID     string                   `json:"resumeId"`
	Name         string                   `json:"name"`
	Description  *string                  `json:"description,omitempty"`
	Role         *string                  `json:"role,omitempty"`
	Screenshots  []string                 `json:"screenshots"`
	Link         *string                  `json:"link,omitempty"`
	GitLink      *string                  `json:"gitLink,omitempty"`
	Stack        []ProjectTechnologyView  `json:"stack"`
	Contributors []ProjectContributorView `json:"contributors"`
}

type ProjectTechnologyView struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}

type ProjectContributorView struct {
	Name     string              `json:"name"`
	Surname  string              `json:"surname"`
	Role     string              `json:"role"`
	Contacts ProjectContactsView `json:"contacts"`
}

type ProjectContactsView struct {
	Phone    *string `json:"phone,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}
