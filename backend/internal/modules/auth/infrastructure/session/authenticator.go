package session

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/auth/application"
	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/application/queries"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/platform/http/middleware"
)

const CookieName = "__Host-session"

type Authenticator struct {
	currentSession *queries.CurrentSessionHandler
	cookieName     string
	observer       ports.AuthObserver
}

func NewAuthenticator(currentSession *queries.CurrentSessionHandler, observers ...ports.AuthObserver) *Authenticator {
	var observer ports.AuthObserver
	if len(observers) > 0 {
		observer = observers[0]
	}
	return &Authenticator{currentSession: currentSession, cookieName: CookieName, observer: observer}
}

func (a *Authenticator) Authenticate(ctx context.Context, request *http.Request) (uuid.UUID, error) {
	if a == nil || a.currentSession == nil {
		return uuid.Nil, fmt.Errorf("session authenticator dependencies are nil")
	}
	if request == nil {
		return uuid.Nil, fmt.Errorf("session authenticator request is nil")
	}

	cookieName := strings.TrimSpace(a.cookieName)
	if cookieName == "" {
		return uuid.Nil, fmt.Errorf("session authenticator cookie name is empty")
	}
	cookie, err := request.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			if a.observer != nil {
				a.observer.SessionRejected(ctx, application.FailureReasonSessionMissing)
			}
			return uuid.Nil, middleware.ErrUnauthenticated
		}
		return uuid.Nil, fmt.Errorf("read auth session cookie: %w", err)
	}

	rawToken := strings.TrimSpace(cookie.Value)
	if rawToken == "" {
		if a.observer != nil {
			a.observer.SessionRejected(ctx, application.FailureReasonSessionMissing)
		}
		return uuid.Nil, middleware.ErrUnauthenticated
	}
	token, err := base64.RawURLEncoding.DecodeString(rawToken)
	if err != nil || len(token) == 0 {
		if a.observer != nil {
			a.observer.SessionRejected(ctx, application.FailureReasonSessionMissing)
		}
		return uuid.Nil, middleware.ErrUnauthenticated
	}

	rawHash := sha256.Sum256(token)
	hash, err := valueobjects.NewSessionTokenHash(rawHash[:])
	if err != nil {
		if a.observer != nil {
			a.observer.SessionRejected(ctx, application.FailureReasonSessionMissing)
		}
		return uuid.Nil, middleware.ErrUnauthenticated
	}

	view, err := a.currentSession.Handle(ctx, hash)
	if err != nil {
		if isUnauthenticatedSessionError(err) {
			if a.observer != nil {
				a.observer.SessionRejected(ctx, sessionFailureReason(err))
			}
			return uuid.Nil, middleware.ErrUnauthenticated
		}
		return uuid.Nil, fmt.Errorf("authenticate session: %w", err)
	}
	if view.UserID.IsZero() {
		return uuid.Nil, fmt.Errorf("authenticate session: returned empty user id")
	}
	return view.UserID.UUID(), nil
}

func sessionFailureReason(err error) string {
	if errors.Is(err, application.ErrSessionRevoked) {
		return application.FailureReasonSessionRevoked
	}
	if errors.Is(err, application.ErrSessionIdleExpired) {
		return application.FailureReasonSessionIdleExpired
	}
	if errors.Is(err, application.ErrSessionExpired) {
		return application.FailureReasonSessionExpired
	}
	return application.FailureReasonSessionMissing
}

func isUnauthenticatedSessionError(err error) bool {
	return errors.Is(err, application.ErrSessionInvalid) ||
		errors.Is(err, application.ErrSessionExpired) ||
		errors.Is(err, application.ErrSessionIdleExpired) ||
		errors.Is(err, application.ErrSessionRevoked)
}
