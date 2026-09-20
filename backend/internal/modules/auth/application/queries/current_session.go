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
	logger        ports.WarningLogger
}

func NewCurrentSessionHandler(sessions ports.SessionRepository, uow ports.UnitOfWork, idleTTL, touchInterval time.Duration, loggers ...ports.WarningLogger) *CurrentSessionHandler {
	var logger ports.WarningLogger
	if len(loggers) > 0 {
		logger = loggers[0]
	}
	return &CurrentSessionHandler{
		sessions: sessions, uow: uow, idleTTL: idleTTL, touchInterval: touchInterval, logger: logger,
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
	if session == nil {
		return CurrentSessionView{}, fmt.Errorf("find active session: repository returned nil session")
	}
	now := h.now().UTC()
	if session.IsExpired(now) {
		return CurrentSessionView{}, application.ErrSessionExpired
	}
	if session.IsIdleExpired(now, h.idleTTL) {
		return CurrentSessionView{}, application.ErrSessionIdleExpired
	}
	if now.Sub(session.LastSeenAt()) >= h.touchInterval {
		if err := ctx.Err(); err != nil {
			return CurrentSessionView{}, err
		}
		if touchErr := h.uow.Do(ctx, func(tx context.Context) error {
			return h.sessions.Touch(tx, session.ID(), now)
		}); touchErr != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return CurrentSessionView{}, ctxErr
			}
			if h.logger != nil {
				h.logger.Warn("failed to update auth session last_seen_at", touchErr)
			}
		} else if touchErr = session.Touch(now); touchErr != nil {
			if h.logger != nil {
				h.logger.Warn("failed to update in-memory auth session activity", touchErr)
			}
		}
	}
	return CurrentSessionView{UserID: session.UserID(), SessionID: session.ID(), CreatedAt: session.CreatedAt(), LastSeenAt: session.LastSeenAt(), ExpiresAt: session.ExpiresAt()}, nil
}
