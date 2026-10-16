package commands

import (
	"context"
	"fmt"
	"time"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CreatePublicVisitInput struct {
	LinkID   sharedvalueobjects.ID
	ResumeID sharedvalueobjects.ID
	Resume   resumequeries.PublicResumeView
	Visitor  ports.PublicVisitInput
}

type CreatePublicVisitHandler struct {
	views      ports.ViewRepository
	identity   ports.VisitorIdentity
	classifier ports.TrafficClassifier
	layout     ports.LayoutSnapshotBuilder
	clock      func() time.Time
}

func NewCreatePublicVisitHandler(
	views ports.ViewRepository,
	identity ports.VisitorIdentity,
	classifier ports.TrafficClassifier,
	layout ports.LayoutSnapshotBuilder,
) *CreatePublicVisitHandler {
	return &CreatePublicVisitHandler{views: views, identity: identity, classifier: classifier, layout: layout, clock: func() time.Time { return time.Now().UTC() }}
}

func (h *CreatePublicVisitHandler) Handle(ctx context.Context, input CreatePublicVisitInput) (*domain.View, ports.VisitorIdentityResult, error) {
	if h == nil || h.views == nil || h.identity == nil || h.classifier == nil || h.layout == nil {
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
	layout, err := h.layout.Build(input.Resume)
	if err != nil {
		return nil, ports.VisitorIdentityResult{}, fmt.Errorf("create public visit: build layout: %w", err)
	}
	view, err := domain.NewView(domain.ViewParams{
		ID: sharedvalueobjects.NewID(), LinkID: input.LinkID, VisitorID: identity.VisitorID, IPHash: identity.IPHash,
		Device: classification.Device, Browser: classification.Browser, OS: classification.OS,
		ReferrerHost: input.Visitor.Referrer, Traffic: classification.Traffic, TrafficReason: classification.Reason,
		Layout: layout, StartedAt: now, CreatedAt: now,
	})
	if err != nil {
		return nil, ports.VisitorIdentityResult{}, fmt.Errorf("create public visit: build view: %w", err)
	}
	if err := h.views.CreateView(ctx, view); err != nil {
		return nil, ports.VisitorIdentityResult{}, fmt.Errorf("create public visit: persist view: %w", err)
	}
	identity.ViewID = view.ID()
	return view, identity, nil
}
