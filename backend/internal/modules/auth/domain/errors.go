package domain

import "errors"

var (
	ErrInvalidProvider         = errors.New("invalid auth provider")
	ErrInvalidIdentity         = errors.New("invalid auth identity")
	ErrInvalidSession          = errors.New("invalid auth session")
	ErrInvalidAuthTransaction  = errors.New("invalid auth transaction")
	ErrSessionRevoked          = errors.New("auth session revoked")
	ErrSessionExpired          = errors.New("auth session expired")
	ErrAuthTransactionExpired  = errors.New("auth transaction expired")
	ErrAuthTransactionConsumed = errors.New("auth transaction consumed")
	ErrEntityNotFound          = errors.New("auth entity not found")
	ErrDuplicateEntity         = errors.New("auth duplicate entity")
)
