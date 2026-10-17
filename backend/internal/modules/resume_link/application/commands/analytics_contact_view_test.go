package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestRecordContactViewEventPersistsOnlyForMatchingViewAndLink(t *testing.T) {
	ctx := context.Background()
	viewID := sharedvalueobjects.NewID()
	linkID := sharedvalueobjects.NewID()
	view := newCommandTestView(t, viewID, linkID)
	events := &analyticsEventRepositoryStub{}
	uow := &unitOfWorkStub{}
	handler := NewRecordPublicInteractionEventHandler(events, &viewRepositoryStub{view: view}, uow)

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
		&analyticsEventRepositoryStub{},
		&viewRepositoryStub{view: newCommandTestView(t, viewID, linkID)},
		&unitOfWorkStub{},
	)

	err := handler.RecordContactView(context.Background(), viewID, linkID, "not-a-contact")
	require.ErrorIs(t, err, domain.ErrInvalidAnalyticsEvent)
}
