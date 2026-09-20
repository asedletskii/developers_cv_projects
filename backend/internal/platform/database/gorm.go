package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGORM(dsn string) (*gorm.DB, error) {
	return NewGORMWithLogger(dsn, zap.NewNop(), 500*time.Millisecond)
}

func NewGORMWithLogger(dsn string, log *zap.Logger, slowQueryThreshold time.Duration) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newGORMLogger(log, slowQueryThreshold)})
	if err != nil {
		return nil, fmt.Errorf("open gorm: %w", err)
	}
	return db, nil
}
