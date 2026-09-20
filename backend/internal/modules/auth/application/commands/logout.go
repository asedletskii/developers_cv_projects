package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
)

type LogoutHandler struct {
	sessions ports.SessionRepository
	uow      ports.UnitOfWork
}

func NewLogoutHandler(sessions ports.SessionRepository, uow ports.UnitOfWork) *LogoutHandler {
	return &LogoutHandler{sessions: sessions, uow: uow}
}

func (h *LogoutHandler) Handle(ctx context.Context, tokenHash valueobjects.SessionTokenHash, now time.Time) error {
	if h == nil || h.sessions == nil || h.uow == nil {
		return fmt.Errorf("logout handler dependencies are nil")
	}
	if now.IsZero() {
		return fmt.Errorf("logout: timestamp must be set")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		session, err := h.sessions.FindActiveByTokenHash(tx, tokenHash)
		if err != nil {
			if errors.Is(err, domain.ErrEntityNotFound) {
				return nil
			}
			return err
		}
		if session == nil {
			return fmt.Errorf("logout: session repository returned nil session")
		}
		if err := session.Revoke(now); err != nil {
			return err
		}
		return h.sessions.Revoke(tx, session.ID(), now)
	})
}
