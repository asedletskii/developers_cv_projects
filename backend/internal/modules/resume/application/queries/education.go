package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
)

type ListEducationsQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetEducationQuery struct {
	ResumeID, EducationID valueobjects.ID
	UserID                useridentity.ID
}
type EducationView struct {
	ID          string  `json:"id"`
	ResumeID    string  `json:"resumeId"`
	Location    string  `json:"location"`
	University  string  `json:"university"`
	Grade       int     `json:"grade"`
	ProgramName string  `json:"programName"`
	ProfileName *string `json:"profileName,omitempty"`
	Faculty     string  `json:"faculty"`
	Course      *int    `json:"course,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
}
type EducationsReader interface {
	ListEducations(context.Context, ListEducationsQuery) ([]EducationView, error)
	GetEducation(context.Context, GetEducationQuery) (EducationView, error)
}
type ListEducationsQueryHandler struct{ reader EducationsReader }

func NewListEducationsQueryHandler(reader EducationsReader) *ListEducationsQueryHandler {
	return &ListEducationsQueryHandler{reader: reader}
}
func (h *ListEducationsQueryHandler) Handle(ctx context.Context, q ListEducationsQuery) ([]EducationView, error) {
	if h == nil || h.reader == nil {
		return nil, fmt.Errorf("list educations query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("invalid list educations query")
	}
	return h.reader.ListEducations(ctx, q)
}

type GetEducationQueryHandler struct{ reader EducationsReader }

func NewGetEducationQueryHandler(reader EducationsReader) *GetEducationQueryHandler {
	return &GetEducationQueryHandler{reader: reader}
}
func (h *GetEducationQueryHandler) Handle(ctx context.Context, q GetEducationQuery) (EducationView, error) {
	if h == nil || h.reader == nil {
		return EducationView{}, fmt.Errorf("get education query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.EducationID.IsZero() || q.UserID.IsZero() {
		return EducationView{}, fmt.Errorf("invalid get education query")
	}
	return h.reader.GetEducation(ctx, q)
}
