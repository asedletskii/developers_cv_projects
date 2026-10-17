package repository

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/domain"
)

func TestMapProfileDBErrorMapsInvalidDataToProfileValidation(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "check violation", err: &pgconn.PgError{Code: "23514"}},
		{name: "too long", err: &pgconn.PgError{Code: "22001"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(t, mapProfileDBError(tt.err), domain.ErrInvalidProfile)
		})
	}
}
