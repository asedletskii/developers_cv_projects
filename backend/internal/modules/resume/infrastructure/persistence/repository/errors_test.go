package repository

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
)

func TestMapDBError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "gorm duplicate", err: gorm.ErrDuplicatedKey, want: domain.ErrDuplicateEntity},
		{name: "postgres duplicate", err: &pgconn.PgError{Code: postgresUniqueViolationCode}, want: domain.ErrDuplicateEntity},
		{name: "wrapped postgres duplicate", err: errors.Join(errors.New("write failed"), &pgconn.PgError{Code: postgresUniqueViolationCode}), want: domain.ErrDuplicateEntity},
		{name: "other postgres error", err: &pgconn.PgError{Code: "23503"}, want: nil},
		{name: "other error", err: errors.New("database unavailable"), want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapDBError(tt.err)
			if tt.want == nil {
				require.Error(t, got)
				require.ErrorIs(t, got, tt.err)
				return
			}
			require.ErrorIs(t, got, tt.want)
		})
	}
}
