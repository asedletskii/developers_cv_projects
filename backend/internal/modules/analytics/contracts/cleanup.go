package contracts

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type AnalyticsLinkCleanup interface {
	DeleteByLinkIDs(context.Context, []sharedvalueobjects.ID) error
}
