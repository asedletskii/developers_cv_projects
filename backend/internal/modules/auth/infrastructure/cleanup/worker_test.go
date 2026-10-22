package cleanup

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/application/commands"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/platform/telemetry"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type cleanupTransactionStub struct{}

func (cleanupTransactionStub) Create(context.Context, *domain.AuthTransaction) error { return nil }
func (cleanupTransactionStub) FindByStateHash(context.Context, valueobjects.StateHash) (*domain.AuthTransaction, error) {
	return nil, nil
}
func (cleanupTransactionStub) FindByStateHashForUpdate(context.Context, valueobjects.StateHash) (*domain.AuthTransaction, error) {
	return nil, nil
}
func (cleanupTransactionStub) Consume(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (cleanupTransactionStub) DeleteExpiredBatch(context.Context, time.Time, int) (int, error) {
	return 0, nil
}

type cleanupSessionStub struct {
	batches []int
	calls   int
}

func (s *cleanupSessionStub) Create(context.Context, *domain.Session) error { return nil }
func (s *cleanupSessionStub) FindActiveByTokenHash(context.Context, valueobjects.SessionTokenHash) (*domain.Session, error) {
	return nil, nil
}
func (s *cleanupSessionStub) Revoke(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (s *cleanupSessionStub) RevokeAllByUserID(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (s *cleanupSessionStub) Touch(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (s *cleanupSessionStub) DeleteExpiredBatch(_ context.Context, _ time.Time, limit int) (int, error) {
	s.calls++
	if limit <= 0 {
		return 0, errors.New("limit must be positive")
	}
	if len(s.batches) == 0 {
		return 0, nil
	}
	count := s.batches[0]
	s.batches = s.batches[1:]
	return count, nil
}

type cleanupUOWStub struct{}

func (cleanupUOWStub) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

func TestRunSessionsProcessesBatchesUntilPartialBatch(t *testing.T) {
	sessions := &cleanupSessionStub{batches: []int{cleanupBatchSize, cleanupBatchSize, 12}}
	handler := commands.NewCleanupExpiredHandler(cleanupTransactionStub{}, sessions, cleanupUOWStub{})
	metrics, err := telemetry.NewMetrics(true)
	require.NoError(t, err)
	worker := NewWorker(handler, nil, metrics)

	worker.runSessions(context.Background())

	require.Equal(t, 3, sessions.calls)
	response := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	require.Contains(t, response.Body.String(), `developers_cv_auth_session_cleanup_total{reason="expired"} 1012`)
}
