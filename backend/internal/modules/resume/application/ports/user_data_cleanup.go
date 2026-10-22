package ports

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type UserDataCleanupRepository interface {
	ListIDsForUpdateByUserID(context.Context, sharedvalueobjects.ID) ([]sharedvalueobjects.ID, error)
	DeleteByUserID(context.Context, sharedvalueobjects.ID) error
}
