package database

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrUniqueViolation     = errors.New("database unique constraint violation")
	ErrForeignKeyViolation = errors.New("database foreign key constraint violation")
	ErrConstraintViolation = errors.New("database constraint violation")
	ErrInvalidData         = errors.New("database rejected invalid data")
)

type sqlStateError interface {
	SQLState() string
}

type classifiedError struct {
	kind  error
	cause error
}

func (e classifiedError) Error() string        { return e.cause.Error() }
func (e classifiedError) Unwrap() error        { return e.cause }
func (e classifiedError) Is(target error) bool { return target == e.kind }

func classify(kind, cause error) error {
	return classifiedError{kind: kind, cause: cause}
}

func MapPostgresError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return classify(ErrUniqueViolation, err)
	}

	var state sqlStateError
	if !errors.As(err, &state) {
		return err
	}

	switch state.SQLState() {
	case "23505": // unique_violation
		return classify(ErrUniqueViolation, err)
	case "23503": // foreign_key_violation
		return classify(ErrForeignKeyViolation, err)
	case "23502", "23514": // not_null_violation, check_violation
		return classify(ErrConstraintViolation, err)
	case "22001", "22021", "22P05": // string truncation, encoding and invalid text representation
		return classify(ErrInvalidData, err)
	default:
		return err
	}
}
