package contracts

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeOwnershipQuery struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
}

type ResumeOwnershipChecker interface {
	OwnsResume(context.Context, ResumeOwnershipQuery) (bool, error)
}

type OwnedResumeLock struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
}

type OwnedResumeLocker interface {
	LockOwnedResumeForUpdate(context.Context, OwnedResumeLock) (bool, error)
}
