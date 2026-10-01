package queries

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrProfileNotFound = errors.New("profile not found")
)
