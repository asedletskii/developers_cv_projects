package commands

import (
	"context"
	"fmt"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type DeleteResumeLinksHandler struct {
	repository ports.ResumeLinkCleanupRepository
	analytics  analyticscontracts.AnalyticsLinkCleanup
}

func NewDeleteResumeLinksHandler(
	repository ports.ResumeLinkCleanupRepository,
	analytics analyticscontracts.AnalyticsLinkCleanup,
) *DeleteResumeLinksHandler {
	return &DeleteResumeLinksHandler{repository: repository, analytics: analytics}
}

func (h *DeleteResumeLinksHandler) DeleteByResumeIDs(
	ctx context.Context,
	resumeIDs []sharedvalueobjects.ID,
) error {
	if h == nil || h.repository == nil || h.analytics == nil {
		return fmt.Errorf("delete resume links: handler dependencies are nil")
	}
	if len(resumeIDs) == 0 {
		return nil
	}
	for _, resumeID := range resumeIDs {
		if resumeID.IsZero() {
			return fmt.Errorf("delete resume links: resume id is empty")
		}
	}
	linkIDs, err := h.repository.ListIDsForUpdateByResumeIDs(ctx, resumeIDs)
	if err != nil {
		return fmt.Errorf("delete resume links: %w", err)
	}
	if len(linkIDs) == 0 {
		return nil
	}
	if err := h.analytics.DeleteByLinkIDs(ctx, linkIDs); err != nil {
		return fmt.Errorf("delete resume links: clean analytics: %w", err)
	}
	if err := h.repository.DeleteByResumeIDs(ctx, resumeIDs); err != nil {
		return fmt.Errorf("delete resume links: %w", err)
	}
	return nil
}
