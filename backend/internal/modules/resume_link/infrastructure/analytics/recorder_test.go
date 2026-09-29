package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/domain/valueobjects"
)

type repositoryStub struct {
	err   error
	event *domain.ResumeLinkEvent
}

func (s *repositoryStub) Create(context.Context, *domain.ResumeLink) error { return nil }
func (s *repositoryStub) GetByID(context.Context, valueobjects.ID) (*domain.ResumeLink, error) {
	return nil, nil
}
func (s *repositoryStub) ReplaceStatus(context.Context, valueobjects.ID, bool, int64) error {
	return nil
}

func (s *repositoryStub) Delete(context.Context, valueobjects.ID) error { return nil }
func (s *repositoryStub) RecordEvent(_ context.Context, event *domain.ResumeLinkEvent) error {
	s.event = event
	return s.err
}

func TestRecorderPersistsEvent(t *testing.T) {
	repository := &repositoryStub{}
	recorder := NewRecorder(repository, zaptest.NewLogger(t))
	event, err := domain.NewResumeLinkEvent(domain.ResumeLinkEventParams{
		ID: valueobjects.NewID(), LinkID: valueobjects.NewID(), Type: domain.EventTypeView,
		OccurredAt: testTime(),
	})
	require.NoError(t, err)

	require.NoError(t, recorder.Record(context.Background(), event))
	require.Same(t, event, repository.event)
}

func TestRecorderReturnsAndLogsPersistenceFailure(t *testing.T) {
	want := errors.New("database unavailable")
	repository := &repositoryStub{err: want}
	recorder := NewRecorder(repository, zaptest.NewLogger(t))
	event, err := domain.NewResumeLinkEvent(domain.ResumeLinkEventParams{
		ID: valueobjects.NewID(), LinkID: valueobjects.NewID(), Type: domain.EventTypeContactView,
		OccurredAt: testTime(),
	})
	require.NoError(t, err)

	require.ErrorIs(t, recorder.Record(context.Background(), event), want)
}

func testTime() (t time.Time) { return time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC) }
