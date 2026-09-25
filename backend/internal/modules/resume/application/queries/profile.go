package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"fmt"
)

type GetResumeProfileQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type ResumeProfileReader interface {
	GetResumeProfile(context.Context, GetResumeProfileQuery) (ResumeProfileView, error)
}
type GetResumeProfileQueryHandler struct{ reader ResumeProfileReader }

func NewGetResumeProfileQueryHandler(reader ResumeProfileReader) *GetResumeProfileQueryHandler {
	return &GetResumeProfileQueryHandler{reader: reader}
}
func (h *GetResumeProfileQueryHandler) Handle(ctx context.Context, q GetResumeProfileQuery) (ResumeProfileView, error) {
	if h == nil || h.reader == nil {
		return ResumeProfileView{}, fmt.Errorf("get resume profile query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return ResumeProfileView{}, fmt.Errorf("invalid get resume profile query")
	}
	return h.reader.GetResumeProfile(ctx, q)
}

type ResumeProfileView struct {
	ID          string   `json:"id"`
	ResumeID    string   `json:"resumeId"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Patronymic  *string  `json:"patronymic,omitempty"`
	Position    string   `json:"position"`
	PhotoPath   string   `json:"photoPath"`
	DateOfBirth string   `json:"dateOfBirth"`
	Location    string   `json:"location"`
	Relocation  bool     `json:"relocation"`
	WorkFormat  []string `json:"workFormat"`
}
