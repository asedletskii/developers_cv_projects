package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

	"developers-cv/backend/internal/platform/config"
)

type Server struct {
	server         *http.Server
	internalServer *http.Server
	cfg            config.HTTPConfig
	log            *zap.Logger
}

func NewServer(cfg config.HTTPConfig, handler http.Handler, log *zap.Logger, internalHandler ...http.Handler) *Server {
	s := &Server{
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), Handler: handler,
			ReadTimeout: cfg.ReadTimeout, ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			WriteTimeout: cfg.WriteTimeout, IdleTimeout: cfg.IdleTimeout, MaxHeaderBytes: 1 << 20,
		},
		cfg: cfg, log: log,
	}
	if len(internalHandler) > 0 && internalHandler[0] != nil && cfg.MetricsAddr != "" {
		s.internalServer = &http.Server{
			Addr: cfg.MetricsAddr, Handler: internalHandler[0],
			ReadTimeout: cfg.ReadTimeout, ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			WriteTimeout: cfg.WriteTimeout, IdleTimeout: cfg.IdleTimeout, MaxHeaderBytes: 1 << 20,
		}
	}
	return s
}

func (s *Server) Run() error {
	s.log.Info("http server starting", zap.String("address", s.server.Addr))
	errCh := make(chan error, 2)
	serve := func(server *http.Server, tlsEnabled bool, cert, key string) {
		var err error
		if tlsEnabled {
			server.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = server.ListenAndServeTLS(cert, key)
		} else {
			err = server.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("listen and serve %s: %w", server.Addr, err)
		} else {
			errCh <- nil
		}
	}
	go serve(s.server, s.cfg.TLS.Enabled, s.cfg.TLS.CertFile, s.cfg.TLS.KeyFile)
	if s.internalServer != nil {
		s.log.Info("internal observability server starting", zap.String("address", s.internalServer.Addr))
		go serve(s.internalServer, false, "", "")
	}
	err := <-errCh
	if err != nil {
		if closeErr := s.server.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close public HTTP server after listener failure: %w", closeErr))
		}
		if s.internalServer != nil {
			if closeErr := s.internalServer.Close(); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("close internal HTTP server after listener failure: %w", closeErr))
			}
		}
	}
	return err
}

func effectiveShutdownTimeout(configured time.Duration) time.Duration {
	if configured <= 0 {
		return 10 * time.Second
	}
	return configured
}

func (s *Server) Shutdown(ctx context.Context) error {
	timeout := effectiveShutdownTimeout(s.cfg.ShutdownTimeout)
	shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	shutdown := func(server *http.Server) {
		defer wg.Done()
		errs <- server.Shutdown(shutdownCtx)
	}
	wg.Add(1)
	go shutdown(s.server)
	if s.internalServer != nil {
		wg.Add(1)
		go shutdown(s.internalServer)
	}
	wg.Wait()
	close(errs)
	var result error
	for err := range errs {
		result = errors.Join(result, err)
	}
	return result
}
