package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
)

type ListExperiencesQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetExperienceQuery struct {
	ResumeID, ExperienceID valueobjects.ID
	UserID                 useridentity.ID
}
type ExperienceTechnologyView struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}
type ExperienceView struct {
	ID           string                     `json:"id"`
	ResumeID     string                     `json:"resumeId"`
	Company      string                     `json:"company"`
	Position     string                     `json:"position"`
	Location     *string                    `json:"location,omitempty"`
	StartDate    string                     `json:"startDate"`
	EndDate      *string                    `json:"endDate,omitempty"`
	Description  *string                    `json:"description,omitempty"`
	Achievements []string                   `json:"achievements"`
	Technologies []ExperienceTechnologyView `json:"technologies"`
}
type ExperiencesReader interface {
	ListExperiences(context.Context, ListExperiencesQuery) ([]ExperienceView, error)
	GetExperience(context.Context, GetExperienceQuery) (ExperienceView, error)
}

type ExperiencesSource interface {
	ReadExperiences(context.Context, valueobjects.ID) ([]ExperienceView, error)
}
type ListExperiencesQueryHandler struct{ reader ExperiencesReader }

func NewListExperiencesQueryHandler(reader ExperiencesReader) *ListExperiencesQueryHandler {
	return &ListExperiencesQueryHandler{reader: reader}
}
func (h *ListExperiencesQueryHandler) Handle(ctx context.Context, q ListExperiencesQuery) ([]ExperienceView, error) {
	if h == nil || h.reader == nil {
		return nil, fmt.Errorf("list experiences query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("invalid list experiences query")
	}
	return h.reader.ListExperiences(ctx, q)
}

type GetExperienceQueryHandler struct{ reader ExperiencesReader }

func NewGetExperienceQueryHandler(reader ExperiencesReader) *GetExperienceQueryHandler {
	return &GetExperienceQueryHandler{reader: reader}
}
func (h *GetExperienceQueryHandler) Handle(ctx context.Context, q GetExperienceQuery) (ExperienceView, error) {
	if h == nil || h.reader == nil {
		return ExperienceView{}, fmt.Errorf("get experience query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.ExperienceID.IsZero() || q.UserID.IsZero() {
		return ExperienceView{}, fmt.Errorf("invalid get experience query")
	}
	return h.reader.GetExperience(ctx, q)
}
