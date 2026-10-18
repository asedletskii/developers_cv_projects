package commands

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/analytics/application/ports"
	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type RecordPublicInteractionEventHandler struct {
	events ports.EventRepository
	views  ports.ViewRepository
	uow    ports.UnitOfWork
	clock  func() time.Time
}

func NewRecordPublicInteractionEventHandler(
	events ports.EventRepository,
	views ports.ViewRepository,
	uow ports.UnitOfWork,
) *RecordPublicInteractionEventHandler {
	return &RecordPublicInteractionEventHandler{
		events: events,
		views:  views,
		uow:    uow,
		clock:  func() time.Time { return time.Now().UTC() },
	}
}

func (h *RecordPublicInteractionEventHandler) RecordContactView(
	ctx context.Context,
	viewID, linkID sharedvalueobjects.ID,
	contact string,
) error {
	if !domain.IsValidAnalyticsContactKey(contact) {
		return fmt.Errorf("%w: invalid contact view event input", domain.ErrInvalidAnalyticsEvent)
	}
	return h.record(ctx, viewID, linkID, domain.EventTypeContactView, map[string]string{"contact": contact})
}

func (h *RecordPublicInteractionEventHandler) RecordDownload(
	ctx context.Context,
	viewID, linkID sharedvalueobjects.ID,
) error {
	return h.record(ctx, viewID, linkID, domain.EventTypeDownload, nil)
}

func (h *RecordPublicInteractionEventHandler) record(
	ctx context.Context,
	viewID, linkID sharedvalueobjects.ID,
	eventType domain.AnalyticsEventType,
	payload map[string]string,
) error {
	if h == nil || h.events == nil || h.views == nil || h.uow == nil {
		return fmt.Errorf("record public interaction event: dependencies are nil")
	}
	if viewID.IsZero() || linkID.IsZero() {
		return fmt.Errorf("%w: invalid public interaction event input", domain.ErrInvalidAnalyticsEvent)
	}
	view, err := h.views.GetView(ctx, viewID)
	if err != nil {
		return err
	}
	if view == nil || view.LinkID() != linkID {
		return nil
	}

	event, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
		ViewID:     viewID,
		LinkID:     linkID,
		Type:       eventType,
		Payload:    payload,
		OccurredAt: h.clock().UTC(),
	})
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(txCtx context.Context) error {
		return h.events.RecordEvents(txCtx, []*domain.AnalyticsEvent{event})
	})
}
