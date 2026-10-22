package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	analyticshttp "developers-cv/backend/internal/modules/analytics/transport/http"
	authhttp "developers-cv/backend/internal/modules/auth/transport/http"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	resumelinkhttp "developers-cv/backend/internal/modules/resume_link/transport/http"
	userhttp "developers-cv/backend/internal/modules/user/transport/http"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/health"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
	"developers-cv/backend/internal/platform/telemetry"
)

type Router struct {
	engine             *gin.Engine
	internal           *gin.Engine
	health             *health.Service
	metrics            *telemetry.Metrics
	corsAllowedOrigins []string
}

func NewRouter(
	cfg config.Config,
	logger *zap.Logger,
	healthService *health.Service,
	metrics *telemetry.Metrics,
	resumeHandler *resumehttp.Handler,
	userHandler *userhttp.Handler,
	resumeLinkHandler *resumelinkhttp.Handler,
	authHandler *authhttp.Handler,
	authenticate middleware.AuthenticateFunc,
	optionalAuthenticate middleware.AuthenticateFunc,
	analyticsHandlers ...*analyticshttp.Handler,
) (*Router, error) {
	if len(analyticsHandlers) > 1 {
		return nil, fmt.Errorf("expected at most one Analytics HTTP handler, got %d", len(analyticsHandlers))
	}
	var analyticsHandler *analyticshttp.Handler
	if len(analyticsHandlers) == 1 {
		analyticsHandler = analyticsHandlers[0]
	}
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	if err := engine.SetTrustedProxies(cfg.HTTP.TrustedProxies); err != nil {
		return nil, fmt.Errorf("configure trusted proxies: %w", err)
	}

	engine.Use(
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Recovery(logger),
		middleware.SecurityHeaders(cfg.HTTP.TLS.Enabled),
		middleware.CORS(cfg.CORS),
		middleware.BodyLimit(cfg.HTTP.MaxBodyBytes),
		telemetry.HTTPMiddleware(cfg.Telemetry.ServiceName),
	)
	if metrics != nil {
		engine.Use(metrics.Middleware())
	}
	if cfg.HTTP.RateLimit.Enabled {
		engine.Use(middleware.NewRateLimiter(cfg.HTTP.RateLimit.Rate, cfg.HTTP.RateLimit.Burst).Middleware())
	}

	var authStartRateLimiter, authCallbackRateLimiter gin.HandlerFunc
	var contactViewRateLimiter, analyticsEventsRateLimiter gin.HandlerFunc
	if cfg.HTTP.ContactViewRateLimit.Enabled {
		contactViewRateLimiter = middleware.NewRateLimiter(
			cfg.HTTP.ContactViewRateLimit.Rate/60,
			cfg.HTTP.ContactViewRateLimit.Burst,
		).Middleware()
	}
	if cfg.HTTP.AnalyticsEventsRateLimit.Enabled {
		analyticsEventsRateLimiter = middleware.NewRateLimiter(
			cfg.HTTP.AnalyticsEventsRateLimit.Rate/60,
			cfg.HTTP.AnalyticsEventsRateLimit.Burst,
		).Middleware()
	}

	if cfg.Auth.RateLimit.Enabled {
		authStartRateLimiter = middleware.NewRateLimiter(
			cfg.Auth.RateLimit.StartRatePerMinute/60,
			cfg.Auth.RateLimit.StartBurst,
		).Middleware()
		authCallbackRateLimiter = middleware.NewRateLimiter(
			cfg.Auth.RateLimit.CallbackRatePerMinute/60,
			cfg.Auth.RateLimit.CallbackBurst,
		).Middleware()
	}

	internal := gin.New()
	internal.Use(middleware.RequestID(), middleware.Logger(logger), middleware.Recovery(logger))
	if metrics != nil {
		internal.Use(metrics.Middleware())
	}
	router := &Router{
		engine: engine, internal: internal, health: healthService, metrics: metrics,
		corsAllowedOrigins: append([]string(nil), cfg.CORS.AllowedOrigins...),
	}
	router.registerRoutes(
		resumeHandler,
		userHandler,
		resumeLinkHandler,
		analyticsHandler,
		authHandler,
		authenticate,
		authStartRateLimiter,
		authCallbackRateLimiter,
		contactViewRateLimiter,
		analyticsEventsRateLimiter,
		optionalAuthenticate,
	)
	return router, nil
}

func (r *Router) Handler() http.Handler {
	return r.engine
}

func (r *Router) registerRoutes(
	resumeHandler *resumehttp.Handler,
	userHandler *userhttp.Handler,
	resumeLinkHandler *resumelinkhttp.Handler,
	analyticsHandler *analyticshttp.Handler,
	authHandler *authhttp.Handler,
	authenticate middleware.AuthenticateFunc,
	authStartRateLimiter gin.HandlerFunc,
	authCallbackRateLimiter gin.HandlerFunc,
	contactViewRateLimiter gin.HandlerFunc,
	analyticsEventsRateLimiter gin.HandlerFunc,
	optionalAuthenticate middleware.AuthenticateFunc,
) {
	engine := r.engine
	engine.GET("/health", r.healthHandler)

	api := engine.Group("/api/v1")
	csrf := middleware.CSRFProtection(r.corsAllowedOrigins)
	auth := api.Group("/auth")
	auth.Use(csrf)
	if authHandler != nil {
		authhttp.RegisterPublicRoutes(authHandler, auth, authStartRateLimiter, authCallbackRateLimiter)
		authhttp.RegisterLogoutRoute(authHandler, auth)
		authProtected := auth.Group("")
		authProtected.Use(middleware.Authentication(authenticate))
		authhttp.RegisterProtectedRoutes(authHandler, authProtected)
	}
	lk := api.Group("/lk")
	lk.Use(middleware.Authentication(authenticate))
	lk.Use(csrf)
	userhttp.RegisterRoutes(userHandler, lk)
	resumehttp.RegisterRoutes(resumeHandler, lk)
	public := api.Group("")
	public.Use(middleware.OptionalAuthentication(optionalAuthenticate))
	anonymous := api.Group("")
	resumelinkhttp.RegisterRoutes(resumeLinkHandler, public, anonymous, lk, contactViewRateLimiter)
	analyticshttp.RegisterRoutes(analyticsHandler, anonymous, lk, analyticsEventsRateLimiter)
	engine.NoRoute(func(c *gin.Context) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Route not found"))
	})
	r.internal.GET("/ready", r.readyHandler)
	if r.metrics != nil {
		r.internal.GET("/metrics", gin.WrapH(r.metrics.Handler()))
	}
}

func (r *Router) InternalHandler() http.Handler {
	return r.internal
}

func (r *Router) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (r *Router) readyHandler(c *gin.Context) {
	results := r.health.Check(c.Request.Context())
	checks := make(map[string]string, len(results))
	ready := true
	for name, err := range results {
		if err != nil {
			ready = false
			checks[name] = "failed"
			continue
		}
		checks[name] = "ok"
	}
	if !ready {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "checks": checks})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready", "checks": checks})
}
