package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
)

type GetSkillsQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type SkillsReader interface {
	GetSkills(context.Context, GetSkillsQuery) (SkillsView, error)
}

type SkillsSource interface {
	ReadSkills(context.Context, valueobjects.ID) (SkillsView, error)
}
type GetSkillsQueryHandler struct{ reader SkillsReader }

func NewGetSkillsQueryHandler(reader SkillsReader) *GetSkillsQueryHandler {
	return &GetSkillsQueryHandler{reader: reader}
}
func (h *GetSkillsQueryHandler) Handle(ctx context.Context, q GetSkillsQuery) (SkillsView, error) {
	if h == nil || h.reader == nil {
		return SkillsView{}, fmt.Errorf("get skills query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return SkillsView{}, fmt.Errorf("invalid get skills query")
	}
	return h.reader.GetSkills(ctx, q)
}

type SkillsView struct {
	ID         string   `json:"id"`
	ResumeID   string   `json:"resumeId"`
	HardSkills []string `json:"hardSkills"`
	SoftSkills []string `json:"softSkills"`
}
