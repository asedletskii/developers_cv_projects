package valueobjects

import (
	"errors"

	"github.com/google/uuid"

	shared "developers-cv/backend/internal/shared/domain/valueobjects"
)

var ErrInvalidID = errors.New("invalid user id")

type ID = shared.ID

func NewID() ID {
	return shared.NewID()
}

func IDFromUUID(value uuid.UUID) (ID, error) {
	id, err := shared.IDFromUUID(value)
	if err != nil {
		return ID{}, ErrInvalidID
	}
	return id, nil
}

func ParseID(value string) (ID, error) {
	id, err := shared.ParseID(value)
	if err != nil {
		return ID{}, ErrInvalidID
	}
	return id, nil
}
