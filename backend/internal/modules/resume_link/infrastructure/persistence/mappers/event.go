package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
)

func ResumeLinkEventToDB(event *domain.ResumeLinkEvent) (gormmodel.ResumeLinkEvent, error) {
	if event == nil {
		return gormmodel.ResumeLinkEvent{}, fmt.Errorf("resume link event is nil")
	}
	if event.ID().IsZero() || event.LinkID().IsZero() || !event.Type().IsValid() || event.OccurredAt().IsZero() {
		return gormmodel.ResumeLinkEvent{}, fmt.Errorf("invalid resume link event")
	}
	return gormmodel.ResumeLinkEvent{
		ID: event.ID().UUID(), LinkID: event.LinkID().UUID(), EventType: string(event.Type()), OccurredAt: event.OccurredAt().UTC(),
	}, nil
}
