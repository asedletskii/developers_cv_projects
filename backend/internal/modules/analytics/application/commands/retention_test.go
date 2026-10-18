package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCleanupAnalyticsHandlerUsesConfiguredRetentionWindows(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	repo := &analyticsRetentionRepositoryStub{clicks: 3, views: 7}
	uow := &unitOfWorkStub{}
	h := NewCleanupAnalyticsHandler(repo, uow, AnalyticsRetentionPolicy{ClickRetentionDays: 90, ViewRetentionDays: 365})

	result, err := h.Handle(context.Background(), now, 500)

	require.NoError(t, err)
	require.Equal(t, AnalyticsRetentionResult{Clicks: 3, Views: 7}, result)
	require.Equal(t, 2, uow.calls)
	require.Equal(t, now.Add(-90*24*time.Hour), repo.cutoffs[0])
	require.Equal(t, now.AddDate(0, 0, -365), repo.cutoffs[1])
}

func TestCleanupAnalyticsHandlerPropagatesClickFailure(t *testing.T) {
	expected := errors.New("db unavailable")
	h := NewCleanupAnalyticsHandler(&analyticsRetentionRepositoryStub{clickErr: expected}, &unitOfWorkStub{}, AnalyticsRetentionPolicy{ClickRetentionDays: 90, ViewRetentionDays: 365})

	_, err := h.Handle(context.Background(), time.Now().UTC(), 500)

	require.ErrorIs(t, err, expected)
}

func TestCleanupAnalyticsHandlerRejectsInvalidInput(t *testing.T) {
	h := NewCleanupAnalyticsHandler(&analyticsRetentionRepositoryStub{}, &unitOfWorkStub{}, AnalyticsRetentionPolicy{ClickRetentionDays: 90, ViewRetentionDays: 365})

	_, err := h.Handle(context.Background(), time.Time{}, 500)
	require.Error(t, err)
	_, err = h.Handle(context.Background(), time.Now().UTC(), 0)
	require.Error(t, err)
}

func TestCleanupAnalyticsHandlerRejectsInvalidRetentionPolicy(t *testing.T) {
	h := NewCleanupAnalyticsHandler(
		&analyticsRetentionRepositoryStub{},
		&unitOfWorkStub{},
		AnalyticsRetentionPolicy{ClickRetentionDays: 0, ViewRetentionDays: 365},
	)

	_, err := h.Handle(context.Background(), time.Now().UTC(), 100)

	require.ErrorContains(t, err, "retention days")
}
