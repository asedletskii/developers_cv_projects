//go:build integration

package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func newAnalyticsView(t *testing.T, linkID sharedvalueobjects.ID) *domain.View {
	t.Helper()
	layout, err := domain.NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionAbout, resumedomain.SectionExperience, resumedomain.SectionProjects},
		map[resumedomain.SectionID]string{
			resumedomain.SectionAbout:      "0123456789abcdef",
			resumedomain.SectionExperience: "fedcba9876543210",
			resumedomain.SectionProjects:   "0011223344556677",
		},
	)
	require.NoError(t, err)
	view, err := domain.NewView(domain.ViewParams{
		ID: sharedvalueobjects.NewID(), LinkID: linkID, VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
		Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux, Traffic: domain.TrafficHuman,
		Layout: layout, StartedAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	return view
}

func TestResumeLinkAnalyticsBatchIsMonotonic(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "Analytics Batch"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	view := newAnalyticsView(t, link.ID())
	require.NoError(t, repo.CreateView(ctx, view))

	clientEventID := sharedvalueobjects.NewID()
	layoutHash, ok := view.Layout().HashFor(resumedomain.SectionAbout)
	require.True(t, ok)
	click, err := domain.NewAnalyticsClick(domain.AnalyticsClickParams{
		ID: sharedvalueobjects.NewID(), ViewID: view.ID(), LinkID: link.ID(), Section: resumedomain.SectionAbout,
		XPct: 4210, YPct: 1733, ViewportW: 1440, LayoutHash: layoutHash, Target: "download", OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	event, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
		ID: sharedvalueobjects.NewID(), ViewID: view.ID(), LinkID: link.ID(), ClientEventID: &clientEventID,
		Type: domain.EventTypeDownload, OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.ApplyBatch(txCtx, view.ID(), link.ID(), domain.ViewStateUpdate{
			Seq: 2, ActiveMS: 7000, MaxSection: resumedomain.SectionExperience,
			ExitSection: resumedomain.SectionExperience, SeenAt: time.Now().UTC(), Final: false,
		}, []*domain.AnalyticsEvent{event}, []*domain.AnalyticsClick{click})
	}))

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.ApplyBatch(txCtx, view.ID(), link.ID(), domain.ViewStateUpdate{
			Seq: 1, ActiveMS: 100, MaxSection: resumedomain.SectionAbout,
			ExitSection: resumedomain.SectionAbout, SeenAt: time.Now().UTC(), Final: false,
		}, nil, nil)
	}))

	loaded, err := repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.Equal(t, int64(7000), loaded.State().ActiveMS())
	require.Equal(t, int64(2), loaded.State().LastSeq())
	require.Equal(t, resumedomain.SectionExperience, loaded.State().MaxSection())
	require.True(t, loaded.State().EndedAt().IsZero())

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.ApplyBatch(txCtx, view.ID(), link.ID(), domain.ViewStateUpdate{
			Seq: 2, ActiveMS: 7000, SeenAt: time.Now().UTC(), Final: true,
		}, nil, []*domain.AnalyticsClick{click})
	}))

	loaded, err = repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.False(t, loaded.State().EndedAt().IsZero())

	var count int
	require.NoError(t, gormDB.Raw(
		"SELECT count(*) FROM resume_link_events WHERE view_id = ? AND client_event_id = ?",
		view.ID().UUID(), clientEventID.UUID(),
	).Scan(&count).Error)
	require.Equal(t, 1, count)
	var clickCount int
	require.NoError(t, gormDB.Raw(
		"SELECT count(*) FROM resume_link_clicks WHERE view_id = ?",
		view.ID().UUID(),
	).Scan(&clickCount).Error)
	require.Equal(t, 1, clickCount)
}

func TestResumeLinkAnalyticsResumeStatisticsAggregatesAllLinks(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	uow := database.NewUnitOfWork(gormDB)
	links := make([]*domain.ResumeLink, 0, 2)
	for _, name := range []string{"HH", "Telegram"} {
		link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: name})
		require.NoError(t, err)
		require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
		links = append(links, link)
	}

	first := newAnalyticsView(t, links[0].ID())
	second := newAnalyticsView(t, links[1].ID())
	require.NoError(t, repo.CreateView(ctx, first))
	require.NoError(t, repo.CreateView(ctx, second))
	startedAt := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	require.NoError(t, gormDB.Exec(`UPDATE resume_link_views SET started_at = ? WHERE id = ?`, startedAt, first.ID().UUID()).Error)
	require.NoError(t, gormDB.Exec(`UPDATE resume_link_views SET started_at = ? WHERE id = ?`, startedAt.Add(5*time.Minute), second.ID().UUID()).Error)

	period, err := queries.NewAnalyticsPeriod("2026-10-01", "2026-10-05", "UTC", startedAt)
	require.NoError(t, err)
	result, err := reader.GetResumeStatistics(ctx, queries.ResumeAnalyticsQuery{ResumeID: resumeID, Period: period})
	require.NoError(t, err)
	require.Len(t, result.Links, 2)
	require.Equal(t, links[0].ID().String(), result.Links[0].ID)
	require.Equal(t, links[1].ID().String(), result.Links[1].ID)
	require.Equal(t, int64(1), result.Links[0].Views)
	require.Equal(t, int64(1), result.Links[1].Views)
	require.Equal(t, int64(2), result.Views)
}

func TestResumeLinkAnalyticsVisitsPaginationUsesBothLinkIDBindings(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "Visits"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	first := newAnalyticsView(t, link.ID())
	second := newAnalyticsView(t, link.ID())
	require.NoError(t, repo.CreateView(ctx, first))
	require.NoError(t, repo.CreateView(ctx, second))

	firstStartedAt := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	secondStartedAt := firstStartedAt.Add(5 * time.Minute)
	require.NoError(t, gormDB.Exec(`UPDATE resume_link_views SET started_at = ? WHERE id = ?`, firstStartedAt, first.ID().UUID()).Error)
	require.NoError(t, gormDB.Exec(`UPDATE resume_link_views SET started_at = ? WHERE id = ?`, secondStartedAt, second.ID().UUID()).Error)

	page, err := reader.GetVisits(ctx, queries.LinkVisitsQuery{LinkID: link.ID(), Limit: 1})
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	require.Equal(t, second.ID().String(), page.Items[0].ID)
	require.NotNil(t, page.NextCursor)

	nextPage, err := reader.GetVisits(ctx, queries.LinkVisitsQuery{LinkID: link.ID(), Limit: 1, Cursor: *page.NextCursor})
	require.NoError(t, err)
	require.Len(t, nextPage.Items, 1)
	require.Equal(t, first.ID().String(), nextPage.Items[0].ID)
	require.Nil(t, nextPage.NextCursor)
}

func TestResumeLinkAnalyticsBatchRejectsChangedPublicLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "Lifecycle"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
	view := newAnalyticsView(t, link.ID())
	require.NoError(t, repo.CreateView(ctx, view))
	require.NoError(t, gormDB.Exec(`UPDATE resume_links SET status = 'inactive' WHERE id = ?`, link.ID().UUID()).Error)
	err = uow.Do(ctx, func(txCtx context.Context) error {
		return repo.ApplyBatch(txCtx, view.ID(), link.ID(), domain.ViewStateUpdate{Seq: 1, SeenAt: time.Now().UTC()}, nil, nil)
	})
	require.ErrorIs(t, err, domain.ErrResumeLinkInactive)

	require.NoError(t, gormDB.Exec(`UPDATE resume_links SET status = 'active' WHERE id = ?`, link.ID().UUID()).Error)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'draft' WHERE id = ?`, resumeID.UUID()).Error)
	err = uow.Do(ctx, func(txCtx context.Context) error {
		return repo.ApplyBatch(txCtx, view.ID(), link.ID(), domain.ViewStateUpdate{Seq: 2, SeenAt: time.Now().UTC()}, nil, nil)
	})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}

func TestResumeLinkAnalyticsFoundationPersistence(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "Analytics"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	view := newAnalyticsView(t, link.ID())
	require.NoError(t, repo.CreateView(ctx, view))
	loaded, err := repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.Equal(t, view.ID(), loaded.ID())
	require.Equal(t, view.LinkID(), loaded.LinkID())
	require.Equal(t, view.Layout().Order(), loaded.Layout().Order())
	require.Equal(t, view.CreatedAt(), loaded.CreatedAt())

	now := time.Now().UTC()
	require.NoError(t, repo.UpdateViewState(ctx, view.ID(), domain.ViewStateUpdate{
		Seq: 1, ActiveMS: 7000, MaxSection: resumedomain.SectionExperience, ExitSection: resumedomain.SectionExperience, SeenAt: now, Final: true,
	}))
	loaded, err = repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.Equal(t, int64(7000), loaded.State().ActiveMS())
	require.Equal(t, resumedomain.SectionExperience, loaded.State().MaxSection())
	require.False(t, loaded.State().EndedAt().IsZero())

	clientEventID := sharedvalueobjects.NewID()
	event, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
		ID: sharedvalueobjects.NewID(), ViewID: view.ID(), LinkID: link.ID(), ClientEventID: &clientEventID,
		Type: domain.EventTypeSectionView, Payload: map[string]string{"section": "experience"}, OccurredAt: now,
	})
	require.NoError(t, err)
	require.NoError(t, repo.RecordEvents(ctx, []*domain.AnalyticsEvent{event, event}))

	var count int
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM resume_link_events WHERE view_id = ?", view.ID().UUID()).Scan(&count).Error)
	require.Equal(t, 1, count)

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Delete(txCtx, link.ID()) }))
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM resume_link_views WHERE id = ?", view.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM resume_link_events WHERE view_id = ?", view.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
}

func TestResumeLinkAnalyticsRetentionDeletesExpiredDataInBatches(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ID: sharedvalueobjects.NewID(), ResumeID: resumeID, Name: "Analytics Retention"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	oldView := newAnalyticsView(t, link.ID())
	oldStartedAt := time.Now().UTC().AddDate(-1, -1, 0)
	oldViewModel := map[string]any{
		"id": oldView.ID().UUID(), "link_id": link.ID().UUID(), "visitor_id": oldView.VisitorID().UUID(),
		"ip_hash": make([]byte, 32), "device": "desktop", "browser": "chrome", "os": "linux",
		"traffic": "human", "layout": `{"order":["about"],"hashes":{"about":"0123456789abcdef"}}`,
		"started_at": oldStartedAt, "active_ms": 0, "last_state_seq": 0, "created_at": oldStartedAt,
	}
	require.NoError(t, gormDB.Table("resume_link_views").Create(oldViewModel).Error)

	eventID := sharedvalueobjects.NewID().UUID()
	require.NoError(t, gormDB.Exec(`INSERT INTO resume_link_events (id, view_id, link_id, client_event_id, event_type, payload, occurred_at) VALUES (?, ?, ?, ?, ?, ?::jsonb, ?)`,
		eventID, oldView.ID().UUID(), link.ID().UUID(), nil, "download", `{}`, oldStartedAt).Error)
	require.NoError(t, gormDB.Exec(`INSERT INTO resume_link_clicks (view_id, link_id, section, x_pct, y_pct, viewport_w, layout_hash, occurred_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		oldView.ID().UUID(), link.ID().UUID(), "about", 100, 200, 1440, "0123456789abcdef", oldStartedAt).Error)

	oldClickCutoff := time.Now().UTC().Add(-90 * 24 * time.Hour)
	oldViewCutoff := time.Now().UTC().AddDate(-1, 0, 0)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		count, err := repo.DeleteExpiredClicks(txCtx, oldClickCutoff, 500)
		require.NoError(t, err)
		require.Equal(t, 1, count)
		return nil
	}))
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		count, err := repo.DeleteExpiredViews(txCtx, oldViewCutoff, 500)
		require.NoError(t, err)
		require.Equal(t, 1, count)
		return nil
	}))

	var count int64
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM resume_link_views WHERE id = ?", oldView.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM resume_link_events WHERE view_id = ?", oldView.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM resume_link_clicks WHERE view_id = ?", oldView.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
}
