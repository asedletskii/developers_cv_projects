package commands

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/auth/application/ports"
)

type CleanupResult struct {
	Transactions int
	Sessions     int
}

type CleanupExpiredHandler struct {
	transactions ports.AuthTransactionRepository
	sessions     ports.SessionRepository
	uow          ports.UnitOfWork
}

func NewCleanupExpiredHandler(
	transactions ports.AuthTransactionRepository,
	sessions ports.SessionRepository,
	uow ports.UnitOfWork,
) *CleanupExpiredHandler {
	return &CleanupExpiredHandler{transactions: transactions, sessions: sessions, uow: uow}
}

func (h *CleanupExpiredHandler) Handle(ctx context.Context, now time.Time, batchSize int) (CleanupResult, error) {
	transactions, err := h.CleanupTransactions(ctx, now, batchSize)
	if err != nil {
		return CleanupResult{}, err
	}
	sessions, err := h.CleanupSessions(ctx, now, batchSize)
	if err != nil {
		return CleanupResult{}, err
	}
	return CleanupResult{Transactions: transactions, Sessions: sessions}, nil
}

func (h *CleanupExpiredHandler) CleanupTransactions(ctx context.Context, now time.Time, batchSize int) (int, error) {
	if err := h.validate(batchSize); err != nil {
		return 0, err
	}
	var deleted int
	err := h.uow.Do(ctx, func(tx context.Context) error {
		count, err := h.transactions.DeleteExpiredBatch(tx, now, batchSize)
		deleted = count
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("cleanup expired auth transactions: %w", err)
	}
	return deleted, nil
}

func (h *CleanupExpiredHandler) CleanupSessions(ctx context.Context, now time.Time, batchSize int) (int, error) {
	if err := h.validate(batchSize); err != nil {
		return 0, err
	}
	var deleted int
	err := h.uow.Do(ctx, func(tx context.Context) error {
		count, err := h.sessions.DeleteExpiredBatch(tx, now, batchSize)
		deleted = count
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("cleanup expired auth sessions: %w", err)
	}
	return deleted, nil
}

func (h *CleanupExpiredHandler) validate(batchSize int) error {
	if h == nil || h.transactions == nil || h.sessions == nil || h.uow == nil {
		return fmt.Errorf("cleanup expired auth data handler dependencies are nil")
	}
	if batchSize <= 0 {
		return fmt.Errorf("cleanup expired auth data: batch size must be positive")
	}
	return nil
}
