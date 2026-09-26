package postgres

import (
	"context"

	"gorm.io/gorm"
)

type HealthChecker struct{ db *gorm.DB }

func NewHealthChecker(db *gorm.DB) *HealthChecker { return &HealthChecker{db: db} }

func (c *HealthChecker) Name() string { return "postgres" }

func (c *HealthChecker) Check(ctx context.Context) error {
	if c == nil || c.db == nil {
		return gorm.ErrInvalidDB
	}
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
