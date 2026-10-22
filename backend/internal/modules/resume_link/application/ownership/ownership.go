package ownership

import (
	"context"
	"fmt"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	application "developers-cv/backend/internal/modules/resume_link/application"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func EnsureResumeOwnedForUpdate(
	ctx context.Context,
	guard resumecontracts.OwnedResumeLocker,
	resumeID, userID sharedvalueobjects.ID,
) error {
	if guard == nil {
		return fmt.Errorf("resume ownership guard is nil")
	}
	if resumeID.IsZero() || userID.IsZero() {
		return fmt.Errorf("invalid resume ownership query")
	}
	owned, err := guard.LockOwnedResumeForUpdate(ctx, resumecontracts.OwnedResumeLock{
		ResumeID: resumeID,
		UserID:   userID,
	})
	if err != nil {
		return fmt.Errorf("lock owned resume: %w", application.NormalizeProviderError(err))
	}
	if !owned {
		return application.ErrResumeNotFound
	}
	return nil
}

func EnsureResumeOwned(ctx context.Context, reader resumecontracts.ResumeOwnershipChecker, query resumecontracts.ResumeOwnershipQuery) error {
	if reader == nil {
		return fmt.Errorf("resume ownership reader is nil")
	}
	if query.ResumeID.IsZero() || query.UserID.IsZero() {
		return fmt.Errorf("invalid resume ownership query")
	}
	owned, err := reader.OwnsResume(ctx, query)
	if err != nil {
		return fmt.Errorf("check resume ownership: %w", application.NormalizeProviderError(err))
	}
	if !owned {
		return application.ErrResumeNotFound
	}
	return nil
}
