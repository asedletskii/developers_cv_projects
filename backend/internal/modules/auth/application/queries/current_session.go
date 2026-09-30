package queries

import (
	"context"
	"errors"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/auth/application"
	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CurrentSessionView struct {
	UserID     sharedvalueobjects.ID
	SessionID  sharedvalueobjects.ID
	CreatedAt  time.Time
	LastSeenAt time.Time
	ExpiresAt  time.Time
}

type CurrentSessionHandler struct {
	sessions      ports.SessionRepository
	uow           ports.UnitOfWork
	idleTTL       time.Duration
	touchInterval time.Duration
	now           func() time.Time
}

func NewCurrentSessionHandler(sessions ports.SessionRepository, uow ports.UnitOfWork, idleTTL, touchInterval time.Duration) *CurrentSessionHandler {
	return &CurrentSessionHandler{
		sessions: sessions, uow: uow, idleTTL: idleTTL, touchInterval: touchInterval,
		now: func() time.Time { return time.Now().UTC() },
	}
}

func (h *CurrentSessionHandler) Handle(ctx context.Context, tokenHash valueobjects.SessionTokenHash) (CurrentSessionView, error) {
	if h == nil || h.sessions == nil || h.uow == nil {
		return CurrentSessionView{}, fmt.Errorf("current session handler dependencies are nil")
	}
	if h.idleTTL <= 0 || h.touchInterval <= 0 {
		return CurrentSessionView{}, fmt.Errorf("%w: invalid current session configuration", application.ErrInvalidAuthFlow)
	}

	session, err := h.sessions.FindActiveByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return CurrentSessionView{}, application.ErrSessionInvalid
		}
		return CurrentSessionView{}, err
	}
	now := h.now().UTC()
	if session.IsExpired(now) {
		return CurrentSessionView{}, application.ErrSessionExpired
	}
	if session.IsIdleExpired(now, h.idleTTL) {
		return CurrentSessionView{}, application.ErrSessionIdleExpired
	}
	if now.Sub(session.LastSeenAt()) >= h.touchInterval {
		if err := h.uow.Do(ctx, func(tx context.Context) error {
			return h.sessions.Touch(tx, session.ID(), now)
		}); err != nil {
			return CurrentSessionView{}, err
		}
		if err = session.Touch(now); err != nil {
			return CurrentSessionView{}, err
		}
	}
	return CurrentSessionView{UserID: session.UserID(), SessionID: session.ID(), CreatedAt: session.CreatedAt(), LastSeenAt: session.LastSeenAt(), ExpiresAt: session.ExpiresAt()}, nil
}
