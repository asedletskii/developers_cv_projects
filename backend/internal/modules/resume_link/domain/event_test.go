package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/shared/domain/valueobjects"
)

func validResumeLinkEventParams() ResumeLinkEventParams {
	return ResumeLinkEventParams{
		ID:         valueobjects.NewID(),
		LinkID:     valueobjects.NewID(),
		Type:       EventTypeView,
		OccurredAt: time.Date(2026, 9, 28, 10, 0, 0, 0, time.FixedZone("test", 2*60*60)),
	}
}

func TestNewResumeLinkEvent(t *testing.T) {
	for _, eventType := range []EventType{EventTypeView, EventTypeContactView} {
		params := validResumeLinkEventParams()
		params.Type = eventType

		event, err := NewResumeLinkEvent(params)
		require.NoError(t, err)
		require.Equal(t, params.ID, event.ID())
		require.Equal(t, params.LinkID, event.LinkID())
		require.Equal(t, eventType, event.Type())
		require.Equal(t, params.OccurredAt.UTC(), event.OccurredAt())
	}
}

func TestNewResumeLinkEventRejectsUnknownType(t *testing.T) {
	params := validResumeLinkEventParams()
	params.Type = EventType("unknown")

	_, err := NewResumeLinkEvent(params)
	require.ErrorIs(t, err, ErrInvalidResumeLinkEvent)
}

func TestNewResumeLinkEventRejectsZeroIdentity(t *testing.T) {
	params := validResumeLinkEventParams()
	params.ID = valueobjects.ID{}
	_, err := NewResumeLinkEvent(params)
	require.ErrorIs(t, err, ErrInvalidResumeLinkEvent)

	params = validResumeLinkEventParams()
	params.LinkID = valueobjects.ID{}
	_, err = NewResumeLinkEvent(params)
	require.ErrorIs(t, err, ErrInvalidResumeLinkEvent)
}

func TestNewResumeLinkEventRejectsZeroTimestamp(t *testing.T) {
	params := validResumeLinkEventParams()
	params.OccurredAt = time.Time{}

	_, err := NewResumeLinkEvent(params)
	require.ErrorIs(t, err, ErrInvalidResumeLinkEvent)
}
