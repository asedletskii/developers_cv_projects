package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/user/application/ports"
	"developers-cv/backend/internal/modules/user/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CreateUserHandler struct {
	repository ports.UserRepository
}

func NewCreateUserHandler(repository ports.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{repository: repository}
}

func (h *CreateUserHandler) Handle(ctx context.Context) (sharedvalueobjects.ID, error) {
	if h == nil || h.repository == nil {
		return sharedvalueobjects.ID{}, fmt.Errorf("create user handler dependencies are nil")
	}

	user, err := domain.NewUser(domain.UserParams{ID: sharedvalueobjects.NewID()})
	if err != nil {
		return sharedvalueobjects.ID{}, err
	}
	if err := h.repository.Create(ctx, user); err != nil {
		return sharedvalueobjects.ID{}, err
	}
	return user.ID(), nil
}

// CreateUser satisfies auth.application.ports.UserProvisioner without coupling
// the user module to the auth module.
func (h *CreateUserHandler) CreateUser(ctx context.Context) (sharedvalueobjects.ID, error) {
	return h.Handle(ctx)
}
