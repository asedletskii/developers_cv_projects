package commands

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type contactViewEventRepositoryStub struct {
	events []*domain.AnalyticsEvent
	err    error
}

func (s *contactViewEventRepositoryStub) RecordEvents(_ context.Context, events []*domain.AnalyticsEvent) error {
	s.events = append(s.events, events...)
	return s.err
}

type contactViewEventViewsStub struct {
	view *domain.View
	err  error
}

func (s contactViewEventViewsStub) CreateView(context.Context, *domain.View) error { return nil }
func (s contactViewEventViewsStub) UpdateViewState(context.Context, sharedvalueobjects.ID, domain.ViewStateUpdate) error {
	return nil
}
func (s contactViewEventViewsStub) GetView(context.Context, sharedvalueobjects.ID) (*domain.View, error) {
	return s.view, s.err
}

func TestRecordContactViewEventPersistsOnlyForMatchingViewAndLink(t *testing.T) {
	ctx := context.Background()
	viewID := sharedvalueobjects.NewID()
	linkID := sharedvalueobjects.NewID()
	view := newContactViewEventTestView(t, viewID, linkID)
	events := &contactViewEventRepositoryStub{}
	uow := &fakeUnitOfWork{}
	handler := NewRecordPublicInteractionEventHandler(events, contactViewEventViewsStub{view: view}, uow)

	err := handler.RecordContactView(ctx, viewID, linkID, "email")
	require.NoError(t, err)
	require.True(t, uow.called)
	require.Len(t, events.events, 1)
	require.Equal(t, domain.EventTypeContactView, events.events[0].Type())
	require.Equal(t, map[string]string{"contact": "email"}, events.events[0].Payload())

	err = handler.RecordDownload(ctx, viewID, linkID)
	require.NoError(t, err)
	require.Len(t, events.events, 2)
	require.Equal(t, domain.EventTypeDownload, events.events[1].Type())
	require.Empty(t, events.events[1].Payload())

	err = handler.RecordContactView(ctx, viewID, sharedvalueobjects.NewID(), "email")
	require.NoError(t, err)
	require.Len(t, events.events, 2)
}

func TestRecordContactViewEventRejectsInvalidContact(t *testing.T) {
	viewID := sharedvalueobjects.NewID()
	linkID := sharedvalueobjects.NewID()
	handler := NewRecordPublicInteractionEventHandler(
		&contactViewEventRepositoryStub{},
		contactViewEventViewsStub{view: newContactViewEventTestView(t, viewID, linkID)},
		&fakeUnitOfWork{},
	)

	err := handler.RecordContactView(context.Background(), viewID, linkID, "not-a-contact")
	require.ErrorIs(t, err, domain.ErrInvalidAnalyticsEvent)
}

func newContactViewEventTestView(t *testing.T, viewID, linkID sharedvalueobjects.ID) *domain.View {
	t.Helper()
	layout, err := domain.NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionAbout},
		map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"},
	)
	require.NoError(t, err)
	view, err := domain.RestoreView(domain.ViewSnapshot{
		ID: viewID, LinkID: linkID, VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
		Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux,
		Traffic: domain.TrafficHuman, Layout: layout, StartedAt: time.Now().UTC(), CreatedAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	return view
}
