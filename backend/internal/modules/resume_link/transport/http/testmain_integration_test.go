//go:build integration

package http

import (
	"os"
	"testing"

	"developers-cv/backend/internal/testutil/pgtest"
)

func TestMain(m *testing.M) {
	os.Exit(pgtest.Main(m))
}
