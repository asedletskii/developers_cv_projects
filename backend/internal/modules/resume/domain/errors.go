package domain

import "errors"

var (
	ErrInvalidPerson   = errors.New("invalid person")
	ErrInvalidEntity   = errors.New("invalid entity")
	ErrPersonInvariant = errors.New("person invariant violation")
	ErrEntityNotFound  = errors.New("entity not found")
	ErrDuplicateEntity = errors.New("duplicate entity")
)
