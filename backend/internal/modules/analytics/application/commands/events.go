package commands

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/analytics/application/ports"
	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const (
	MaxAnalyticsEventsPerBatch = 100
	MaxAnalyticsClicksPerBatch = 200
	MaxAnalyticsEventsPerView  = domain.MaxAnalyticsEventsPerView
)

type AnalyticsEventInput struct {
	ClientEventID     uuid.UUID
	Type              domain.AnalyticsEventType
	Payload           map[string]string
	ClientTimestampMS int64
}

type AnalyticsClickInput struct {
	Section           string
	XPct              int
	YPct              int
	ViewportW         int
	Target            string
	ClientTimestampMS int64
}

type AnalyticsBatchInput struct {
	ViewID      sharedvalueobjects.ID
	LinkID      sharedvalueobjects.ID
	Seq         int64
	ActiveMS    int64
	MaxSection  string
	ExitSection string
	Final       bool
	Events      []AnalyticsEventInput
	Clicks      []AnalyticsClickInput
}

type IngestAnalyticsBatchHandler struct {
	repository ports.AnalyticsBatchRepository
	uow        ports.UnitOfWork
	observer   ports.AnalyticsObserver
	clock      func() time.Time
}

func NewIngestAnalyticsBatchHandler(
	repository ports.AnalyticsBatchRepository,
	uow ports.UnitOfWork,
	observer ports.AnalyticsObserver,
) *IngestAnalyticsBatchHandler {
	return &IngestAnalyticsBatchHandler{
		repository: repository,
		uow:        uow,
		observer:   observer,
		clock:      func() time.Time { return time.Now().UTC() },
	}
}

func (h *IngestAnalyticsBatchHandler) Handle(ctx context.Context, input AnalyticsBatchInput) error {
	if h == nil {
		return fmt.Errorf("ingest analytics batch: handler is nil")
	}
	if h.repository == nil || h.uow == nil {
		h.observeBatchOutcome("failure", "unknown")
		return fmt.Errorf("ingest analytics batch: dependencies are nil")
	}
	if input.ViewID.IsZero() || input.LinkID.IsZero() || input.Seq <= 0 || input.ActiveMS < 0 {
		h.observeBatchOutcome("rejected", "invalid_batch")
		return fmt.Errorf("%w: invalid analytics batch", domain.ErrInvalidAnalyticsEvent)
	}
	if len(input.Events) > MaxAnalyticsEventsPerBatch || len(input.Clicks) > MaxAnalyticsClicksPerBatch {
		h.observeBatchOutcome("rejected", "batch_limit")
		return fmt.Errorf("%w: too many analytics inputs", domain.ErrInvalidAnalyticsEvent)
	}

	orderedInput := append([]AnalyticsEventInput(nil), input.Events...)
	sort.SliceStable(orderedInput, func(i, j int) bool {
		return orderedInput[i].ClientTimestampMS < orderedInput[j].ClientTimestampMS
	})
	orderedClicks := append([]AnalyticsClickInput(nil), input.Clicks...)
	sort.SliceStable(orderedClicks, func(i, j int) bool {
		return orderedClicks[i].ClientTimestampMS < orderedClicks[j].ClientTimestampMS
	})
	now := h.clock().UTC()

	needsPrintHistory := false
	for _, item := range orderedInput {
		if item.Type == domain.EventTypePrint {
			needsPrintHistory = true
			break
		}
	}
	outcome, reason := "accepted", "success"
	err := h.uow.Do(ctx, func(txCtx context.Context) error {
		snapshot, err := h.repository.LoadAnalyticsBatchForUpdate(
			txCtx, input.ViewID, input.LinkID,
			len(orderedInput) > 0, len(orderedClicks) > 0, needsPrintHistory,
		)
		if err != nil {
			outcome, reason = classifyAnalyticsBatchError(err), classifyAnalyticsBatchReason(err)
			if outcome == "failure" {
				reason = "load_failed"
			}
			return err
		}
		view := snapshot.View
		if view == nil || view.ID() != input.ViewID || view.LinkID() != input.LinkID {
			outcome, reason = "failure", "repository_mismatch"
			return fmt.Errorf("ingest analytics batch: repository returned a mismatched view")
		}
		update := domain.ViewStateUpdate{
			Seq: input.Seq, ActiveMS: input.ActiveMS,
			MaxSection: domain.AnalyticsSectionID(input.MaxSection), ExitSection: domain.AnalyticsSectionID(input.ExitSection),
			SeenAt: now, Final: input.Final,
		}
		decision, err := view.EvaluateAnalyticsBatch(snapshot.LinkContext.LinkActive, snapshot.LinkContext.ResumePublished, now, update)
		if err != nil {
			outcome, reason = classifyAnalyticsBatchError(err), classifyAnalyticsBatchReason(err)
			return err
		}
		if decision.Ignore {
			outcome = "ignored"
			if view.Traffic() == domain.TrafficBot {
				reason = "bot_traffic"
			} else {
				reason = "stale_sequence"
			}
			return nil
		}
		if err := view.ApplyState(update); err != nil {
			outcome, reason = "rejected", "invalid_state"
			return fmt.Errorf("ingest analytics batch: apply state: %w", err)
		}

		events := make([]*domain.AnalyticsEvent, 0, len(orderedInput))
		if !decision.SameSeq {
			seen := make(map[uuid.UUID]struct{}, len(orderedInput))
			for index, item := range orderedInput {
				if item.Type == domain.EventTypeContactView {
					h.observeDroppedInput("event", "contact_view_routed_separately")
					continue
				}
				if item.ClientEventID == uuid.Nil {
					h.observeDroppedInput("event", "missing_client_id")
					continue
				}
				if _, exists := seen[item.ClientEventID]; exists {
					h.observeDroppedInput("event", "duplicate_client_id")
					continue
				}
				seen[item.ClientEventID] = struct{}{}

				clientID, err := sharedvalueobjects.IDFromUUID(item.ClientEventID)
				if err != nil {
					h.observeDroppedInput("event", "invalid_client_id")
					continue
				}
				payload := item.Payload
				if item.Type == domain.EventTypeTechFilter && view.Layout().HasTechnologyCatalog() {
					canonical, ok := view.Layout().NormalizeTechnology(payload["technology"])
					if !ok {
						h.observeDroppedInput("event", "technology_not_in_resume")
						continue
					}
					payload = make(map[string]string, len(item.Payload))
					for key, value := range item.Payload {
						payload[key] = value
					}
					payload["technology"] = canonical
				}
				event, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
					ViewID: input.ViewID, LinkID: input.LinkID,
					ClientEventID: &clientID, Type: item.Type, Payload: payload,
					OccurredAt: now.Add(time.Duration(index) * time.Nanosecond),
				})
				if err != nil {
					h.observeDroppedInput("event", "invalid_payload_or_type")
					continue
				}
				events = append(events, event)
			}
			before := len(events)
			events = domain.SelectAnalyticsEventsForStorage(events, snapshot.StoredEventCount, snapshot.ExistingPrints)
			if snapshot.StoredEventCount+int64(before) > int64(domain.MaxAnalyticsEventsPerView) {
				dropped := before - len(events)
				if dropped > 0 {
					h.observeDroppedInputN("event", "per_view_limit", dropped)
				}
			} else if dropped := before - len(events); dropped > 0 {
				h.observeDroppedInputN("event", "print_deduplicated", dropped)
			}
		}

		clicks := make([]*domain.AnalyticsClick, 0, len(orderedClicks))
		if !decision.SameSeq {
			for index, item := range orderedClicks {
				section := domain.AnalyticsSectionID(item.Section)
				layoutHash, ok := view.Layout().HashFor(section)
				if !ok {
					h.observeDroppedInput("click", "section_not_in_layout")
					continue
				}
				click, err := domain.NewAnalyticsClick(domain.AnalyticsClickParams{
					ViewID: input.ViewID, LinkID: input.LinkID, Section: section,
					XPct: item.XPct, YPct: item.YPct, ViewportW: item.ViewportW, LayoutHash: layoutHash, Target: item.Target,
					OccurredAt: now.Add(time.Duration(index) * time.Nanosecond),
				})
				if err != nil {
					h.observeDroppedInput("click", "invalid_click")
					continue
				}
				clicks = append(clicks, click)
			}
			before := len(clicks)
			clicks = domain.SelectAnalyticsClicksForStorage(clicks, snapshot.StoredClickCount)
			if dropped := before - len(clicks); dropped > 0 {
				h.observeDroppedInputN("click", "per_view_limit", dropped)
			}
		}
		if decision.SameSeq {
			outcome, reason = "ignored", "same_sequence"
		}
		if err := h.repository.SaveAnalyticsBatch(txCtx, view, events, clicks); err != nil {
			outcome, reason = "failure", "save_failed"
			return err
		}
		return nil
	})
	if err != nil && (outcome == "accepted" || outcome == "ignored") {
		outcome, reason = "failure", "transaction_failed"
	}
	h.observeBatchOutcome(outcome, reason)
	return err
}

func classifyAnalyticsBatchError(err error) string {
	if errors.Is(err, domain.ErrInvalidAnalyticsEvent) || errors.Is(err, domain.ErrInvalidViewState) ||
		errors.Is(err, domain.ErrAnalyticsLinkInactive) || errors.Is(err, domain.ErrAnalyticsSubjectNotFound) ||
		errors.Is(err, domain.ErrAnalyticsViewExpired) {
		return "rejected"
	}
	return "failure"
}

func classifyAnalyticsBatchReason(err error) string {
	if errors.Is(err, domain.ErrInvalidAnalyticsEvent) {
		return "invalid_batch"
	}
	if errors.Is(err, domain.ErrInvalidViewState) {
		return "invalid_state"
	}
	if errors.Is(err, domain.ErrAnalyticsLinkInactive) || errors.Is(err, domain.ErrAnalyticsSubjectNotFound) || errors.Is(err, domain.ErrAnalyticsViewExpired) {
		return "view_rejected"
	}
	return "load_failed"
}

func (h *IngestAnalyticsBatchHandler) RecordDroppedTransportInput(kind, reason string) {
	h.observeDroppedInput(kind, reason)
}

func (h *IngestAnalyticsBatchHandler) RecordRejectedTransportBatch(reason string) {
	h.observeBatchOutcome("rejected", reason)
}

func (h *IngestAnalyticsBatchHandler) observeDroppedInput(kind, reason string) {
	h.observeDroppedInputN(kind, reason, 1)
}

func (h *IngestAnalyticsBatchHandler) observeDroppedInputN(kind, reason string, count int) {
	if h == nil || h.observer == nil || count <= 0 {
		return
	}
	for range count {
		h.observer.RecordDroppedAnalyticsInput(kind, reason)
	}
}

func (h *IngestAnalyticsBatchHandler) observeBatchOutcome(outcome, reason string) {
	if h == nil || h.observer == nil {
		return
	}
	h.observer.RecordAnalyticsBatchOutcome(outcome, reason)
}
