package repository

import (
	"errors"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
)

const postgresUniqueViolationCode = "23505"

func mapDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domain.ErrDuplicateEntity
	}
	var sqlStateError interface{ SQLState() string }
	if errors.As(err, &sqlStateError) && sqlStateError.SQLState() == postgresUniqueViolationCode {
		return domain.ErrDuplicateEntity
	}
	return err
}
