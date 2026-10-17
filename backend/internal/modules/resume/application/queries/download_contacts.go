package queries

import (
	"context"
	"fmt"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type DownloadContactsQuery struct {
	ResumeID sharedvalueobjects.ID
}

type DownloadContactsView struct {
	Phone    *string `json:"phone,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type DownloadContactsReader interface {
	GetContactsForDownload(context.Context, DownloadContactsQuery) (DownloadContactsView, error)
}

type DownloadContactsQueryHandler struct {
	reader DownloadContactsReader
}

func NewDownloadContactsQueryHandler(reader DownloadContactsReader) *DownloadContactsQueryHandler {
	return &DownloadContactsQueryHandler{reader: reader}
}

func (h *DownloadContactsQueryHandler) Handle(
	ctx context.Context,
	query DownloadContactsQuery,
) (DownloadContactsView, error) {
	if h == nil || h.reader == nil {
		return DownloadContactsView{}, fmt.Errorf("download contacts query handler dependencies are nil")
	}
	if query.ResumeID.IsZero() {
		return DownloadContactsView{}, fmt.Errorf("invalid download contacts query")
	}
	return h.reader.GetContactsForDownload(ctx, query)
}
