package queries

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	ViewID  sharedvalueobjects.ID
}

type ResumeContactViewQuery interface {
	Handle(context.Context, resumequeries.ContactViewQuery) (resumequeries.ContactView, error)
}

type ContactViewThroughLinkQueryHandler struct {
	links  ResumeLinkPublicReader
	resume ResumeContactViewQuery
	events ports.AnalyticsEventRepository
	views  ports.ViewRepository
}

func NewContactViewThroughLinkQueryHandler(links ResumeLinkPublicReader, resume ResumeContactViewQuery) *ContactViewThroughLinkQueryHandler {
	return &ContactViewThroughLinkQueryHandler{links: links, resume: resume}
}

func NewContactViewThroughLinkQueryHandlerWithAnalytics(links ResumeLinkPublicReader, resume ResumeContactViewQuery, events ports.AnalyticsEventRepository, views ports.ViewRepository) *ContactViewThroughLinkQueryHandler {
	return &ContactViewThroughLinkQueryHandler{links: links, resume: resume, events: events, views: views}
}

func (h *ContactViewThroughLinkQueryHandler) Handle(ctx context.Context, query ContactViewThroughLinkQuery) (resumequeries.ContactView, error) {
	if h == nil || h.links == nil || h.resume == nil {
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
	if !query.ViewID.IsZero() && h.events != nil && h.views != nil {
		analyticsView, viewErr := h.views.GetView(ctx, query.ViewID)

		if viewErr == nil && analyticsView.LinkID() == query.LinkID {
			h.recordContactViewEvent(ctx, query.ViewID, query.LinkID, query.Contact)
		}
	}
	return view, nil
}

func (h *ContactViewThroughLinkQueryHandler) recordContactViewEvent(
	ctx context.Context,
	viewID sharedvalueobjects.ID,
	linkID sharedvalueobjects.ID,
	contact resumedomain.ContactKey,
) {
	event, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
		ID:         sharedvalueobjects.NewID(),
		ViewID:     viewID,
		LinkID:     linkID,
		Type:       domain.EventTypeContactView,
		Payload:    map[string]string{"contact": string(contact)},
		OccurredAt: time.Now().UTC(),
	})
	if err != nil {
		return
	}

	if err := h.events.RecordEvents(ctx, []*domain.AnalyticsEvent{event}); err != nil {
		return
	}
}
