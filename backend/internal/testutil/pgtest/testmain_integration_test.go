//go:build integration

package pgtest

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(Main(m))
}
