package ports

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type UserProvisioner interface {
	CreateUser(ctx context.Context) (sharedvalueobjects.ID, error)
}
