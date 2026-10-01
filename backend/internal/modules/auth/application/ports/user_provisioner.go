package ports

import (
	"context"

	valueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type UserProvisioner interface {
	CreateUser(ctx context.Context) (valueobjects.ID, error)
}
