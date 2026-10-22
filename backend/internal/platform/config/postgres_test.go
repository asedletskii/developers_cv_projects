package config

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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
		MaxIdleConns:    2,
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
	setValidConfigEnvironment(t)
	t.Setenv("POSTGRES_SSLMODE", "invalid")

	_, err := parse()

	require.Error(t, err)
	require.Contains(t, err.Error(), "POSTGRES_SSLMODE")
}

func TestPostgresConnectionStringIncludesSessionTimeouts(t *testing.T) {
	cfg := PostgresConfig{
		Host: "localhost", Port: 5432, Database: "developers_cv", User: "developers", Password: "secret", SSLMode: "require",
		StatementTimeout: 15 * time.Second, LockTimeout: 3 * time.Second, IdleInTransactionTimeout: time.Minute,
	}

	parsed, err := url.Parse(cfg.ConnectionString())
	require.NoError(t, err)
	require.Equal(t, "-c statement_timeout=15000 -c lock_timeout=3000 -c idle_in_transaction_session_timeout=60000", parsed.Query().Get("options"))
}
