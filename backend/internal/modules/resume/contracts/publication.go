package contracts

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumePublicationQuery struct {
	ResumeID sharedvalueobjects.ID
}

type ResumePublicationState struct {
	Exists    bool
	Published bool
	Status    string
}

type ResumePublicationReader interface {
	GetPublicationState(context.Context, ResumePublicationQuery) (ResumePublicationState, error)
}

type ResumePublicationLocker interface {
	LockPublicationState(context.Context, ResumePublicationQuery) (ResumePublicationState, error)
}
