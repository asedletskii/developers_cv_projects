package repository

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
)

func TestMapDBError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "gorm duplicate", err: gorm.ErrDuplicatedKey, want: ErrDuplicateAnalyticsRecord},
		{name: "postgres duplicate", err: &pgconn.PgError{Code: "23505"}, want: ErrDuplicateAnalyticsRecord},
		{name: "foreign key violation", err: &pgconn.PgError{Code: "23503"}, want: analyticscontracts.ErrAnalyticsSubjectNotFound},
		{name: "check violation", err: &pgconn.PgError{Code: "23514"}, want: ErrInvalidAnalyticsData},
		{name: "too long", err: &pgconn.PgError{Code: "22001"}, want: ErrInvalidAnalyticsData},
		{name: "other error", err: errors.New("database unavailable")},
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
