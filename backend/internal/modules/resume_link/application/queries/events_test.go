package queries

import (
	"context"

	"developers-cv/backend/internal/modules/resume_link/domain"
)

type linkEventRecorderStub struct {
	event *domain.ResumeLinkEvent
	err   error
}

func (s *linkEventRecorderStub) Record(_ context.Context, event *domain.ResumeLinkEvent) error {
	s.event = event
	return s.err
}
