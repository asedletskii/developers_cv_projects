package cleanup

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/application/commands"
)

type retentionRepositoryStub struct {
	clicksDeleted int
	clickCalls    int
}

func (s *retentionRepositoryStub) DeleteExpiredClicks(ctx context.Context, _ time.Time, _ int) (int, error) {
	s.clickCalls++
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	return s.clicksDeleted, nil
}

func (s *retentionRepositoryStub) DeleteExpiredViews(ctx context.Context, _ time.Time, _ int) (int, error) {
	return 0, ctx.Err()
}

type retentionUnitOfWorkStub struct{}

func (retentionUnitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestWorkerRunStopsWhenItsTimeBudgetExpires(t *testing.T) {
	repository := &retentionRepositoryStub{clicksDeleted: 1}
	handler := commands.NewCleanupAnalyticsHandler(repository, retentionUnitOfWorkStub{}, commands.AnalyticsRetentionPolicy{ClickRetentionDays: 90, ViewRetentionDays: 365})
	worker := &Worker{
		cleanup:        handler,
		maxRunDuration: 10 * time.Millisecond,
		batchSize:      1,
	}

	err := worker.run(context.Background())
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Greater(t, repository.clickCalls, 0)
}

func TestWorkerRunReturnsRepositoryFailure(t *testing.T) {
	expected := errors.New("database unavailable")
	worker := &Worker{
		cleanup: commands.NewCleanupAnalyticsHandler(
			&failingRetentionRepository{err: expected}, retentionUnitOfWorkStub{},
			commands.AnalyticsRetentionPolicy{ClickRetentionDays: 90, ViewRetentionDays: 365},
		),
		maxRunDuration: time.Second,
		batchSize:      1,
	}

	require.ErrorIs(t, worker.run(context.Background()), expected)
}

type failingRetentionRepository struct {
	err error
}

func (s *failingRetentionRepository) DeleteExpiredClicks(context.Context, time.Time, int) (int, error) {
	return 0, s.err
}

func (s *failingRetentionRepository) DeleteExpiredViews(context.Context, time.Time, int) (int, error) {
	return 0, s.err
}
