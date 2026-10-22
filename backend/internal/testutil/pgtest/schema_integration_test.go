//go:build integration

package pgtest

import (
	"sort"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

var moduleTablesBySchema = map[string][]string{
	"users": {
		"users",
		"user_profiles",
		"user_profile_localizations",
	},
	"auth": {
		"user_identities",
		"auth_sessions",
		"auth_transactions",
	},
	"resume": {
		"resumes",
		"resume_overviews",
		"resume_overview_sections",
		"resume_profiles",
		"resume_profile_work_formats",
		"resume_contacts",
		"resume_contact_hidden",
		"resume_contacts_text",
		"resume_skills",
		"resume_skill_hard",
		"resume_skill_soft",
		"resume_spoken_languages",
		"resume_technologies",
		"resume_experiences",
		"resume_experience_achievements",
		"resume_experience_stack",
		"resume_educations",
		"resume_courses",
		"resume_activities",
		"resume_projects",
		"resume_project_screenshots",
		"resume_project_technologies",
		"resume_project_contributors",
		"resume_experience_tasks",
	},
	"resume_link": {
		"resume_links",
	},
	"analytics": {
		"analytics_views",
		"analytics_events",
		"analytics_clicks",
		"analytics_visitors",
	},
}

func TestTablesLiveInOwnedSchemas(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := NewDatabase(t)
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	rows, err := conn.Query(ctx, `
		SELECT table_schema, table_name
		FROM information_schema.tables
		WHERE table_type = 'BASE TABLE'
		  AND table_schema = ANY($1)
		ORDER BY table_schema, table_name
	`, moduleSchemas())
	require.NoError(t, err)
	defer rows.Close()

	actual := make(map[string][]string)
	for rows.Next() {
		var schema, table string
		require.NoError(t, rows.Scan(&schema, &table))
		actual[schema] = append(actual[schema], table)
	}
	require.NoError(t, rows.Err())

	expected := make(map[string][]string, len(moduleTablesBySchema))
	for schema, tables := range moduleTablesBySchema {
		expected[schema] = append([]string(nil), tables...)
		sort.Strings(expected[schema])
	}
	require.Equal(t, expected, actual)
}

func TestNoModuleTablesInPublic(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := NewDatabase(t)
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	rows, err := conn.Query(ctx, `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		  AND table_type = 'BASE TABLE'
		  AND table_name <> 'schema_migrations'
		ORDER BY table_name
	`)
	require.NoError(t, err)
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		require.NoError(t, rows.Scan(&table))
		tables = append(tables, table)
	}
	require.NoError(t, rows.Err())
	require.Empty(t, tables, "module tables must not live in public")
}

func TestNoCrossSchemaForeignKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := NewDatabase(t)
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	rows, err := conn.Query(ctx, `
		SELECT
			n.nspname AS source_schema,
			t.relname AS source_table,
			fn.nspname AS target_schema,
			ft.relname AS target_table,
			c.conname AS constraint_name
		FROM pg_constraint c
		JOIN pg_class t ON t.oid = c.conrelid
		JOIN pg_namespace n ON n.oid = t.relnamespace
		JOIN pg_class ft ON ft.oid = c.confrelid
		JOIN pg_namespace fn ON fn.oid = ft.relnamespace
		WHERE c.contype = 'f'
		  AND n.nspname <> fn.nspname
		ORDER BY n.nspname, t.relname, c.conname
	`)
	require.NoError(t, err)
	defer rows.Close()

	var actual []string
	for rows.Next() {
		var sourceSchema, sourceTable, targetSchema, targetTable, constraint string
		require.NoError(t, rows.Scan(&sourceSchema, &sourceTable, &targetSchema, &targetTable, &constraint))
		actual = append(actual, crossSchemaFKKey(sourceSchema, sourceTable, targetSchema, targetTable, constraint))
	}
	require.NoError(t, rows.Err())
	require.Empty(t, actual, "cross-schema foreign keys must not exist")
}

func TestModuleForeignKeysAndCleanupIndexesRemain(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := NewDatabase(t)
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	rows, err := conn.Query(ctx, `
		SELECT n.nspname || '.' || t.relname || ' -> ' || fn.nspname || '.' || ft.relname
		FROM pg_constraint c
		JOIN pg_class t ON t.oid = c.conrelid
		JOIN pg_namespace n ON n.oid = t.relnamespace
		JOIN pg_class ft ON ft.oid = c.confrelid
		JOIN pg_namespace fn ON fn.oid = ft.relnamespace
		WHERE c.contype = 'f'
		  AND n.nspname = fn.nspname
	`)
	require.NoError(t, err)
	defer rows.Close()
	var foreignKeys []string
	for rows.Next() {
		var foreignKey string
		require.NoError(t, rows.Scan(&foreignKey))
		foreignKeys = append(foreignKeys, foreignKey)
	}
	require.NoError(t, rows.Err())
	require.Contains(t, foreignKeys, "users.user_profiles -> users.users")
	require.Contains(t, foreignKeys, "resume.resume_overviews -> resume.resumes")
	require.Contains(t, foreignKeys, "analytics.analytics_events -> analytics.analytics_views")

	rows, err = conn.Query(ctx, `
		SELECT schemaname || '.' || indexname
		FROM pg_indexes
		WHERE indexname IN (
			'idx_resumes_user_id',
			'idx_resume_links_resume_id',
			'auth_sessions_user_id_idx'
		)
	`)
	require.NoError(t, err)
	defer rows.Close()
	var indexes []string
	for rows.Next() {
		var index string
		require.NoError(t, rows.Scan(&index))
		indexes = append(indexes, index)
	}
	require.NoError(t, rows.Err())
	require.ElementsMatch(t, []string{
		"resume.idx_resumes_user_id",
		"resume_link.idx_resume_links_resume_id",
		"auth.auth_sessions_user_id_idx",
	}, indexes)
}

func TestOrphanChecksOnEmptyDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := NewDatabase(t)
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	RequireNoOrphans(t, ctx, conn)
}

func TestAnalyticsViewsResumeSchema(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := NewDatabase(t)
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	var dataType, nullable string
	require.NoError(t, conn.QueryRow(ctx, `
		SELECT data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'analytics'
		  AND table_name = 'analytics_views'
		  AND column_name = 'resume_id'
	`).Scan(&dataType, &nullable))
	require.Equal(t, "uuid", dataType)
	require.Equal(t, "NO", nullable)

	rows, err := conn.Query(ctx, `
		SELECT indexname
		FROM pg_indexes
		WHERE schemaname = 'analytics'
		  AND tablename = 'analytics_views'
		  AND indexname IN (
		      'idx_analytics_views_resume_started',
		      'idx_analytics_views_resume_visitor_started'
		  )
		ORDER BY indexname
	`)
	require.NoError(t, err)
	defer rows.Close()
	var indexes []string
	for rows.Next() {
		var name string
		require.NoError(t, rows.Scan(&name))
		indexes = append(indexes, name)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []string{
		"idx_analytics_views_resume_started",
		"idx_analytics_views_resume_visitor_started",
	}, indexes)
}

func moduleSchemas() []string {
	schemas := make([]string, 0, len(moduleTablesBySchema))
	for schema := range moduleTablesBySchema {
		schemas = append(schemas, schema)
	}
	sort.Strings(schemas)
	return schemas
}

func crossSchemaFKKey(sourceSchema, sourceTable, targetSchema, targetTable, constraint string) string {
	return sourceSchema + "." + sourceTable + " -> " + targetSchema + "." + targetTable + " (" + constraint + ")"
}
