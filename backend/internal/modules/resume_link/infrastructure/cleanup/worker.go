package cleanup

import (
	"context"
	"time"

	"go.uber.org/zap"

	"developers-cv/backend/internal/modules/resume_link/application/commands"
)

const (
	retentionCleanupInterval = 24 * time.Hour
	retentionCleanupBatch    = 500
)

type Worker struct {
	cleanup   *commands.CleanupAnalyticsHandler
	logger    *zap.Logger
	interval  time.Duration
	batchSize int
}

func NewWorker(cleanup *commands.CleanupAnalyticsHandler, logger *zap.Logger) *Worker {
	return &Worker{cleanup: cleanup, logger: logger, interval: retentionCleanupInterval, batchSize: retentionCleanupBatch}
}

func (w *Worker) Run(ctx context.Context) {
	if w == nil || w.cleanup == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.run(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.run(ctx)
		}
	}
}

func (w *Worker) run(ctx context.Context) {
	startedAt := time.Now()
	var total commands.AnalyticsRetentionResult
	for {
		result, err := w.cleanup.Handle(ctx, time.Now().UTC(), w.batchSize)
		if err != nil {
			if w.logger != nil {
				w.logger.Error("analytics retention cleanup failed",
					zap.Error(err),
					zap.Int("clicks_deleted", total.Clicks),
					zap.Int("views_deleted", total.Views),
					zap.Duration("duration", time.Since(startedAt)),
				)
			}
			return
		}
		total.Clicks += result.Clicks
		total.Views += result.Views
		if result.Clicks == 0 && result.Views == 0 {
			break
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
	if w.logger != nil && (total.Clicks > 0 || total.Views > 0) {
		w.logger.Info("analytics retention cleanup completed",
			zap.Int("clicks_deleted", total.Clicks),
			zap.Int("views_deleted", total.Views),
			zap.Duration("duration", time.Since(startedAt)),
		)
	}
}
