package readers

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"developers-cv/backend/internal/platform/database"
)

func resolveDB(ctx context.Context, fallback *gorm.DB, operation string) (*gorm.DB, error) {
	if fallback == nil {
		return nil, fmt.Errorf("%s: reader dependencies are nil", operation)
	}
	return database.DBFromContext(ctx, fallback).WithContext(ctx), nil
}

func yearMonth(value time.Time) string {
	return value.Format("2006-01")
}

func optionalYearMonth(value *time.Time) *string {
	if value == nil {
		return nil
	}
	result := yearMonth(*value)
	return &result
}

func optionalYear(value *time.Time) *string {
	if value == nil {
		return nil
	}
	result := value.Format("2006")
	return &result
}
