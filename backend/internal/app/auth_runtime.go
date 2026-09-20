package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	analyticsadapters "developers-cv/backend/internal/modules/analytics/infrastructure/analytics"
	authports "developers-cv/backend/internal/modules/auth/application/ports"
	authqueries "developers-cv/backend/internal/modules/auth/application/queries"
	authobservability "developers-cv/backend/internal/modules/auth/infrastructure/observability"
	authsession "developers-cv/backend/internal/modules/auth/infrastructure/session"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/http/middleware"

	"go.uber.org/zap"

	"developers-cv/backend/internal/platform/telemetry"
)

type runtimeAuthenticator interface {
	Authenticate(context.Context, *http.Request) (uuid.UUID, error)
}

func newAuthenticate(authenticator runtimeAuthenticator) middleware.AuthenticateFunc {
	return func(ctx context.Context, request *http.Request) (uuid.UUID, error) {
		userID, err := authenticator.Authenticate(ctx, request)
		if err != nil {
			if errors.Is(err, middleware.ErrUnauthenticated) || isMockUnauthenticated(err) {
				return uuid.Nil, middleware.ErrUnauthenticated
			}
			return uuid.Nil, err
		}
		return userID, nil
	}
}

func newRuntimeAuthenticator(cfg config.AuthConfig, currentSession *authqueries.CurrentSessionHandler, reader userqueries.UserReader, observer authports.AuthObserver) (runtimeAuthenticator, error) {
	switch cfg.Mode {
	case "mock":
		return newMockAuthenticator(reader)
	case "telegram":
		return authsession.NewAuthenticator(currentSession, cfg.Cookie.Name, observer), nil
	default:
		return nil, errors.New("unsupported authentication mode")
	}
}

func newAuthObserver(log *zap.Logger, metrics *telemetry.Metrics, cfg config.Config) authports.AuthObserver {
	return authobservability.NewObserver(log, metrics, cfg.Telemetry.ServiceName)
}

type optionalAuthenticateFunc middleware.AuthenticateFunc

func newOptionalAuthenticate(cfg config.AuthConfig, authenticator runtimeAuthenticator) optionalAuthenticateFunc {
	return optionalAuthenticateFunc(func(ctx context.Context, request *http.Request) (uuid.UUID, error) {
		if request == nil {
			return uuid.Nil, middleware.ErrUnauthenticated
		}
		hasCredentials := false
		switch cfg.Mode {
		case "mock":
			hasCredentials = hasMockCredentials(request)
		case "telegram":
			_, err := request.Cookie(cfg.Cookie.Name)
			hasCredentials = err == nil
		}
		if !hasCredentials {
			return uuid.Nil, middleware.ErrUnauthenticated
		}
		return newAuthenticate(authenticator)(ctx, request)
	})
}

func newAnalyticsVisitorIdentifier(cfg config.Config) (*analyticsadapters.VisitorIdentifier, error) {
	return analyticsadapters.NewVisitorIdentifier(cfg.Analytics.VisitorSecret, cfg.App.Env == "production")
}

func newAnalyticsTrafficClassifier(cfg config.Config) *analyticsadapters.TrafficClassifier {
	return analyticsadapters.NewTrafficClassifier(cfg.Analytics.SuspectBurst, cfg.Analytics.SuspectWindow)
}
