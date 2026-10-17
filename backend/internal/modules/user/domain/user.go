package domain

import (
	"fmt"
	"time"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type User struct {
	id        sharedvalueobjects.ID
	createdAt time.Time
	updatedAt time.Time
}

type UserParams struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserSnapshot struct {
	ID        sharedvalueobjects.ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(params UserParams) (*User, error) {
	id := sharedvalueobjects.NewID()

	createdAt := params.CreatedAt.UTC()
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	updatedAt := params.UpdatedAt.UTC()
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	if err := validateUser(id, createdAt, updatedAt); err != nil {
		return nil, err
	}

	return &User{
		id:        id,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func RestoreUser(snapshot UserSnapshot) (*User, error) {
	if err := validateUser(snapshot.ID, snapshot.CreatedAt, snapshot.UpdatedAt); err != nil {
		return nil, err
	}
	return &User{id: snapshot.ID, createdAt: snapshot.CreatedAt, updatedAt: snapshot.UpdatedAt}, nil
}

func validateUser(id sharedvalueobjects.ID, createdAt, updatedAt time.Time) error {
	if id.IsZero() {
		return fmt.Errorf("%w: id is empty", ErrInvalidUser)
	}
	if createdAt.IsZero() || updatedAt.IsZero() {
		return fmt.Errorf("%w: timestamps must be set", ErrInvalidUser)
	}
	if updatedAt.Before(createdAt) {
		return fmt.Errorf("%w: updated_at is before created_at", ErrInvalidUser)
	}
	return nil
}

func (u *User) ID() sharedvalueobjects.ID {
	return u.id
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
