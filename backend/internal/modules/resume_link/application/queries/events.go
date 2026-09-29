package queries

import (
	"context"
	"time"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/domain/valueobjects"
)

func recordLinkEvent(ctx context.Context, recorder ports.LinkEventRecorder, linkID valueobjects.ID, eventType domain.EventType) {
	event, err := domain.NewResumeLinkEvent(domain.ResumeLinkEventParams{
		ID:         valueobjects.NewID(),
		LinkID:     linkID,
		Type:       eventType,
		OccurredAt: time.Now().UTC(),
	})
	if err != nil {
		return
	}
	if err = recorder.Record(ctx, event); err != nil {
		return
	}
}
