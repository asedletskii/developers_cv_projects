package commands

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/analytics/application/ports"
	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type unitOfWorkStub struct {
	called bool
	calls  int
	err    error
}

func (u *unitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	u.called = true
	u.calls++
	if u.err != nil {
		return u.err
	}
	return fn(ctx)
}

type viewRepositoryStub struct {
	created *domain.View
	view    *domain.View
	getErr  error
}

func (r *viewRepositoryStub) CreateView(_ context.Context, view *domain.View) error {
	r.created = view
	return nil
}

func (r *viewRepositoryStub) GetView(context.Context, sharedvalueobjects.ID) (*domain.View, error) {
	return r.view, r.getErr
}

func (r *viewRepositoryStub) UpdateViewState(context.Context, sharedvalueobjects.ID, domain.ViewStateUpdate) error {
	return nil
}

type analyticsEventRepositoryStub struct {
	events []*domain.AnalyticsEvent
	err    error
}

func (r *analyticsEventRepositoryStub) RecordEvents(_ context.Context, events []*domain.AnalyticsEvent) error {
	r.events = append(r.events, events...)
	return r.err
}

type analyticsRetentionRepositoryStub struct {
	clicks   int
	views    int
	clickErr error
	viewErr  error
	cutoffs  []time.Time
}

func (r *analyticsRetentionRepositoryStub) DeleteExpiredClicks(_ context.Context, cutoff time.Time, _ int) (int, error) {
	r.cutoffs = append(r.cutoffs, cutoff)
	return r.clicks, r.clickErr
}

func (r *analyticsRetentionRepositoryStub) DeleteExpiredViews(_ context.Context, cutoff time.Time, _ int) (int, error) {
	r.cutoffs = append(r.cutoffs, cutoff)
	return r.views, r.viewErr
}

type analyticsBatchRepositoryStub struct {
	called       bool
	err          error
	events       int
	received     []*domain.AnalyticsEvent
	clicks       []*domain.AnalyticsClick
	technologies []string
	lastSeq      int64
}

func (r *analyticsBatchRepositoryStub) LoadAnalyticsBatchForUpdate(
	_ context.Context,
	viewID, linkID sharedvalueobjects.ID,
	_, _, _ bool,
) (ports.AnalyticsBatchSnapshot, error) {
	layout, err := domain.NewLayoutSnapshot(
		[]domain.AnalyticsSectionID{domain.AnalyticsSectionAbout},
		map[domain.AnalyticsSectionID]string{domain.AnalyticsSectionAbout: "0123456789abcdef"},
		r.technologies,
	)
	if err != nil {
		return ports.AnalyticsBatchSnapshot{}, err
	}
	now := time.Now().UTC()
	view, err := domain.RestoreView(domain.ViewSnapshot{
		ID: viewID, LinkID: linkID, ResumeID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
		Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux,
		Traffic: domain.TrafficHuman, Layout: layout, StartedAt: now, CreatedAt: now,
		State: domain.ViewStateSnapshot{LastSeq: r.lastSeq},
	})
	if err != nil {
		return ports.AnalyticsBatchSnapshot{}, err
	}
	return ports.AnalyticsBatchSnapshot{
		View: view,
		LinkContext: ports.LinkAnalyticsContext{
			LinkID: linkID, ResumeID: sharedvalueobjects.NewID(), LinkActive: true, ResumePublished: true,
		},
	}, nil
}

func (r *analyticsBatchRepositoryStub) SaveAnalyticsBatch(
	_ context.Context,
	_ *domain.View,
	events []*domain.AnalyticsEvent,
	clicks []*domain.AnalyticsClick,
) error {
	r.called = true
	r.events = len(events)
	r.received = append([]*domain.AnalyticsEvent(nil), events...)
	r.clicks = append([]*domain.AnalyticsClick(nil), clicks...)
	return r.err
}

type analyticsInputObserverStub struct {
	dropped  map[string]int
	outcomes map[string]int
}

func (s *analyticsInputObserverStub) RecordAnalyticsBatchOutcome(outcome, reason string) {
	if s.outcomes == nil {
		s.outcomes = make(map[string]int)
	}
	s.outcomes[outcome+":"+reason]++
}

func (s *analyticsInputObserverStub) RecordDroppedAnalyticsInput(kind, reason string) {
	if s.dropped == nil {
		s.dropped = make(map[string]int)
	}
	s.dropped[kind+":"+reason]++
}

type analyticsSubjectReaderStub struct {
	lifecycle *analyticscontracts.LinkLifecycle
	err       error
	refs      []analyticscontracts.ResumeLinkRef
	refsErr   error
}

func (s *analyticsSubjectReaderStub) LoadLinkLifecycle(_ context.Context, linkID sharedvalueobjects.ID) (analyticscontracts.LinkLifecycle, error) {
	if s.err != nil {
		return analyticscontracts.LinkLifecycle{}, s.err
	}
	if s.lifecycle != nil {
		return *s.lifecycle, nil
	}
	return analyticscontracts.LinkLifecycle{
		LinkID: linkID, ResumeID: sharedvalueobjects.NewID(),
		LinkActive: true, ResumePublished: true,
	}, nil
}

func (s *analyticsSubjectReaderStub) ListResumeLinks(context.Context, sharedvalueobjects.ID) ([]analyticscontracts.ResumeLinkRef, error) {
	return s.refs, s.refsErr
}

type visitorIdentityStub struct{}

func (visitorIdentityStub) Identify(_ time.Time, _ analyticscontracts.PublicVisitInput) (analyticscontracts.VisitorIdentityResult, error) {
	return analyticscontracts.VisitorIdentityResult{VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32)}, nil
}

type trafficClassifierStub struct{}

func (trafficClassifierStub) Classify(context.Context, ports.TrafficClassificationInput) (ports.TrafficClassificationResult, error) {
	return ports.TrafficClassificationResult{
		Traffic: domain.TrafficHuman,
		Device:  domain.DeviceDesktop,
		Browser: domain.BrowserChrome,
		OS:      domain.OSMacOS,
	}, nil
}

type layoutSnapshotBuilderStub struct{}

func (layoutSnapshotBuilderStub) Build(analyticscontracts.ResumeAnalyticsLayout) (domain.LayoutSnapshot, error) {
	return domain.NewLayoutSnapshot(
		[]domain.AnalyticsSectionID{domain.AnalyticsSectionAbout},
		map[domain.AnalyticsSectionID]string{domain.AnalyticsSectionAbout: "0123456789abcdef"},
	)
}

func newCommandTestView(t *testing.T, viewID, linkID sharedvalueobjects.ID) *domain.View {
	t.Helper()
	layout, err := domain.NewLayoutSnapshot(
		[]domain.AnalyticsSectionID{domain.AnalyticsSectionAbout},
		map[domain.AnalyticsSectionID]string{domain.AnalyticsSectionAbout: "0123456789abcdef"},
	)
	require.NoError(t, err)
	now := time.Now().UTC()
	view, err := domain.RestoreView(domain.ViewSnapshot{
		ID: viewID, LinkID: linkID, ResumeID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
		Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux,
		Traffic: domain.TrafficHuman, Layout: layout, StartedAt: now, CreatedAt: now,
	})
	require.NoError(t, err)
	return view
}
