package commands

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func validateDeleteCommandIDs(resource string, resumeID, entityID valueobjects.ID, userID useridentity.ID) error {
	if resumeID.IsZero() {
		return fmt.Errorf("%s resume id is required: %w", resource, domain.ErrInvalidEntity)
	}
	if entityID.IsZero() {
		return fmt.Errorf("%s id is required: %w", resource, domain.ErrInvalidEntity)
	}
	if userID.IsZero() {
		return fmt.Errorf("%s user id is required: %w", resource, domain.ErrInvalidEntity)
	}
	return nil
}
