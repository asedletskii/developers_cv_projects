package contracts

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeLinkCleanup interface {
	DeleteByResumeIDs(context.Context, []sharedvalueobjects.ID) error
}
