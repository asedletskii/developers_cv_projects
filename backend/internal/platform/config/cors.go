package config

import "os"

func parseCORSConfig(authMode string) (CORSConfig, error) {
	credentials, err := envBool("CORS_ALLOW_CREDENTIALS", true)
	if err != nil {
		return CORSConfig{}, err
	}

	return CORSConfig{
		AllowedOrigins:      splitCSV(os.Getenv("CORS_ALLOWED_ORIGINS")),
		AllowCredentials:    credentials,
		AllowMockAuthHeader: authMode == "mock",
	}, nil
}
