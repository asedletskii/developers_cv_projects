package domain

import (
	"fmt"
	"time"

	"developers-cv/backend/internal/shared/domain/valueobjects"
)

type EventType string

const (
	EventTypeView        EventType = "view"
	EventTypeContactView EventType = "contact_view"
)

func (v EventType) IsValid() bool {
	switch v {
	case EventTypeView, EventTypeContactView:
		return true
	default:
		return false
	}
}

type ResumeLinkEventParams struct {
	ID         valueobjects.ID
	LinkID     valueobjects.ID
	Type       EventType
	OccurredAt time.Time
}

type ResumeLinkEvent struct {
	id         valueobjects.ID
	linkID     valueobjects.ID
	typeValue  EventType
	occurredAt time.Time
}

func NewResumeLinkEvent(params ResumeLinkEventParams) (*ResumeLinkEvent, error) {
	if params.ID.IsZero() || params.LinkID.IsZero() {
		return nil, fmt.Errorf("%w: identity is empty", ErrInvalidResumeLinkEvent)
	}
	if !params.Type.IsValid() {
		return nil, fmt.Errorf("%w: invalid event type %q", ErrInvalidResumeLinkEvent, params.Type)
	}
	if params.OccurredAt.IsZero() {
		return nil, fmt.Errorf("%w: occurred_at is empty", ErrInvalidResumeLinkEvent)
	}

	return &ResumeLinkEvent{
		id:         params.ID,
		linkID:     params.LinkID,
		typeValue:  params.Type,
		occurredAt: params.OccurredAt.UTC(),
	}, nil
}

func (e *ResumeLinkEvent) ID() valueobjects.ID {
	if e == nil {
		return valueobjects.ID{}
	}
	return e.id
}

func (e *ResumeLinkEvent) LinkID() valueobjects.ID {
	if e == nil {
		return valueobjects.ID{}
	}
	return e.linkID
}

func (e *ResumeLinkEvent) Type() EventType {
	if e == nil {
		return ""
	}
	return e.typeValue
}

func (e *ResumeLinkEvent) OccurredAt() time.Time {
	if e == nil {
		return time.Time{}
	}
	return e.occurredAt
}
