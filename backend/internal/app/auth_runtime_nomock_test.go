//go:build !devauth

package app

import (
	"testing"

	"developers-cv/backend/internal/platform/config"

	"github.com/stretchr/testify/require"
)

func TestNewRuntimeAuthenticatorRejectsMockWhenDevauthTagIsAbsent(t *testing.T) {
	_, err := newRuntimeAuthenticator(config.AuthConfig{Mode: "mock"}, nil, nil, nil)
	require.ErrorContains(t, err, "mock authentication is not compiled into this binary")
}
