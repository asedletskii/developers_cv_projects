package ports

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type UserDataCleanupRepository interface {
	DeleteByUserID(context.Context, sharedvalueobjects.ID) error
}
