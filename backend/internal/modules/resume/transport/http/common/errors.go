package common

import (
	"errors"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func IsValidationError(err error) bool {
	return errors.Is(err, domain.ErrInvalidResume) ||
		errors.Is(err, domain.ErrInvalidEntity) ||
		errors.Is(err, domain.ErrResumeInvariant) ||
		errors.Is(err, domain.ErrDuplicateEntity) ||
		errors.Is(err, valueobjects.ErrInvalidValue)
}
