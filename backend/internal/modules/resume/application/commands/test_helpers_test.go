package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
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

func testDate(t *testing.T) valueobjects.Date {
	t.Helper()
	v, err := valueobjects.ParseDate("2020-01-02")
	require.NoError(t, err)
	return v
}

func testURL(t *testing.T) valueobjects.ExternalURL {
	t.Helper()
	v, err := valueobjects.NewExternalURL("https://example.com/item")
	require.NoError(t, err)
	return v
}

func ptr[T any](value T) *T {
	return &value
}
