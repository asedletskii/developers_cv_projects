package ports

import (
	"context"
	"errors"
)

const MaxTokenLength = 2048

var ErrTokenVerificationFailed = errors.New("anti-bot token verification failed")

type TokenVerifier interface {
	Verify(ctx context.Context, token string, remoteIP string) error
}
