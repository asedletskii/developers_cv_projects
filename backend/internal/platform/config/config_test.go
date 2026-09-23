package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadRequiresOriginsWhenCredentialsEnabled(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "true")

	_, err := Load()

	require.Error(t, err)
	require.Contains(t, err.Error(), "CORS_ALLOWED_ORIGINS")
}

func TestLoadParsesConfiguration(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("HTTP_MAX_BODY_BYTES", "2048")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:4173")
	t.Setenv("LOG_FORMAT", "console")

	cfg, err := Load()

	require.NoError(t, err)
	require.Equal(t, 9090, cfg.HTTP.Port)
	require.Equal(t, int64(2048), cfg.HTTP.MaxBodyBytes)
	require.Equal(
		t,
		[]string{"http://localhost:5173", "http://localhost:4173"},
		cfg.CORS.AllowedOrigins,
	)
	require.Equal(t, "console", cfg.Log.Format)
}

func TestLoadDotEnvDoesNotOverrideEnvironment(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	err := os.WriteFile(
		path,
		[]byte(
			"TEST_DOTENV_VALUE=from-file\n"+
				"TEST_DOTENV_QUOTED=\"hello world\"\n",
		),
		0o600,
	)
	require.NoError(t, err)

	t.Setenv("TEST_DOTENV_VALUE", "from-environment")

	err = loadDotEnv(path)
	require.NoError(t, err)

	require.Equal(t, "from-environment", os.Getenv("TEST_DOTENV_VALUE"))
	require.Equal(t, "hello world", os.Getenv("TEST_DOTENV_QUOTED"))
}

func TestLoadRejectsInvalidEnvironmentValue(t *testing.T) {
	t.Setenv("HTTP_PORT", "not-a-port")

	_, err := parse()

	require.Error(t, err)
	require.Contains(t, err.Error(), "HTTP_PORT")
}

func TestLoadRejectsWildcardCORSOrigin(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "*")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "true")

	_, err := parse()

	require.Error(t, err)
	require.Contains(t, err.Error(), "must not contain *")
}

func TestPostgresConnectionStringEscapesCredentials(t *testing.T) {
	t.Setenv("TEST_POSTGRES_PASSWORD", "p@ss:w?rd#1")

	cfg := PostgresConfig{
		Host:            "localhost",
		Port:            5432,
		Database:        "developers_cv",
		User:            "developers",
		Password:        os.Getenv("TEST_POSTGRES_PASSWORD"),
		SSLMode:         "require",
		MaxConns:        10,
		MinConns:        2,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}

	require.Equal(
		t,
		"postgres://developers:p%40ss%3Aw%3Frd%231@localhost:5432/developers_cv?sslmode=require",
		cfg.ConnectionString(),
	)
}

func TestLoadRejectsInvalidPostgresSSLMode(t *testing.T) {
	t.Setenv("POSTGRES_SSLMODE", "invalid")

	_, err := parse()

	require.Error(t, err)
	require.Contains(t, err.Error(), "POSTGRES_SSLMODE")
}

func TestLoadParsesAuthMode(t *testing.T) {
	t.Setenv("AUTH_MODE", "mock")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")

	cfg, err := parse()

	require.NoError(t, err)
	require.Equal(t, "mock", cfg.Auth.Mode)
}

func TestLoadRejectsUnsupportedAuthMode(t *testing.T) {
	t.Setenv("AUTH_MODE", "jwt")

	_, err := parse()

	require.Error(t, err)
	require.Contains(t, err.Error(), "AUTH_MODE")
}
