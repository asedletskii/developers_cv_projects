package mappers

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
)

func AnalyticsEventToDB(event *domain.AnalyticsEvent) (gormmodel.ResumeLinkAnalyticsEvent, error) {
	if event == nil {
		return gormmodel.ResumeLinkAnalyticsEvent{}, fmt.Errorf("analytics event is nil")
	}
	payload, err := json.Marshal(event.Payload())
	if err != nil {
		return gormmodel.ResumeLinkAnalyticsEvent{}, fmt.Errorf("marshal event payload: %w", err)
	}
	var clientEventID *uuid.UUID
	if id := event.ClientEventID(); id != nil {
		value := id.UUID()
		clientEventID = &value
	}
	return gormmodel.ResumeLinkAnalyticsEvent{
		ID:            event.ID().UUID(),
		ViewID:        event.ViewID().UUID(),
		LinkID:        event.LinkID().UUID(),
		ClientEventID: clientEventID,
		EventType:     string(event.Type()),
		Payload:       payload,
		OccurredAt:    event.OccurredAt(),
	}, nil
}
