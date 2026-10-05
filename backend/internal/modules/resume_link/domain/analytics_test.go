package domain

import (
	"testing"
	"time"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"

	"github.com/stretchr/testify/require"
)

func testLayout(t *testing.T) LayoutSnapshot {
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

func testView(t *testing.T) *View {
	t.Helper()
	now := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	view, err := NewView(ViewParams{
		ID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(),
		IPHash: make([]byte, 32), Device: DeviceDesktop, Browser: BrowserChrome, OS: OSLinux,
		Traffic: TrafficHuman, Layout: testLayout(t), StartedAt: now, CreatedAt: now,
	})
	require.NoError(t, err)
	return view
}

func TestViewStateApplyIsMonotonic(t *testing.T) {
	view := testView(t)
	first := time.Date(2026, 10, 5, 10, 0, 5, 0, time.UTC)
	require.NoError(t, view.ApplyState(ViewStateUpdate{Seq: 1, ActiveMS: 7000, MaxSection: resumedomain.SectionExperience, ExitSection: resumedomain.SectionExperience, SeenAt: first}))
	require.Equal(t, int64(7000), view.State().ActiveMS())
	require.Equal(t, resumedomain.SectionExperience, view.State().MaxSection())
	require.Equal(t, first, view.State().ConfirmedAt())

	second := first.Add(time.Second)
	require.NoError(t, view.ApplyState(ViewStateUpdate{Seq: 2, ActiveMS: 3000, MaxSection: resumedomain.SectionAbout, ExitSection: resumedomain.SectionProjects, SeenAt: second, Final: true}))
	require.Equal(t, int64(7000), view.State().ActiveMS())
	require.Equal(t, resumedomain.SectionExperience, view.State().MaxSection())
	require.Equal(t, resumedomain.SectionProjects, view.State().ExitSection())
	require.Equal(t, second, view.State().EndedAt())

	require.NoError(t, view.ApplyState(ViewStateUpdate{Seq: 3, ActiveMS: 8000, MaxSection: resumedomain.SectionProjects, SeenAt: second.Add(time.Second)}))
	require.True(t, view.State().EndedAt().IsZero())
	require.Equal(t, int64(8000), view.State().ActiveMS())

	stateBeforeStale := view.State()
	require.NoError(t, view.ApplyState(ViewStateUpdate{Seq: 2, ActiveMS: 1, MaxSection: resumedomain.SectionAbout, ExitSection: resumedomain.SectionAbout, SeenAt: second.Add(2 * time.Second), Final: true}))
	require.Equal(t, stateBeforeStale.ActiveMS(), view.State().ActiveMS())
	require.Equal(t, stateBeforeStale.ExitSection(), view.State().ExitSection())
	require.True(t, view.State().EndedAt().IsZero())
}

func TestNewAnalyticsEventValidatesAndNormalizesPayload(t *testing.T) {
	clientEventID := sharedvalueobjects.NewID()
	event, err := NewAnalyticsEvent(AnalyticsEventParams{
		ID: sharedvalueobjects.NewID(), ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
		ClientEventID: &clientEventID, Type: EventTypeContactView,
		Payload: map[string]string{"contact": "email", "ignored": "value", "long": "ignored"}, OccurredAt: time.Now(),
	})
	require.NoError(t, err)
	require.Equal(t, map[string]string{"contact": "email"}, event.Payload())

	_, err = NewAnalyticsEvent(AnalyticsEventParams{ID: sharedvalueobjects.NewID(), ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Type: AnalyticsEventType("view"), OccurredAt: time.Now()})
	require.ErrorIs(t, err, ErrInvalidAnalyticsEvent)
}

func TestLayoutSnapshotRejectsTechStackAndInvalidHash(t *testing.T) {
	_, err := NewLayoutSnapshot([]resumedomain.SectionID{resumedomain.SectionTechStack}, map[resumedomain.SectionID]string{resumedomain.SectionTechStack: "0123456789abcdef"})
	require.ErrorIs(t, err, ErrInvalidLayoutSnapshot)
	_, err = NewLayoutSnapshot([]resumedomain.SectionID{resumedomain.SectionAbout}, map[resumedomain.SectionID]string{resumedomain.SectionAbout: "invalid"})
	require.ErrorIs(t, err, ErrInvalidLayoutSnapshot)
}

func TestAnalyticsClickValidation(t *testing.T) {
	now := time.Now().UTC()
	click, err := NewAnalyticsClick(AnalyticsClickParams{
		ID: sharedvalueobjects.NewID(), ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
		Section: resumedomain.SectionAbout, XPct: 4210, YPct: 1733, ViewportW: 1440, LayoutHash: "0123456789abcdef",
		Target: "download", OccurredAt: now,
	})
	require.NoError(t, err)
	require.Equal(t, 4210, click.XPct())
	require.Equal(t, "0123456789abcdef", click.LayoutHash())

	_, err = NewAnalyticsClick(AnalyticsClickParams{
		ID: sharedvalueobjects.NewID(), ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
		Section: resumedomain.SectionAbout, XPct: 10001, YPct: 0, ViewportW: 1440, LayoutHash: "0123456789abcdef", OccurredAt: now,
	})
	require.ErrorIs(t, err, ErrInvalidAnalyticsClick)
}
