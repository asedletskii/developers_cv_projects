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
	sessionCleanupMaxDuration  = 30 * time.Second
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
	w.runSessions(ctx)
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
	cleanupCtx, cancel := context.WithTimeout(ctx, sessionCleanupMaxDuration)
	defer cancel()

	var totalDeleted int
	for {
		if err := cleanupCtx.Err(); err != nil {
			if ctx.Err() == nil && w.logger != nil {
				w.logger.Warn("auth session cleanup stopped after reaching its time budget", zap.Int("sessions_deleted", totalDeleted), zap.Duration("duration", time.Since(startedAt)))
			}
			break
		}

		count, err := w.cleanup.CleanupSessions(cleanupCtx, time.Now().UTC(), w.batchSize)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			if w.metrics != nil {
				w.metrics.RecordAuthCleanupFailure("sessions")
			}
			if w.logger != nil {
				w.logger.Error("auth session cleanup failed", zap.Error(err), zap.Int("sessions_deleted", totalDeleted), zap.Duration("duration", time.Since(startedAt)))
			}
			return
		}
		totalDeleted += count
		if count < w.batchSize {
			break
		}
	}

	if totalDeleted > 0 {
		if w.metrics != nil {
			w.metrics.RecordAuthSessionCleanup("expired", totalDeleted)
		}
		if w.logger != nil {
			w.logger.Info("auth session cleanup completed", zap.Int("sessions_deleted", totalDeleted), zap.Duration("duration", time.Since(startedAt)))
		}
	}
}
