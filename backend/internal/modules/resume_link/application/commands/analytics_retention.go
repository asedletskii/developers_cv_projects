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

type AnalyticsRetentionPolicy struct {
	ClickRetentionDays int
	ViewRetentionDays  int
}

type CleanupAnalyticsHandler struct {
	repository ports.AnalyticsRetentionRepository
	uow        ports.UnitOfWork
	retention  AnalyticsRetentionPolicy
}

func NewCleanupAnalyticsHandler(
	repository ports.AnalyticsRetentionRepository,
	uow ports.UnitOfWork,
	retention AnalyticsRetentionPolicy,
) *CleanupAnalyticsHandler {
	return &CleanupAnalyticsHandler{repository: repository, uow: uow, retention: retention}
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
	cutoff := now.UTC().AddDate(0, 0, -h.retention.ClickRetentionDays)
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
	cutoff := now.UTC().AddDate(0, 0, -h.retention.ViewRetentionDays)
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
	if h.retention.ClickRetentionDays <= 0 || h.retention.ClickRetentionDays > 3650 ||
		h.retention.ViewRetentionDays <= 0 || h.retention.ViewRetentionDays > 3650 {
		return fmt.Errorf("cleanup analytics: retention days must be between 1 and 3650")
	}
	if batchSize <= 0 {
		return fmt.Errorf("cleanup analytics: batch size must be positive")
	}
	return nil
}
