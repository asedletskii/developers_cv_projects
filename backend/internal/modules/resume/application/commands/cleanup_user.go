package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	resumelinkcontracts "developers-cv/backend/internal/modules/resume_link/contracts"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type DeleteUserResumesHandler struct {
	repository ports.UserDataCleanupRepository
	links      resumelinkcontracts.ResumeLinkCleanup
}

func NewDeleteUserResumesHandler(
	repository ports.UserDataCleanupRepository,
	links resumelinkcontracts.ResumeLinkCleanup,
) *DeleteUserResumesHandler {
	return &DeleteUserResumesHandler{repository: repository, links: links}
}

func (h *DeleteUserResumesHandler) DeleteByUserID(ctx context.Context, userID sharedvalueobjects.ID) error {
	if h == nil || h.repository == nil || h.links == nil {
		return fmt.Errorf("delete user resumes: handler dependencies are nil")
	}
	if userID.IsZero() {
		return fmt.Errorf("delete user resumes: user id is empty")
	}
	resumeIDs, err := h.repository.ListIDsForUpdateByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete user resumes: %w", err)
	}
	if len(resumeIDs) > 0 {
		if err := h.links.DeleteByResumeIDs(ctx, resumeIDs); err != nil {
			return fmt.Errorf("delete user resumes: clean links: %w", err)
		}
	}
	if err := h.repository.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("delete user resumes: %w", err)
	}
	return nil
}
