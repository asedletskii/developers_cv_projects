package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ContactViewThroughLinkQuery struct {
	LinkID  sharedvalueobjects.ID
	Contact resumedomain.ContactKey
}

type ResumeContactViewQuery interface {
	Handle(context.Context, resumequeries.ContactViewQuery) (resumequeries.ContactView, error)
}

type ContactViewThroughLinkQueryHandler struct {
	links  ResumeLinkPublicReader
	resume ResumeContactViewQuery
	events ports.LinkEventRecorder
}

func NewContactViewThroughLinkQueryHandler(links ResumeLinkPublicReader, resume ResumeContactViewQuery, events ports.LinkEventRecorder) *ContactViewThroughLinkQueryHandler {
	return &ContactViewThroughLinkQueryHandler{links: links, resume: resume, events: events}
}

func (h *ContactViewThroughLinkQueryHandler) Handle(ctx context.Context, query ContactViewThroughLinkQuery) (resumequeries.ContactView, error) {
	if h == nil || h.links == nil || h.resume == nil || h.events == nil {
		return resumequeries.ContactView{}, fmt.Errorf("contact view through link query handler dependencies are nil")
	}
	if query.LinkID.IsZero() || !query.Contact.IsValid() {
		return resumequeries.ContactView{}, fmt.Errorf("invalid contact view through link query")
	}

	link, err := h.links.GetByID(ctx, PublicResumeLinkQuery{LinkID: query.LinkID})
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return resumequeries.ContactView{}, ErrResumeLinkNotFound
		}
		return resumequeries.ContactView{}, err
	}
	if domain.ResumeLinkStatus(link.Status) == domain.ResumeLinkStatusInactive {
		return resumequeries.ContactView{}, domain.ErrResumeLinkInactive
	}
	if domain.ResumeLinkStatus(link.Status) != domain.ResumeLinkStatusActive {
		return resumequeries.ContactView{}, ErrResumeLinkNotFound
	}

	resumeUUID, err := uuid.Parse(link.ResumeID)
	if err != nil {
		return resumequeries.ContactView{}, fmt.Errorf("resolve resume link: invalid resume id: %w", err)
	}
	resumeID, err := sharedvalueobjects.IDFromUUID(resumeUUID)
	if err != nil {
		return resumequeries.ContactView{}, fmt.Errorf("resolve resume link: invalid resume id: %w", err)
	}

	view, err := h.resume.Handle(ctx, resumequeries.ContactViewQuery{ResumeID: resumeID, Contact: query.Contact})
	if err != nil {
		return resumequeries.ContactView{}, err
	}
	recordLinkEvent(ctx, h.events, query.LinkID, domain.EventTypeContactView)
	return view, nil
}
