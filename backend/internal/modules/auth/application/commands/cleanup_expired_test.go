package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCleanupExpiredHandlerCleansBothTablesInSeparateBoundedTransactions(t *testing.T) {
	transactions := &cleanupTransactionRepositoryStub{deleted: 3}
	sessions := &cleanupSessionRepositoryStub{deleted: 7}
	uow := &unitOfWorkStub{}
	h := NewCleanupExpiredHandler(transactions, sessions, uow)

	result, err := h.Handle(context.Background(), time.Now().UTC(), 500)

	require.NoError(t, err)
	require.Equal(t, CleanupResult{Transactions: 3, Sessions: 7}, result)
	require.Equal(t, 2, uow.calls)
}

func TestCleanupExpiredHandlerRejectsInvalidBatchSize(t *testing.T) {
	h := NewCleanupExpiredHandler(&cleanupTransactionRepositoryStub{}, &cleanupSessionRepositoryStub{}, &unitOfWorkStub{})
	_, err := h.Handle(context.Background(), time.Now().UTC(), 0)
	require.Error(t, err)
}

func TestCleanupExpiredHandlerPropagatesFailure(t *testing.T) {
	expected := errors.New("db unavailable")
	h := NewCleanupExpiredHandler(&cleanupTransactionRepositoryStub{err: expected}, &cleanupSessionRepositoryStub{}, &unitOfWorkStub{})
	_, err := h.CleanupTransactions(context.Background(), time.Now().UTC(), 500)
	require.ErrorIs(t, err, expected)
}
