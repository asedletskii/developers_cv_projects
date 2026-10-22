package repository

import (
	"errors"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	"developers-cv/backend/internal/platform/database"
)

var (
	ErrDuplicateAnalyticsRecord = errors.New("duplicate analytics record")
	ErrInvalidAnalyticsData     = errors.New("invalid analytics data")
)

func mapDBError(err error) error {
	mapped := database.MapPostgresError(err)
	if errors.Is(mapped, database.ErrUniqueViolation) {
		return ErrDuplicateAnalyticsRecord
	}
	if errors.Is(mapped, database.ErrForeignKeyViolation) {
		return analyticscontracts.ErrAnalyticsSubjectNotFound
	}
	if errors.Is(mapped, database.ErrConstraintViolation) || errors.Is(mapped, database.ErrInvalidData) {
		return ErrInvalidAnalyticsData
	}
	return mapped
}
