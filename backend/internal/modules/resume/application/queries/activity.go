package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
)

type ListActivitiesQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetActivityQuery struct {
	ResumeID, ActivityID valueobjects.ID
	UserID               useridentity.ID
}
type ActivityView struct {
	ID          string  `json:"id"`
	ResumeID    string  `json:"resumeId"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Date        *string `json:"date,omitempty"`
	Link        *string `json:"link,omitempty"`
}
type ActivitiesReader interface {
	ListActivities(context.Context, ListActivitiesQuery) ([]ActivityView, error)
	GetActivity(context.Context, GetActivityQuery) (ActivityView, error)
}
type ListActivitiesQueryHandler struct{ reader ActivitiesReader }

func NewListActivitiesQueryHandler(reader ActivitiesReader) *ListActivitiesQueryHandler {
	return &ListActivitiesQueryHandler{reader: reader}
}
func (h *ListActivitiesQueryHandler) Handle(ctx context.Context, q ListActivitiesQuery) ([]ActivityView, error) {
	if h == nil || h.reader == nil {
		return nil, fmt.Errorf("list activities query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("invalid list activities query")
	}
	return h.reader.ListActivities(ctx, q)
}

type GetActivityQueryHandler struct{ reader ActivitiesReader }

func NewGetActivityQueryHandler(reader ActivitiesReader) *GetActivityQueryHandler {
	return &GetActivityQueryHandler{reader: reader}
}
func (h *GetActivityQueryHandler) Handle(ctx context.Context, q GetActivityQuery) (ActivityView, error) {
	if h == nil || h.reader == nil {
		return ActivityView{}, fmt.Errorf("get activity query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.ActivityID.IsZero() || q.UserID.IsZero() {
		return ActivityView{}, fmt.Errorf("invalid get activity query")
	}
	return h.reader.GetActivity(ctx, q)
}
