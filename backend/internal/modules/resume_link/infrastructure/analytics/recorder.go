package analytics

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
)

type Recorder struct {
	repository ports.ResumeLinkRepository
	logger     *zap.Logger
}

func NewRecorder(repository ports.ResumeLinkRepository, logger *zap.Logger) *Recorder {
	return &Recorder{repository: repository, logger: logger}
}

func (r *Recorder) Record(ctx context.Context, event *domain.ResumeLinkEvent) error {
	if r == nil || r.repository == nil {
		return fmt.Errorf("record resume link analytics: recorder dependencies are nil")
	}
	if event == nil {
		return fmt.Errorf("record resume link analytics: event is nil")
	}

	if err := r.repository.RecordEvent(ctx, event); err != nil {
		if r.logger != nil {
			r.logger.Error(
				"resume link analytics event recording failed",
				zap.String("link_id", event.LinkID().String()),
				zap.String("event_type", string(event.Type())),
				zap.Error(err),
			)
		}
		return err
	}
	return nil
}
