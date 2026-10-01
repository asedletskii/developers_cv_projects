package commands

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/auth/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LogoutAllHandler struct {
	sessions ports.SessionRepository
	uow      ports.UnitOfWork
}

func NewLogoutAllHandler(sessions ports.SessionRepository, uow ports.UnitOfWork) *LogoutAllHandler {
	return &LogoutAllHandler{sessions: sessions, uow: uow}
}

func (h *LogoutAllHandler) Handle(ctx context.Context, userID sharedvalueobjects.ID, now time.Time) error {
	if h == nil || h.sessions == nil || h.uow == nil {
		return fmt.Errorf("logout all handler dependencies are nil")
	}
	if userID.IsZero() {
		return fmt.Errorf("invalid logout all command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.sessions.RevokeAllByUserID(tx, userID, now)
	})
}
