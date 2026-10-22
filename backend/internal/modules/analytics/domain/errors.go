package domain

import "errors"

var (
	ErrInvalidView             = errors.New("invalid analytics view")
	ErrInvalidViewState        = errors.New("invalid analytics view state")
	ErrInvalidAnalyticsSection = errors.New("invalid analytics section")
	ErrInvalidLayoutSnapshot   = errors.New("invalid analytics layout snapshot")
	ErrInvalidAnalyticsEvent   = errors.New("invalid analytics event")
	ErrInvalidAnalyticsClick   = errors.New("invalid analytics click")
	ErrAnalyticsViewExpired    = errors.New("analytics view expired")

	ErrAnalyticsLinkInactive = errors.New("resume link is inactive")
)
