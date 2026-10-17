package ownership

import (
	"context"
	"fmt"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
)

func EnsureResumeOwned(ctx context.Context, reader ports.ResumeOwnershipReader, query ports.ResumeOwnershipQuery) error {
	if reader == nil {
		return fmt.Errorf("resume ownership reader is nil")
	}
	if query.ResumeID.IsZero() || query.UserID.IsZero() {
		return fmt.Errorf("invalid resume ownership query")
	}
	owned, err := reader.OwnsResume(ctx, query)
	if err != nil {
		return fmt.Errorf("check resume ownership: %w", err)
	}
	if !owned {
		return resumedomain.ErrEntityNotFound
	}
	return nil
}
