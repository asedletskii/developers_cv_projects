package application

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/application/ports"
)

func TestFailureReasonNeverUsesRawErrorText(t *testing.T) {
	sensitive := errors.New("authorization code=secret state=secret nonce=secret session_token=secret")

	reason := FailureReason(sensitive)

	require.Equal(t, FailureReasonDatabaseError, reason)
	require.NotContains(t, reason, "authorization")
	require.NotContains(t, reason, "secret")
}

func TestFailureReasonKeepsOIDCSecurityCategoriesStable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "exchange", err: ports.ErrTokenExchangeFailed, want: FailureReasonTokenExchangeFailed},
		{name: "id token", err: ports.ErrIDTokenInvalid, want: FailureReasonIDTokenInvalid},
		{name: "nonce", err: ports.ErrNonceMismatch, want: FailureReasonNonceMismatch},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, FailureReason(tt.err))
		})
	}
}
