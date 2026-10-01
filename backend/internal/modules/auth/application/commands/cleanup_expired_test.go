package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type cleanupTransactionRepo struct {
	deleted int
	err     error
}

func (r *cleanupTransactionRepo) Create(context.Context, *domain.AuthTransaction) error { return nil }
func (r *cleanupTransactionRepo) FindByStateHash(context.Context, valueobjects.StateHash) (*domain.AuthTransaction, error) {
	return nil, domain.ErrEntityNotFound
}
func (r *cleanupTransactionRepo) FindByStateHashForUpdate(context.Context, valueobjects.StateHash) (*domain.AuthTransaction, error) {
	return nil, domain.ErrEntityNotFound
}
func (r *cleanupTransactionRepo) Consume(context.Context, identity.ID, time.Time) error {
	return nil
}
func (r *cleanupTransactionRepo) DeleteExpiredBatch(context.Context, time.Time, int) (int, error) {
	return r.deleted, r.err
}

type cleanupSessionRepo struct {
	deleted int
	err     error
}

func (r *cleanupSessionRepo) Create(context.Context, *domain.Session) error { return nil }
func (r *cleanupSessionRepo) FindActiveByTokenHash(context.Context, valueobjects.SessionTokenHash) (*domain.Session, error) {
	return nil, domain.ErrEntityNotFound
}
func (r *cleanupSessionRepo) Revoke(context.Context, identity.ID, time.Time) error {
	return nil
}
func (r *cleanupSessionRepo) RevokeAllByUserID(context.Context, identity.ID, time.Time) error {
	return nil
}
func (r *cleanupSessionRepo) Touch(context.Context, identity.ID, time.Time) error {
	return nil
}
func (r *cleanupSessionRepo) DeleteExpiredBatch(context.Context, time.Time, int) (int, error) {
	return r.deleted, r.err
}

type cleanupUOW struct {
	calls int
	err   error
}

func (u *cleanupUOW) Do(ctx context.Context, fn func(context.Context) error) error {
	u.calls++
	if u.err != nil {
		return u.err
	}
	return fn(ctx)
}

func TestCleanupExpiredHandlerCleansBothTablesInSeparateBoundedTransactions(t *testing.T) {
	transactions := &cleanupTransactionRepo{deleted: 3}
	sessions := &cleanupSessionRepo{deleted: 7}
	uow := &cleanupUOW{}
	h := NewCleanupExpiredHandler(transactions, sessions, uow)

	result, err := h.Handle(context.Background(), time.Now().UTC(), 500)

	require.NoError(t, err)
	require.Equal(t, CleanupResult{Transactions: 3, Sessions: 7}, result)
	require.Equal(t, 2, uow.calls)
}

func TestCleanupExpiredHandlerRejectsInvalidBatchSize(t *testing.T) {
	h := NewCleanupExpiredHandler(&cleanupTransactionRepo{}, &cleanupSessionRepo{}, &cleanupUOW{})
	_, err := h.Handle(context.Background(), time.Now().UTC(), 0)
	require.Error(t, err)
}

func TestCleanupExpiredHandlerPropagatesFailure(t *testing.T) {
	expected := errors.New("db unavailable")
	h := NewCleanupExpiredHandler(&cleanupTransactionRepo{err: expected}, &cleanupSessionRepo{}, &cleanupUOW{})
	_, err := h.CleanupTransactions(context.Background(), time.Now().UTC(), 500)
	require.ErrorIs(t, err, expected)
}
