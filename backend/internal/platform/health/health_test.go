package health

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeChecker struct {
	name string
	err  error
}

func (f fakeChecker) Name() string                { return f.name }
func (f fakeChecker) Check(context.Context) error { return f.err }

func TestServiceCheck(t *testing.T) {
	expectedErr := errors.New("not ready")
	service := NewService(fakeChecker{name: "database", err: expectedErr}, fakeChecker{name: "storage"})

	results := service.Check(context.Background())

	require.Len(t, results, 2)
	require.ErrorIs(t, results["database"], expectedErr)
	require.NoError(t, results["storage"])
}
