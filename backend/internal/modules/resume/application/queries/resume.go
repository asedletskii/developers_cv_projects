package queries

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GetResumeQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type ResumeReader interface {
	GetOwned(context.Context, GetResumeQuery) (ResumeView, error)
}

type ResumeView struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	Language     string              `json:"language"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
	Version      *int64              `json:"version,omitempty"`
	Overview     ResumeOverviewView  `json:"overview"`
	Profile      *ResumeProfileView  `json:"profile,omitempty"`
	Contacts     *ResumeContactsView `json:"contacts,omitempty"`
	Skills       *SkillsView         `json:"skills,omitempty"`
	Languages    []LanguageView      `json:"languages"`
	Technologies []TechnologyView    `json:"technologies"`
	Experiences  []ExperienceView    `json:"experiences"`
	Educations   []EducationView     `json:"educations"`
	Courses      []CourseView        `json:"courses"`
	Activities   []ActivityView      `json:"activities"`
	Projects     []ProjectView       `json:"projects"`
}

type GetResumeQueryHandler struct{ reader ResumeReader }

func NewGetResumeQueryHandler(reader ResumeReader) *GetResumeQueryHandler {
	return &GetResumeQueryHandler{reader: reader}
}
func (h *GetResumeQueryHandler) Handle(ctx context.Context, q GetResumeQuery) (ResumeView, error) {
	if h == nil || h.reader == nil {
		return ResumeView{}, fmt.Errorf("get resume query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return ResumeView{}, fmt.Errorf("invalid get resume query")
	}
	return h.reader.GetOwned(ctx, q)
}
