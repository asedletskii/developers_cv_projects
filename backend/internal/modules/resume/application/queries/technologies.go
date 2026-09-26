package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ListTechnologiesQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type GetTechnologyQuery struct {
	ResumeID     valueobjects.ID
	TechnologyID valueobjects.ID
	UserID       useridentity.ID
}

type TechnologiesReader interface {
	ListTechnologies(context.Context, ListTechnologiesQuery) ([]TechnologyView, error)
	GetTechnology(context.Context, GetTechnologyQuery) (TechnologyView, error)
}

type TechnologiesSource interface {
	ReadTechnologies(context.Context, valueobjects.ID) ([]TechnologyView, error)
}

type ListTechnologiesQueryHandler struct{ reader TechnologiesReader }

func NewListTechnologiesQueryHandler(reader TechnologiesReader) *ListTechnologiesQueryHandler {
	return &ListTechnologiesQueryHandler{reader: reader}
}

func (h *ListTechnologiesQueryHandler) Handle(ctx context.Context, q ListTechnologiesQuery) ([]TechnologyView, error) {
	if h == nil || h.reader == nil {
		return nil, fmt.Errorf("list technologies query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("invalid list technologies query")
	}
	return h.reader.ListTechnologies(ctx, q)
}

type GetTechnologyQueryHandler struct{ reader TechnologiesReader }

func NewGetTechnologyQueryHandler(reader TechnologiesReader) *GetTechnologyQueryHandler {
	return &GetTechnologyQueryHandler{reader: reader}
}

func (h *GetTechnologyQueryHandler) Handle(ctx context.Context, q GetTechnologyQuery) (TechnologyView, error) {
	if h == nil || h.reader == nil {
		return TechnologyView{}, fmt.Errorf("get technology query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.TechnologyID.IsZero() || q.UserID.IsZero() {
		return TechnologyView{}, fmt.Errorf("invalid get technology query")
	}
	return h.reader.GetTechnology(ctx, q)
}

type TechnologyView struct {
	ID       string  `json:"id"`
	ResumeID string  `json:"resumeId"`
	Name     string  `json:"name"`
	Grade    *string `json:"grade,omitempty"`
	Position int     `json:"position"`
}
