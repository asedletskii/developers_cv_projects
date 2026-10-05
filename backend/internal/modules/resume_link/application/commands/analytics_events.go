package commands

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const (
	MaxAnalyticsEventsPerBatch = 100
	MaxAnalyticsClicksPerBatch = 200
	MaxAnalyticsEventsPerView  = 1000
)

type AnalyticsEventInput struct {
	ClientEventID uuid.UUID
	Type          domain.AnalyticsEventType
	Payload       map[string]string
	At            int64
}

type AnalyticsClickInput struct {
	Section string
	XPct    int
	YPct    int
	VW      int
	Target  string
	At      int64
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
	views      ports.ViewRepository
	uow        ports.UnitOfWork
	clock      func() time.Time
}

func NewIngestAnalyticsBatchHandler(repository ports.AnalyticsBatchRepository, views ports.ViewRepository, uow ports.UnitOfWork) *IngestAnalyticsBatchHandler {
	return &IngestAnalyticsBatchHandler{
		repository: repository,
		views:      views,
		uow:        uow,
		clock:      func() time.Time { return time.Now().UTC() },
	}
}

func (h *IngestAnalyticsBatchHandler) Handle(ctx context.Context, input AnalyticsBatchInput) error {
	if h == nil || h.repository == nil || h.views == nil || h.uow == nil {
		return fmt.Errorf("ingest analytics batch: dependencies are nil")
	}
	if input.ViewID.IsZero() || input.LinkID.IsZero() || input.Seq <= 0 || input.ActiveMS < 0 {
		return fmt.Errorf("%w: invalid analytics batch", domain.ErrInvalidAnalyticsEvent)
	}
	if len(input.Events) > MaxAnalyticsEventsPerBatch {
		return fmt.Errorf("%w: too many events", domain.ErrInvalidAnalyticsEvent)
	}
	if len(input.Clicks) > MaxAnalyticsClicksPerBatch {
		return fmt.Errorf("%w: too many clicks", domain.ErrInvalidAnalyticsEvent)
	}

	view, err := h.views.GetView(ctx, input.ViewID)
	if err != nil {
		return err
	}
	events := make([]*domain.AnalyticsEvent, 0, len(input.Events))
	seen := make(map[uuid.UUID]struct{}, len(input.Events))
	orderedInput := append([]AnalyticsEventInput(nil), input.Events...)
	sort.SliceStable(orderedInput, func(i, j int) bool { return orderedInput[i].At < orderedInput[j].At })
	now := h.clock().UTC()
	for index, item := range orderedInput {
		if item.Type == domain.EventTypeContactView || item.ClientEventID == uuid.Nil {
			continue
		}
		if _, exists := seen[item.ClientEventID]; exists {
			continue
		}
		seen[item.ClientEventID] = struct{}{}

		clientID, err := sharedvalueobjects.IDFromUUID(item.ClientEventID)
		if err != nil {
			continue
		}
		event, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
			ID: sharedvalueobjects.NewID(), ViewID: input.ViewID, LinkID: input.LinkID,
			ClientEventID: &clientID, Type: item.Type, Payload: item.Payload, OccurredAt: now.Add(time.Duration(index) * time.Nanosecond),
		})
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	clicks := make([]*domain.AnalyticsClick, 0, len(input.Clicks))
	orderedClicks := append([]AnalyticsClickInput(nil), input.Clicks...)
	sort.SliceStable(orderedClicks, func(i, j int) bool { return orderedClicks[i].At < orderedClicks[j].At })
	for index, item := range orderedClicks {
		section := resumedomain.SectionID(item.Section)
		layoutHash, ok := view.Layout().HashFor(section)
		if !ok {
			continue
		}
		click, err := domain.NewAnalyticsClick(domain.AnalyticsClickParams{
			ID: sharedvalueobjects.NewID(), ViewID: input.ViewID, LinkID: input.LinkID, Section: section,
			XPct: item.XPct, YPct: item.YPct, ViewportW: item.VW, LayoutHash: layoutHash, Target: item.Target,
			OccurredAt: now.Add(time.Duration(index) * time.Nanosecond),
		})
		if err != nil {
			continue
		}
		clicks = append(clicks, click)
	}

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		return h.repository.ApplyBatch(txCtx, input.ViewID, input.LinkID, domain.ViewStateUpdate{
			Seq: input.Seq, ActiveMS: input.ActiveMS,
			MaxSection: resumedomain.SectionID(input.MaxSection), ExitSection: resumedomain.SectionID(input.ExitSection),
			SeenAt: now, Final: input.Final,
		}, events, clicks)
	})
}
