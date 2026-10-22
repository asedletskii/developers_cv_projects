package domain

import (
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Session struct {
	id         sharedvalueobjects.ID
	userID     sharedvalueobjects.ID
	tokenHash  valueobjects.SessionTokenHash
	createdAt  time.Time
	lastSeenAt time.Time
	expiresAt  time.Time
	revokedAt  *time.Time
	userAgent  *string
	ipAddress  *string
}

type SessionParams struct {
	UserID     sharedvalueobjects.ID
	TokenHash  valueobjects.SessionTokenHash
	CreatedAt  time.Time
	LastSeenAt time.Time
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	UserAgent  *string
	IPAddress  *string
}

type SessionSnapshot struct {
	ID         sharedvalueobjects.ID
	UserID     sharedvalueobjects.ID
	TokenHash  valueobjects.SessionTokenHash
	CreatedAt  time.Time
	LastSeenAt time.Time
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	UserAgent  *string
	IPAddress  *string
}

func NewSession(params SessionParams) (*Session, error) {
	id := sharedvalueobjects.NewID()
	createdAt := params.CreatedAt.UTC()
	lastSeenAt := params.LastSeenAt.UTC()
	expiresAt := params.ExpiresAt.UTC()
	var revokedAtInput *time.Time
	if params.RevokedAt != nil {
		value := params.RevokedAt.UTC()
		revokedAtInput = &value
	}
	revokedAt, err := validateSession(id, params.UserID, params.TokenHash, createdAt, lastSeenAt, expiresAt, revokedAtInput)
	if err != nil {
		return nil, err
	}
	return &Session{
		id: id, userID: params.UserID, tokenHash: params.TokenHash,
		createdAt: createdAt, lastSeenAt: lastSeenAt, expiresAt: expiresAt,
		revokedAt: revokedAt, userAgent: cloneString(params.UserAgent), ipAddress: cloneString(params.IPAddress),
	}, nil
}

func RestoreSession(snapshot SessionSnapshot) (*Session, error) {
	revokedAt, err := validateSession(
		snapshot.ID, snapshot.UserID, snapshot.TokenHash,
		snapshot.CreatedAt, snapshot.LastSeenAt, snapshot.ExpiresAt, snapshot.RevokedAt,
	)
	if err != nil {
		return nil, err
	}
	return &Session{
		id: snapshot.ID, userID: snapshot.UserID, tokenHash: snapshot.TokenHash,
		createdAt: snapshot.CreatedAt, lastSeenAt: snapshot.LastSeenAt, expiresAt: snapshot.ExpiresAt,
		revokedAt: revokedAt, userAgent: cloneString(snapshot.UserAgent), ipAddress: cloneString(snapshot.IPAddress),
	}, nil
}

func validateSession(
	id, userID sharedvalueobjects.ID,
	tokenHash valueobjects.SessionTokenHash,
	createdAt, lastSeenAt, expiresAt time.Time,
	revokedAt *time.Time,
) (*time.Time, error) {
	if id.IsZero() || userID.IsZero() {
		return nil, fmt.Errorf("%w: invalid session identity", ErrInvalidSession)
	}
	if createdAt.IsZero() || lastSeenAt.IsZero() || expiresAt.IsZero() {
		return nil, fmt.Errorf("%w: session timestamps must be set", ErrInvalidSession)
	}
	if lastSeenAt.Before(createdAt) || !expiresAt.After(createdAt) {
		return nil, fmt.Errorf("%w: invalid session timestamp order", ErrInvalidSession)
	}
	var restoredRevokedAt *time.Time
	if revokedAt != nil {
		value := *revokedAt
		if value.IsZero() || value.Before(createdAt) {
			return nil, fmt.Errorf("%w: invalid revoked timestamp", ErrInvalidSession)
		}
		restoredRevokedAt = &value
	}
	return restoredRevokedAt, nil
}

func (s *Session) ID() sharedvalueobjects.ID                { return s.id }
func (s *Session) UserID() sharedvalueobjects.ID            { return s.userID }
func (s *Session) TokenHash() valueobjects.SessionTokenHash { return s.tokenHash }
func (s *Session) CreatedAt() time.Time                     { return s.createdAt }
func (s *Session) LastSeenAt() time.Time                    { return s.lastSeenAt }
func (s *Session) ExpiresAt() time.Time                     { return s.expiresAt }
func (s *Session) RevokedAt() *time.Time                    { return cloneTime(s.revokedAt) }
func (s *Session) UserAgent() *string                       { return cloneString(s.userAgent) }
func (s *Session) IPAddress() *string                       { return cloneString(s.ipAddress) }

func (s *Session) IsRevoked() bool              { return s != nil && s.revokedAt != nil }
func (s *Session) IsExpired(now time.Time) bool { return s == nil || !now.UTC().Before(s.expiresAt) }
func (s *Session) IsIdleExpired(now time.Time, idleTTL time.Duration) bool {
	if s == nil || idleTTL < 0 {
		return true
	}
	return !now.UTC().Before(s.lastSeenAt.Add(idleTTL))
}
func (s *Session) IsActive(now time.Time, idleTTL time.Duration) bool {
	if s == nil || s.IsRevoked() || s.IsExpired(now) {
		return false
	}
	return !s.IsIdleExpired(now, idleTTL)
}

func (s *Session) Touch(now time.Time) error {
	if s == nil {
		return fmt.Errorf("%w: session is nil", ErrInvalidSession)
	}
	now = now.UTC()
	if now.Before(s.lastSeenAt) {
		return fmt.Errorf("%w: last_seen_at cannot move backwards", ErrInvalidSession)
	}
	s.lastSeenAt = now
	return nil
}

func (s *Session) Revoke(now time.Time) error {
	if s == nil {
		return fmt.Errorf("%w: session is nil", ErrInvalidSession)
	}
	if s.revokedAt != nil {
		return nil
	}
	now = now.UTC()
	if now.IsZero() || now.Before(s.createdAt) {
		return fmt.Errorf("%w: invalid revoke timestamp", ErrInvalidSession)
	}
	s.revokedAt = &now
	return nil
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
