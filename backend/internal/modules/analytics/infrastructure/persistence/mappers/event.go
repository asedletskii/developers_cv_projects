package mappers

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/analytics/domain"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
)

func EventToDB(event *domain.AnalyticsEvent) (gormmodel.AnalyticsEvent, error) {
	if event == nil {
		return gormmodel.AnalyticsEvent{}, fmt.Errorf("event is nil")
	}
	payload, err := json.Marshal(event.Payload())
	if err != nil {
		return gormmodel.AnalyticsEvent{}, fmt.Errorf("marshal event payload: %w", err)
	}
	var clientEventID *uuid.UUID
	if id := event.ClientEventID(); id != nil {
		value := id.UUID()
		clientEventID = &value
	}
	return gormmodel.AnalyticsEvent{
		ID:            event.ID().UUID(),
		ViewID:        event.ViewID().UUID(),
		LinkID:        event.LinkID().UUID(),
		ClientEventID: clientEventID,
		EventType:     string(event.Type()),
		Payload:       payload,
		OccurredAt:    event.OccurredAt(),
	}, nil
}
