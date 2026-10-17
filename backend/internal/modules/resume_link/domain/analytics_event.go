package domain

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const MaxAnalyticsEventsPerView = 1000

type AnalyticsEventType string

const (
	EventTypeContactView   AnalyticsEventType = "contact_view"
	EventTypeSectionView   AnalyticsEventType = "section_view"
	EventTypeDownload      AnalyticsEventType = "download"
	EventTypeTechFilter    AnalyticsEventType = "tech_filter"
	EventTypeExternalClick AnalyticsEventType = "external_click"
	EventTypePrint         AnalyticsEventType = "print"
	EventTypeContactCopy   AnalyticsEventType = "contact_copy"
	EventTypeContactClick  AnalyticsEventType = "contact_click"
)

func (v AnalyticsEventType) IsValid() bool {
	switch v {
	case EventTypeContactView, EventTypeSectionView,
		EventTypeDownload, EventTypeTechFilter,
		EventTypeExternalClick, EventTypePrint,
		EventTypeContactCopy, EventTypeContactClick:
		return true
	default:
		return false
	}
}

type AnalyticsEventParams struct {
	ViewID        sharedvalueobjects.ID
	LinkID        sharedvalueobjects.ID
	ClientEventID *sharedvalueobjects.ID
	Type          AnalyticsEventType
	Payload       map[string]string
	OccurredAt    time.Time
}

type AnalyticsEvent struct {
	id            sharedvalueobjects.ID
	viewID        sharedvalueobjects.ID
	linkID        sharedvalueobjects.ID
	clientEventID *sharedvalueobjects.ID
	eventType     AnalyticsEventType
	payload       map[string]string
	occurredAt    time.Time
}

func NewAnalyticsEvent(params AnalyticsEventParams) (*AnalyticsEvent, error) {
	if params.ViewID.IsZero() || params.LinkID.IsZero() || !params.Type.IsValid() ||
		params.OccurredAt.IsZero() {
		return nil, fmt.Errorf("%w: invalid identity, type or timestamp", ErrInvalidAnalyticsEvent)
	}
	payload, err := normalizeEventPayload(params.Type, params.Payload)
	if err != nil {
		return nil, err
	}
	var clientEventID *sharedvalueobjects.ID
	if params.ClientEventID != nil {
		if params.ClientEventID.IsZero() {
			return nil, fmt.Errorf("%w: client event id is empty", ErrInvalidAnalyticsEvent)
		}
		id := *params.ClientEventID
		clientEventID = &id
	}
	return &AnalyticsEvent{
		id:            sharedvalueobjects.NewID(),
		viewID:        params.ViewID,
		linkID:        params.LinkID,
		clientEventID: clientEventID,
		eventType:     params.Type,
		payload:       payload,
		occurredAt:    params.OccurredAt.UTC(),
	}, nil
}

func normalizeEventPayload(
	eventType AnalyticsEventType,
	input map[string]string,
) (map[string]string, error) {
	if (eventType == EventTypeDownload || eventType == EventTypePrint) && len(input) > 0 {
		return nil, fmt.Errorf("%w: payload is not allowed", ErrInvalidAnalyticsEvent)
	}
	allowed := map[AnalyticsEventType]map[string]bool{
		EventTypeContactView:   {"contact": true},
		EventTypeSectionView:   {"section": true},
		EventTypeDownload:      {},
		EventTypeTechFilter:    {"technology": true, "action": true},
		EventTypeExternalClick: {"kind": true, "target": true},
		EventTypePrint:         {},
		EventTypeContactCopy:   {"contact": true},
		EventTypeContactClick:  {"contact": true},
	}[eventType]
	result := make(map[string]string)
	for key, value := range input {
		if allowed[key] {
			result[key] = truncateEventValue(value)
		}
	}
	switch eventType {
	case EventTypeContactView, EventTypeContactCopy, EventTypeContactClick:
		if _, err := resumedomain.ParseContactKey(result["contact"]); err != nil {
			return nil, fmt.Errorf("%w: invalid contact", ErrInvalidAnalyticsEvent)
		}
	case EventTypeSectionView:
		section, err := resumedomain.ParseSectionID(result["section"])
		if err != nil || section == resumedomain.SectionTechStack {
			return nil, fmt.Errorf("%w: invalid section", ErrInvalidAnalyticsEvent)
		}
	case EventTypeTechFilter:
		if strings.TrimSpace(result["technology"]) == "" ||
			(result["action"] != "select" && result["action"] != "clear") {
			return nil, fmt.Errorf("%w: invalid technology filter", ErrInvalidAnalyticsEvent)
		}
	case EventTypeExternalClick:
		if strings.TrimSpace(result["kind"]) == "" || strings.TrimSpace(result["target"]) == "" {
			return nil, fmt.Errorf("%w: invalid external click", ErrInvalidAnalyticsEvent)
		}
	case EventTypeDownload, EventTypePrint:
		if len(result) != 0 {
			return nil, fmt.Errorf("%w: payload is not allowed", ErrInvalidAnalyticsEvent)
		}
	}
	return result, nil
}

const maxAnalyticsEventValueRunes = 64

func truncateEventValue(value string) string {
	value = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, value)
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) > maxAnalyticsEventValueRunes {
		runes = runes[:maxAnalyticsEventValueRunes]
	}
	return string(runes)
}

func (e *AnalyticsEvent) ID() sharedvalueobjects.ID     { return e.id }
func (e *AnalyticsEvent) ViewID() sharedvalueobjects.ID { return e.viewID }
func (e *AnalyticsEvent) LinkID() sharedvalueobjects.ID { return e.linkID }
func (e *AnalyticsEvent) ClientEventID() *sharedvalueobjects.ID {
	if e.clientEventID == nil {
		return nil
	}
	id := *e.clientEventID
	return &id
}
func (e *AnalyticsEvent) Type() AnalyticsEventType { return e.eventType }
func (e *AnalyticsEvent) Payload() map[string]string {
	result := make(map[string]string, len(e.payload))
	for k, v := range e.payload {
		result[k] = v
	}
	return result
}
func (e *AnalyticsEvent) OccurredAt() time.Time { return e.occurredAt }
