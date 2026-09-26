package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ListLanguagesQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type GetLanguageQuery struct {
	ResumeID   valueobjects.ID
	LanguageID valueobjects.ID
	UserID     useridentity.ID
}

type LanguagesReader interface {
	ListLanguages(context.Context, ListLanguagesQuery) ([]LanguageView, error)
	GetLanguage(context.Context, GetLanguageQuery) (LanguageView, error)
}

type ListLanguagesQueryHandler struct{ reader LanguagesReader }

func NewListLanguagesQueryHandler(reader LanguagesReader) *ListLanguagesQueryHandler {
	return &ListLanguagesQueryHandler{reader: reader}
}

func (h *ListLanguagesQueryHandler) Handle(ctx context.Context, q ListLanguagesQuery) ([]LanguageView, error) {
	if h == nil || h.reader == nil {
		return nil, fmt.Errorf("list languages query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("invalid list languages query")
	}
	return h.reader.ListLanguages(ctx, q)
}

type GetLanguageQueryHandler struct{ reader LanguagesReader }

func NewGetLanguageQueryHandler(reader LanguagesReader) *GetLanguageQueryHandler {
	return &GetLanguageQueryHandler{reader: reader}
}

func (h *GetLanguageQueryHandler) Handle(ctx context.Context, q GetLanguageQuery) (LanguageView, error) {
	if h == nil || h.reader == nil {
		return LanguageView{}, fmt.Errorf("get language query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.LanguageID.IsZero() || q.UserID.IsZero() {
		return LanguageView{}, fmt.Errorf("invalid get language query")
	}
	return h.reader.GetLanguage(ctx, q)
}

type LanguageView struct {
	ID       string  `json:"id"`
	ResumeID string  `json:"resumeId"`
	Language string  `json:"language"`
	Grade    *string `json:"grade,omitempty"`
	Position int     `json:"position"`
}
