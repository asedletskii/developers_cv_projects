package cleanup

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/platform/telemetry"
)

const (
	retentionCleanupInterval   = 24 * time.Hour
	retentionCleanupRetry      = 5 * time.Minute
	retentionCleanupMaxRunTime = 5 * time.Minute
	retentionCleanupBatch      = 500
)

type Worker struct {
	cleanup        *commands.CleanupAnalyticsHandler
	logger         *zap.Logger
	metrics        *telemetry.Metrics
	interval       time.Duration
	retryInterval  time.Duration
	maxRunDuration time.Duration
	batchSize      int
}

func NewWorker(cleanup *commands.CleanupAnalyticsHandler, logger *zap.Logger, metrics *telemetry.Metrics) *Worker {
	return &Worker{
		cleanup:        cleanup,
		logger:         logger,
		metrics:        metrics,
		interval:       retentionCleanupInterval,
		retryInterval:  retentionCleanupRetry,
		maxRunDuration: retentionCleanupMaxRunTime,
		batchSize:      retentionCleanupBatch,
	}
}

func (w *Worker) Run(ctx context.Context) {
	if w == nil || w.cleanup == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	interval := w.interval
	if interval <= 0 {
		interval = retentionCleanupInterval
	}
	retryInterval := w.retryInterval
	if retryInterval <= 0 {
		retryInterval = retentionCleanupRetry
	}

	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}

		if err := w.run(ctx); err != nil {
			timer.Reset(retryInterval)
			continue
		}
		timer.Reset(interval)
	}
}

func (w *Worker) run(ctx context.Context) error {
	if w == nil || w.cleanup == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	maxRunDuration := w.maxRunDuration
	if maxRunDuration <= 0 {
		maxRunDuration = retentionCleanupMaxRunTime
	}
	batchSize := w.batchSize
	if batchSize <= 0 {
		batchSize = retentionCleanupBatch
	}

	startedAt := time.Now()
	runCtx, cancel := context.WithTimeout(ctx, maxRunDuration)
	defer cancel()

	var total commands.AnalyticsRetentionResult
	outcome := "success"
	defer func() {
		if w.metrics != nil {
			w.metrics.RecordAnalyticsRetentionCleanup(outcome, total.Clicks, total.Views, time.Since(startedAt))
		}
	}()
	for {
		if err := runCtx.Err(); err != nil {
			outcome = "failure"
			w.logFailure(err, total, startedAt)
			return err
		}

		result, err := w.cleanup.Handle(runCtx, time.Now().UTC(), batchSize)
		if err != nil {
			outcome = "failure"
			w.logFailure(err, total, startedAt)
			return err
		}
		total.Clicks += result.Clicks
		total.Views += result.Views
		if result.Clicks == 0 && result.Views == 0 {
			break
		}
	}

	if w.logger != nil && (total.Clicks > 0 || total.Views > 0) {
		w.logger.Info("analytics retention cleanup completed",
			zap.Int("clicks_deleted", total.Clicks),
			zap.Int("views_deleted", total.Views),
			zap.Duration("duration", time.Since(startedAt)),
		)
	}
	return nil
}

func (w *Worker) logFailure(err error, total commands.AnalyticsRetentionResult, startedAt time.Time) {
	if w.logger == nil || err == nil {
		return
	}
	fields := []zap.Field{
		zap.Error(err),
		zap.Int("clicks_deleted", total.Clicks),
		zap.Int("views_deleted", total.Views),
		zap.Duration("duration", time.Since(startedAt)),
	}
	if errors.Is(err, context.DeadlineExceeded) {
		fields = append(fields, zap.Bool("run_budget_exceeded", true))
	}
	w.logger.Error("analytics retention cleanup failed", fields...)
}
