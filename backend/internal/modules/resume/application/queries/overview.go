package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
)

type GetOverviewQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type OverviewReader interface {
	GetOverview(context.Context, GetOverviewQuery) (ResumeOverviewView, error)
}
type GetOverviewQueryHandler struct{ reader OverviewReader }

func NewGetOverviewQueryHandler(r OverviewReader) *GetOverviewQueryHandler {
	return &GetOverviewQueryHandler{reader: r}
}
func (h *GetOverviewQueryHandler) Handle(ctx context.Context, q GetOverviewQuery) (ResumeOverviewView, error) {
	if h == nil || h.reader == nil {
		return ResumeOverviewView{}, fmt.Errorf("get overview query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return ResumeOverviewView{}, fmt.Errorf("invalid get overview query")
	}
	return h.reader.GetOverview(ctx, q)
}

type ResumeOverviewView struct {
	Presentation OverviewPresentationView `json:"presentation"`
}
type OverviewPresentationView struct {
	Theme        *OverviewThemeSettings    `json:"theme,omitempty"`
	Greeting     OverviewGreeting          `json:"greeting"`
	Availability *OverviewAvailability     `json:"availability,omitempty"`
	ResumeLink   *string                   `json:"resumeLink,omitempty"`
	Sections     []OverviewSectionSettings `json:"sections"`
}
type OverviewGreeting struct {
	First           string `json:"first"`
	Second          string `json:"second"`
	SecondHighlight string `json:"secondHighlight"`
	Third           string `json:"third"`
}
type OverviewAvailability struct {
	Status string  `json:"status"`
	From   *string `json:"from,omitempty"`
	Note   *string `json:"note,omitempty"`
}
type OverviewThemeSettings struct {
	Primary   string `json:"primary,omitempty"`
	Secondary string `json:"secondary,omitempty"`
}
type OverviewSectionSettings struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}
