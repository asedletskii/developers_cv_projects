package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"developers-cv/backend/internal/platform/config"
)

type Server struct {
	server *http.Server
	cfg    config.HTTPConfig
	log    *zap.Logger
}

func NewServer(cfg config.HTTPConfig, handler http.Handler, log *zap.Logger) *Server {
	return &Server{
		server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:           handler,
			ReadTimeout:       cfg.ReadTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			MaxHeaderBytes:    1 << 20,
		},
		cfg: cfg,
		log: log,
	}
}

func (s *Server) Run() error {
	s.log.Info("http server starting", zap.String("address", s.server.Addr))
	var err error
	if s.cfg.TLS.Enabled {
		s.server.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		err = s.server.ListenAndServeTLS(s.cfg.TLS.CertFile, s.cfg.TLS.KeyFile)
	} else {
		err = s.server.ListenAndServe()
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(shutdownCtx)
}
