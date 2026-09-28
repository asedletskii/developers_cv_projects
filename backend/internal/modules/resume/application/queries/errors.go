package queries

import "errors"

var ErrResumeNotFound = errors.New("resume not found")

var (
	ErrContactNotFound  = errors.New("contact not found")
	ErrContactNotHidden = errors.New("contact is not hidden")
)
