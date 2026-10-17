package database

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type testSQLError struct{ state string }

func (e testSQLError) Error() string    { return "postgres error " + e.state }
func (e testSQLError) SQLState() string { return e.state }

func TestMapPostgresError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "nil", err: nil, want: nil},
		{name: "gorm duplicate", err: gorm.ErrDuplicatedKey, want: ErrUniqueViolation},
		{name: "unique", err: testSQLError{state: "23505"}, want: ErrUniqueViolation},
		{name: "foreign key", err: testSQLError{state: "23503"}, want: ErrForeignKeyViolation},
		{name: "check", err: testSQLError{state: "23514"}, want: ErrConstraintViolation},
		{name: "not null", err: testSQLError{state: "23502"}, want: ErrConstraintViolation},
		{name: "too long", err: testSQLError{state: "22001"}, want: ErrInvalidData},
		{name: "invalid text", err: testSQLError{state: "22P05"}, want: ErrInvalidData},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapPostgresError(tt.err)
			if tt.want == nil {
				require.NoError(t, got)
				return
			}
			require.ErrorIs(t, got, tt.want)
		})
	}
}

func TestMapPostgresErrorPreservesUnknownError(t *testing.T) {
	original := errors.New("connection refused")
	require.ErrorIs(t, MapPostgresError(original), original)
}
