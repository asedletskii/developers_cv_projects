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
	ID        sharedvalueobjects.ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(params UserParams) (*User, error) {
	if params.ID.IsZero() {
		return nil, fmt.Errorf("%w: id is empty", ErrInvalidUser)
	}

	now := time.Now().UTC()
	createdAt := params.CreatedAt.UTC()
	if createdAt.IsZero() {
		createdAt = now
	}
	updatedAt := params.UpdatedAt.UTC()
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	if updatedAt.Before(createdAt) {
		return nil, fmt.Errorf("%w: updated_at is before created_at", ErrInvalidUser)
	}

	return &User{
		id:        params.ID,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
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
