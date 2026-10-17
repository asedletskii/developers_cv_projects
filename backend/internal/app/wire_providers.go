package app

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	authcommands "developers-cv/backend/internal/modules/auth/application/commands"
	authports "developers-cv/backend/internal/modules/auth/application/ports"
	authqueries "developers-cv/backend/internal/modules/auth/application/queries"
	authhttp "developers-cv/backend/internal/modules/auth/transport/http"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	resumelinkhttp "developers-cv/backend/internal/modules/resume_link/transport/http"
	userhttp "developers-cv/backend/internal/modules/user/transport/http"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/database"
	"developers-cv/backend/internal/platform/health"
	httpserver "developers-cv/backend/internal/platform/http"
	"developers-cv/backend/internal/platform/http/middleware"
	"developers-cv/backend/internal/platform/logger"
	"developers-cv/backend/internal/platform/telemetry"
)

func newLogger(cfg config.Config) (*zap.Logger, error) {
	return logger.New(cfg.Log)
}

func newDatabase(cfg config.Config, log *zap.Logger) (*gorm.DB, func(), error) {
	db, err := database.NewGORMFromConfig(cfg, log)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := database.Close(db); err != nil {
			log.Error("close database after failed application initialization", zap.Error(err))
		}
	}
	return db, cleanup, nil
}

func newTelemetry(cfg config.Config, log *zap.Logger) (*telemetry.Service, func(), error) {
	service, err := telemetry.New(cfg.Telemetry, log)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if service == nil {
			return
		}
		timeout := cfg.HTTP.ShutdownTimeout
		if timeout <= 0 {
			timeout = 10 * time.Second
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := service.Shutdown(ctx); err != nil {
			log.Error("shutdown telemetry after failed application initialization", zap.Error(err))
		}
	}
	return service, cleanup, nil
}

func newMetrics(cfg config.Config) (*telemetry.Metrics, error) {
	return telemetry.NewMetrics(cfg.Prometheus.Enabled)
}

type authApplicationLogger struct {
	logger *zap.Logger
}

func (l authApplicationLogger) Warn(message string, err error) {
	if l.logger == nil {
		return
	}
	l.logger.Warn(message, zap.Error(err))
}

func newCurrentSessionHandler(sessions authports.SessionRepository, uow authports.UnitOfWork, cfg config.AuthConfig, logger *zap.Logger) *authqueries.CurrentSessionHandler {
	return authqueries.NewCurrentSessionHandler(
		sessions,
		uow,
		cfg.Session.IdleTTL,
		cfg.Session.TouchInterval,
		authApplicationLogger{logger: logger},
	)
}

func newStartTelegramLoginHandler(
	repository authports.AuthTransactionRepository,
	oidc authports.TelegramOIDC,
	tokenGenerator authports.TokenGenerator,
	uow authports.UnitOfWork,
	cfg config.AuthConfig,
) *authcommands.StartTelegramLoginHandler {
	return authcommands.NewStartTelegramLoginHandler(repository, oidc, tokenGenerator, uow, cfg.Telegram.RedirectURL, cfg.Transaction.TTL)
}

func newHTTPHandler(router *httpserver.Router) http.Handler {
	return router.Handler()
}

func newHTTPServer(cfg config.Config, handler http.Handler, log *zap.Logger, router *httpserver.Router) *httpserver.Server {
	return httpserver.NewServer(cfg.HTTP, handler, log, router.InternalHandler())
}

func newApplication(
	server *httpserver.Server,
	telemetryService *telemetry.Service,
	log *zap.Logger,
	db *gorm.DB,
	cfg config.Config,
	workers []backgroundWorker,
) *Application {
	return NewApplication(server, telemetryService, log, db, cfg.HTTP.ShutdownTimeout, workers)
}

func newRouter(
	cfg config.Config,
	log *zap.Logger,
	healthService *health.Service,
	metrics *telemetry.Metrics,
	resumeHandler *resumehttp.Handler,
	userHandler *userhttp.Handler,
	resumeLinkHandler *resumelinkhttp.Handler,
	authHandler *authhttp.Handler,
	authenticate middleware.AuthenticateFunc,
	optionalAuthenticate optionalAuthenticateFunc,
) (*httpserver.Router, error) {
	return httpserver.NewRouter(cfg, log, healthService, metrics, resumeHandler, userHandler, resumeLinkHandler, authHandler, authenticate, middleware.AuthenticateFunc(optionalAuthenticate))
}

func newAuthHandler(
	start *authcommands.StartTelegramLoginHandler,
	complete *authcommands.CompleteTelegramLoginHandler,
	logout *authcommands.LogoutHandler,
	observer authports.AuthObserver,
	tokenGenerator authports.TokenGenerator,
	cfg config.AuthConfig,
) *authhttp.Handler {
	return authhttp.NewHandler(authhttp.HandlerDependencies{
		StartTelegramLogin:   start,
		CompleteTelegram:     complete,
		Logout:               logout,
		Observer:             observer,
		TokenGenerator:       tokenGenerator,
		Cookie:               cookieConfigFromRuntime(cfg),
		BrowserBindingMaxAge: int(cfg.Transaction.TTL.Seconds()),
	})
}
