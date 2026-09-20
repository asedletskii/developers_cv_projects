package repository

import (
	"errors"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/platform/database"
)

func mapDBError(err error) error {
	mapped := database.MapPostgresError(err)
	if errors.Is(mapped, database.ErrUniqueViolation) {
		return domain.ErrDuplicateEntity
	}
	if errors.Is(mapped, database.ErrForeignKeyViolation) {
		return domain.ErrEntityNotFound
	}
	if errors.Is(mapped, database.ErrConstraintViolation) || errors.Is(mapped, database.ErrInvalidData) {
		return domain.ErrInvalidEntity
	}
	return mapped
}
