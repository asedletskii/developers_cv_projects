package health

import (
	"context"
	"sync"
)

type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

type Service struct {
	mu     sync.RWMutex
	checks []Checker
}

func NewService(checks ...Checker) *Service {
	return &Service{checks: append([]Checker(nil), checks...)}
}

func (s *Service) Check(ctx context.Context) map[string]error {
	s.mu.RLock()
	checks := append([]Checker(nil), s.checks...)
	s.mu.RUnlock()

	results := make(map[string]error, len(checks))
	for _, check := range checks {
		results[check.Name()] = check.Check(ctx)
	}
	return results
}
