package pgtest

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	testcontainers "github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	image        = "postgres:17-alpine"
	databaseName = "developers_cv"
	username     = "developers_cv"
	password     = "test-password"
)

var (
	baseDSN string
	nextDB  atomic.Uint64
)

func Main(m *testing.M) int {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	container, err := testpostgres.Run(
		ctx,
		image,
		testpostgres.WithDatabase(databaseName),
		testpostgres.WithUsername(username),
		testpostgres.WithPassword(password),
		testpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start integration PostgreSQL: %v\n", err)
		return 1
	}

	baseDSN, err = container.ConnectionString(ctx, "sslmode=disable")
	if err == nil {
		err = applyMigrations(ctx, baseDSN)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "prepare integration PostgreSQL: %v\n", err)
		if terminateErr := testcontainers.TerminateContainer(container); terminateErr != nil {
			fmt.Fprintf(os.Stderr, "terminate PostgreSQL after preparation failure: %v\n", terminateErr)
		}
		return 1
	}

	code := m.Run()
	if err := testcontainers.TerminateContainer(container); err != nil {
		fmt.Fprintf(os.Stderr, "terminate integration PostgreSQL: %v\n", err)
		if code == 0 {
			code = 1
		}
	}
	return code
}

func NewDatabase(t *testing.T) (context.Context, string) {
	t.Helper()
	if baseDSN == "" {
		t.Fatal("pgtest.Main(m) must be called from TestMain for integration tests")
	}

	ctx := context.Background()
	name := fmt.Sprintf("developers_cv_test_%d", nextDB.Add(1))
	adminDSN, err := withDatabase(baseDSN, "postgres")
	if err != nil {
		t.Fatalf("build PostgreSQL admin DSN: %v", err)
	}
	admin, err := pgx.Connect(ctx, adminDSN)
	if err != nil {
		t.Fatalf("connect to PostgreSQL admin database: %v", err)
	}
	_, err = admin.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", name, databaseName))
	closeErr := admin.Close(ctx)
	if err != nil {
		t.Fatalf("create isolated PostgreSQL database: %v", err)
	}
	if closeErr != nil {
		t.Fatalf("close PostgreSQL admin connection: %v", closeErr)
	}

	dsn, err := withDatabase(baseDSN, name)
	if err != nil {
		t.Fatalf("build isolated PostgreSQL DSN: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		conn, err := pgx.Connect(cleanupCtx, adminDSN)
		if err != nil {
			t.Errorf("connect to drop isolated PostgreSQL database %q: %v", name, err)
			return
		}
		_, dropErr := conn.Exec(cleanupCtx, fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", name))
		closeErr := conn.Close(cleanupCtx)
		if dropErr != nil {
			t.Errorf("drop isolated PostgreSQL database %q: %v", name, dropErr)
		}
		if closeErr != nil {
			t.Errorf("close PostgreSQL cleanup connection: %v", closeErr)
		}
	})
	return ctx, dsn
}

func withDatabase(dsn, name string) (string, error) {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}
	parsed.Path = "/" + name
	return parsed.String(), nil
}

func applyMigrations(ctx context.Context, dsn string) (retErr error) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("resolve pgtest source directory")
	}
	migrationsDir := filepath.Join(filepath.Dir(sourceFile), "../../../db/migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	if len(names) == 0 {
		return fmt.Errorf("no up migrations found in %s", migrationsDir)
	}

	root, err := os.OpenRoot(migrationsDir)
	if err != nil {
		return fmt.Errorf("open migrations directory: %w", err)
	}
	defer func() {
		if err := root.Close(); err != nil {
			retErr = errors.Join(retErr, fmt.Errorf("close migrations directory: %w", err))
		}
	}()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("connect to migration database: %w", err)
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := conn.Close(closeCtx); err != nil {
			retErr = errors.Join(retErr, fmt.Errorf("close migration database connection: %w", err))
		}
	}()

	for _, name := range names {
		content, err := root.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err := conn.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("apply migration %s: %w", name, err)
		}
	}
	return nil
}
