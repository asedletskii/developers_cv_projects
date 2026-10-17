package config

import "time"

func parseTurnstileConfig() (TurnstileConfig, error) {
	enabled, err := envBool("TURNSTILE_ENABLED", false)
	if err != nil {
		return TurnstileConfig{}, err
	}
	timeout, err := envDuration("TURNSTILE_TIMEOUT", 5*time.Second)
	if err != nil {
		return TurnstileConfig{}, err
	}

	return TurnstileConfig{
		Enabled:        enabled,
		SecretKey:      envString("TURNSTILE_SECRET_KEY", ""),
		ExpectedHost:   envString("TURNSTILE_EXPECTED_HOSTNAME", ""),
		ExpectedAction: envString("TURNSTILE_EXPECTED_ACTION", "contact_view"),
		Timeout:        timeout,
	}, nil
}
