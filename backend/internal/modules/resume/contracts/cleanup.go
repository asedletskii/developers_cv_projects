package contracts

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type UserDataCleanup interface {
	DeleteByUserID(context.Context, sharedvalueobjects.ID) error
}
