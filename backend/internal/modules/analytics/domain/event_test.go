package domain

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"

	"github.com/stretchr/testify/require"
)

func TestNewAnalyticsEventValidatesAndNormalizesPayload(t *testing.T) {
	clientEventID := sharedvalueobjects.NewID()
	event, err := NewAnalyticsEvent(AnalyticsEventParams{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
		ClientEventID: &clientEventID, Type: EventTypeContactView,
		Payload:    map[string]string{"contact": "email", "ignored": "value", "long": "ignored"},
		OccurredAt: time.Now(),
	})
	require.NoError(t, err)
	require.Equal(t, map[string]string{"contact": "email"}, event.Payload())

	_, err = NewAnalyticsEvent(AnalyticsEventParams{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
		Type: AnalyticsEventType("view"), OccurredAt: time.Now(),
	})
	require.ErrorIs(t, err, ErrInvalidAnalyticsEvent)
}

func TestAnalyticsEventCopiesPayload(t *testing.T) {
	payload := map[string]string{"contact": "email"}
	event, err := NewAnalyticsEvent(AnalyticsEventParams{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
		Type: EventTypeContactView, Payload: payload, OccurredAt: time.Now(),
	})
	require.NoError(t, err)
	require.False(t, event.ID().IsZero())
	payload["contact"] = "telegram"
	require.Equal(t, "email", event.Payload()["contact"])
}

func TestAnalyticsEventsWithoutPayloadRejectNonEmptyPayload(t *testing.T) {
	for _, eventType := range []AnalyticsEventType{EventTypeDownload, EventTypePrint} {
		_, err := NewAnalyticsEvent(AnalyticsEventParams{
			ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
			Type: eventType, Payload: map[string]string{"unexpected": "value"}, OccurredAt: time.Now(),
		})
		require.ErrorIs(t, err, ErrInvalidAnalyticsEvent)
	}
}

func TestAnalyticsEventPayloadRemovesControlCharactersAndTruncatesByRunes(t *testing.T) {
	value := strings.Repeat("я", 70) + "\x00\n"
	event, err := NewAnalyticsEvent(AnalyticsEventParams{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
		Type: EventTypeExternalClick,
		Payload: map[string]string{
			"kind":   "external\x00-link\n",
			"target": value,
		},
		OccurredAt: time.Now(),
	})
	require.NoError(t, err)
	payload := event.Payload()
	require.Equal(t, "external-link", payload["kind"])
	require.Equal(t, 64, utf8.RuneCountInString(payload["target"]))
	require.Equal(t, strings.Repeat("я", 64), payload["target"])
	require.True(t, utf8.ValidString(payload["target"]))
}
