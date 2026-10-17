//go:build integration

package persistence_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type analyticsViewRepository interface {
	CreateView(context.Context, *domain.View) error
}

func createAnalyticsView(
	ctx context.Context,
	db *gorm.DB,
	repo analyticsViewRepository,
	view *domain.View,
) error {
	return database.NewUnitOfWork(db).Do(ctx, func(txCtx context.Context) error {
		return repo.CreateView(txCtx, view)
	})
}

func applyAnalyticsBatchInTransaction(
	ctx context.Context,
	repo interface {
		LoadAnalyticsBatchForUpdate(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID, bool, bool, bool) (domain.AnalyticsBatchSnapshot, error)
		SaveAnalyticsBatch(context.Context, *domain.View, []*domain.AnalyticsEvent, []*domain.AnalyticsClick) error
	},
	viewID, linkID sharedvalueobjects.ID,
	update domain.ViewStateUpdate,
	events []*domain.AnalyticsEvent,
	clicks []*domain.AnalyticsClick,
) error {
	needsPrintHistory := false
	for _, event := range events {
		if event != nil && event.Type() == domain.EventTypePrint {
			needsPrintHistory = true
			break
		}
	}
	snapshot, err := repo.LoadAnalyticsBatchForUpdate(ctx, viewID, linkID, len(events) > 0, len(clicks) > 0, needsPrintHistory)
	if err != nil {
		return err
	}
	now := update.SeenAt
	if now.IsZero() {
		now = time.Now().UTC()
	}
	decision, err := snapshot.View.EvaluateAnalyticsBatch(snapshot.LinkStatus, snapshot.ResumeStatus, now, update)
	if err != nil {
		return err
	}
	if decision.Ignore {
		return nil
	}
	if err := snapshot.View.ApplyState(update); err != nil {
		return err
	}
	if decision.SameSeq {
		events = nil
		clicks = nil
	} else {
		events = domain.SelectAnalyticsEventsForStorage(events, snapshot.StoredEventCount, snapshot.ExistingPrints)
		clicks = domain.SelectAnalyticsClicksForStorage(clicks, snapshot.StoredClickCount)
	}
	return repo.SaveAnalyticsBatch(ctx, snapshot.View, events, clicks)
}

func newAnalyticsView(t *testing.T, linkID sharedvalueobjects.ID, technologyCatalog ...[]string) *domain.View {
	t.Helper()
	layout, err := domain.NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionAbout, resumedomain.SectionExperience, resumedomain.SectionProjects},
		map[resumedomain.SectionID]string{
			resumedomain.SectionAbout:      "0123456789abcdef",
			resumedomain.SectionExperience: "fedcba9876543210",
			resumedomain.SectionProjects:   "0011223344556677",
		},
		technologyCatalog...,
	)
	require.NoError(t, err)
	view, err := domain.NewView(domain.ViewParams{
		LinkID: linkID, VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
		Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux, Traffic: domain.TrafficHuman,
		Layout: layout, StartedAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	return view
}

func TestResumeLinkAnalyticsSanitizesEventPayloadBeforeJSONBInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Analytics JSONB payload",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
	view := newAnalyticsView(t, link.ID(), []string{"Go", "PostgreSQL"})
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	loadedView, err := repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.True(t, loadedView.Layout().HasTechnologyCatalog())
	require.Equal(t, []string{"Go", "PostgreSQL"}, loadedView.Layout().Technologies())

	unsafeEvent, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
		ViewID: view.ID(), LinkID: link.ID(),
		Type: domain.EventTypeExternalClick,
		Payload: map[string]string{
			"kind": "external\x00-link", "target": "https://example.test/путь\x00\n",
		},
		OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	validEvent, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
		ViewID: view.ID(), LinkID: link.ID(),
		Type: domain.EventTypeDownload, OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.RecordEvents(txCtx, []*domain.AnalyticsEvent{unsafeEvent, validEvent})
	}))

	var payload string
	require.NoError(t, gormDB.Raw(
		"SELECT payload::text FROM resume_link_events WHERE id = ?", unsafeEvent.ID().UUID(),
	).Scan(&payload).Error)
	require.NotContains(t, payload, "\\u0000")
	require.Contains(t, payload, "https://example.test/путь")
	require.NotContains(t, payload, "\\n")
	var eventCount int64
	require.NoError(t, gormDB.Raw(
		"SELECT count(*) FROM resume_link_events WHERE view_id = ?", view.ID().UUID(),
	).Scan(&eventCount).Error)
	require.EqualValues(t, 2, eventCount, "sanitized payload and unrelated valid event should both persist")
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

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Analytics Batch"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	view := newAnalyticsView(t, link.ID())
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))

	clientEventID := sharedvalueobjects.NewID()
	layoutHash, ok := view.Layout().HashFor(resumedomain.SectionAbout)
	require.True(t, ok)
	click, err := domain.NewAnalyticsClick(domain.AnalyticsClickParams{
		ViewID: view.ID(), LinkID: link.ID(), Section: resumedomain.SectionAbout,
		XPct: 4210, YPct: 1733, ViewportW: 1440, LayoutHash: layoutHash, Target: "download", OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	event, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
		ViewID: view.ID(), LinkID: link.ID(), ClientEventID: &clientEventID,
		Type: domain.EventTypeDownload, OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), domain.ViewStateUpdate{
			Seq: 2, ActiveMS: 7000, MaxSection: resumedomain.SectionExperience,
			ExitSection: resumedomain.SectionExperience, SeenAt: time.Now().UTC(), Final: false,
		}, []*domain.AnalyticsEvent{event}, []*domain.AnalyticsClick{click})
	}))

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), domain.ViewStateUpdate{
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
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), domain.ViewStateUpdate{
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

func TestResumeLinkAnalyticsPrintEventsAreDeduplicatedPerBatchAndAcrossBatches(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Print deduplication"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
	view := newAnalyticsView(t, link.ID())
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))

	newPrintEvent := func() *domain.AnalyticsEvent {
		t.Helper()
		clientEventID := sharedvalueobjects.NewID()
		event, eventErr := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
			ViewID: view.ID(), LinkID: link.ID(), ClientEventID: &clientEventID,
			Type: domain.EventTypePrint, OccurredAt: time.Now().UTC(),
		})
		require.NoError(t, eventErr)
		return event
	}

	apply := func(seq int64, events ...*domain.AnalyticsEvent) {
		t.Helper()
		err := uow.Do(ctx, func(txCtx context.Context) error {
			return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), domain.ViewStateUpdate{
				Seq: seq, SeenAt: time.Now().UTC(),
			}, events, nil)
		})
		require.NoError(t, err)
	}
	apply(1, newPrintEvent(), newPrintEvent())

	var count int64
	require.NoError(t, gormDB.Model(&gormmodel.ResumeLinkAnalyticsEvent{}).
		Where("view_id = ? AND event_type = ?", view.ID().UUID(), string(domain.EventTypePrint)).Count(&count).Error)
	require.EqualValues(t, 1, count, "only one print should be stored from a single batch")

	apply(2, newPrintEvent())
	require.NoError(t, gormDB.Model(&gormmodel.ResumeLinkAnalyticsEvent{}).
		Where("view_id = ? AND event_type = ?", view.ID().UUID(), string(domain.EventTypePrint)).Count(&count).Error)
	require.EqualValues(t, 1, count, "a print inside the two-second window should be deduplicated across batches")
}

func TestResumeLinkAnalyticsBatchEnforcesPerViewEventAndClickCaps(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Analytics caps",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
	view := newAnalyticsView(t, link.ID())
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))

	now := time.Now().UTC()
	seedEvents := make([]gormmodel.ResumeLinkAnalyticsEvent, domain.MaxAnalyticsEventsPerView-1)
	for i := range seedEvents {
		clientID := sharedvalueobjects.NewID().UUID()
		seedEvents[i] = gormmodel.ResumeLinkAnalyticsEvent{
			ID: sharedvalueobjects.NewID().UUID(), ViewID: view.ID().UUID(), LinkID: link.ID().UUID(),
			ClientEventID: &clientID, EventType: string(domain.EventTypeDownload), Payload: []byte("{}"), OccurredAt: now,
		}
	}
	require.NoError(t, gormDB.Create(&seedEvents).Error)

	seedClicks := make([]gormmodel.ResumeLinkClick, domain.MaxAnalyticsClicksPerView-1)
	for i := range seedClicks {
		target := "button"
		seedClicks[i] = gormmodel.ResumeLinkClick{
			ViewID: view.ID().UUID(), LinkID: link.ID().UUID(), Section: "about",
			XPct: 100, YPct: 200, ViewportW: 1280, LayoutHash: "0123456789abcdef",
			Target: &target, OccurredAt: now,
		}
	}
	require.NoError(t, gormDB.Create(&seedClicks).Error)

	newDownload := func() *domain.AnalyticsEvent {
		clientID := sharedvalueobjects.NewID()
		event, eventErr := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
			ViewID: view.ID(), LinkID: link.ID(), ClientEventID: &clientID,
			Type: domain.EventTypeDownload, OccurredAt: now,
		})
		require.NoError(t, eventErr)
		return event
	}
	newClick := func() *domain.AnalyticsClick {
		click, clickErr := domain.NewAnalyticsClick(domain.AnalyticsClickParams{
			ViewID: view.ID(), LinkID: link.ID(),
			Section: resumedomain.SectionAbout, XPct: 100, YPct: 200, ViewportW: 1280,
			LayoutHash: "0123456789abcdef", Target: "button", OccurredAt: now,
		})
		require.NoError(t, clickErr)
		return click
	}
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), domain.ViewStateUpdate{
			Seq: 1, SeenAt: now,
		}, []*domain.AnalyticsEvent{newDownload(), newDownload()}, []*domain.AnalyticsClick{newClick(), newClick()})
	}))

	var eventCount, clickCount int64
	require.NoError(t, gormDB.Model(&gormmodel.ResumeLinkAnalyticsEvent{}).
		Where("view_id = ? AND client_event_id IS NOT NULL", view.ID().UUID()).Count(&eventCount).Error)
	require.NoError(t, gormDB.Model(&gormmodel.ResumeLinkClick{}).
		Where("view_id = ?", view.ID().UUID()).Count(&clickCount).Error)
	require.EqualValues(t, domain.MaxAnalyticsEventsPerView, eventCount)
	require.EqualValues(t, domain.MaxAnalyticsClicksPerView, clickCount)
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
		link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: name})
		require.NoError(t, err)
		require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
		links = append(links, link)
	}

	first := newAnalyticsView(t, links[0].ID())
	second := newAnalyticsView(t, links[1].ID())
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, first))
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, second))
	startedAt := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	require.NoError(t, gormDB.Exec(`UPDATE resume_link_views SET started_at = ? WHERE id = ?`, startedAt, first.ID().UUID()).Error)
	require.NoError(t, gormDB.Exec(`UPDATE resume_link_views SET started_at = ? WHERE id = ?`, startedAt.Add(5*time.Minute), second.ID().UUID()).Error)

	period, err := queries.NewAnalyticsPeriod("2026-10-01", "2026-10-05", "UTC", startedAt)
	require.NoError(t, err)
	linkStats, err := reader.GetLinkStatistics(ctx, queries.LinkAnalyticsQuery{
		LinkID: links[0].ID(),
		Period: period,
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), linkStats.Views)

	result, err := reader.GetResumeStatistics(ctx, queries.ResumeAnalyticsQuery{ResumeID: resumeID, Period: period})
	require.NoError(t, err)
	require.Len(t, result.Links, 2)
	require.Equal(t, links[0].ID().String(), result.Links[0].ID)
	require.Equal(t, links[1].ID().String(), result.Links[1].ID)
	require.Equal(t, int64(1), result.Links[0].Views)
	require.Equal(t, int64(1), result.Links[1].Views)
	require.Equal(t, int64(2), result.Views)
}

func TestResumeLinkAnalyticsReferrerStatisticsAreCappedAndDeterministic(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Referrer cap",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
	startedAt := time.Date(2026, 10, 10, 10, 0, 0, 0, time.UTC)
	for index := 0; index < 25; index++ {
		host := fmt.Sprintf("host-%02d.example", index)
		view, viewErr := domain.NewView(domain.ViewParams{
			LinkID: link.ID(), VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
			Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux, Traffic: domain.TrafficHuman,
			ReferrerHost: host, Layout: newAnalyticsView(t, link.ID()).Layout(), StartedAt: startedAt.Add(time.Duration(index) * time.Minute),
		})
		require.NoError(t, viewErr)
		require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	}
	period, err := queries.NewAnalyticsPeriod("2026-10-10", "2026-10-11", "UTC", startedAt)
	require.NoError(t, err)
	stats, err := reader.GetLinkStatistics(ctx, queries.LinkAnalyticsQuery{LinkID: link.ID(), Period: period})
	require.NoError(t, err)
	require.Len(t, stats.Referrers, 20)
	require.Equal(t, "host-00.example", *stats.Referrers[0].Host)
	require.Equal(t, "host-19.example", *stats.Referrers[19].Host)
}

func TestResumeLinkAnalyticsReadMetricsUseConfirmationState(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Confirmation state",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	startedAt := time.Date(2026, 10, 10, 10, 0, 0, 0, time.UTC)
	unconfirmed := newAnalyticsView(t, link.ID())
	confirmed := newAnalyticsView(t, link.ID())
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, unconfirmed))
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, confirmed))
	require.NoError(t, gormDB.Exec(
		`UPDATE resume_link_views SET started_at = ?, active_ms = 9000, confirmed_at = NULL WHERE id = ?`,
		startedAt, unconfirmed.ID().UUID(),
	).Error)
	confirmedAt := startedAt.Add(time.Second)
	require.NoError(t, gormDB.Exec(
		`UPDATE resume_link_views SET started_at = ?, active_ms = 0, confirmed_at = ? WHERE id = ?`,
		startedAt.Add(time.Minute), confirmedAt, confirmed.ID().UUID(),
	).Error)
	period, err := queries.NewAnalyticsPeriod("2026-10-10", "2026-10-11", "UTC", startedAt)
	require.NoError(t, err)

	stats, err := reader.GetLinkStatistics(ctx, queries.LinkAnalyticsQuery{LinkID: link.ID(), Period: period})
	require.NoError(t, err)
	require.EqualValues(t, 2, stats.Views)
	require.EqualValues(t, 1, stats.Reads)
	require.Zero(t, stats.MedianActiveSeconds)

	daily, err := reader.GetDaily(ctx, queries.LinkAnalyticsQuery{LinkID: link.ID(), Period: period})
	require.NoError(t, err)
	require.EqualValues(t, 1, daily.Days[0].Reads)

	sections, err := reader.GetSections(ctx, queries.LinkSectionsQuery{
		LinkID: link.ID(), Period: period,
		SectionOrder: []string{"about", "experience", "projects"},
	})
	require.NoError(t, err)
	require.EqualValues(t, 1, sections.Reads)
	require.NotNil(t, sections.Sections)

	visits, err := reader.GetVisits(ctx, queries.VisitsReadQuery{LinkID: link.ID(), Limit: 10})
	require.NoError(t, err)
	readByID := make(map[string]bool, len(visits.Items))
	for _, visit := range visits.Items {
		readByID[visit.ID] = visit.IsRead
	}
	require.False(t, readByID[unconfirmed.ID().String()])
	require.True(t, readByID[confirmed.ID().String()])
}

func TestResumeLinkAnalyticsReturningVisitorMetricsKeepTheirCalendarDaySemantics(t *testing.T) {
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
	for _, name := range []string{"Primary", "Secondary"} {
		link, err := domain.NewResumeLink(domain.ResumeLinkParams{
			ResumeID: resumeID, Name: name,
		})
		require.NoError(t, err)
		require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
		links = append(links, link)
	}

	base := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	visitorID := sharedvalueobjects.NewID()
	first := newAnalyticsViewForVisitor(t, links[0].ID(), visitorID, base)
	sameDayReturn := newAnalyticsViewForVisitor(t, links[0].ID(), visitorID, base.Add(30*time.Minute))
	nextDayOtherLink := newAnalyticsViewForVisitor(t, links[1].ID(), visitorID, base.Add(24*time.Hour))
	for _, view := range []*domain.View{first, sameDayReturn, nextDayOtherLink} {
		require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	}

	period, err := queries.NewAnalyticsPeriod("2026-10-01", "2026-10-03", "UTC", base)
	require.NoError(t, err)

	visits, err := reader.GetVisits(ctx, queries.VisitsReadQuery{LinkID: links[0].ID(), Limit: 10})
	require.NoError(t, err)
	require.Len(t, visits.Items, 2)
	require.Equal(t, sameDayReturn.ID().String(), visits.Items[0].ID)
	require.True(t, visits.Items[0].Visitor.IsReturning)
	require.False(t, visits.Items[1].Visitor.IsReturning)

	linkStats, err := reader.GetLinkStatistics(ctx, queries.LinkAnalyticsQuery{LinkID: links[0].ID(), Period: period})
	require.NoError(t, err)
	require.EqualValues(t, 0, linkStats.ReturningVisitors)

	resumeStats, err := reader.GetResumeStatistics(ctx, queries.ResumeAnalyticsQuery{ResumeID: resumeID, Period: period})
	require.NoError(t, err)
	require.EqualValues(t, 1, resumeStats.ReturningVisitors)

	daily, err := reader.GetDaily(ctx, queries.LinkAnalyticsQuery{LinkID: links[0].ID(), Period: period})
	require.NoError(t, err)
	require.Len(t, daily.Days, 3)
	require.EqualValues(t, 1, daily.Days[0].UniqueVisitors)
	require.EqualValues(t, 0, daily.Days[1].UniqueVisitors)
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

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Visits"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	first := newAnalyticsView(t, link.ID())
	second := newAnalyticsView(t, link.ID())
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, first))
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, second))

	firstStartedAt := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	secondStartedAt := firstStartedAt.Add(5 * time.Minute)
	require.NoError(t, gormDB.Exec(`UPDATE resume_link_views SET started_at = ? WHERE id = ?`, firstStartedAt, first.ID().UUID()).Error)
	require.NoError(t, gormDB.Exec(`UPDATE resume_link_views SET started_at = ? WHERE id = ?`, secondStartedAt, second.ID().UUID()).Error)

	page, err := reader.GetVisits(ctx, queries.VisitsReadQuery{LinkID: link.ID(), Limit: 1})
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	require.Equal(t, second.ID().String(), page.Items[0].ID)
	require.NotNil(t, page.NextCursor)

	cursorTime, cursorID, err := queries.DecodeVisitCursor(*page.NextCursor)
	require.NoError(t, err)
	nextPage, err := reader.GetVisits(ctx, queries.VisitsReadQuery{
		LinkID: link.ID(), Limit: 1, CursorTime: cursorTime, CursorID: cursorID,
	})
	require.NoError(t, err)
	require.Len(t, nextPage.Items, 1)
	require.Equal(t, first.ID().String(), nextPage.Items[0].ID)
	require.Nil(t, nextPage.NextCursor)
}

func TestResumeLinkAnalyticsVisitorOrdinalAndFirstSeenSurviveRetention(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Visitor history",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	base := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	visitorA := sharedvalueobjects.NewID()
	visitorB := sharedvalueobjects.NewID()
	firstA := newAnalyticsViewForVisitor(t, link.ID(), visitorA, base.Add(time.Minute))
	firstB := newAnalyticsViewForVisitor(t, link.ID(), visitorB, base)
	secondA := newAnalyticsViewForVisitor(t, link.ID(), visitorA, base.Add(2*time.Minute))

	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, firstA))
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, firstB))
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, secondA))

	page, err := reader.GetVisits(ctx, queries.VisitsReadQuery{LinkID: link.ID(), Limit: 10})
	require.NoError(t, err)
	require.Len(t, page.Items, 3)
	require.Equal(t, int64(2), page.Items[0].Visitor.Ordinal)
	require.True(t, page.Items[0].Visitor.IsReturning)
	require.Equal(t, firstA.ID().String(), page.Items[1].ID)
	require.Equal(t, int64(2), page.Items[1].Visitor.Ordinal)
	require.False(t, page.Items[1].Visitor.IsReturning)
	require.Equal(t, int64(1), page.Items[2].Visitor.Ordinal)

	var deleted int
	err = uow.Do(ctx, func(txCtx context.Context) error {
		var deleteErr error
		deleted, deleteErr = repo.DeleteExpiredViews(txCtx, base.Add(90*time.Second), 100)
		return deleteErr
	})
	require.NoError(t, err)
	require.Equal(t, 2, deleted)
	remaining, err := reader.GetVisits(ctx, queries.VisitsReadQuery{LinkID: link.ID(), Limit: 10})
	require.NoError(t, err)
	require.Len(t, remaining.Items, 1)
	require.Equal(t, secondA.ID().String(), remaining.Items[0].ID)
	require.Equal(t, int64(2), remaining.Items[0].Visitor.Ordinal)
	require.True(t, remaining.Items[0].Visitor.IsReturning)
}

func TestResumeLinkAnalyticsConcurrentVisitorRegistrationKeepsDenseOrdinal(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Concurrent visitors",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	const visitorCount = 8
	base := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	views := make([]*domain.View, 0, visitorCount)
	for i := visitorCount - 1; i >= 0; i-- {
		views = append(views, newAnalyticsViewForVisitor(
			t, link.ID(), sharedvalueobjects.NewID(), base.Add(time.Duration(i)*time.Minute),
		))
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(views))
	for _, view := range views {
		view := view
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- createAnalyticsView(ctx, gormDB, repo, view)
		}()
	}
	wg.Wait()
	close(errCh)
	for createErr := range errCh {
		require.NoError(t, createErr)
	}

	page, err := reader.GetVisits(ctx, queries.VisitsReadQuery{LinkID: link.ID(), Limit: visitorCount})
	require.NoError(t, err)
	require.Len(t, page.Items, visitorCount)
	for i, visit := range page.Items {
		require.Equal(t, int64(visitorCount-i), visit.Visitor.Ordinal)
		require.False(t, visit.Visitor.IsReturning)
	}
}

func newAnalyticsViewForVisitor(
	t *testing.T,
	linkID sharedvalueobjects.ID,
	visitorID sharedvalueobjects.ID,
	startedAt time.Time,
) *domain.View {
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
		LinkID: linkID, VisitorID: visitorID, IPHash: make([]byte, 32),
		Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux, Traffic: domain.TrafficHuman,
		Layout: layout, StartedAt: startedAt, CreatedAt: startedAt,
	})
	require.NoError(t, err)
	return view
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

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Lifecycle"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
	view := newAnalyticsView(t, link.ID())
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	require.NoError(t, gormDB.Exec(`UPDATE resume_links SET status = 'inactive' WHERE id = ?`, link.ID().UUID()).Error)
	err = uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), domain.ViewStateUpdate{Seq: 1, SeenAt: time.Now().UTC()}, nil, nil)
	})
	require.ErrorIs(t, err, domain.ErrResumeLinkInactive)

	require.NoError(t, gormDB.Exec(`UPDATE resume_links SET status = 'active' WHERE id = ?`, link.ID().UUID()).Error)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'draft' WHERE id = ?`, resumeID.UUID()).Error)
	err = uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), domain.ViewStateUpdate{Seq: 2, SeenAt: time.Now().UTC()}, nil, nil)
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

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Analytics"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	view := newAnalyticsView(t, link.ID())
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	loaded, err := repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.Equal(t, view.ID(), loaded.ID())
	require.Equal(t, view.LinkID(), loaded.LinkID())
	require.Equal(t, view.Layout().Order(), loaded.Layout().Order())
	require.WithinDuration(t, view.CreatedAt(), loaded.CreatedAt(), time.Microsecond)

	now := time.Now().UTC()
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.UpdateViewState(txCtx, view.ID(), domain.ViewStateUpdate{
			Seq: 1, ActiveMS: 7000, MaxSection: resumedomain.SectionExperience, ExitSection: resumedomain.SectionExperience, SeenAt: now, Final: true,
		})
	}))
	loaded, err = repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.Equal(t, int64(7000), loaded.State().ActiveMS())
	require.Equal(t, resumedomain.SectionExperience, loaded.State().MaxSection())
	require.False(t, loaded.State().EndedAt().IsZero())

	clientEventID := sharedvalueobjects.NewID()
	event, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
		ViewID: view.ID(), LinkID: link.ID(), ClientEventID: &clientEventID,
		Type: domain.EventTypeSectionView, Payload: map[string]string{"section": "experience"}, OccurredAt: now,
	})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.RecordEvents(txCtx, []*domain.AnalyticsEvent{event, event})
	}))

	var count int
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM resume_link_events WHERE view_id = ?", view.ID().UUID()).Scan(&count).Error)
	require.Equal(t, 1, count)

	require.NoError(t, gormDB.Raw("SELECT count(*) FROM resume_link_visitors WHERE link_id = ?", link.ID().UUID()).Scan(&count).Error)
	require.EqualValues(t, 1, count)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Delete(txCtx, link.ID()) }))
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM resume_link_views WHERE id = ?", view.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM resume_link_events WHERE view_id = ?", view.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM resume_link_visitors WHERE link_id = ?", link.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count, "visitor pseudonymous metadata must be deleted with its link")
}

func TestResumeLinkAnalyticsRetentionDeletesExpiredDataInBatches(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Analytics Retention"})
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
	require.NoError(t, gormDB.Exec(
		`INSERT INTO resume_link_events
		 (id, view_id, link_id, client_event_id, event_type, payload, occurred_at)
		 VALUES (?, ?, ?, ?, ?, ?::jsonb, ?)`,
		eventID, oldView.ID().UUID(), link.ID().UUID(), nil,
		"download", `{}`, oldStartedAt,
	).Error)
	require.NoError(t, gormDB.Exec(
		`INSERT INTO resume_link_clicks
		 (view_id, link_id, section, x_pct, y_pct, viewport_w, layout_hash, occurred_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		oldView.ID().UUID(), link.ID().UUID(), "about", 100, 200,
		1440, "0123456789abcdef", oldStartedAt,
	).Error)

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

func TestResumeLinkAnalyticsClickReaderUsesGORMFiltersAndCountsStaleClicks(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Click map",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.Create(txCtx, link)
	}))

	view := newAnalyticsView(t, link.ID())
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	currentHash, ok := view.Layout().HashFor(resumedomain.SectionAbout)
	require.True(t, ok)
	clicks := make([]*domain.AnalyticsClick, 0, 2)
	for _, layoutHash := range []string{currentHash, "fedcba9876543210"} {
		click, err := domain.NewAnalyticsClick(domain.AnalyticsClickParams{
			ViewID:     view.ID(),
			LinkID:     link.ID(),
			Section:    resumedomain.SectionAbout,
			XPct:       4210,
			YPct:       1733,
			ViewportW:  1440,
			LayoutHash: layoutHash,
			Target:     "download",
			OccurredAt: time.Now().UTC(),
		})
		require.NoError(t, err)
		clicks = append(clicks, click)
	}
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.RecordClicks(txCtx, clicks)
	}))

	currentLayout, err := reader.GetClicks(ctx, queries.LinkClicksQuery{
		LinkID:      link.ID(),
		Section:     string(resumedomain.SectionAbout),
		Breakpoint:  "lg",
		CurrentHash: currentHash,
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), currentLayout.Total)
	require.Equal(t, int64(1), currentLayout.StaleTotal)
	require.Len(t, currentLayout.Cells, 1)
	require.Equal(t, int64(8), currentLayout.Cells[0].X)
	require.Equal(t, int64(3), currentLayout.Cells[0].Y)

	allLayouts, err := reader.GetClicks(ctx, queries.LinkClicksQuery{
		LinkID:       link.ID(),
		Section:      string(resumedomain.SectionAbout),
		Breakpoint:   "lg",
		IncludeStale: true,
		CurrentHash:  currentHash,
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), allLayouts.Total)
	require.Zero(t, allLayouts.StaleTotal)
	require.Len(t, allLayouts.Cells, 1)
	require.Equal(t, int64(2), allLayouts.Cells[0].Count)
}

func TestResumeLinkAnalyticsRejectsEventsAndClicksForAnotherLink(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	firstLink, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "First",
	})
	require.NoError(t, err)
	secondLink, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Second",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, firstLink) }))
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, secondLink) }))

	view := newAnalyticsView(t, firstLink.ID())
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))

	event, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
		ViewID: view.ID(), LinkID: secondLink.ID(),
		Type: domain.EventTypeDownload, OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	require.ErrorIs(t, repo.RecordEvents(ctx, []*domain.AnalyticsEvent{event}), database.ErrTransactionRequired)
	require.ErrorIs(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.RecordEvents(txCtx, []*domain.AnalyticsEvent{event}) }), domain.ErrEntityNotFound)
	err = uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), firstLink.ID(), domain.ViewStateUpdate{
			Seq: 1, SeenAt: time.Now().UTC(),
		}, []*domain.AnalyticsEvent{event}, nil)
	})
	require.Error(t, err)

	layoutHash, ok := view.Layout().HashFor(resumedomain.SectionAbout)
	require.True(t, ok)
	click, err := domain.NewAnalyticsClick(domain.AnalyticsClickParams{
		ViewID: view.ID(), LinkID: secondLink.ID(),
		Section: resumedomain.SectionAbout, XPct: 100, YPct: 200, ViewportW: 1440,
		LayoutHash: layoutHash, OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	require.ErrorIs(t, repo.RecordClicks(ctx, []*domain.AnalyticsClick{click}), database.ErrTransactionRequired)
	require.ErrorIs(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.RecordClicks(txCtx, []*domain.AnalyticsClick{click}) }), domain.ErrEntityNotFound)

	var eventCount, clickCount int64
	require.NoError(t, gormDB.Raw(
		`SELECT count(*) FROM resume_link_events WHERE view_id = ?`, view.ID().UUID(),
	).Scan(&eventCount).Error)
	require.NoError(t, gormDB.Raw(
		`SELECT count(*) FROM resume_link_clicks WHERE view_id = ?`, view.ID().UUID(),
	).Scan(&clickCount).Error)
	require.Zero(t, eventCount)
	require.Zero(t, clickCount)
}

func TestResumeLinkAnalyticsRetentionCleanupLockIsTransactionScoped(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, _, gormDB, repo, _ := newDatabase(t)
	uow := database.NewUnitOfWork(gormDB)
	locked := make(chan struct{})
	release := make(chan struct{})
	firstDone := make(chan error, 1)

	go func() {
		firstDone <- uow.Do(ctx, func(txCtx context.Context) error {
			count, err := repo.DeleteExpiredClicks(txCtx, time.Now().UTC(), 1)
			if err != nil {
				return err
			}
			if count != 0 {
				return fmt.Errorf("expected no expired clicks, got %d", count)
			}
			close(locked)
			<-release
			return nil
		})
	}()
	defer func() {
		select {
		case <-release:
		default:
			close(release)
		}
	}()

	select {
	case <-locked:
	case <-time.After(5 * time.Second):
		t.Fatal("first cleanup transaction did not acquire the advisory lock")
	}

	var acquired bool
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		tx, err := database.DBFromTransactionContext(txCtx)
		if err != nil {
			return err
		}
		return tx.Raw("SELECT pg_try_advisory_xact_lock(?)", int64(741209882133)).Scan(&acquired).Error
	}))
	require.False(t, acquired, "another transaction must not acquire the cleanup lock while the first transaction is open")

	close(release)
	select {
	case err := <-firstDone:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("first cleanup transaction did not finish after releasing the lock")
	}

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		tx, err := database.DBFromTransactionContext(txCtx)
		if err != nil {
			return err
		}
		return tx.Raw("SELECT pg_try_advisory_xact_lock(?)", int64(741209882133)).Scan(&acquired).Error
	}))
	require.True(t, acquired, "transaction-scoped cleanup lock must be released after commit")
}
