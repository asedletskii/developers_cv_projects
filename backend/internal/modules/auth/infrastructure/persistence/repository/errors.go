package repository

import (
	"errors"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/platform/database"
)

func mapDBError(err error) error {
	mapped := database.MapPostgresError(err)
	if errors.Is(mapped, database.ErrUniqueViolation) {
		return domain.ErrDuplicateEntity
	}
	return mapped
}
