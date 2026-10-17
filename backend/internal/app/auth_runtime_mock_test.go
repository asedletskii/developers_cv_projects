//go:build devauth

package app

import (
	"testing"

	authmock "developers-cv/backend/internal/modules/auth/infrastructure/mock"
	"developers-cv/backend/internal/platform/config"

	"github.com/stretchr/testify/require"
)

func TestNewRuntimeAuthenticatorSupportsMockWithDevauthTag(t *testing.T) {
	authenticator, err := newRuntimeAuthenticator(config.AuthConfig{Mode: "mock"}, nil, nil, nil)
	require.NoError(t, err)
	require.IsType(t, &authmock.MockAuthenticator{}, authenticator)
}
