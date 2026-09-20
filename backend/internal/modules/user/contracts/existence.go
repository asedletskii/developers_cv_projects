package contracts

import (
	"context"
	"errors"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

var ErrUserNotFound = errors.New("user not found")

type UserExistenceChecker interface {
	EnsureUserExists(context.Context, sharedvalueobjects.ID) error
}
