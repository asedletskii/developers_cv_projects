package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	authports "developers-cv/backend/internal/modules/auth/application/ports"
	authqueries "developers-cv/backend/internal/modules/auth/application/queries"
	authmock "developers-cv/backend/internal/modules/auth/infrastructure/mock"
	authobservability "developers-cv/backend/internal/modules/auth/infrastructure/observability"
	authsession "developers-cv/backend/internal/modules/auth/infrastructure/session"
	resumelinkanalytics "developers-cv/backend/internal/modules/resume_link/infrastructure/analytics"
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
			if errors.Is(err, middleware.ErrUnauthenticated) || errors.Is(err, authmock.ErrUnauthenticated) {
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
		return authmock.NewMockAuthenticator(reader), nil
	case "telegram":
		return authsession.NewAuthenticator(currentSession, cfg.Cookie.Name, observer), nil
	default:
		return nil, errors.New("unsupported authentication mode")
	}
}

func newAuthObserver(log *zap.Logger, metrics *telemetry.Metrics, serviceName string) authports.AuthObserver {
	return authobservability.NewObserver(log, metrics, serviceName)
}

func newOptionalAuthenticate(cfg config.AuthConfig, authenticator runtimeAuthenticator) middleware.AuthenticateFunc {
	return func(ctx context.Context, request *http.Request) (uuid.UUID, error) {
		if request == nil {
			return uuid.Nil, middleware.ErrUnauthenticated
		}
		hasCredentials := false
		switch cfg.Mode {
		case "mock":
			hasCredentials = request.Header.Get(authmock.MockUserIDHeader) != ""
		case "telegram":
			_, err := request.Cookie(cfg.Cookie.Name)
			hasCredentials = err == nil
		}
		if !hasCredentials {
			return uuid.Nil, middleware.ErrUnauthenticated
		}
		return newAuthenticate(authenticator)(ctx, request)
	}
}

func newAnalyticsVisitorIdentifier(cfg config.Config) (*resumelinkanalytics.VisitorIdentifier, error) {
	return resumelinkanalytics.NewVisitorIdentifier(cfg.Analytics.VisitorSecret, cfg.App.Env == "production")
}

func newAnalyticsTrafficClassifier(cfg config.Config) *resumelinkanalytics.TrafficClassifier {
	return resumelinkanalytics.NewTrafficClassifier(cfg.Analytics.SuspectBurst, cfg.Analytics.SuspectWindow)
}
