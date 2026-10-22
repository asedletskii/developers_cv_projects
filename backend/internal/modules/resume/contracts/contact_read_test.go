package contracts

import (
	"testing"

	"github.com/stretchr/testify/require"

	contactkey "developers-cv/backend/internal/shared/domain/contact"
)

func TestContactKeyAliasesSharedKey(t *testing.T) {
	require.Equal(t, contactkey.Phone, ContactKeyPhone)
	require.Equal(t, contactkey.Telegram, ContactKeyTelegram)
	require.Equal(t, contactkey.Email, ContactKeyEmail)
}
