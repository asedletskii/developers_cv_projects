package commands

import (
	"context"

	"developers-cv/backend/internal/modules/user/domain"
)

type createUserRepositoryStub struct {
	user *domain.User
	err  error
}

func (s *createUserRepositoryStub) Create(_ context.Context, user *domain.User) error {
	s.user = user
	return s.err
}

type profileRepositoryStub struct {
	profile *domain.Profile
	err     error
}

func (s *profileRepositoryStub) ReplaceProfile(_ context.Context, profile *domain.Profile) error {
	s.profile = profile
	return s.err
}

type unitOfWorkStub struct{}

func (unitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func stringPtr(value string) *string { return &value }
