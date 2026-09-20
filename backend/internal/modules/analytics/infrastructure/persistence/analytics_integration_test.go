//go:build integration

package persistence_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	analyticscommands "developers-cv/backend/internal/modules/analytics/application/commands"
	analyticsports "developers-cv/backend/internal/modules/analytics/application/ports"
	analyticsqueries "developers-cv/backend/internal/modules/analytics/application/queries"
	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	analyticsdomain "developers-cv/backend/internal/modules/analytics/domain"
	analyticsgormmodel "developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type analyticsViewRepository interface {
	CreateView(context.Context, *analyticsdomain.View) error
}

func createAnalyticsView(
	ctx context.Context,
	db *gorm.DB,
	repo analyticsViewRepository,
	view *analyticsdomain.View,
) error {
	return database.NewUnitOfWork(db).Do(ctx, func(txCtx context.Context) error {
		return repo.CreateView(txCtx, view)
	})
}

func applyAnalyticsBatchInTransaction(
	ctx context.Context,
	repo *analyticsTestRepository,
	viewID, linkID sharedvalueobjects.ID,
	update analyticsdomain.ViewStateUpdate,
	events []*analyticsdomain.AnalyticsEvent,
	clicks []*analyticsdomain.AnalyticsClick,
) error {
	eventInputs := make([]analyticscommands.AnalyticsEventInput, 0, len(events))
	for _, event := range events {
		if event == nil {
			continue
		}
		var clientEventID uuid.UUID
		if id := event.ClientEventID(); id != nil {
			clientEventID = id.UUID()
		}
		eventInputs = append(eventInputs, analyticscommands.AnalyticsEventInput{
			ClientEventID: clientEventID, Type: event.Type(), Payload: event.Payload(),
			ClientTimestampMS: event.OccurredAt().UnixMilli(),
		})
	}
	clickInputs := make([]analyticscommands.AnalyticsClickInput, 0, len(clicks))
	for _, click := range clicks {
		if click == nil {
			continue
		}
		clickInputs = append(clickInputs, analyticscommands.AnalyticsClickInput{
			Section: string(click.Section()), XPct: click.XPct(), YPct: click.YPct(),
			ViewportW: click.ViewportW(), Target: click.Target(),
			ClientTimestampMS: click.OccurredAt().UnixMilli(),
		})
	}
	handler := analyticscommands.NewIngestAnalyticsBatchHandler(repo, txContextUnitOfWork{}, repo.Subjects, nil)
	return handler.Handle(ctx, analyticscommands.AnalyticsBatchInput{
		ViewID: viewID, LinkID: linkID, Seq: update.Seq, ActiveMS: update.ActiveMS,
		MaxSection: string(update.MaxSection), ExitSection: string(update.ExitSection),
		Final: update.Final, Events: eventInputs, Clicks: clickInputs,
	})
}

type txContextUnitOfWork struct{}

func (txContextUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type signalingSubjectReader struct {
	inner       analyticscontracts.AnalyticsSubjectReader
	linkContext chan analyticsports.LinkAnalyticsContext
	locked      chan struct{}
	release     chan struct{}
	once        sync.Once
}

func (s *signalingSubjectReader) LoadLinkLifecycle(ctx context.Context, linkID sharedvalueobjects.ID) (analyticscontracts.LinkLifecycle, error) {
	lifecycle, err := s.inner.LoadLinkLifecycle(ctx, linkID)
	if err != nil {
		return analyticscontracts.LinkLifecycle{}, err
	}
	s.once.Do(func() {
		s.linkContext <- analyticsports.LinkAnalyticsContext{
			LinkID: lifecycle.LinkID, ResumeID: lifecycle.ResumeID,
			LinkActive: lifecycle.LinkActive, ResumePublished: lifecycle.ResumePublished,
		}
		close(s.locked)
		<-s.release
	})
	return lifecycle, nil
}

func (s *signalingSubjectReader) ListResumeLinks(ctx context.Context, resumeID sharedvalueobjects.ID) ([]analyticscontracts.ResumeLinkRef, error) {
	return s.inner.ListResumeLinks(ctx, resumeID)
}

func newAnalyticsView(t *testing.T, link *domain.ResumeLink, technologyCatalog ...[]string) *analyticsdomain.View {
	t.Helper()
	layout, err := analyticsdomain.NewLayoutSnapshot(
		[]analyticsdomain.AnalyticsSectionID{analyticsdomain.AnalyticsSectionAbout, analyticsdomain.AnalyticsSectionExperience, analyticsdomain.AnalyticsSectionProjects},
		map[analyticsdomain.AnalyticsSectionID]string{
			analyticsdomain.AnalyticsSectionAbout:      "0123456789abcdef",
			analyticsdomain.AnalyticsSectionExperience: "fedcba9876543210",
			analyticsdomain.AnalyticsSectionProjects:   "0011223344556677",
		},
		technologyCatalog...,
	)
	require.NoError(t, err)
	view, err := analyticsdomain.NewView(analyticsdomain.ViewParams{
		LinkID: link.ID(), ResumeID: link.ResumeID(), VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
		Device: analyticsdomain.DeviceDesktop, Browser: analyticsdomain.BrowserChrome, OS: analyticsdomain.OSLinux, Traffic: analyticsdomain.TrafficHuman,
		Layout: layout, StartedAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	return view
}

func TestAnalyticsSanitizesEventPayloadBeforeJSONBInsert(t *testing.T) {
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
	view := newAnalyticsView(t, link, []string{"Go", "PostgreSQL"})
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	loadedView, err := repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.True(t, loadedView.Layout().HasTechnologyCatalog())
	require.Equal(t, []string{"Go", "PostgreSQL"}, loadedView.Layout().Technologies())

	unsafeEvent, err := analyticsdomain.NewAnalyticsEvent(analyticsdomain.AnalyticsEventParams{
		ViewID: view.ID(), LinkID: link.ID(),
		Type: analyticsdomain.EventTypeExternalClick,
		Payload: map[string]string{
			"kind": "external\x00-link", "target": "https://example.test/путь\x00\n",
		},
		OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	validEvent, err := analyticsdomain.NewAnalyticsEvent(analyticsdomain.AnalyticsEventParams{
		ViewID: view.ID(), LinkID: link.ID(),
		Type: analyticsdomain.EventTypeDownload, OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.RecordEvents(txCtx, []*analyticsdomain.AnalyticsEvent{unsafeEvent, validEvent})
	}))

	var payload string
	require.NoError(t, gormDB.Raw(
		"SELECT payload::text FROM analytics.analytics_events WHERE id = ?", unsafeEvent.ID().UUID(),
	).Scan(&payload).Error)
	require.NotContains(t, payload, "\\u0000")
	require.Contains(t, payload, "https://example.test/путь")
	require.NotContains(t, payload, "\\n")
	var eventCount int64
	require.NoError(t, gormDB.Raw(
		"SELECT count(*) FROM analytics.analytics_events WHERE view_id = ?", view.ID().UUID(),
	).Scan(&eventCount).Error)
	require.EqualValues(t, 2, eventCount, "sanitized payload and unrelated valid event should both persist")
}

func TestAnalyticsBatchIsMonotonic(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Analytics Batch"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	view := newAnalyticsView(t, link)
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))

	clientEventID := sharedvalueobjects.NewID()
	layoutHash, ok := view.Layout().HashFor(analyticsdomain.AnalyticsSectionAbout)
	require.True(t, ok)
	click, err := analyticsdomain.NewAnalyticsClick(analyticsdomain.AnalyticsClickParams{
		ViewID: view.ID(), LinkID: link.ID(), Section: analyticsdomain.AnalyticsSectionAbout,
		XPct: 4210, YPct: 1733, ViewportW: 1440, LayoutHash: layoutHash, Target: "download", OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	event, err := analyticsdomain.NewAnalyticsEvent(analyticsdomain.AnalyticsEventParams{
		ViewID: view.ID(), LinkID: link.ID(), ClientEventID: &clientEventID,
		Type: analyticsdomain.EventTypeDownload, OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), analyticsdomain.ViewStateUpdate{
			Seq: 2, ActiveMS: 7000, MaxSection: analyticsdomain.AnalyticsSectionExperience,
			ExitSection: analyticsdomain.AnalyticsSectionExperience, SeenAt: time.Now().UTC(), Final: false,
		}, []*analyticsdomain.AnalyticsEvent{event}, []*analyticsdomain.AnalyticsClick{click})
	}))

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), analyticsdomain.ViewStateUpdate{
			Seq: 1, ActiveMS: 100, MaxSection: analyticsdomain.AnalyticsSectionAbout,
			ExitSection: analyticsdomain.AnalyticsSectionAbout, SeenAt: time.Now().UTC(), Final: false,
		}, nil, nil)
	}))

	loaded, err := repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.Equal(t, int64(7000), loaded.State().ActiveMS())
	require.Equal(t, int64(2), loaded.State().LastSeq())
	require.Equal(t, analyticsdomain.AnalyticsSectionExperience, loaded.State().MaxSection())
	require.True(t, loaded.State().EndedAt().IsZero())

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), analyticsdomain.ViewStateUpdate{
			Seq: 2, ActiveMS: 7000, SeenAt: time.Now().UTC(), Final: true,
		}, nil, []*analyticsdomain.AnalyticsClick{click})
	}))

	loaded, err = repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.False(t, loaded.State().EndedAt().IsZero())

	var count int
	require.NoError(t, gormDB.Raw(
		"SELECT count(*) FROM analytics.analytics_events WHERE view_id = ? AND client_event_id = ?",
		view.ID().UUID(), clientEventID.UUID(),
	).Scan(&count).Error)
	require.Equal(t, 1, count)
	var clickCount int
	require.NoError(t, gormDB.Raw(
		"SELECT count(*) FROM analytics.analytics_clicks WHERE view_id = ?",
		view.ID().UUID(),
	).Scan(&clickCount).Error)
	require.Equal(t, 1, clickCount)
}

func TestAnalyticsBatchSerializesWithLifecycleChanges(t *testing.T) {
	tests := []struct {
		name          string
		table         string
		getID         func(sharedvalueobjects.ID, sharedvalueobjects.ID) sharedvalueobjects.ID
		updateSQL     string
		expectedError error
	}{
		{
			name:          "link deactivation",
			table:         "resume_link.resume_links",
			getID:         func(linkID, _ sharedvalueobjects.ID) sharedvalueobjects.ID { return linkID },
			updateSQL:     "UPDATE resume_link.resume_links SET status = 'inactive' WHERE id = $1",
			expectedError: analyticsdomain.ErrAnalyticsLinkInactive,
		},
		{
			name:          "resume unpublish",
			table:         "resume.resumes",
			getID:         func(_, resumeID sharedvalueobjects.ID) sharedvalueobjects.ID { return resumeID },
			updateSQL:     "UPDATE resume.resumes SET status = 'draft' WHERE id = $1",
			expectedError: analyticscontracts.ErrAnalyticsSubjectNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if testing.Short() {
				t.Skip("integration test")
			}

			ctx, dsn, gormDB, repo, _ := newDatabase(t)
			owner := sharedvalueobjects.NewID()
			resumeID := sharedvalueobjects.NewID()
			insertUserAndResume(t, ctx, dsn, owner, resumeID)
			require.NoError(t, gormDB.Exec(
				"UPDATE resume.resumes SET status = 'opened' WHERE id = ?", resumeID.UUID(),
			).Error)

			link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Lifecycle race"})
			require.NoError(t, err)
			uow := database.NewUnitOfWork(gormDB)
			require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

			view := newAnalyticsView(t, link)
			require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))

			locked := make(chan struct{})
			linkContext := make(chan analyticsports.LinkAnalyticsContext, 1)
			release := make(chan struct{})
			var releaseOnce sync.Once
			unblock := func() { releaseOnce.Do(func() { close(release) }) }
			defer unblock()
			repo.Subjects = &signalingSubjectReader{
				inner: repo.Subjects, linkContext: linkContext, locked: locked, release: release,
			}
			batchResult := make(chan error, 1)
			go func() {
				batchResult <- uow.Do(ctx, func(txCtx context.Context) error {
					return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), analyticsdomain.ViewStateUpdate{
						Seq: 1, SeenAt: time.Now().UTC(),
					}, nil, nil)
				})
			}()

			select {
			case <-locked:
			case <-time.After(5 * time.Second):
				t.Fatal("analytics batch did not acquire lifecycle locks")
			}
			lifecycleContext := <-linkContext
			require.Equal(t, link.ID(), lifecycleContext.LinkID)
			require.Equal(t, resumeID, lifecycleContext.ResumeID)
			require.True(t, lifecycleContext.LinkActive)
			require.True(t, lifecycleContext.ResumePublished)

			conn, err := pgx.Connect(ctx, dsn)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
			var lockedID string
			lockErr := conn.QueryRow(ctx,
				fmt.Sprintf("SELECT id FROM %s WHERE id = $1 FOR UPDATE NOWAIT", test.table),
				test.getID(link.ID(), resumeID).UUID(),
			).Scan(&lockedID)
			var postgresErr *pgconn.PgError
			require.ErrorAs(t, lockErr, &postgresErr)
			require.Equal(t, "55P03", postgresErr.Code, "lifecycle row must stay locked through the batch transaction")

			unblock()
			select {
			case err := <-batchResult:
				require.NoError(t, err)
			case <-time.After(5 * time.Second):
				t.Fatal("analytics batch transaction did not finish")
			}

			_, err = conn.Exec(ctx, test.updateSQL, test.getID(link.ID(), resumeID).UUID())
			require.NoError(t, err)
			err = uow.Do(ctx, func(txCtx context.Context) error {
				return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), analyticsdomain.ViewStateUpdate{
					Seq: 2, SeenAt: time.Now().UTC(),
				}, nil, nil)
			})
			require.ErrorIs(t, err, test.expectedError)
		})
	}
}

func TestDeletingResumeLinkCascadesAnalyticsHistory(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Delete history"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	view := newAnalyticsView(t, link)
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	event, err := analyticsdomain.NewAnalyticsEvent(analyticsdomain.AnalyticsEventParams{
		ViewID: view.ID(), LinkID: link.ID(), Type: analyticsdomain.EventTypeDownload, OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.RecordEvents(txCtx, []*analyticsdomain.AnalyticsEvent{event})
	}))
	hash, ok := view.Layout().HashFor(analyticsdomain.AnalyticsSectionAbout)
	require.True(t, ok)
	click, err := analyticsdomain.NewAnalyticsClick(analyticsdomain.AnalyticsClickParams{
		ViewID: view.ID(), LinkID: link.ID(), Section: analyticsdomain.AnalyticsSectionAbout,
		XPct: 100, YPct: 100, ViewportW: 1280, LayoutHash: hash, OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.RecordClicks(txCtx, []*analyticsdomain.AnalyticsClick{click})
	}))
	require.NoError(t, gormDB.Exec(
		"INSERT INTO analytics.analytics_visitors (link_id, visitor_id, first_seen_at, first_view_id) VALUES (?, ?, ?, ?)",
		link.ID().UUID(), sharedvalueobjects.NewID().UUID(), time.Now().UTC(), view.ID().UUID(),
	).Error)

	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		if err := repo.DeleteByLinkIDs(txCtx, []sharedvalueobjects.ID{link.ID()}); err != nil {
			return err
		}
		return repo.Delete(txCtx, link.ID())
	}))

	var views, events, clicks, visitors int64
	require.NoError(t, gormDB.Model(&analyticsgormmodel.AnalyticsView{}).Where("link_id = ?", link.ID().UUID()).Count(&views).Error)
	require.NoError(t, gormDB.Model(&analyticsgormmodel.AnalyticsEvent{}).Where("link_id = ?", link.ID().UUID()).Count(&events).Error)
	require.NoError(t, gormDB.Model(&analyticsgormmodel.AnalyticsClick{}).Where("link_id = ?", link.ID().UUID()).Count(&clicks).Error)
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM analytics.analytics_visitors WHERE link_id = ?", link.ID().UUID()).Scan(&visitors).Error)
	require.Zero(t, views)
	require.Zero(t, events)
	require.Zero(t, clicks)
	require.Zero(t, visitors)
}

func TestAnalyticsPrintEventsAreDeduplicatedPerBatchAndAcrossBatches(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Print deduplication"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
	view := newAnalyticsView(t, link)
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))

	newPrintEvent := func() *analyticsdomain.AnalyticsEvent {
		t.Helper()
		clientEventID := sharedvalueobjects.NewID()
		event, eventErr := analyticsdomain.NewAnalyticsEvent(analyticsdomain.AnalyticsEventParams{
			ViewID: view.ID(), LinkID: link.ID(), ClientEventID: &clientEventID,
			Type: analyticsdomain.EventTypePrint, OccurredAt: time.Now().UTC(),
		})
		require.NoError(t, eventErr)
		return event
	}

	apply := func(seq int64, events ...*analyticsdomain.AnalyticsEvent) {
		t.Helper()
		err := uow.Do(ctx, func(txCtx context.Context) error {
			return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), analyticsdomain.ViewStateUpdate{
				Seq: seq, SeenAt: time.Now().UTC(),
			}, events, nil)
		})
		require.NoError(t, err)
	}
	apply(1, newPrintEvent(), newPrintEvent())

	var count int64
	require.NoError(t, gormDB.Model(&analyticsgormmodel.AnalyticsEvent{}).
		Where("view_id = ? AND event_type = ?", view.ID().UUID(), string(analyticsdomain.EventTypePrint)).Count(&count).Error)
	require.EqualValues(t, 1, count, "only one print should be stored from a single batch")

	apply(2, newPrintEvent())
	require.NoError(t, gormDB.Model(&analyticsgormmodel.AnalyticsEvent{}).
		Where("view_id = ? AND event_type = ?", view.ID().UUID(), string(analyticsdomain.EventTypePrint)).Count(&count).Error)
	require.EqualValues(t, 1, count, "a print inside the two-second window should be deduplicated across batches")
}

func TestAnalyticsBatchEnforcesPerViewEventAndClickCaps(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Analytics caps",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
	view := newAnalyticsView(t, link)
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))

	now := time.Now().UTC()
	seedEvents := make([]analyticsgormmodel.AnalyticsEvent, analyticsdomain.MaxAnalyticsEventsPerView-1)
	for i := range seedEvents {
		clientID := sharedvalueobjects.NewID().UUID()
		seedEvents[i] = analyticsgormmodel.AnalyticsEvent{
			ID: sharedvalueobjects.NewID().UUID(), ViewID: view.ID().UUID(), LinkID: link.ID().UUID(),
			ClientEventID: &clientID, EventType: string(analyticsdomain.EventTypeDownload), Payload: []byte("{}"), OccurredAt: now,
		}
	}
	require.NoError(t, gormDB.Create(&seedEvents).Error)

	seedClicks := make([]analyticsgormmodel.AnalyticsClick, analyticsdomain.MaxAnalyticsClicksPerView-1)
	for i := range seedClicks {
		target := "button"
		seedClicks[i] = analyticsgormmodel.AnalyticsClick{
			ViewID: view.ID().UUID(), LinkID: link.ID().UUID(), Section: "about",
			XPct: 100, YPct: 200, ViewportW: 1280, LayoutHash: "0123456789abcdef",
			Target: &target, OccurredAt: now,
		}
	}
	require.NoError(t, gormDB.Create(&seedClicks).Error)

	newDownload := func() *analyticsdomain.AnalyticsEvent {
		clientID := sharedvalueobjects.NewID()
		event, eventErr := analyticsdomain.NewAnalyticsEvent(analyticsdomain.AnalyticsEventParams{
			ViewID: view.ID(), LinkID: link.ID(), ClientEventID: &clientID,
			Type: analyticsdomain.EventTypeDownload, OccurredAt: now,
		})
		require.NoError(t, eventErr)
		return event
	}
	newClick := func() *analyticsdomain.AnalyticsClick {
		click, clickErr := analyticsdomain.NewAnalyticsClick(analyticsdomain.AnalyticsClickParams{
			ViewID: view.ID(), LinkID: link.ID(),
			Section: analyticsdomain.AnalyticsSectionAbout, XPct: 100, YPct: 200, ViewportW: 1280,
			LayoutHash: "0123456789abcdef", Target: "button", OccurredAt: now,
		})
		require.NoError(t, clickErr)
		return click
	}
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), analyticsdomain.ViewStateUpdate{
			Seq: 1, SeenAt: now,
		}, []*analyticsdomain.AnalyticsEvent{newDownload(), newDownload()}, []*analyticsdomain.AnalyticsClick{newClick(), newClick()})
	}))

	var eventCount, clickCount int64
	require.NoError(t, gormDB.Model(&analyticsgormmodel.AnalyticsEvent{}).
		Where("view_id = ? AND client_event_id IS NOT NULL", view.ID().UUID()).Count(&eventCount).Error)
	require.NoError(t, gormDB.Model(&analyticsgormmodel.AnalyticsClick{}).
		Where("view_id = ?", view.ID().UUID()).Count(&clickCount).Error)
	require.EqualValues(t, analyticsdomain.MaxAnalyticsEventsPerView, eventCount)
	require.EqualValues(t, analyticsdomain.MaxAnalyticsClicksPerView, clickCount)
}

func TestAnalyticsResumeStatisticsAggregatesAllLinks(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	uow := database.NewUnitOfWork(gormDB)
	links := make([]*domain.ResumeLink, 0, 3)
	for _, name := range []string{"HH", "Telegram", "No views"} {
		link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: name})
		require.NoError(t, err)
		require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
		links = append(links, link)
	}

	first := newAnalyticsView(t, links[0])
	second := newAnalyticsView(t, links[1])
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, first))
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, second))
	startedAt := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	require.NoError(t, gormDB.Exec(`
		UPDATE analytics.analytics_views
		SET started_at = ?, confirmed_at = ?, active_ms = 12000,
		    device = 'mobile', referrer_host = 'referrer.example'
		WHERE id = ?`, startedAt, startedAt.Add(2*time.Minute), first.ID().UUID()).Error)
	require.NoError(t, gormDB.Exec(`UPDATE analytics.analytics_views SET started_at = ? WHERE id = ?`, startedAt.Add(5*time.Minute), second.ID().UUID()).Error)

	period, err := analyticsqueries.NewAnalyticsPeriod("2026-10-01", "2026-10-05", "UTC", startedAt)
	require.NoError(t, err)
	linkStats, err := reader.GetLinkStatistics(ctx, analyticsqueries.LinkAnalyticsQuery{
		LinkID: links[0].ID(),
		Period: period,
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), linkStats.Views)

	result, err := reader.GetResumeStatistics(ctx, analyticsqueries.ResumeAnalyticsQuery{ResumeID: resumeID, Period: period})
	require.NoError(t, err)
	require.Len(t, result.Links, 3)
	require.Equal(t, links[0].ID().String(), result.Links[0].ID)
	require.Equal(t, links[1].ID().String(), result.Links[1].ID)
	require.Equal(t, links[2].ID().String(), result.Links[2].ID)
	require.Equal(t, "No views", result.Links[2].Name)
	require.Equal(t, int64(1), result.Links[0].Views)
	require.Equal(t, int64(1), result.Links[0].Reads)
	require.Equal(t, int64(1), result.Links[1].Views)
	require.Zero(t, result.Links[1].Reads)
	require.Zero(t, result.Links[2].Views)
	require.Nil(t, result.Links[2].LastViewAt)
	require.Equal(t, int64(2), result.Views)
	require.Equal(t, int64(1), result.Reads)
	require.Equal(t, int64(2), result.UniqueVisitors)
	require.Equal(t, int64(12), result.MedianActiveSeconds)
	require.Equal(t, int64(1), result.Devices["mobile"])
	require.Equal(t, int64(1), result.Devices["desktop"])
	require.Len(t, result.Referrers, 2)
	referrers := make(map[string]int64, len(result.Referrers))
	for _, referrer := range result.Referrers {
		host := ""
		if referrer.Host != nil {
			host = *referrer.Host
		}
		referrers[host] = referrer.Views
	}
	require.Equal(t, int64(1), referrers["referrer.example"])
	require.Equal(t, int64(1), referrers[""])
	require.Zero(t, result.Filtered.Bot)
	require.Zero(t, result.Filtered.Owner)
	require.Zero(t, result.Filtered.Suspect)
	require.NotNil(t, result.FirstViewAt)
	require.NotNil(t, result.LastViewAt)
	require.True(t, result.FirstViewAt.Equal(startedAt))
	require.True(t, result.LastViewAt.Equal(startedAt.Add(5*time.Minute)))

	require.NoError(t, gormDB.Exec(`UPDATE resume_link.resume_links SET name = 'HH updated' WHERE id = ?`, links[0].ID().UUID()).Error)
	result, err = reader.GetResumeStatistics(ctx, analyticsqueries.ResumeAnalyticsQuery{ResumeID: resumeID, Period: period})
	require.NoError(t, err)
	require.Equal(t, "HH updated", result.Links[0].Name)
}

func TestAnalyticsReferrerStatisticsAreCappedAndDeterministic(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Referrer cap",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
	startedAt := time.Date(2026, 10, 10, 10, 0, 0, 0, time.UTC)
	for index := 0; index < 25; index++ {
		host := fmt.Sprintf("host-%02d.example", index)
		view, viewErr := analyticsdomain.NewView(analyticsdomain.ViewParams{
			LinkID: link.ID(), ResumeID: link.ResumeID(), VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
			Device: analyticsdomain.DeviceDesktop, Browser: analyticsdomain.BrowserChrome, OS: analyticsdomain.OSLinux, Traffic: analyticsdomain.TrafficHuman,
			ReferrerHost: host, Layout: newAnalyticsView(t, link).Layout(), StartedAt: startedAt.Add(time.Duration(index) * time.Minute),
		})
		require.NoError(t, viewErr)
		require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	}
	period, err := analyticsqueries.NewAnalyticsPeriod("2026-10-10", "2026-10-11", "UTC", startedAt)
	require.NoError(t, err)
	stats, err := reader.GetLinkStatistics(ctx, analyticsqueries.LinkAnalyticsQuery{LinkID: link.ID(), Period: period})
	require.NoError(t, err)
	require.Len(t, stats.Referrers, 20)
	require.Equal(t, "host-00.example", *stats.Referrers[0].Host)
	require.Equal(t, "host-19.example", *stats.Referrers[19].Host)
}

func TestAnalyticsReadMetricsUseConfirmationState(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Confirmation state",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	startedAt := time.Date(2026, 10, 10, 10, 0, 0, 0, time.UTC)
	unconfirmed := newAnalyticsView(t, link)
	confirmed := newAnalyticsView(t, link)
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, unconfirmed))
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, confirmed))
	require.NoError(t, gormDB.Exec(
		`UPDATE analytics.analytics_views SET started_at = ?, active_ms = 9000, confirmed_at = NULL WHERE id = ?`,
		startedAt, unconfirmed.ID().UUID(),
	).Error)
	confirmedAt := startedAt.Add(time.Second)
	require.NoError(t, gormDB.Exec(
		`UPDATE analytics.analytics_views SET started_at = ?, active_ms = 0, confirmed_at = ? WHERE id = ?`,
		startedAt.Add(time.Minute), confirmedAt, confirmed.ID().UUID(),
	).Error)
	period, err := analyticsqueries.NewAnalyticsPeriod("2026-10-10", "2026-10-11", "UTC", startedAt)
	require.NoError(t, err)

	stats, err := reader.GetLinkStatistics(ctx, analyticsqueries.LinkAnalyticsQuery{LinkID: link.ID(), Period: period})
	require.NoError(t, err)
	require.EqualValues(t, 2, stats.Views)
	require.EqualValues(t, 1, stats.Reads)
	require.Zero(t, stats.MedianActiveSeconds)

	daily, err := reader.GetDaily(ctx, analyticsqueries.LinkAnalyticsQuery{LinkID: link.ID(), Period: period})
	require.NoError(t, err)
	require.EqualValues(t, 1, daily.Days[0].Reads)

	sections, err := reader.GetSections(ctx, analyticsqueries.LinkSectionsQuery{
		LinkID: link.ID(), Period: period,
		SectionOrder: []string{"about", "experience", "projects"},
	})
	require.NoError(t, err)
	require.EqualValues(t, 1, sections.Reads)
	require.NotNil(t, sections.Sections)

	visits, err := reader.GetVisits(ctx, analyticsqueries.VisitsReadQuery{LinkID: link.ID(), Limit: 10})
	require.NoError(t, err)
	readByID := make(map[string]bool, len(visits.Items))
	for _, visit := range visits.Items {
		readByID[visit.ID] = visit.IsRead
	}
	require.False(t, readByID[unconfirmed.ID().String()])
	require.True(t, readByID[confirmed.ID().String()])
}

func TestAnalyticsReturningVisitorMetricsKeepTheirCalendarDaySemantics(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

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
	first := newAnalyticsViewForVisitor(t, links[0], visitorID, base)
	sameDayReturn := newAnalyticsViewForVisitor(t, links[0], visitorID, base.Add(30*time.Minute))
	nextDayOtherLink := newAnalyticsViewForVisitor(t, links[1], visitorID, base.Add(24*time.Hour))
	for _, view := range []*analyticsdomain.View{first, sameDayReturn, nextDayOtherLink} {
		require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	}

	period, err := analyticsqueries.NewAnalyticsPeriod("2026-10-01", "2026-10-03", "UTC", base)
	require.NoError(t, err)

	visits, err := reader.GetVisits(ctx, analyticsqueries.VisitsReadQuery{LinkID: links[0].ID(), Limit: 10})
	require.NoError(t, err)
	require.Len(t, visits.Items, 2)
	require.Equal(t, sameDayReturn.ID().String(), visits.Items[0].ID)
	require.True(t, visits.Items[0].Visitor.IsReturning)
	require.False(t, visits.Items[1].Visitor.IsReturning)

	linkStats, err := reader.GetLinkStatistics(ctx, analyticsqueries.LinkAnalyticsQuery{LinkID: links[0].ID(), Period: period})
	require.NoError(t, err)
	require.EqualValues(t, 0, linkStats.ReturningVisitors)

	resumeStats, err := reader.GetResumeStatistics(ctx, analyticsqueries.ResumeAnalyticsQuery{ResumeID: resumeID, Period: period})
	require.NoError(t, err)
	require.EqualValues(t, 1, resumeStats.ReturningVisitors)

	daily, err := reader.GetDaily(ctx, analyticsqueries.LinkAnalyticsQuery{LinkID: links[0].ID(), Period: period})
	require.NoError(t, err)
	require.Len(t, daily.Days, 3)
	require.EqualValues(t, 1, daily.Days[0].UniqueVisitors)
	require.EqualValues(t, 0, daily.Days[1].UniqueVisitors)
}

func TestAnalyticsVisitsPaginationUsesBothLinkIDBindings(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Visits"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	first := newAnalyticsView(t, link)
	second := newAnalyticsView(t, link)
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, first))
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, second))

	firstStartedAt := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	secondStartedAt := firstStartedAt.Add(5 * time.Minute)
	require.NoError(t, gormDB.Exec(`UPDATE analytics.analytics_views SET started_at = ? WHERE id = ?`, firstStartedAt, first.ID().UUID()).Error)
	require.NoError(t, gormDB.Exec(`UPDATE analytics.analytics_views SET started_at = ? WHERE id = ?`, secondStartedAt, second.ID().UUID()).Error)

	page, err := reader.GetVisits(ctx, analyticsqueries.VisitsReadQuery{LinkID: link.ID(), Limit: 1})
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	require.Equal(t, second.ID().String(), page.Items[0].ID)
	require.NotNil(t, page.Items[0].Events)
	require.Empty(t, page.Items[0].Events)
	require.NotNil(t, page.NextCursor)

	cursorTime, cursorID, err := analyticsqueries.DecodeVisitCursor(*page.NextCursor)
	require.NoError(t, err)
	nextPage, err := reader.GetVisits(ctx, analyticsqueries.VisitsReadQuery{
		LinkID: link.ID(), Limit: 1, CursorTime: cursorTime, CursorID: cursorID,
	})
	require.NoError(t, err)
	require.Len(t, nextPage.Items, 1)
	require.Equal(t, first.ID().String(), nextPage.Items[0].ID)
	require.Nil(t, nextPage.NextCursor)
}

func TestAnalyticsVisitorOrdinalAndFirstSeenSurviveRetention(t *testing.T) {
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
	firstA := newAnalyticsViewForVisitor(t, link, visitorA, base.Add(time.Minute))
	firstB := newAnalyticsViewForVisitor(t, link, visitorB, base)
	secondA := newAnalyticsViewForVisitor(t, link, visitorA, base.Add(2*time.Minute))

	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, firstA))
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, firstB))
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, secondA))

	page, err := reader.GetVisits(ctx, analyticsqueries.VisitsReadQuery{LinkID: link.ID(), Limit: 10})
	require.NoError(t, err)
	require.Len(t, page.Items, 3)
	require.NotNil(t, page.Items[0].Events, "visits without events must expose an empty array")
	require.Empty(t, page.Items[0].Events)
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
	remaining, err := reader.GetVisits(ctx, analyticsqueries.VisitsReadQuery{LinkID: link.ID(), Limit: 10})
	require.NoError(t, err)
	require.Len(t, remaining.Items, 1)
	require.Equal(t, secondA.ID().String(), remaining.Items[0].ID)
	require.Equal(t, int64(2), remaining.Items[0].Visitor.Ordinal)
	require.True(t, remaining.Items[0].Visitor.IsReturning)
}

func TestAnalyticsConcurrentVisitorRegistrationKeepsDenseOrdinal(t *testing.T) {
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
	views := make([]*analyticsdomain.View, 0, visitorCount)
	for i := visitorCount - 1; i >= 0; i-- {
		views = append(views, newAnalyticsViewForVisitor(
			t, link, sharedvalueobjects.NewID(), base.Add(time.Duration(i)*time.Minute),
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

	page, err := reader.GetVisits(ctx, analyticsqueries.VisitsReadQuery{LinkID: link.ID(), Limit: visitorCount})
	require.NoError(t, err)
	require.Len(t, page.Items, visitorCount)
	for i, visit := range page.Items {
		require.Equal(t, int64(visitorCount-i), visit.Visitor.Ordinal)
		require.False(t, visit.Visitor.IsReturning)
	}
}

func newAnalyticsViewForVisitor(
	t *testing.T,
	link *domain.ResumeLink,
	visitorID sharedvalueobjects.ID,
	startedAt time.Time,
) *analyticsdomain.View {
	t.Helper()
	layout, err := analyticsdomain.NewLayoutSnapshot(
		[]analyticsdomain.AnalyticsSectionID{analyticsdomain.AnalyticsSectionAbout, analyticsdomain.AnalyticsSectionExperience, analyticsdomain.AnalyticsSectionProjects},
		map[analyticsdomain.AnalyticsSectionID]string{
			analyticsdomain.AnalyticsSectionAbout:      "0123456789abcdef",
			analyticsdomain.AnalyticsSectionExperience: "fedcba9876543210",
			analyticsdomain.AnalyticsSectionProjects:   "0011223344556677",
		},
	)
	require.NoError(t, err)
	view, err := analyticsdomain.NewView(analyticsdomain.ViewParams{
		LinkID: link.ID(), ResumeID: link.ResumeID(), VisitorID: visitorID, IPHash: make([]byte, 32),
		Device: analyticsdomain.DeviceDesktop, Browser: analyticsdomain.BrowserChrome, OS: analyticsdomain.OSLinux, Traffic: analyticsdomain.TrafficHuman,
		Layout: layout, StartedAt: startedAt, CreatedAt: startedAt,
	})
	require.NoError(t, err)
	return view
}

func TestAnalyticsBatchRejectsChangedPublicLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Lifecycle"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))
	view := newAnalyticsView(t, link)
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	require.NoError(t, gormDB.Exec(`UPDATE resume_link.resume_links SET status = 'inactive' WHERE id = ?`, link.ID().UUID()).Error)
	err = uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), analyticsdomain.ViewStateUpdate{Seq: 1, SeenAt: time.Now().UTC()}, nil, nil)
	})
	require.ErrorIs(t, err, analyticsdomain.ErrAnalyticsLinkInactive)

	require.NoError(t, gormDB.Exec(`UPDATE resume_link.resume_links SET status = 'active' WHERE id = ?`, link.ID().UUID()).Error)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'draft' WHERE id = ?`, resumeID.UUID()).Error)
	err = uow.Do(ctx, func(txCtx context.Context) error {
		return applyAnalyticsBatchInTransaction(txCtx, repo, view.ID(), link.ID(), analyticsdomain.ViewStateUpdate{Seq: 2, SeenAt: time.Now().UTC()}, nil, nil)
	})
	require.ErrorIs(t, err, analyticscontracts.ErrAnalyticsSubjectNotFound)
}

func TestAnalyticsFoundationPersistence(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: resumeID, Name: "Analytics"})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return repo.Create(txCtx, link) }))

	view := newAnalyticsView(t, link)
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	loaded, err := repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.Equal(t, view.ID(), loaded.ID())
	require.Equal(t, view.LinkID(), loaded.LinkID())
	require.Equal(t, view.Layout().Order(), loaded.Layout().Order())
	require.WithinDuration(t, view.CreatedAt(), loaded.CreatedAt(), time.Microsecond)

	now := time.Now().UTC()
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.UpdateViewState(txCtx, view.ID(), analyticsdomain.ViewStateUpdate{
			Seq: 1, ActiveMS: 7000, MaxSection: analyticsdomain.AnalyticsSectionExperience, ExitSection: analyticsdomain.AnalyticsSectionExperience, SeenAt: now, Final: true,
		})
	}))
	loaded, err = repo.GetView(ctx, view.ID())
	require.NoError(t, err)
	require.Equal(t, int64(7000), loaded.State().ActiveMS())
	require.Equal(t, analyticsdomain.AnalyticsSectionExperience, loaded.State().MaxSection())
	require.False(t, loaded.State().EndedAt().IsZero())

	clientEventID := sharedvalueobjects.NewID()
	event, err := analyticsdomain.NewAnalyticsEvent(analyticsdomain.AnalyticsEventParams{
		ViewID: view.ID(), LinkID: link.ID(), ClientEventID: &clientEventID,
		Type: analyticsdomain.EventTypeSectionView, Payload: map[string]string{"section": "experience"}, OccurredAt: now,
	})
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.RecordEvents(txCtx, []*analyticsdomain.AnalyticsEvent{event, event})
	}))

	var count int
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM analytics.analytics_events WHERE view_id = ?", view.ID().UUID()).Scan(&count).Error)
	require.Equal(t, 1, count)

	require.NoError(t, gormDB.Raw("SELECT count(*) FROM analytics.analytics_visitors WHERE link_id = ?", link.ID().UUID()).Scan(&count).Error)
	require.EqualValues(t, 1, count)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		if err := repo.DeleteByLinkIDs(txCtx, []sharedvalueobjects.ID{link.ID()}); err != nil {
			return err
		}
		return repo.Delete(txCtx, link.ID())
	}))
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM analytics.analytics_views WHERE id = ?", view.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM analytics.analytics_events WHERE view_id = ?", view.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM analytics.analytics_visitors WHERE link_id = ?", link.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count, "visitor pseudonymous metadata must be deleted with its link")
}

func TestAnalyticsRetentionDeletesExpiredDataInBatches(t *testing.T) {
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

	oldView := newAnalyticsView(t, link)
	oldStartedAt := time.Now().UTC().AddDate(-1, -1, 0)
	oldViewModel := map[string]any{
		"id": oldView.ID().UUID(), "link_id": link.ID().UUID(), "resume_id": link.ResumeID().UUID(), "visitor_id": oldView.VisitorID().UUID(),
		"ip_hash": make([]byte, 32), "device": "desktop", "browser": "chrome", "os": "linux",
		"traffic": "human", "layout": `{"order":["about"],"hashes":{"about":"0123456789abcdef"}}`,
		"started_at": oldStartedAt, "active_ms": 0, "last_state_seq": 0, "created_at": oldStartedAt,
	}
	require.NoError(t, gormDB.Table("analytics.analytics_views").Create(oldViewModel).Error)

	eventID := sharedvalueobjects.NewID().UUID()
	require.NoError(t, gormDB.Exec(
		`INSERT INTO analytics.analytics_events
		 (id, view_id, link_id, client_event_id, event_type, payload, occurred_at)
		 VALUES (?, ?, ?, ?, ?, ?::jsonb, ?)`,
		eventID, oldView.ID().UUID(), link.ID().UUID(), nil,
		"download", `{}`, oldStartedAt,
	).Error)
	require.NoError(t, gormDB.Exec(
		`INSERT INTO analytics.analytics_clicks
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
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM analytics.analytics_views WHERE id = ?", oldView.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM analytics.analytics_events WHERE view_id = ?", oldView.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
	require.NoError(t, gormDB.Raw("SELECT count(*) FROM analytics.analytics_clicks WHERE view_id = ?", oldView.ID().UUID()).Scan(&count).Error)
	require.Zero(t, count)
}

func TestAnalyticsClickReaderUsesGORMFiltersAndCountsStaleClicks(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, reader := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: resumeID, Name: "Click map",
	})
	require.NoError(t, err)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.Create(txCtx, link)
	}))

	view := newAnalyticsView(t, link)
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))
	currentHash, ok := view.Layout().HashFor(analyticsdomain.AnalyticsSectionAbout)
	require.True(t, ok)
	clicks := make([]*analyticsdomain.AnalyticsClick, 0, 2)
	for _, layoutHash := range []string{currentHash, "fedcba9876543210"} {
		click, err := analyticsdomain.NewAnalyticsClick(analyticsdomain.AnalyticsClickParams{
			ViewID:     view.ID(),
			LinkID:     link.ID(),
			Section:    analyticsdomain.AnalyticsSectionAbout,
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

	currentLayout, err := reader.GetClicks(ctx, analyticsqueries.LinkClicksQuery{
		LinkID:      link.ID(),
		Section:     string(analyticsdomain.AnalyticsSectionAbout),
		Breakpoint:  "lg",
		CurrentHash: currentHash,
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), currentLayout.Total)
	require.Equal(t, int64(1), currentLayout.StaleTotal)
	require.Len(t, currentLayout.Cells, 1)
	require.Equal(t, int64(8), currentLayout.Cells[0].X)
	require.Equal(t, int64(3), currentLayout.Cells[0].Y)

	allLayouts, err := reader.GetClicks(ctx, analyticsqueries.LinkClicksQuery{
		LinkID:       link.ID(),
		Section:      string(analyticsdomain.AnalyticsSectionAbout),
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

func TestAnalyticsRejectsEventsAndClicksForAnotherLink(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, gormDB, repo, _ := newDatabase(t)
	owner := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	insertUserAndResume(t, ctx, dsn, owner, resumeID)
	require.NoError(t, gormDB.Exec(`UPDATE resume.resumes SET status = 'opened' WHERE id = ?`, resumeID.UUID()).Error)

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

	view := newAnalyticsView(t, firstLink)
	require.NoError(t, createAnalyticsView(ctx, gormDB, repo, view))

	event, err := analyticsdomain.NewAnalyticsEvent(analyticsdomain.AnalyticsEventParams{
		ViewID: view.ID(), LinkID: secondLink.ID(),
		Type: analyticsdomain.EventTypeDownload, OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	require.ErrorIs(t, repo.RecordEvents(ctx, []*analyticsdomain.AnalyticsEvent{event}), database.ErrTransactionRequired)
	require.ErrorIs(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.RecordEvents(txCtx, []*analyticsdomain.AnalyticsEvent{event})
	}), analyticscontracts.ErrAnalyticsSubjectNotFound)
	err = uow.Do(ctx, func(txCtx context.Context) error {
		snapshot, loadErr := repo.LoadAnalyticsBatchForUpdate(
			txCtx, view.ID(), firstLink.ID(), true, false, false,
		)
		if loadErr != nil {
			return loadErr
		}
		return repo.SaveAnalyticsBatch(txCtx, snapshot.View, []*analyticsdomain.AnalyticsEvent{event}, nil)
	})
	require.Error(t, err)

	layoutHash, ok := view.Layout().HashFor(analyticsdomain.AnalyticsSectionAbout)
	require.True(t, ok)
	click, err := analyticsdomain.NewAnalyticsClick(analyticsdomain.AnalyticsClickParams{
		ViewID: view.ID(), LinkID: secondLink.ID(),
		Section: analyticsdomain.AnalyticsSectionAbout, XPct: 100, YPct: 200, ViewportW: 1440,
		LayoutHash: layoutHash, OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	require.ErrorIs(t, repo.RecordClicks(ctx, []*analyticsdomain.AnalyticsClick{click}), database.ErrTransactionRequired)
	require.ErrorIs(t, uow.Do(ctx, func(txCtx context.Context) error {
		return repo.RecordClicks(txCtx, []*analyticsdomain.AnalyticsClick{click})
	}), analyticscontracts.ErrAnalyticsSubjectNotFound)

	var eventCount, clickCount int64
	require.NoError(t, gormDB.Raw(
		`SELECT count(*) FROM analytics.analytics_events WHERE view_id = ?`, view.ID().UUID(),
	).Scan(&eventCount).Error)
	require.NoError(t, gormDB.Raw(
		`SELECT count(*) FROM analytics.analytics_clicks WHERE view_id = ?`, view.ID().UUID(),
	).Scan(&clickCount).Error)
	require.Zero(t, eventCount)
	require.Zero(t, clickCount)
}

func TestAnalyticsRetentionCleanupLockIsTransactionScoped(t *testing.T) {
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
