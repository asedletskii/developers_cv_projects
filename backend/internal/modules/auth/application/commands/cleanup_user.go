package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/auth/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type DeleteUserDataHandler struct {
	repository ports.UserDataCleanupRepository
}

func NewDeleteUserDataHandler(repository ports.UserDataCleanupRepository) *DeleteUserDataHandler {
	return &DeleteUserDataHandler{repository: repository}
}

func (h *DeleteUserDataHandler) DeleteByUserID(ctx context.Context, userID sharedvalueobjects.ID) error {
	if h == nil || h.repository == nil {
		return fmt.Errorf("delete user auth data: handler dependencies are nil")
	}
	if userID.IsZero() {
		return fmt.Errorf("delete user auth data: user id is empty")
	}
	if err := h.repository.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("delete user auth data: %w", err)
	}
	return nil
}
