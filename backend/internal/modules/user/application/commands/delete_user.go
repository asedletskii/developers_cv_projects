package commands

import (
	"context"
	"fmt"

	authcontracts "developers-cv/backend/internal/modules/auth/contracts"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/user/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type DeleteUserCommand struct {
	UserID sharedvalueobjects.ID
}

type DeleteUserHandler struct {
	users   ports.UserDeleteRepository
	auth    authcontracts.UserDataCleanup
	resumes resumecontracts.UserDataCleanup
	uow     ports.UnitOfWork
}

func NewDeleteUserHandler(
	users ports.UserDeleteRepository,
	auth authcontracts.UserDataCleanup,
	resumes resumecontracts.UserDataCleanup,
	uow ports.UnitOfWork,
) *DeleteUserHandler {
	return &DeleteUserHandler{users: users, auth: auth, resumes: resumes, uow: uow}
}

func (h *DeleteUserHandler) Handle(ctx context.Context, c DeleteUserCommand) error {
	if h == nil || h.users == nil || h.auth == nil || h.resumes == nil || h.uow == nil {
		return fmt.Errorf("delete user handler dependencies are nil")
	}
	if c.UserID.IsZero() {
		return fmt.Errorf("delete user: user id is empty")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		exists, err := h.users.LockForDelete(tx, c.UserID)
		if err != nil {
			return fmt.Errorf("delete user: lock user: %w", err)
		}
		if !exists {
			return nil
		}
		if err := h.auth.DeleteByUserID(tx, c.UserID); err != nil {
			return fmt.Errorf("delete user: clean auth data: %w", err)
		}
		if err := h.resumes.DeleteByUserID(tx, c.UserID); err != nil {
			return fmt.Errorf("delete user: clean resume data: %w", err)
		}
		if err := h.users.Delete(tx, c.UserID); err != nil {
			return fmt.Errorf("delete user: %w", err)
		}
		return nil
	})
}
