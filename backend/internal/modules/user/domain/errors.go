package domain

import "errors"

var (
	ErrInvalidUser    = errors.New("invalid user")
	ErrInvalidProfile = errors.New("invalid profile")
	ErrEntityNotFound = errors.New("user not found")
)
