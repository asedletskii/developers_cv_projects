package queries

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/shared/domain/valueobjects"
)

type linkEventRecorderStub struct {
	event *domain.ResumeLinkEvent
	err   error
}

func (s *linkEventRecorderStub) Record(_ context.Context, event *domain.ResumeLinkEvent) error {
	s.event = event
	return s.err
}

type contextRecorderStub struct {
	ctx         context.Context
	ctxErr      error
	deadline    time.Time
	hasDeadline bool
}

func (r *contextRecorderStub) Record(ctx context.Context, _ *domain.ResumeLinkEvent) error {
	r.ctx = ctx
	r.ctxErr = ctx.Err()
	r.deadline, r.hasDeadline = ctx.Deadline()
	return nil
}

func TestRecordLinkEventUsesDetachedBoundedContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	recorder := &contextRecorderStub{}
	recordLinkEvent(ctx, recorder, valueobjects.NewID(), domain.EventTypeView)

	require.NotNil(t, recorder.ctx)
	require.NoError(t, recorder.ctxErr)
	require.True(t, recorder.hasDeadline)
	require.WithinDuration(t, time.Now().Add(linkEventRecordTimeout), recorder.deadline, 100*time.Millisecond)
}
