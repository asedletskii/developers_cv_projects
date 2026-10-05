package commands

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
)

type AnalyticsRetentionResult struct {
	Clicks int
	Views  int
}

type CleanupAnalyticsHandler struct {
	repository ports.AnalyticsRetentionRepository
	uow        ports.UnitOfWork
}

func NewCleanupAnalyticsHandler(repository ports.AnalyticsRetentionRepository, uow ports.UnitOfWork) *CleanupAnalyticsHandler {
	return &CleanupAnalyticsHandler{repository: repository, uow: uow}
}

func (h *CleanupAnalyticsHandler) Handle(ctx context.Context, now time.Time, batchSize int) (AnalyticsRetentionResult, error) {
	clicks, err := h.CleanupClicks(ctx, now, batchSize)
	if err != nil {
		return AnalyticsRetentionResult{}, err
	}
	views, err := h.CleanupViews(ctx, now, batchSize)
	if err != nil {
		return AnalyticsRetentionResult{}, err
	}
	return AnalyticsRetentionResult{Clicks: clicks, Views: views}, nil
}

func (h *CleanupAnalyticsHandler) CleanupClicks(ctx context.Context, now time.Time, batchSize int) (int, error) {
	if err := h.validate(now, batchSize); err != nil {
		return 0, err
	}
	cutoff := now.UTC().Add(-90 * 24 * time.Hour)
	var deleted int
	if err := h.uow.Do(ctx, func(tx context.Context) error {
		count, err := h.repository.DeleteExpiredClicks(tx, cutoff, batchSize)
		deleted = count
		return err
	}); err != nil {
		return 0, fmt.Errorf("cleanup expired analytics clicks: %w", err)
	}
	return deleted, nil
}

func (h *CleanupAnalyticsHandler) CleanupViews(ctx context.Context, now time.Time, batchSize int) (int, error) {
	if err := h.validate(now, batchSize); err != nil {
		return 0, err
	}
	cutoff := now.UTC().AddDate(-1, 0, 0)
	var deleted int
	if err := h.uow.Do(ctx, func(tx context.Context) error {
		count, err := h.repository.DeleteExpiredViews(tx, cutoff, batchSize)
		deleted = count
		return err
	}); err != nil {
		return 0, fmt.Errorf("cleanup expired analytics views: %w", err)
	}
	return deleted, nil
}

func (h *CleanupAnalyticsHandler) validate(now time.Time, batchSize int) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("cleanup analytics dependencies are nil")
	}
	if now.IsZero() {
		return fmt.Errorf("cleanup analytics: now must be set")
	}
	if batchSize <= 0 {
		return fmt.Errorf("cleanup analytics: batch size must be positive")
	}
	return nil
}
