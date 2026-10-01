package ports

import (
	"context"
	"time"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	FindActiveByTokenHash(ctx context.Context, tokenHash valueobjects.SessionTokenHash) (*domain.Session, error)
	Revoke(ctx context.Context, sessionID identity.ID, at time.Time) error
	RevokeAllByUserID(ctx context.Context, userID identity.ID, at time.Time) error
	Touch(ctx context.Context, sessionID identity.ID, at time.Time) error
	DeleteExpiredBatch(ctx context.Context, now time.Time, limit int) (int, error)
}
