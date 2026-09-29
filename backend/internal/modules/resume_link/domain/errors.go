package domain

import "errors"

var (
	ErrConcurrencyConflict     = errors.New("resume link concurrency conflict")
	ErrInvalidResumeLink       = errors.New("invalid resume link")
	ErrInvalidResumeLinkName   = errors.New("invalid resume link name")
	ErrInvalidResumeLinkEvent  = errors.New("invalid resume link event")
	ErrEntityNotFound          = errors.New("resume link not found")
	ErrDuplicateResumeLink     = errors.New("duplicate resume link")
	ErrResumeLinkLimitExceeded = errors.New("resume link limit exceeded")
)
