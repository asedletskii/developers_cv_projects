package database

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"developers-cv/backend/internal/platform/config"
)

func NewGORMFromConfig(cfg config.Config, log *zap.Logger) (*gorm.DB, error) {
	db, err := NewGORMWithLogger(cfg.Postgres.ConnectionString(), log, cfg.Postgres.SlowQueryThreshold)
	if err != nil {
		return nil, err
	}
	if err := ConfigurePool(db, cfg.Postgres); err != nil {
		if closeErr := Close(db); closeErr != nil {
			return nil, errors.Join(
				err,
				fmt.Errorf("close gorm after pool configuration failure: %w", closeErr),
			)
		}
		return nil, err
	}
	return db, nil
}

func ConfigurePool(db *gorm.DB, cfg config.PostgresConfig) error {
	if db == nil {
		return fmt.Errorf("configure gorm pool: database is nil")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxConnLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.MaxConnIdleTime)

	return nil
}
