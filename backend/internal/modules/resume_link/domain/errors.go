package domain

import "errors"

var (
	ErrConcurrencyConflict     = errors.New("resume link concurrency conflict")
	ErrInvalidResumeLink       = errors.New("invalid resume link")
	ErrInvalidResumeLinkName   = errors.New("invalid resume link name")
	ErrInvalidView             = errors.New("invalid resume link view")
	ErrInvalidViewState        = errors.New("invalid resume link view state")
	ErrInvalidLayoutSnapshot   = errors.New("invalid analytics layout snapshot")
	ErrInvalidAnalyticsEvent   = errors.New("invalid analytics event")
	ErrInvalidAnalyticsClick   = errors.New("invalid analytics click")
	ErrAnalyticsViewExpired    = errors.New("analytics view expired")
	ErrEntityNotFound          = errors.New("resume link not found")
	ErrResumeLinkInactive      = errors.New("resume link is inactive")
	ErrDuplicateResumeLink     = errors.New("duplicate resume link")
	ErrResumeLinkLimitExceeded = errors.New("resume link limit exceeded")
)
