package application

import "errors"

var (
	ErrInvalidAuthFlow            = errors.New("invalid auth flow")
	ErrStateInvalid               = errors.New("invalid auth state")
	ErrBrowserBindingInvalid      = errors.New("invalid auth browser binding")
	ErrProviderUnavailable        = errors.New("auth provider unavailable")
	ErrIdentityVerificationFailed = errors.New("identity verification failed")
	ErrSessionInvalid             = errors.New("invalid auth session")
	ErrSessionExpired             = errors.New("auth session expired")
	ErrSessionIdleExpired         = errors.New("auth session idle expired")
	ErrSessionRevoked             = errors.New("auth session revoked")
	ErrConcurrentIdentityCreation = errors.New("concurrent identity creation")
)

type identityVerificationError struct {
	cause error
}

func (e identityVerificationError) Error() string {
	return ErrIdentityVerificationFailed.Error()
}

func (e identityVerificationError) Unwrap() []error {
	return []error{ErrIdentityVerificationFailed, e.cause}
}

func WrapIdentityVerificationError(err error) error {
	if err == nil {
		return ErrIdentityVerificationFailed
	}
	return identityVerificationError{cause: err}
}
