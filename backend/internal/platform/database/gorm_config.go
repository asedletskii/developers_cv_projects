package database

import (
	"developers-cv/backend/internal/platform/config"

	"gorm.io/gorm"
)

func NewGORMFromConfig(cfg config.Config) (*gorm.DB, error) {
	return NewGORM(cfg.Postgres.ConnectionString())
}
