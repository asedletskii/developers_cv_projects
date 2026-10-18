package commands

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/analytics/application/ports"
	"developers-cv/backend/internal/modules/analytics/domain"
)

type CreatePublicVisitInput = ports.CreatePublicVisitInput

type CreatePublicVisitHandler struct {
	views      ports.ViewRepository
	uow        ports.UnitOfWork
	identity   ports.VisitorIdentity
	classifier ports.TrafficClassifier
	layout     ports.LayoutSnapshotBuilder
	clock      func() time.Time
}

func NewCreatePublicVisitHandler(
	views ports.ViewRepository,
	uow ports.UnitOfWork,
	identity ports.VisitorIdentity,
	classifier ports.TrafficClassifier,
	layout ports.LayoutSnapshotBuilder,
) *CreatePublicVisitHandler {
	return &CreatePublicVisitHandler{
		views: views, uow: uow, identity: identity, classifier: classifier, layout: layout,
		clock: func() time.Time { return time.Now().UTC() },
	}
}

func (h *CreatePublicVisitHandler) Handle(ctx context.Context, input CreatePublicVisitInput) (*domain.View, ports.VisitorIdentityResult, error) {
	if h == nil || h.views == nil || h.uow == nil || h.identity == nil || h.classifier == nil || h.layout == nil {
		return nil, ports.VisitorIdentityResult{}, fmt.Errorf("create public visit: dependencies are nil")
	}
	if input.LinkID.IsZero() {
		return nil, ports.VisitorIdentityResult{}, fmt.Errorf("create public visit: invalid input")
	}
	now := h.clock().UTC()
	identity, err := h.identity.Identify(now, input.Visitor)
	if err != nil {
		return nil, ports.VisitorIdentityResult{}, fmt.Errorf("create public visit: identify visitor: %w", err)
	}
	classification, err := h.classifier.Classify(ctx, ports.TrafficClassificationInput{
		LinkID: input.LinkID, OwnerID: input.Visitor.OwnerID, UserAgent: input.Visitor.UserAgent, IPHash: identity.IPHash, Now: now,
	})
	if err != nil {
		return nil, ports.VisitorIdentityResult{}, fmt.Errorf("create public visit: classify traffic: %w", err)
	}
	layout, err := h.layout.Build(input.ResumeLayout)
	if err != nil {
		return nil, ports.VisitorIdentityResult{}, fmt.Errorf("create public visit: build layout: %w", err)
	}
	view, err := domain.NewView(domain.ViewParams{
		LinkID: input.LinkID, VisitorID: identity.VisitorID, IPHash: identity.IPHash,
		Device: classification.Device, Browser: classification.Browser, OS: classification.OS,
		ReferrerHost: input.Visitor.Referrer, Traffic: classification.Traffic, TrafficReason: classification.Reason,
		Layout: layout, StartedAt: now, CreatedAt: now,
	})
	if err != nil {
		return nil, ports.VisitorIdentityResult{}, fmt.Errorf("create public visit: build view: %w", err)
	}
	if err := h.uow.Do(ctx, func(txCtx context.Context) error {
		return h.views.CreateView(txCtx, view)
	}); err != nil {
		return nil, ports.VisitorIdentityResult{}, fmt.Errorf("create public visit: persist view: %w", err)
	}
	identity.ViewID = view.ID()
	return view, identity, nil
}
