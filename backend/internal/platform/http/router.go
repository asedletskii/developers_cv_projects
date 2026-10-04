package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

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
	engine  *gin.Engine
	health  *health.Service
	metrics *telemetry.Metrics
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
) (*Router, error) {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	if err := engine.SetTrustedProxies(cfg.HTTP.TrustedProxies); err != nil {
		return nil, fmt.Errorf("configure trusted proxies: %w", err)
	}

	engine.Use(
		middleware.RequestID(),
		middleware.Recovery(logger),
		middleware.SecurityHeaders(cfg.HTTP.TLS.Enabled),
		middleware.CORS(cfg.CORS),
		middleware.BodyLimit(cfg.HTTP.MaxBodyBytes),
		telemetry.HTTPMiddleware(cfg.Telemetry.ServiceName),
		middleware.Logger(logger),
	)
	if metrics != nil {
		engine.Use(metrics.Middleware())
	}
	if cfg.HTTP.RateLimit.Enabled {
		engine.Use(middleware.NewRateLimiter(cfg.HTTP.RateLimit.Rate, cfg.HTTP.RateLimit.Burst).Middleware())
	}

	var authStartRateLimiter, authCallbackRateLimiter gin.HandlerFunc
	var contactViewRateLimiter gin.HandlerFunc
	if cfg.HTTP.ContactViewRateLimit.Enabled {
		contactViewRateLimiter = middleware.NewRateLimiter(
			cfg.HTTP.ContactViewRateLimit.Rate/60,
			cfg.HTTP.ContactViewRateLimit.Burst,
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

	router := &Router{engine: engine, health: healthService, metrics: metrics}
	router.registerRoutes(
		resumeHandler,
		userHandler,
		resumeLinkHandler,
		authHandler,
		authenticate,
		authStartRateLimiter,
		authCallbackRateLimiter,
		contactViewRateLimiter,
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
	authHandler *authhttp.Handler,
	authenticate middleware.AuthenticateFunc,
	authStartRateLimiter gin.HandlerFunc,
	authCallbackRateLimiter gin.HandlerFunc,
	contactViewRateLimiter gin.HandlerFunc,
) {
	engine := r.engine
	engine.GET("/health", r.healthHandler)
	engine.GET("/ready", r.readyHandler)
	if r.metrics != nil {
		engine.GET("/metrics", gin.WrapH(r.metrics.Handler()))
	}

	api := engine.Group("/api/v1")
	auth := api.Group("/auth")
	if authHandler != nil {
		authhttp.RegisterPublicRoutes(authHandler, auth, authStartRateLimiter, authCallbackRateLimiter)
		authProtected := auth.Group("")
		authProtected.Use(middleware.Authentication(authenticate))
		authhttp.RegisterProtectedRoutes(authHandler, authProtected)
	}
	lk := api.Group("/lk")
	lk.Use(middleware.Authentication(authenticate))
	userhttp.RegisterRoutes(userHandler, lk)
	resumehttp.RegisterRoutes(resumeHandler, lk)
	resumelinkhttp.RegisterRoutes(resumeLinkHandler, api, lk, contactViewRateLimiter)
	engine.NoRoute(func(c *gin.Context) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Route not found"))
	})
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
