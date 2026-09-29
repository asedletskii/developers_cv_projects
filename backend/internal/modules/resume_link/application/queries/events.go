package queries

import (
	"context"
	"time"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/domain/valueobjects"
)

const linkEventRecordTimeout = 2 * time.Second

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

	analyticsCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), linkEventRecordTimeout)
	defer cancel()
	if err := recorder.Record(analyticsCtx, event); err != nil {
		// Infrastructure recorder is responsible for logging persistence failures.
		return
	}
}
