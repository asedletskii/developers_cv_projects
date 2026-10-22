//go:build devauth

package app

import (
	"testing"

	"developers-cv/backend/internal/devauth"
	"developers-cv/backend/internal/platform/config"

	"github.com/stretchr/testify/require"
)

func TestNewRuntimeAuthenticatorSupportsMockWithDevauthTag(t *testing.T) {
	authenticator, err := newRuntimeAuthenticator(config.AuthConfig{Mode: "mock"}, nil, nil, nil)
	require.NoError(t, err)
	require.IsType(t, &devauth.MockAuthenticator{}, authenticator)
}
