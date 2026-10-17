package domain

import (
	"strings"
	"testing"
	"time"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"

	"github.com/stretchr/testify/require"
)

func testViewLayout(t *testing.T) LayoutSnapshot {
	t.Helper()
	layout, err := NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionAbout, resumedomain.SectionExperience, resumedomain.SectionProjects},
		map[resumedomain.SectionID]string{
			resumedomain.SectionAbout:      "0123456789abcdef",
			resumedomain.SectionExperience: "fedcba9876543210",
			resumedomain.SectionProjects:   "0011223344556677",
		},
	)
	require.NoError(t, err)
	return layout
}

func newTestView(t *testing.T) *View {
	t.Helper()
	now := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	view, err := NewView(ViewParams{
		LinkID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(),
		IPHash: make([]byte, 32), Device: DeviceDesktop, Browser: BrowserChrome, OS: OSLinux,
		Traffic: TrafficHuman, Layout: testViewLayout(t), StartedAt: now, CreatedAt: now,
	})
	require.NoError(t, err)
	return view
}

func TestNewViewCopiesIPHashAndDefaultsCreatedAt(t *testing.T) {
	now := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	ipHash := make([]byte, 32)
	view, err := NewView(ViewParams{
		LinkID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(),
		IPHash: ipHash, Device: DeviceDesktop, Browser: BrowserChrome, OS: OSLinux,
		Traffic: TrafficHuman, Layout: testViewLayout(t), StartedAt: now,
	})
	require.NoError(t, err)
	require.False(t, view.ID().IsZero())
	require.False(t, view.CreatedAt().IsZero())
	require.Equal(t, view.StartedAt(), view.CreatedAt())

	ipHash[0] = 1
	require.Equal(t, byte(0), view.IPHash()[0])
}

func TestRestoreViewPreservesStoredStrings(t *testing.T) {
	now := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	view, err := RestoreView(ViewSnapshot{
		ID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(),
		IPHash: make([]byte, 32), Device: DeviceDesktop, Browser: BrowserChrome, OS: OSLinux,
		ReferrerHost: " example.com ", Traffic: TrafficHuman, TrafficReason: " retained ",
		Layout: testViewLayout(t), StartedAt: now, CreatedAt: now,
	})

	require.NoError(t, err)
	require.Equal(t, " example.com ", view.ReferrerHost())
	require.Equal(t, " retained ", view.TrafficReason())
}

func TestRestoreViewRejectsMissingCreationTimestamp(t *testing.T) {
	_, err := RestoreView(ViewSnapshot{
		ID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(),
		IPHash: make([]byte, 32), Device: DeviceDesktop, Browser: BrowserChrome, OS: OSLinux,
		Traffic: TrafficHuman, Layout: testViewLayout(t),
		StartedAt: time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC),
	})

	require.ErrorIs(t, err, ErrInvalidView)
}

func TestReconstructViewRestoresSnapshot(t *testing.T) {
	now := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	view, err := ReconstructView(ViewSnapshot{
		ID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(),
		IPHash: make([]byte, 32), Device: DeviceMobile, Browser: BrowserSafari, OS: OSIOS,
		ReferrerHost: "example.com", Traffic: TrafficHuman, TrafficReason: "classified",
		Layout: testViewLayout(t), StartedAt: now, CreatedAt: now,
	})
	require.NoError(t, err)
	require.Equal(t, "example.com", view.ReferrerHost())
	require.Equal(t, DeviceMobile, view.Device())
	require.Equal(t, BrowserSafari, view.Browser())
	require.Equal(t, OSIOS, view.OS())
	require.Equal(t, TrafficHuman, view.Traffic())
	require.Equal(t, now, view.StartedAt())
}

func TestReconstructViewRestoresState(t *testing.T) {
	view := newTestView(t)
	state := ViewState{}
	require.NoError(t, state.Apply(ViewStateUpdate{
		Seq: 1, ActiveMS: 7000, MaxSection: resumedomain.SectionExperience,
		SeenAt: time.Date(2026, 10, 5, 10, 0, 5, 0, time.UTC),
	}, view.Layout()))

	reconstructed, err := ReconstructView(ViewSnapshot{
		ID: view.ID(), LinkID: view.LinkID(), VisitorID: view.VisitorID(),
		IPHash: view.IPHash(), Device: view.Device(), Browser: view.Browser(), OS: view.OS(),
		ReferrerHost: view.ReferrerHost(), Traffic: view.Traffic(), TrafficReason: view.TrafficReason(),
		Layout: view.Layout(), StartedAt: view.StartedAt(), CreatedAt: view.CreatedAt(), State: state,
	})
	require.NoError(t, err)
	require.Equal(t, int64(7000), reconstructed.State().ActiveMS())
	require.Equal(t, resumedomain.SectionExperience, reconstructed.State().MaxSection())
}

func TestNewViewRejectsOverlongReferrerHost(t *testing.T) {
	now := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	_, err := NewView(ViewParams{
		LinkID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(),
		IPHash: make([]byte, 32), Device: DeviceDesktop, Browser: BrowserChrome, OS: OSLinux,
		ReferrerHost: strings.Repeat("a", MaxReferrerHostLength+1),
		Traffic:      TrafficHuman, Layout: testViewLayout(t), StartedAt: now,
	})
	require.ErrorIs(t, err, ErrInvalidView)
}

func TestViewApplyStateRejectsActiveTimeAboveStorageLimit(t *testing.T) {
	view := newTestView(t)
	err := view.ApplyState(ViewStateUpdate{
		Seq: 1, ActiveMS: MaxActiveMS + 1,
		SeenAt: view.StartedAt().Add(13 * time.Hour),
	})
	require.ErrorIs(t, err, ErrInvalidViewState)
}

func TestViewRestoreStateRejectsActiveTimeAboveStorageLimit(t *testing.T) {
	view := newTestView(t)
	err := view.RestoreState(ViewState{activeMS: MaxActiveMS + 1})
	require.ErrorIs(t, err, ErrInvalidViewState)
}

func TestViewRestoreStateRejectsSectionOutsideLayout(t *testing.T) {
	view := newTestView(t)
	err := view.RestoreState(ViewState{maxSection: resumedomain.SectionContacts})
	require.ErrorIs(t, err, ErrInvalidViewState)
}

func TestViewRestoreStateRejectsNegativeSequenceAtomically(t *testing.T) {
	view := newTestView(t)
	before := view.State()
	err := view.RestoreState(ViewState{
		lastSeq: -1,
	})

	require.ErrorIs(t, err, ErrInvalidViewState)
	require.Equal(t, before, view.State())
}

func TestViewApplyStateRejectsActiveTimeBeyondElapsedTime(t *testing.T) {
	view := newTestView(t)
	seenAt := view.StartedAt().Add(10 * time.Second)

	err := view.ApplyState(ViewStateUpdate{
		Seq: 1, ActiveMS: 60_001, SeenAt: seenAt,
	})

	require.ErrorIs(t, err, ErrInvalidViewState)
	require.Equal(t, int64(0), view.State().ActiveMS())
}

func TestViewApplyStateAllowsSmallClockSkew(t *testing.T) {
	view := newTestView(t)
	seenAt := view.StartedAt().Add(10 * time.Second)

	err := view.ApplyState(ViewStateUpdate{
		Seq: 1, ActiveMS: 40_000, SeenAt: seenAt,
	})

	require.NoError(t, err)
	require.Equal(t, int64(40_000), view.State().ActiveMS())
}

func TestViewEvaluateAnalyticsBatchAppliesLifecycleAndSequenceRules(t *testing.T) {
	view := newTestView(t)
	now := view.StartedAt().Add(time.Minute)

	decision, err := view.EvaluateAnalyticsBatch(
		ResumeLinkStatusActive, resumedomain.ResumeStatusOpened, now,
		ViewStateUpdate{Seq: 1, SeenAt: now},
	)
	require.NoError(t, err)
	require.False(t, decision.Ignore)
	require.False(t, decision.SameSeq)

	decision, err = view.EvaluateAnalyticsBatch(
		ResumeLinkStatusInactive, resumedomain.ResumeStatusOpened, now,
		ViewStateUpdate{Seq: 1, SeenAt: now},
	)
	require.ErrorIs(t, err, ErrResumeLinkInactive)

	decision, err = view.EvaluateAnalyticsBatch(
		ResumeLinkStatusActive, resumedomain.ResumeStatusDraft, now,
		ViewStateUpdate{Seq: 1, SeenAt: now},
	)
	require.ErrorIs(t, err, ErrEntityNotFound)

	decision, err = view.EvaluateAnalyticsBatch(
		ResumeLinkStatusActive, resumedomain.ResumeStatusOpened, view.StartedAt().Add(analyticsViewTTL+time.Second),
		ViewStateUpdate{Seq: 1, SeenAt: now},
	)
	require.ErrorIs(t, err, ErrAnalyticsViewExpired)
}

func TestViewEvaluateAnalyticsBatchIgnoresBotsAndStaleSequences(t *testing.T) {
	view := newTestView(t)
	now := view.StartedAt().Add(time.Minute)
	bot, err := NewView(ViewParams{
		LinkID: view.LinkID(), VisitorID: sharedvalueobjects.NewID(),
		IPHash: make([]byte, 32), Device: DeviceDesktop, Browser: BrowserChrome, OS: OSLinux,
		Traffic: TrafficBot, Layout: view.Layout(), StartedAt: view.StartedAt(), CreatedAt: view.CreatedAt(),
	})
	require.NoError(t, err)
	decision, err := bot.EvaluateAnalyticsBatch(
		ResumeLinkStatusActive, resumedomain.ResumeStatusOpened, now, ViewStateUpdate{Seq: 1, SeenAt: now},
	)
	require.NoError(t, err)
	require.True(t, decision.Ignore)

	state := ViewState{}
	require.NoError(t, state.Apply(ViewStateUpdate{Seq: 2, ActiveMS: 1000, SeenAt: now}, view.Layout()))
	require.NoError(t, view.RestoreState(state))
	decision, err = view.EvaluateAnalyticsBatch(
		ResumeLinkStatusActive, resumedomain.ResumeStatusOpened, now, ViewStateUpdate{Seq: 1, SeenAt: now},
	)
	require.NoError(t, err)
	require.True(t, decision.Ignore)
	decision, err = view.EvaluateAnalyticsBatch(
		ResumeLinkStatusActive, resumedomain.ResumeStatusOpened, now, ViewStateUpdate{Seq: 2, Final: true, SeenAt: now},
	)
	require.NoError(t, err)
	require.False(t, decision.Ignore)
	require.True(t, decision.SameSeq)
}
