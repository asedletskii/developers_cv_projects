package config

import "time"

func parsePostgresConfig() (PostgresConfig, error) {
	port, err := envInt("POSTGRES_PORT", 5432)
	if err != nil {
		return PostgresConfig{}, err
	}
	statementTimeout, err := envDuration("POSTGRES_STATEMENT_TIMEOUT", 15*time.Second)
	if err != nil {
		return PostgresConfig{}, err
	}
	lockTimeout, err := envDuration("POSTGRES_LOCK_TIMEOUT", 3*time.Second)
	if err != nil {
		return PostgresConfig{}, err
	}
	idleTimeout, err := envDuration("POSTGRES_IDLE_IN_TRANSACTION_TIMEOUT", time.Minute)
	if err != nil {
		return PostgresConfig{}, err
	}
	slowQueryThreshold, err := envDuration("POSTGRES_SLOW_QUERY_THRESHOLD", 500*time.Millisecond)
	if err != nil {
		return PostgresConfig{}, err
	}
	maxConns, err := envInt("POSTGRES_MAX_CONNS", 10)
	if err != nil {
		return PostgresConfig{}, err
	}
	maxIdleConns, err := envInt("POSTGRES_MAX_IDLE_CONNS", 2)
	if err != nil {
		return PostgresConfig{}, err
	}
	maxLifetime, err := envDuration("POSTGRES_MAX_CONN_LIFETIME", time.Hour)
	if err != nil {
		return PostgresConfig{}, err
	}
	maxIdle, err := envDuration("POSTGRES_MAX_CONN_IDLE_TIME", 30*time.Minute)
	if err != nil {
		return PostgresConfig{}, err
	}

	return PostgresConfig{
		Host:                     envString("POSTGRES_HOST", "localhost"),
		Port:                     port,
		Database:                 envString("POSTGRES_DB", "developers_cv"),
		User:                     envString("POSTGRES_USER", "developers_cv"),
		Password:                 envString("POSTGRES_PASSWORD", ""),
		SSLMode:                  envString("POSTGRES_SSLMODE", "disable"),
		StatementTimeout:         statementTimeout,
		LockTimeout:              lockTimeout,
		IdleInTransactionTimeout: idleTimeout,
		SlowQueryThreshold:       slowQueryThreshold,
		MaxConns:                 maxConns,
		MaxIdleConns:             maxIdleConns,
		MaxConnLifetime:          maxLifetime,
		MaxConnIdleTime:          maxIdle,
	}, nil
}
