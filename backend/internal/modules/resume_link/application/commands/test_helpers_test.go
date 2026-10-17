package commands

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
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

type resumeLinkRepositoryStub struct {
	link         *domain.ResumeLink
	err          error
	getErr       error
	count        int64
	countErr     error
	createErr    error
	deleteErr    error
	delErr       error
	replaceErr   error
	returnNil    bool
	lockCalled   bool
	countCalled  bool
	called       bool
	createCalled bool
	deleted      sharedvalueobjects.ID
	replaced     bool
	active       bool
	version      int64
}

func validCreateCommand() CreateResumeLinkCommand {
	return CreateResumeLinkCommand{UserID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(), Name: "  HH  "}
}

func (r *resumeLinkRepositoryStub) LockResumeForLinkCreation(context.Context, sharedvalueobjects.ID) error {
	r.lockCalled = true
	return nil
}
func (r *resumeLinkRepositoryStub) CountByResumeID(context.Context, sharedvalueobjects.ID) (int64, error) {
	r.countCalled = true
	return r.count, r.countErr
}
func (r *resumeLinkRepositoryStub) Create(_ context.Context, link *domain.ResumeLink) error {
	r.called = true
	r.createCalled = true
	if r.createErr != nil {
		return r.createErr
	}
	if r.err != nil {
		return r.err
	}
	r.link = link
	return nil
}
func (r *resumeLinkRepositoryStub) GetByID(context.Context, sharedvalueobjects.ID) (*domain.ResumeLink, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	if r.err != nil {
		return nil, r.err
	}
	if r.returnNil {
		return nil, nil
	}
	return r.link, nil
}
func (r *resumeLinkRepositoryStub) Delete(_ context.Context, id sharedvalueobjects.ID) error {
	r.deleted = id
	if r.deleteErr != nil {
		return r.deleteErr
	}
	if r.delErr != nil {
		return r.delErr
	}
	return r.err
}
func (r *resumeLinkRepositoryStub) ReplaceStatus(_ context.Context, _ sharedvalueobjects.ID, active bool, version int64) error {
	r.replaced, r.active, r.version = true, active, version
	if r.replaceErr != nil {
		return r.replaceErr
	}
	return r.err
}

type resumeOwnershipReaderStub struct {
	err   error
	query resumequeries.GetResumeQuery
}

func (r *resumeOwnershipReaderStub) GetOwned(_ context.Context, query resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	r.query = query
	if r.err != nil {
		return resumequeries.ResumeView{}, r.err
	}
	return resumequeries.ResumeView{ID: query.ResumeID.String()}, nil
}
func (r *resumeOwnershipReaderStub) OwnsResume(_ context.Context, query resumequeries.GetResumeQuery) (bool, error) {
	r.query = query
	return r.err == nil, r.err
}

type resumeReaderStub struct {
	view   resumequeries.ResumeView
	err    error
	called bool
	query  resumequeries.GetResumeQuery
}

func (r *resumeReaderStub) GetOwned(_ context.Context, query resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	r.called = true
	r.query = query
	if r.err != nil {
		return resumequeries.ResumeView{}, r.err
	}
	return r.view, nil
}
func (r *resumeReaderStub) OwnsResume(_ context.Context, query resumequeries.GetResumeQuery) (bool, error) {
	r.called = true
	r.query = query
	if r.err != nil {
		return false, r.err
	}
	return true, nil
}

type analyticsBatchRepositoryStub struct {
	called       bool
	err          error
	events       int
	received     []*domain.AnalyticsEvent
	technologies []string
}

func (r *analyticsBatchRepositoryStub) LoadAnalyticsBatchForUpdate(_ context.Context, viewID, linkID sharedvalueobjects.ID, _, _, _ bool) (domain.AnalyticsBatchSnapshot, error) {
	layout, err := domain.NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionAbout},
		map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"},
		r.technologies,
	)
	if err != nil {
		return domain.AnalyticsBatchSnapshot{}, err
	}
	now := time.Now().UTC()
	view, err := domain.RestoreView(domain.ViewSnapshot{
		ID: viewID, LinkID: linkID, VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
		Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux,
		Traffic: domain.TrafficHuman, Layout: layout, StartedAt: now, CreatedAt: now,
	})
	if err != nil {
		return domain.AnalyticsBatchSnapshot{}, err
	}
	return domain.AnalyticsBatchSnapshot{View: view, LinkStatus: domain.ResumeLinkStatusActive, ResumeStatus: resumedomain.ResumeStatusOpened}, nil
}
func (r *analyticsBatchRepositoryStub) SaveAnalyticsBatch(_ context.Context, _ *domain.View, events []*domain.AnalyticsEvent, _ []*domain.AnalyticsClick) error {
	r.called = true
	r.events = len(events)
	r.received = append([]*domain.AnalyticsEvent(nil), events...)
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

type visitorIdentityStub struct{}

func (visitorIdentityStub) Identify(_ time.Time, _ ports.PublicVisitInput) (ports.VisitorIdentityResult, error) {
	return ports.VisitorIdentityResult{VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32)}, nil
}

type trafficClassifierStub struct{}

func (trafficClassifierStub) Classify(context.Context, ports.TrafficClassificationInput) (ports.TrafficClassificationResult, error) {
	return ports.TrafficClassificationResult{Traffic: domain.TrafficHuman, Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSMacOS}, nil
}

type layoutSnapshotBuilderStub struct{}

func (layoutSnapshotBuilderStub) Build(resumequeries.PublicResumeView) (domain.LayoutSnapshot, error) {
	return domain.NewLayoutSnapshot([]resumedomain.SectionID{resumedomain.SectionAbout}, map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"})
}

func newCommandTestView(t *testing.T, viewID, linkID sharedvalueobjects.ID) *domain.View {
	t.Helper()
	layout, err := domain.NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionAbout},
		map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"},
	)
	require.NoError(t, err)
	now := time.Now().UTC()
	view, err := domain.RestoreView(domain.ViewSnapshot{
		ID: viewID, LinkID: linkID, VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
		Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux,
		Traffic: domain.TrafficHuman, Layout: layout, StartedAt: now, CreatedAt: now,
	})
	require.NoError(t, err)
	return view
}

func restoreCommandTestLink(t *testing.T, status domain.ResumeLinkStatus, version int64) *domain.ResumeLink {
	t.Helper()
	link, err := domain.RestoreResumeLink(domain.ResumeLinkSnapshot{
		ID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(),
		Name: "HH", Status: status, Version: version,
	})
	require.NoError(t, err)
	return link
}
