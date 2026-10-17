package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWaitForWorkersReturnsWhenWorkersFinish(t *testing.T) {
	done := make(chan struct{})
	close(done)

	require.NoError(t, waitForWorkers(context.Background(), done))
}

func TestWaitForWorkersHonorsShutdownDeadline(t *testing.T) {
	done := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	started := time.Now()
	err := waitForWorkers(ctx, done)

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, time.Since(started), time.Second)
}

func TestWaitForWorkersAcceptsNoWorkers(t *testing.T) {
	require.NoError(t, waitForWorkers(context.Background(), nil))
}

type lifecycleServerStub struct {
	runErr       error
	shutdownErr  error
	shutdownCall int
}

func (s *lifecycleServerStub) Run() error { return s.runErr }

func (s *lifecycleServerStub) Shutdown(context.Context) error {
	s.shutdownCall++
	return s.shutdownErr
}

type lifecycleTelemetryStub struct {
	shutdownCall int
	shutdownErr  error
}

func (s *lifecycleTelemetryStub) Shutdown(context.Context) error {
	s.shutdownCall++
	return s.shutdownErr
}

type cancellationWorkerStub struct {
	done chan struct{}
}

func (w cancellationWorkerStub) Run(ctx context.Context) {
	<-ctx.Done()
	close(w.done)
}

func TestApplicationRunShutsDownResourcesWhenServerStops(t *testing.T) {
	serverErr := errors.New("listener failed")
	server := &lifecycleServerStub{runErr: serverErr}
	telemetry := &lifecycleTelemetryStub{}
	worker := cancellationWorkerStub{done: make(chan struct{})}

	application := NewApplication(server, telemetry, zap.NewNop(), nil, time.Second, []backgroundWorker{worker})
	err := application.Run(context.Background())

	require.ErrorIs(t, err, serverErr)
	require.Equal(t, 1, server.shutdownCall)
	require.Equal(t, 1, telemetry.shutdownCall)
	select {
	case <-worker.done:
	case <-time.After(time.Second):
		t.Fatal("background worker did not stop after application shutdown")
	}
}

func TestApplicationShutdownIsIdempotent(t *testing.T) {
	server := &lifecycleServerStub{}
	telemetry := &lifecycleTelemetryStub{}
	application := NewApplication(server, telemetry, zap.NewNop(), nil, time.Second, nil)

	require.NoError(t, application.Shutdown(context.Background()))
	require.NoError(t, application.Shutdown(context.Background()))
	require.Equal(t, 1, server.shutdownCall)
	require.Equal(t, 1, telemetry.shutdownCall)
}
