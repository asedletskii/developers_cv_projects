package repository

import (
	"errors"

	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/platform/database"
)

func mapDBError(err error) error {
	mapped := database.MapPostgresError(err)
	if errors.Is(mapped, database.ErrUniqueViolation) {
		return domain.ErrDuplicateResumeLink
	}
	if errors.Is(mapped, database.ErrForeignKeyViolation) {
		return domain.ErrEntityNotFound
	}
	if errors.Is(mapped, database.ErrConstraintViolation) || errors.Is(mapped, database.ErrInvalidData) {
		return domain.ErrInvalidResumeLink
	}
	return mapped
}
