package app

import (
	"net/http"

	resumelinkports "developers-cv/backend/internal/modules/resume_link/application/ports"
	resumelinksecurity "developers-cv/backend/internal/modules/resume_link/infrastructure/security"
	"developers-cv/backend/internal/platform/config"
)

func newTurnstileVerifier(cfg config.Config) (resumelinkports.TokenVerifier, error) {
	if !cfg.Turnstile.Enabled {
		return nil, nil
	}

	return resumelinksecurity.NewTurnstileVerifier(resumelinksecurity.Config{
		Secret:         cfg.Turnstile.SecretKey,
		Timeout:        cfg.Turnstile.Timeout,
		ExpectedHost:   cfg.Turnstile.ExpectedHost,
		ExpectedAction: cfg.Turnstile.ExpectedAction,
	}, &http.Client{})
}
