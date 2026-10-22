//go:build integration

package pgtest

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestMigrationsAndSeedsUseSchemas(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := NewDatabase(t)
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	content, err := os.ReadFile(developmentSeedsPath(t))
	require.NoError(t, err)
	seeds := string(content)
	require.Contains(t, seeds, "INSERT INTO users.users")
	require.Contains(t, seeds, "INSERT INTO users.user_profiles")
	require.Contains(t, seeds, "INSERT INTO users.user_profile_localizations")
	for _, unqualified := range []string{
		"INSERT INTO users ",
		"INSERT INTO user_profiles ",
		"INSERT INTO user_profile_localizations ",
	} {
		require.False(t, strings.Contains(seeds, unqualified), "seed must qualify %q", unqualified)
	}

	_, err = conn.Exec(ctx, seeds)
	require.NoError(t, err, "apply development seeds")

	var users, profiles, localizations int
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM users.users").Scan(&users))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM users.user_profiles").Scan(&profiles))
	require.NoError(t, conn.QueryRow(ctx, "SELECT count(*) FROM users.user_profile_localizations").Scan(&localizations))
	require.Equal(t, 2, users)
	require.Equal(t, 2, profiles)
	require.Equal(t, 2, localizations)
}

func developmentSeedsPath(t *testing.T) string {
	t.Helper()
	_, sourceFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "resolve integration test source directory")
	return filepath.Join(filepath.Dir(sourceFile), "../../../db/seeds/development.sql")
}
