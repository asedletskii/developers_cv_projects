package queries

import (
	"context"
	"fmt"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	valueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GetResumeLinkStatisticsQuery struct {
	UserID valueobjects.ID
	LinkID valueobjects.ID
}

type ResumeLinkStatistics struct {
	Views        int64 `json:"views"`
	ContactViews int64 `json:"contactViews"`
}

type ResumeLinkStatisticsReader interface {
	GetStatistics(context.Context, GetResumeLinkStatisticsQuery) (ResumeLinkStatistics, error)
}

type GetResumeLinkStatisticsQueryHandler struct {
	links   ResumeLinkPublicReader
	resumes ResumeOwnershipReader
	reader  ResumeLinkStatisticsReader
}

func NewGetResumeLinkStatisticsQueryHandler(
	links ResumeLinkPublicReader,
	resumes ResumeOwnershipReader,
	reader ResumeLinkStatisticsReader,
) *GetResumeLinkStatisticsQueryHandler {
	return &GetResumeLinkStatisticsQueryHandler{links: links, resumes: resumes, reader: reader}
}

func (h *GetResumeLinkStatisticsQueryHandler) Handle(
	ctx context.Context,
	query GetResumeLinkStatisticsQuery,
) (ResumeLinkStatistics, error) {
	if h == nil || h.links == nil || h.resumes == nil || h.reader == nil {
		return ResumeLinkStatistics{}, fmt.Errorf("get resume link statistics query handler dependencies are nil")
	}
	if query.UserID.IsZero() || query.LinkID.IsZero() {
		return ResumeLinkStatistics{}, fmt.Errorf("invalid get resume link statistics query")
	}

	link, err := h.links.GetByID(ctx, PublicResumeLinkQuery{LinkID: query.LinkID})
	if err != nil {
		return ResumeLinkStatistics{}, err
	}

	resumeID, err := valueobjects.ParseID(link.ResumeID)
	if err != nil {
		return ResumeLinkStatistics{}, fmt.Errorf("resolve resume link: invalid resume id: %w", err)
	}
	if _, err := h.resumes.GetOwned(ctx, resumequeries.GetResumeQuery{
		ResumeID: resumeID,
		UserID:   query.UserID,
	}); err != nil {
		return ResumeLinkStatistics{}, err
	}

	return h.reader.GetStatistics(ctx, query)
}
