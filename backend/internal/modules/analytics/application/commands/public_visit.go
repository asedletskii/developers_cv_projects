package commands

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/analytics/application/ports"
	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	"developers-cv/backend/internal/modules/analytics/domain"
)

type CreatePublicVisitHandler struct {
	views      ports.ViewRepository
	uow        ports.UnitOfWork
	identity   ports.VisitorIdentity
	classifier ports.TrafficClassifier
	layout     ports.LayoutSnapshotBuilder
	subjects   analyticscontracts.AnalyticsSubjectReader
	clock      func() time.Time
}

func NewCreatePublicVisitHandler(
	views ports.ViewRepository,
	uow ports.UnitOfWork,
	identity ports.VisitorIdentity,
	classifier ports.TrafficClassifier,
	layout ports.LayoutSnapshotBuilder,
	subjects analyticscontracts.AnalyticsSubjectReader,
) *CreatePublicVisitHandler {
	return &CreatePublicVisitHandler{
		views: views, uow: uow, identity: identity, classifier: classifier, layout: layout,
		subjects: subjects, clock: func() time.Time { return time.Now().UTC() },
	}
}

func (h *CreatePublicVisitHandler) Handle(ctx context.Context, input analyticscontracts.CreatePublicVisitInput) (analyticscontracts.VisitorIdentityResult, error) {
	if h == nil || h.views == nil || h.uow == nil || h.identity == nil || h.classifier == nil || h.layout == nil || h.subjects == nil {
		return analyticscontracts.VisitorIdentityResult{}, fmt.Errorf("create public visit: dependencies are nil")
	}
	if input.LinkID.IsZero() || input.ResumeID.IsZero() {
		return analyticscontracts.VisitorIdentityResult{}, fmt.Errorf("create public visit: invalid input")
	}
	now := h.clock().UTC()
	identity, err := h.identity.Identify(now, input.Visitor)
	if err != nil {
		return analyticscontracts.VisitorIdentityResult{}, fmt.Errorf("create public visit: identify visitor: %w", err)
	}
	classification, err := h.classifier.Classify(ctx, ports.TrafficClassificationInput{
		LinkID: input.LinkID, OwnerID: input.Visitor.OwnerID, UserAgent: input.Visitor.UserAgent, IPHash: identity.IPHash, Now: now,
	})
	if err != nil {
		return analyticscontracts.VisitorIdentityResult{}, fmt.Errorf("create public visit: classify traffic: %w", err)
	}
	layout, err := h.layout.Build(input.ResumeLayout)
	if err != nil {
		return analyticscontracts.VisitorIdentityResult{}, fmt.Errorf("create public visit: build layout: %w", err)
	}
	view, err := domain.NewView(domain.ViewParams{
		LinkID: input.LinkID, ResumeID: input.ResumeID, VisitorID: identity.VisitorID, IPHash: identity.IPHash,
		Device: classification.Device, Browser: classification.Browser, OS: classification.OS,
		ReferrerHost: input.Visitor.Referrer, Traffic: classification.Traffic, TrafficReason: classification.Reason,
		Layout: layout, StartedAt: now, CreatedAt: now,
	})
	if err != nil {
		return analyticscontracts.VisitorIdentityResult{}, fmt.Errorf("create public visit: build view: %w", err)
	}
	if err := h.uow.Do(ctx, func(txCtx context.Context) error {
		lifecycle, err := h.subjects.LoadLinkLifecycle(txCtx, input.LinkID)
		if err != nil {
			return err
		}
		if !lifecycle.LinkActive {
			return domain.ErrAnalyticsLinkInactive
		}
		if !lifecycle.ResumePublished {
			return analyticscontracts.ErrAnalyticsSubjectNotFound
		}
		if lifecycle.ResumeID != input.ResumeID {
			return fmt.Errorf("create public visit: link resume mismatch")
		}
		return h.views.CreateView(txCtx, view)
	}); err != nil {
		return analyticscontracts.VisitorIdentityResult{}, fmt.Errorf("create public visit: persist view: %w", err)
	}
	identity.ViewID = view.ID()
	return identity, nil
}
