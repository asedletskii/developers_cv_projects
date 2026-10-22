package queries

import (
	"context"
	"fmt"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
)

type DownloadContactsReader interface {
	GetContactsForDownload(context.Context, resumecontracts.DownloadContactsQuery) (resumecontracts.DownloadContactsView, error)
}

type DownloadContactsQueryHandler struct {
	reader DownloadContactsReader
}

func NewDownloadContactsQueryHandler(reader DownloadContactsReader) *DownloadContactsQueryHandler {
	return &DownloadContactsQueryHandler{reader: reader}
}

func (h *DownloadContactsQueryHandler) Handle(
	ctx context.Context,
	query resumecontracts.DownloadContactsQuery,
) (resumecontracts.DownloadContactsView, error) {
	if h == nil || h.reader == nil {
		return resumecontracts.DownloadContactsView{}, fmt.Errorf("download contacts query handler dependencies are nil")
	}
	if query.ResumeID.IsZero() {
		return resumecontracts.DownloadContactsView{}, fmt.Errorf("invalid download contacts query")
	}
	return h.reader.GetContactsForDownload(ctx, query)
}
