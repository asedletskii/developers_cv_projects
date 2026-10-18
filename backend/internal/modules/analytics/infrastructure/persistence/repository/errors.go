package repository

import (
	"errors"

	resumelinkdomain "developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/platform/database"
)

func mapDBError(err error) error {
	mapped := database.MapPostgresError(err)
	if errors.Is(mapped, database.ErrUniqueViolation) {
		return resumelinkdomain.ErrDuplicateResumeLink
	}
	if errors.Is(mapped, database.ErrForeignKeyViolation) {
		return resumelinkdomain.ErrEntityNotFound
	}
	if errors.Is(mapped, database.ErrConstraintViolation) || errors.Is(mapped, database.ErrInvalidData) {
		return resumelinkdomain.ErrInvalidResumeLink
	}
	return mapped
}
