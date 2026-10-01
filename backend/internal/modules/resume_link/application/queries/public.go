package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicResumeLinkQuery struct {
	LinkID sharedvalueobjects.ID
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
	Handle(context.Context, resumequeries.PublicResumeQuery) (resumequeries.PublicResumeView, error)
}

type PublicResumeLinkQueryHandler struct {
	links  ResumeLinkPublicReader
	resume ResumePublicQuery
	events ports.LinkEventRecorder
}

func NewPublicResumeLinkQueryHandler(links ResumeLinkPublicReader, resume ResumePublicQuery, events ports.LinkEventRecorder) *PublicResumeLinkQueryHandler {
	return &PublicResumeLinkQueryHandler{links: links, resume: resume, events: events}
}

func (h *PublicResumeLinkQueryHandler) Handle(ctx context.Context, query PublicResumeLinkQuery) (resumequeries.PublicResumeView, error) {
	if h == nil || h.links == nil || h.resume == nil || h.events == nil {
		return resumequeries.PublicResumeView{}, fmt.Errorf("public resume link query handler dependencies are nil")
	}
	if query.LinkID.IsZero() {
		return resumequeries.PublicResumeView{}, fmt.Errorf("invalid resume link id")
	}

	link, err := h.links.GetByID(ctx, query)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return resumequeries.PublicResumeView{}, ErrResumeLinkNotFound
		}
		return resumequeries.PublicResumeView{}, err
	}
	if domain.ResumeLinkStatus(link.Status) == domain.ResumeLinkStatusInactive {
		return resumequeries.PublicResumeView{}, domain.ErrResumeLinkInactive
	}
	if domain.ResumeLinkStatus(link.Status) != domain.ResumeLinkStatusActive {
		return resumequeries.PublicResumeView{}, ErrResumeLinkNotFound
	}

	resumeUUID, err := uuid.Parse(link.ResumeID)
	if err != nil {
		return resumequeries.PublicResumeView{}, fmt.Errorf("resolve resume link: invalid resume id: %w", err)
	}
	resumeID, err := sharedvalueobjects.IDFromUUID(resumeUUID)
	if err != nil {
		return resumequeries.PublicResumeView{}, fmt.Errorf("resolve resume link: invalid resume id: %w", err)
	}

	view, err := h.resume.Handle(ctx, resumequeries.PublicResumeQuery{ID: resumeID})
	if err != nil {
		return resumequeries.PublicResumeView{}, err
	}
	recordLinkEvent(ctx, h.events, query.LinkID, domain.EventTypeView)
	return view, nil
}
