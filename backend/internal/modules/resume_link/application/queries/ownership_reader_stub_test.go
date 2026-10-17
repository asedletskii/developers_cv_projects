package queries

import (
	"context"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
)

type ownershipReaderStub struct {
	err    error
	viewID string
	query  resumequeries.GetResumeQuery
	calls  int
}

func (s *ownershipReaderStub) GetOwned(_ context.Context, query resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	s.calls++
	s.query = query
	if s.err != nil {
		return resumequeries.ResumeView{}, s.err
	}
	viewID := s.viewID
	if viewID == "" {
		viewID = query.ResumeID.String()
	}
	return resumequeries.ResumeView{ID: viewID}, nil
}

func (s *ownershipReaderStub) OwnsResume(_ context.Context, query resumequeries.GetResumeQuery) (bool, error) {
	s.calls++
	s.query = query
	if s.err != nil {
		return false, s.err
	}
	return s.viewID == "" || s.viewID == query.ResumeID.String(), nil
}
