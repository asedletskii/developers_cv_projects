package errors

import (
	"errors"
	"net/http"
)

type Code string

const (
	CodeBadRequest      Code = "BAD_REQUEST"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeForbidden       Code = "FORBIDDEN"
	CodeNotFound        Code = "NOT_FOUND"
	CodeConflict        Code = "CONFLICT"
	CodeTooMany         Code = "TOO_MANY_REQUESTS"
	CodeInternal        Code = "INTERNAL_ERROR"
	CodeUnavailable     Code = "SERVICE_UNAVAILABLE"
	CodeRequestTooLarge Code = "REQUEST_TOO_LARGE"
)

type AppError struct {
	Code       Code
	Message    string
	Status     int
	Cause      error
	PublicSafe bool
}

func (e *AppError) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + ": " + e.Cause.Error()
}

func (e *AppError) Unwrap() error { return e.Cause }

func New(code Code, status int, message string) *AppError {
	return &AppError{Code: code, Status: status, Message: message, PublicSafe: true}
}

func Internal(cause error) *AppError {
	return &AppError{
		Code:    CodeInternal,
		Status:  http.StatusInternalServerError,
		Message: "Internal server error",
		Cause:   cause,
	}
}

func From(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return Internal(err)
}
