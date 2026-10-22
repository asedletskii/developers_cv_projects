package health

import (
	"context"
	"sync"
	"time"
)

const defaultCheckTimeout = 2 * time.Second

type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

type Service struct {
	mu           sync.RWMutex
	checks       []Checker
	checkTimeout time.Duration
}

func NewService(checks ...Checker) *Service {
	return NewServiceWithTimeout(defaultCheckTimeout, checks...)
}

func NewServiceWithTimeout(timeout time.Duration, checks ...Checker) *Service {
	if timeout <= 0 {
		timeout = defaultCheckTimeout
	}
	return &Service{checks: append([]Checker(nil), checks...), checkTimeout: timeout}
}

func (s *Service) Check(ctx context.Context) map[string]error {
	s.mu.RLock()
	checks := append([]Checker(nil), s.checks...)
	timeout := s.checkTimeout
	s.mu.RUnlock()

	results := make(map[string]error, len(checks))
	for _, check := range checks {
		checkCtx, cancel := context.WithTimeout(ctx, timeout)
		results[check.Name()] = check.Check(checkCtx)
		cancel()
	}
	return results
}
