package repository

import (
	"errors"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume_link/domain"
)

const postgresUniqueViolationCode = "23505"

func mapDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domain.ErrDuplicateResumeLink
	}
	var sqlStateError interface{ SQLState() string }
	if errors.As(err, &sqlStateError) && sqlStateError.SQLState() == postgresUniqueViolationCode {
		return domain.ErrDuplicateResumeLink
	}
	return err
}
