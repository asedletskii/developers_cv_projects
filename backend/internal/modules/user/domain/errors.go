package domain

import "errors"

var (
	ErrInvalidUser    = errors.New("invalid user")
	ErrEntityNotFound = errors.New("user not found")
)
