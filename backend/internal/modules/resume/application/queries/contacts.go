package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
)

type GetResumeContactsQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type ResumeContactsReader interface {
	GetResumeContacts(context.Context, GetResumeContactsQuery) (ResumeContactsView, error)
}

type ResumeContactsSource interface {
	ReadResumeContacts(context.Context, valueobjects.ID) (ResumeContactsView, error)
}
type GetResumeContactsQueryHandler struct{ reader ResumeContactsReader }

func NewGetResumeContactsQueryHandler(reader ResumeContactsReader) *GetResumeContactsQueryHandler {
	return &GetResumeContactsQueryHandler{reader: reader}
}
func (h *GetResumeContactsQueryHandler) Handle(ctx context.Context, q GetResumeContactsQuery) (ResumeContactsView, error) {
	if h == nil || h.reader == nil {
		return ResumeContactsView{}, fmt.Errorf("get resume contacts query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return ResumeContactsView{}, fmt.Errorf("invalid get resume contacts query")
	}
	return h.reader.GetResumeContacts(ctx, q)
}

type ResumeContactsView struct {
	ID             string                  `json:"id"`
	ResumeID       string                  `json:"resumeId"`
	Phone          *string                 `json:"phone,omitempty"`
	Telegram       *string                 `json:"telegram,omitempty"`
	Email          *string                 `json:"email,omitempty"`
	HiddenContacts []string                `json:"hiddenContacts"`
	ContactsText   *ResumeContactsTextView `json:"contactsText,omitempty"`
}
type ResumeContactsTextView struct {
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
}
