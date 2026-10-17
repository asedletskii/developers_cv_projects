package queries

import (
	"context"
	"errors"
	"fmt"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ownership"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
)

func EnsureResumeOwned(ctx context.Context, reader ports.ResumeOwnershipReader, query resumequeries.GetResumeQuery) error {
	err := ownership.EnsureResumeOwned(ctx, reader, query)
	if errors.Is(err, resumedomain.ErrEntityNotFound) {
		return resumequeries.ErrResumeNotFound
	}
	if errors.Is(err, resumequeries.ErrResumeNotFound) {
		return resumequeries.ErrResumeNotFound
	}
	if err != nil {
		return fmt.Errorf("ensure resume ownership: %w", err)
	}
	return nil
}
