package cleanup

import (
	"context"
	"time"

	"go.uber.org/zap"

	"developers-cv/backend/internal/modules/auth/application/commands"
	"developers-cv/backend/internal/platform/telemetry"
)

const (
	transactionCleanupInterval = 5 * time.Minute
	sessionCleanupInterval     = 6 * time.Hour
	cleanupBatchSize           = 500
)

type Worker struct {
	cleanup   *commands.CleanupExpiredHandler
	logger    *zap.Logger
	metrics   *telemetry.Metrics
	interval  time.Duration
	batchSize int
}

func NewWorker(cleanup *commands.CleanupExpiredHandler, logger *zap.Logger, metrics *telemetry.Metrics) *Worker {
	return &Worker{cleanup: cleanup, logger: logger, metrics: metrics, interval: transactionCleanupInterval, batchSize: cleanupBatchSize}
}

func (w *Worker) Run(ctx context.Context) {
	if w == nil || w.cleanup == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	transactionTicker := time.NewTicker(w.interval)
	sessionTicker := time.NewTicker(sessionCleanupInterval)
	defer transactionTicker.Stop()
	defer sessionTicker.Stop()

	w.runTransactions(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-transactionTicker.C:
			w.runTransactions(ctx)
		case <-sessionTicker.C:
			w.runSessions(ctx)
		}
	}
}

func (w *Worker) runTransactions(ctx context.Context) {
	startedAt := time.Now()
	count, err := w.cleanup.CleanupTransactions(ctx, time.Now().UTC(), w.batchSize)
	if err != nil {
		if w.metrics != nil {
			w.metrics.RecordAuthCleanupFailure("transactions")
		}
		if w.logger != nil {
			w.logger.Error("auth cleanup failed", zap.Error(err), zap.Duration("duration", time.Since(startedAt)))
		}
		return
	}
	if count > 0 && w.metrics != nil {
		w.metrics.RecordAuthTransactionCleanup("telegram", "expired", count)
	}
	if w.logger != nil && count > 0 {
		w.logger.Info("auth transaction cleanup completed", zap.Int("transactions_deleted", count), zap.Duration("duration", time.Since(startedAt)))
	}
}

func (w *Worker) runSessions(ctx context.Context) {
	startedAt := time.Now()
	count, err := w.cleanup.CleanupSessions(ctx, time.Now().UTC(), w.batchSize)
	if err != nil {
		if w.metrics != nil {
			w.metrics.RecordAuthCleanupFailure("sessions")
		}
		if w.logger != nil {
			w.logger.Error("auth session cleanup failed", zap.Error(err), zap.Duration("duration", time.Since(startedAt)))
		}
		return
	}
	if w.logger != nil && count > 0 {
		w.logger.Info("auth session cleanup completed", zap.Int("sessions_deleted", count), zap.Duration("duration", time.Since(startedAt)))
	}
}
