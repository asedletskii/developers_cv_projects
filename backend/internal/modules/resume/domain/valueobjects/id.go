package valueobjects

import (
	"github.com/google/uuid"
)

type ID struct {
	value uuid.UUID
}

func NewID() ID {
	return ID{value: uuid.New()}
}

func IDFromUUID(value uuid.UUID) (ID, error) {
	if value == uuid.Nil {
		return ID{}, ErrInvalidValue
	}
	return ID{value: value}, nil
}

func ParseID(value string) (ID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return ID{}, ErrInvalidValue
	}
	return IDFromUUID(parsed)
}

func (id ID) UUID() uuid.UUID {
	return id.value
}

func (id ID) String() string {
	return id.value.String()
}

func (id ID) IsZero() bool {
	return id.value == uuid.Nil
}
