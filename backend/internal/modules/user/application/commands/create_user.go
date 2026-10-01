package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/shared/domain/valueobjects"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
}

type CreateUserHandler struct {
	repository UserRepository
}

func NewCreateUserHandler(repository UserRepository) *CreateUserHandler {
	return &CreateUserHandler{repository: repository}
}

func (h *CreateUserHandler) Handle(ctx context.Context) (valueobjects.ID, error) {
	if h == nil || h.repository == nil {
		return valueobjects.ID{}, fmt.Errorf("create user handler dependencies are nil")
	}

	user, err := domain.NewUser(domain.UserParams{ID: valueobjects.NewID()})
	if err != nil {
		return valueobjects.ID{}, err
	}
	if err := h.repository.Create(ctx, user); err != nil {
		return valueobjects.ID{}, err
	}
	return user.ID(), nil
}

// CreateUser satisfies auth.application.ports.UserProvisioner without coupling
// the user module to the auth module.
func (h *CreateUserHandler) CreateUser(ctx context.Context) (valueobjects.ID, error) {
	return h.Handle(ctx)
}
