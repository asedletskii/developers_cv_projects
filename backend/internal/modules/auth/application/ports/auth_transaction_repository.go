package ports

import (
	"context"
	"time"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type AuthTransactionRepository interface {
	Create(ctx context.Context, transaction *domain.AuthTransaction) error
	FindByStateHash(ctx context.Context, stateHash valueobjects.StateHash) (*domain.AuthTransaction, error)
	FindByStateHashForUpdate(ctx context.Context, stateHash valueobjects.StateHash) (*domain.AuthTransaction, error)
	Consume(ctx context.Context, transactionID identity.ID, at time.Time) error
	DeleteExpiredBatch(ctx context.Context, now time.Time, limit int) (int, error)
}
