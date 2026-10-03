package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeCommandUnitOfWork struct {
	called bool
}

func (u *fakeCommandUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	u.called = true
	return fn(ctx)
}

func testDate(t *testing.T) domain.Date {
	t.Helper()
	v, err := domain.ParseDate("2020-01-02")
	require.NoError(t, err)
	return v
}

func testURL(t *testing.T) domain.ExternalURL {
	t.Helper()
	v, err := domain.NewExternalURL("https://example.com/item")
	require.NoError(t, err)
	return v
}
