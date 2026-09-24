package domain

import "errors"

var (
	ErrInvalidResume   = errors.New("invalid resume")
	ErrInvalidEntity   = errors.New("invalid entity")
	ErrResumeInvariant = errors.New("resume invariant violation")
	ErrEntityNotFound  = errors.New("entity not found")
	ErrDuplicateEntity = errors.New("duplicate entity")
)
