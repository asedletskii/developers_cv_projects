package application

import (
	"errors"
	"fmt"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume_link/domain"
)

var (
	ErrResumeNotFound     = errors.New("resume not found")
	ErrContactNotFound    = errors.New("contact not found")
	ErrContactNotHidden   = errors.New("contact is not hidden")
	ErrResumeLinkNotFound = errors.New("resume link not found")
)

func NormalizeProviderError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, resumecontracts.ErrResumeNotFound):
		return fmt.Errorf("%w: %w", ErrResumeNotFound, err)
	case errors.Is(err, resumecontracts.ErrContactNotFound):
		return fmt.Errorf("%w: %w", ErrContactNotFound, err)
	case errors.Is(err, resumecontracts.ErrContactNotHidden):
		return fmt.Errorf("%w: %w", ErrContactNotHidden, err)
	default:
		return err
	}
}

func NormalizeLinkNotFoundError(err error) error {
	if errors.Is(err, domain.ErrEntityNotFound) {
		return fmt.Errorf("%w: %w", ErrResumeLinkNotFound, err)
	}
	return err
}
