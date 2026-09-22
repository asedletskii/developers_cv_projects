package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromPreservesAppError(t *testing.T) {
	expected := New(CodeNotFound, http.StatusNotFound, "Resource not found")

	actual := From(expected)

	require.Same(t, expected, actual)
}

func TestFromConvertsUnknownError(t *testing.T) {
	actual := From(errors.New("database password leaked"))

	require.Equal(t, CodeInternal, actual.Code)
	require.Equal(t, http.StatusInternalServerError, actual.Status)
	require.Equal(t, "Internal server error", actual.Message)
	require.NotContains(t, actual.Message, "database")
}
