package queries

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type GetResumeRootQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type ResumeRootReader interface {
	GetResumeRoot(context.Context, GetResumeRootQuery) (ResumeRootView, error)
}

type GetResumeRootQueryHandler struct {
	reader ResumeRootReader
}

func NewGetResumeRootQueryHandler(reader ResumeRootReader) *GetResumeRootQueryHandler {
	return &GetResumeRootQueryHandler{reader: reader}
}

func (h *GetResumeRootQueryHandler) Handle(ctx context.Context, query GetResumeRootQuery) (ResumeRootView, error) {
	if h == nil || h.reader == nil {
		return ResumeRootView{}, fmt.Errorf("get resume root query handler dependencies are nil")
	}
	if query.ResumeID.IsZero() || query.UserID.IsZero() {
		return ResumeRootView{}, fmt.Errorf("invalid get resume root query")
	}
	return h.reader.GetResumeRoot(ctx, query)
}

type ResumeRootView struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
	Presentation ResumePresentationView `json:"presentation"`
}

type ResumePresentationView struct {
	Theme        *ResumeRootThemeSettings    `json:"theme,omitempty"`
	Greeting     ResumeRootGreeting          `json:"greeting"`
	Availability *ResumeRootAvailability     `json:"availability,omitempty"`
	ResumeLink   *string                     `json:"resumeLink,omitempty"`
	Sections     []ResumeRootSectionSettings `json:"sections"`
}

type ResumeRootLocalizedText struct {
	RU string `json:"ru"`
	EN string `json:"en"`
}
type ResumeRootGreeting struct {
	First           ResumeRootLocalizedText `json:"first"`
	Second          ResumeRootLocalizedText `json:"second"`
	SecondHighlight ResumeRootLocalizedText `json:"secondHighlight"`
	Third           ResumeRootLocalizedText `json:"third"`
}
type ResumeRootAvailability struct {
	Status string                   `json:"status"`
	From   *string                  `json:"from,omitempty"`
	Note   *ResumeRootLocalizedText `json:"note,omitempty"`
}
type ResumeRootThemeSettings struct {
	Primary   string `json:"primary,omitempty"`
	Secondary string `json:"secondary,omitempty"`
}
type ResumeRootSectionSettings struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}
