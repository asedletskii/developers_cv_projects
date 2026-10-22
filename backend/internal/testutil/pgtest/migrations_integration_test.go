//go:build integration

package pgtest

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestAllMigrationsCanBeRolledBack(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := NewDatabase(t)
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	migrationsDir := migrationDirectory(t)
	root, err := os.OpenRoot(migrationsDir)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, root.Close()) })

	entries, err := os.ReadDir(migrationsDir)
	require.NoError(t, err)
	var names []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".down.sql") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	require.NotEmpty(t, names, "expected at least one down migration")

	for _, name := range names {
		content, err := root.ReadFile(name)
		require.NoError(t, err, "read migration %s", name)
		_, err = conn.Exec(ctx, string(content))
		require.NoError(t, err, "apply rollback migration %s", name)
		if strings.HasPrefix(name, "000003_") {
			for schema, tables := range moduleTablesBySchema {
				for _, table := range tables {
					var actual *string
					require.NoError(t, conn.QueryRow(ctx, "SELECT to_regclass($1)::text", "public."+table).Scan(&actual))
					require.NotNil(t, actual, "rollback of migration 3 should restore public.%s from %s", table, schema)
				}
			}
		}
		if strings.HasPrefix(name, "000002_") {
			for _, table := range []string{
				"resume_link_views", "resume_link_events", "resume_link_clicks", "resume_link_visitors",
			} {
				var actual *string
				require.NoError(t, conn.QueryRow(ctx, "SELECT to_regclass($1)::text", "public."+table).Scan(&actual))
				require.NotNil(t, actual, "rollback of migration 2 should restore %s", table)
			}
		}
	}

	var usersTable *string
	require.NoError(t, conn.QueryRow(ctx, "SELECT to_regclass('public.users')::text").Scan(&usersTable))
	require.Nil(t, usersTable, "rollback should remove the users table")
	var resumesTable *string
	require.NoError(t, conn.QueryRow(ctx, "SELECT to_regclass('public.resumes')::text").Scan(&resumesTable))
	require.Nil(t, resumesTable, "rollback should remove the resumes table")
}

func TestAnalyticsTablesUseAnalyticsOwnedNames(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := NewDatabase(t)
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	for _, table := range []string{
		"analytics_views", "analytics_events", "analytics_clicks", "analytics_visitors",
	} {
		var actual *string
		require.NoError(t, conn.QueryRow(ctx, "SELECT to_regclass($1)::text", "analytics."+table).Scan(&actual))
		require.NotNil(t, actual, "expected table analytics.%s", table)
		require.Equal(t, "analytics."+table, *actual)
	}
}

func migrationDirectory(t *testing.T) string {
	t.Helper()
	_, sourceFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "resolve integration test source directory")
	return filepath.Join(filepath.Dir(sourceFile), "../../../db/migrations")
}

func TestMigrationDirectoryContainsBothDirections(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	migrationsDir := migrationDirectory(t)
	entries, err := os.ReadDir(migrationsDir)
	require.NoError(t, err)
	up := make(map[string]struct{})
	down := make(map[string]struct{})
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			continue
		}
		switch {
		case strings.HasSuffix(name, ".up.sql"):
			up[strings.TrimSuffix(name, ".up.sql")] = struct{}{}
		case strings.HasSuffix(name, ".down.sql"):
			down[strings.TrimSuffix(name, ".down.sql")] = struct{}{}
		}
	}
	require.NotEmpty(t, up, "expected up migrations")
	require.Equal(t, up, down, "every up migration must have a matching down migration")

	for name := range up {
		for _, suffix := range []string{".up.sql", ".down.sql"} {
			content, err := os.ReadFile(filepath.Join(migrationsDir, name+suffix))
			require.NoError(t, err, "read migration %s%s", name, suffix)
			require.NotContains(t, string(content), "schema_migrations", "migrations must not move or modify the migration ledger")
		}
	}
}
