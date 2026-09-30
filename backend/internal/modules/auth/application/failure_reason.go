package application

import (
	"errors"

	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain"
)

const (
	FailureReasonInvalidRequest      = "invalid_request"
	FailureReasonProviderDenied      = "provider_denied"
	FailureReasonProviderUnavailable = "provider_unavailable"
	FailureReasonTokenExchangeFailed = "token_exchange_failed"
	FailureReasonIDTokenInvalid      = "id_token_invalid"
	FailureReasonNonceMismatch       = "nonce_mismatch"
	FailureReasonStateInvalid        = "state_invalid"
	FailureReasonTransactionExpired  = "transaction_expired"
	FailureReasonTransactionConsumed = "transaction_consumed"
	FailureReasonIdentityConflict    = "identity_conflict"
	FailureReasonDatabaseError       = "database_error"
	FailureReasonSessionMissing      = "session_missing"
	FailureReasonSessionRevoked      = "session_revoked"
	FailureReasonSessionExpired      = "session_expired"
	FailureReasonSessionIdleExpired  = "session_idle_expired"
)

func FailureReason(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, ErrProviderUnavailable):
		return FailureReasonProviderUnavailable
	case errors.Is(err, ports.ErrTokenExchangeFailed):
		return FailureReasonTokenExchangeFailed
	case errors.Is(err, ports.ErrNonceMismatch):
		return FailureReasonNonceMismatch
	case errors.Is(err, ports.ErrIDTokenInvalid):
		return FailureReasonIDTokenInvalid
	case errors.Is(err, ErrConcurrentIdentityCreation), errors.Is(err, domain.ErrDuplicateEntity):
		return FailureReasonIdentityConflict
	case errors.Is(err, domain.ErrAuthTransactionExpired):
		return FailureReasonTransactionExpired
	case errors.Is(err, domain.ErrAuthTransactionConsumed):
		return FailureReasonTransactionConsumed
	case errors.Is(err, ErrSessionRevoked):
		return FailureReasonSessionRevoked
	case errors.Is(err, ErrSessionIdleExpired):
		return FailureReasonSessionIdleExpired
	case errors.Is(err, ErrSessionExpired):
		return FailureReasonSessionExpired
	case errors.Is(err, ErrSessionInvalid):
		return FailureReasonSessionMissing
	case errors.Is(err, ErrStateInvalid):
		return FailureReasonStateInvalid
	case errors.Is(err, ErrInvalidAuthFlow):
		return FailureReasonInvalidRequest
	case errors.Is(err, ErrIdentityVerificationFailed):
		return FailureReasonIDTokenInvalid
	default:
		return FailureReasonDatabaseError
	}
}
