package queries

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicResumeLinkQuery struct {
	LinkID valueobjects.ID
}

type ResumeLinkPublicView struct {
	ID       string
	ResumeID string
	Status   string
}

type ResumeLinkPublicReader interface {
	GetByID(context.Context, PublicResumeLinkQuery) (ResumeLinkPublicView, error)
}

type ResumePublicQuery interface {
	Handle(context.Context, resumequeries.PublicResumeQuery) (resumequeries.ResumeView, error)
}

type PublicResumeLinkQueryHandler struct {
	links  ResumeLinkPublicReader
	resume ResumePublicQuery
	events ports.LinkEventRecorder
}

func NewPublicResumeLinkQueryHandler(links ResumeLinkPublicReader, resume ResumePublicQuery, events ports.LinkEventRecorder) *PublicResumeLinkQueryHandler {
	return &PublicResumeLinkQueryHandler{links: links, resume: resume, events: events}
}

func (h *PublicResumeLinkQueryHandler) Handle(ctx context.Context, query PublicResumeLinkQuery) (resumequeries.ResumeView, error) {
	if h == nil || h.links == nil || h.resume == nil || h.events == nil {
		return resumequeries.ResumeView{}, fmt.Errorf("public resume link query handler dependencies are nil")
	}
	if query.LinkID.IsZero() {
		return resumequeries.ResumeView{}, fmt.Errorf("invalid resume link id")
	}

	link, err := h.links.GetByID(ctx, query)
	if err != nil {
		return resumequeries.ResumeView{}, err
	}
	if link.Status == domain.ResumeLinkStatusInactive {
		return resumequeries.ResumeView{}, domain.ErrResumeLinkInactive
	}
	if link.Status != domain.ResumeLinkStatusActive {
		return resumequeries.ResumeView{}, domain.ErrEntityNotFound
	}

	resumeUUID, err := uuid.Parse(link.ResumeID)
	if err != nil {
		return resumequeries.ResumeView{}, fmt.Errorf("resolve resume link: invalid resume id: %w", err)
	}
	resumeID, err := valueobjects.IDFromUUID(resumeUUID)
	if err != nil {
		return resumequeries.ResumeView{}, fmt.Errorf("resolve resume link: invalid resume id: %w", err)
	}

	view, err := h.resume.Handle(ctx, resumequeries.PublicResumeQuery{ID: resumeID})
	if err != nil {
		return resumequeries.ResumeView{}, err
	}
	recordLinkEvent(ctx, h.events, query.LinkID, domain.EventTypeView)
	return view, nil
}
