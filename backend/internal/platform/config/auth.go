package config

import "time"

func parseAuthConfig(appEnv string) (AuthConfig, error) {
	mode := envString("AUTH_MODE", "mock")
	absoluteTTL, err := envDuration("AUTH_SESSION_ABSOLUTE_TTL", 720*time.Hour)
	if err != nil {
		return AuthConfig{}, err
	}
	idleTTL, err := envDuration("AUTH_SESSION_IDLE_TTL", 168*time.Hour)
	if err != nil {
		return AuthConfig{}, err
	}
	touchInterval, err := envDuration("AUTH_SESSION_TOUCH_INTERVAL", 5*time.Minute)
	if err != nil {
		return AuthConfig{}, err
	}
	transactionTTL, err := envDuration("AUTH_TRANSACTION_TTL", 5*time.Minute)
	if err != nil {
		return AuthConfig{}, err
	}

	secureDefault := appEnv == "staging" || appEnv == "production"
	cookieNameDefault := "session"
	if secureDefault {
		cookieNameDefault = "__Host-session"
	}
	cookieSecure, err := envBool("AUTH_COOKIE_SECURE", secureDefault)
	if err != nil {
		return AuthConfig{}, err
	}
	rateEnabled, err := envBool("AUTH_RATE_LIMIT_ENABLED", true)
	if err != nil {
		return AuthConfig{}, err
	}
	startRate, err := envFloat("AUTH_RATE_LIMIT_START_PER_MINUTE", 5)
	if err != nil {
		return AuthConfig{}, err
	}
	startBurst, err := envInt("AUTH_RATE_LIMIT_START_BURST", 5)
	if err != nil {
		return AuthConfig{}, err
	}
	callbackRate, err := envFloat("AUTH_RATE_LIMIT_CALLBACK_PER_MINUTE", 20)
	if err != nil {
		return AuthConfig{}, err
	}
	callbackBurst, err := envInt("AUTH_RATE_LIMIT_CALLBACK_BURST", 20)
	if err != nil {
		return AuthConfig{}, err
	}

	return AuthConfig{
		Mode: mode,
		Telegram: TelegramAuthConfig{
			ClientID:     envString("AUTH_TELEGRAM_CLIENT_ID", ""),
			ClientSecret: envString("AUTH_TELEGRAM_CLIENT_SECRET", ""),
			RedirectURL:  envString("AUTH_TELEGRAM_REDIRECT_URL", "http://localhost:8080/api/v1/auth/telegram/callback"),
			IssuerURL:    envString("AUTH_TELEGRAM_ISSUER_URL", "https://oauth.telegram.org"),
		},
		Session: SessionAuthConfig{
			AbsoluteTTL:   absoluteTTL,
			IdleTTL:       idleTTL,
			TouchInterval: touchInterval,
		},
		Cookie:      CookieAuthConfig{Name: envString("AUTH_COOKIE_NAME", cookieNameDefault), Secure: cookieSecure},
		Transaction: TransactionAuthConfig{TTL: transactionTTL},
		RateLimit: AuthRateLimitConfig{
			Enabled:               rateEnabled,
			StartRatePerMinute:    startRate,
			StartBurst:            startBurst,
			CallbackRatePerMinute: callbackRate,
			CallbackBurst:         callbackBurst,
		},
	}, nil
}
