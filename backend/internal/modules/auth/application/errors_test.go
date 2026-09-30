package application

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapIdentityVerificationErrorPreservesCauseWithoutExposingDetails(t *testing.T) {
	secret := errors.New("telegram client_secret=super-secret")

	err := WrapIdentityVerificationError(secret)

	require.ErrorIs(t, err, ErrIdentityVerificationFailed)
	require.ErrorIs(t, err, secret)
	require.Equal(t, ErrIdentityVerificationFailed.Error(), err.Error())
	require.NotContains(t, err.Error(), "super-secret")
}
