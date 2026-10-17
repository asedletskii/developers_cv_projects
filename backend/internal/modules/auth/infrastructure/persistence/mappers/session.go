package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func SessionToDB(session *domain.Session) (gormmodel.AuthSession, error) {
	if session == nil {
		return gormmodel.AuthSession{}, fmt.Errorf("session is nil")
	}
	return gormmodel.AuthSession{
		ID:         session.ID().UUID(),
		UserID:     session.UserID().UUID(),
		TokenHash:  session.TokenHash().Bytes(),
		CreatedAt:  session.CreatedAt(),
		LastSeenAt: session.LastSeenAt(),
		ExpiresAt:  session.ExpiresAt(),
		RevokedAt:  session.RevokedAt(),
		UserAgent:  session.UserAgent(),
		IPAddress:  session.IPAddress(),
	}, nil
}

func SessionToDomain(model gormmodel.AuthSession) (*domain.Session, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("map session id: %w", err)
	}
	userID, err := sharedvalueobjects.IDFromUUID(model.UserID)
	if err != nil {
		return nil, fmt.Errorf("map session user id: %w", err)
	}
	hash, err := valueobjects.NewSessionTokenHash(model.TokenHash)
	if err != nil {
		return nil, fmt.Errorf("map session token hash: %w", err)
	}
	return domain.RestoreSession(domain.SessionSnapshot{
		ID: id, UserID: userID, TokenHash: hash,
		CreatedAt: model.CreatedAt, LastSeenAt: model.LastSeenAt, ExpiresAt: model.ExpiresAt,
		RevokedAt: model.RevokedAt, UserAgent: model.UserAgent, IPAddress: model.IPAddress,
	})
}
