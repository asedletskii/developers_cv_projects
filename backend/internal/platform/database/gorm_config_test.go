//go:build integration

package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"developers-cv/backend/internal/platform/config"
)

func TestConfigurePool(t *testing.T) {
	db, err := gorm.Open(postgres.Open("postgres://user:password@localhost:5432/developers_cv?sslmode=disable"), &gorm.Config{DryRun: true})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, Close(db)) })

	cfg := config.PostgresConfig{
		MaxConns:        12,
		MaxIdleConns:    4,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 15 * time.Minute,
	}
	require.NoError(t, ConfigurePool(db, cfg))

	sqlDB, err := db.DB()
	require.NoError(t, err)
	stats := sqlDB.Stats()
	require.Equal(t, 12, stats.MaxOpenConnections)
	require.Equal(t, 4, stats.Idle)
}
