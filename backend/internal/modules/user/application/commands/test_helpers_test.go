package commands

import (
	"context"

	"developers-cv/backend/internal/modules/user/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type createUserRepositoryStub struct {
	user   *domain.User
	err    error
	exists bool
}

func (s *createUserRepositoryStub) Create(_ context.Context, user *domain.User) error {
	s.user = user
	return s.err
}

func (s *createUserRepositoryStub) Exists(context.Context, sharedvalueobjects.ID) (bool, error) {
	return s.exists, s.err
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
