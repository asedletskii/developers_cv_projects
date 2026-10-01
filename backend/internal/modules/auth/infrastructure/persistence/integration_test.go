//go:build integration

package persistence_test

import (
	"context"
	"errors"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"crypto/sha256"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func newAuthDatabase(t *testing.T) (context.Context, string, *database.UnitOfWork, *repository.IdentityRepository, *repository.SessionRepository, *repository.AuthTransactionRepository) {
	t.Helper()
	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, container.Terminate(ctx)) })
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyMigrations(t, ctx, dsn)
	db, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(db)) })
	return ctx, dsn, database.NewUnitOfWork(db), repository.NewIdentityRepository(db), repository.NewSessionRepository(db), repository.NewAuthTransactionRepository(db)
}

func applyMigrations(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()
	cfg, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, cfg)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	dir := filepath.Join("..", "..", "..", "..", "..", "db", "migrations")
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	var names []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	for _, name := range names {
		sql, err := os.ReadFile(filepath.Join(dir, name))
		require.NoError(t, err)
		_, err = conn.Exec(ctx, string(sql))
		require.NoErrorf(t, err, "failed to apply %s", name)
	}
}

func insertUser(t *testing.T, ctx context.Context, dsn string, id identity.ID) {
	t.Helper()
	cfg, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, cfg)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, "INSERT INTO users (id) VALUES ($1)", id.UUID())
	require.NoError(t, err)
}

func testIdentity(t *testing.T, userID identity.ID, subject string, at time.Time) *domain.Identity {
	t.Helper()
	providerSubject, err := valueobjects.NewProviderSubject(subject)
	require.NoError(t, err)
	identity, err := domain.NewIdentity(domain.IdentityParams{ID: identity.NewID(), UserID: userID, Provider: domain.ProviderTelegram, ProviderSubject: providerSubject, CreatedAt: at, UpdatedAt: at, LastLoginAt: at})
	require.NoError(t, err)
	return identity
}

func testSession(t *testing.T, userID identity.ID, at time.Time) *domain.Session {
	t.Helper()
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	session, err := domain.NewSession(domain.SessionParams{ID: identity.NewID(), UserID: userID, TokenHash: hash, CreatedAt: at, LastSeenAt: at, ExpiresAt: at.Add(time.Hour)})
	require.NoError(t, err)
	return session
}

func testTransaction(t *testing.T, at time.Time) *domain.AuthTransaction {
	t.Helper()
	state, err := valueobjects.NewStateHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	nonceBytes := make([]byte, valueobjects.HashSize)
	nonceBytes[0] = 1
	nonce, err := valueobjects.NewNonceHash(nonceBytes)
	require.NoError(t, err)
	bindingBytes := sha256.Sum256([]byte("browser-binding"))
	binding, err := valueobjects.NewBrowserBindingHash(bindingBytes[:])
	require.NoError(t, err)
	tx, err := domain.NewAuthTransaction(domain.AuthTransactionParams{ID: identity.NewID(), Provider: domain.ProviderTelegram, StateHash: state, NonceHash: nonce, BrowserBindingHash: binding, CodeVerifier: "verifier", RedirectURI: "https://example.test/callback", CreatedAt: at, ExpiresAt: at.Add(5 * time.Minute)})
	require.NoError(t, err)
	return tx
}

func TestAuthRepositoriesPersistAndCascade(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, uow, identities, sessions, transactions := newAuthDatabase(t)
	userID := identity.NewID()
	insertUser(t, ctx, dsn, userID)
	at := time.Now().UTC().Truncate(time.Microsecond)

	identity := testIdentity(t, userID, "telegram-subject", at)
	session := testSession(t, userID, at)
	transaction := testTransaction(t, at)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		require.NoError(t, identities.Create(tx, identity))
		require.NoError(t, sessions.Create(tx, session))
		return transactions.Create(tx, transaction)
	}))

	loadedIdentity, err := identities.FindByProviderSubject(ctx, domain.ProviderTelegram, identity.ProviderSubject())
	require.NoError(t, err)
	require.Equal(t, identity.ID(), loadedIdentity.ID())
	updatedAt := at.Add(2 * time.Minute)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return identities.UpdateLastLogin(tx, identity.ID(), updatedAt) }))
	loadedIdentity, err = identities.FindByProviderSubject(ctx, domain.ProviderTelegram, identity.ProviderSubject())
	require.NoError(t, err)
	require.Equal(t, updatedAt, loadedIdentity.LastLoginAt())

	duplicate := testIdentity(t, userID, identity.ProviderSubject().String(), updatedAt)
	err = uow.Do(ctx, func(tx context.Context) error { return identities.Create(tx, duplicate) })
	require.ErrorIs(t, err, domain.ErrDuplicateEntity)
	loadedSession, err := sessions.FindActiveByTokenHash(ctx, session.TokenHash())
	require.NoError(t, err)
	require.Equal(t, session.ID(), loadedSession.ID())
	duplicateSession := testSession(t, userID, at)
	err = uow.Do(ctx, func(tx context.Context) error { return sessions.Create(tx, duplicateSession) })
	require.ErrorIs(t, err, domain.ErrDuplicateEntity)
	loadedTransaction, err := transactions.FindByStateHash(ctx, transaction.StateHash())
	require.NoError(t, err)
	require.Equal(t, transaction.ID(), loadedTransaction.ID())
	require.Equal(t, transaction.BrowserBindingHash().Bytes(), loadedTransaction.BrowserBindingHash().Bytes())
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return transactions.Consume(tx, transaction.ID(), at.Add(time.Minute)) }))
	consumed, err := transactions.FindByStateHash(ctx, transaction.StateHash())
	require.NoError(t, err)
	require.NotNil(t, consumed.ConsumedAt())

	duplicateTx := testTransaction(t, at)
	err = uow.Do(ctx, func(tx context.Context) error { return transactions.Create(tx, duplicateTx) })
	require.ErrorIs(t, err, domain.ErrDuplicateEntity)

	var count int64
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return sessions.Revoke(tx, session.ID(), at.Add(time.Minute)) }))
	_, err = sessions.FindActiveByTokenHash(ctx, session.TokenHash())
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	cfg, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, cfg)
	require.NoError(t, err)

	_, err = conn.Exec(ctx, "DELETE FROM users WHERE id = $1", userID.UUID())
	require.NoError(t, err)
	db, err := database.NewGORM(dsn)
	require.NoError(t, err)
	defer database.Close(db)
	require.NoError(t, db.Table("user_identities").Where("user_id = ?", userID.UUID()).Count(&count).Error)
	require.Zero(t, count)
	require.NoError(t, db.Table("auth_sessions").Where("user_id = ?", userID.UUID()).Count(&count).Error)
	require.Zero(t, count)
}

func TestAuthTransactionForUpdateBlocksConcurrentReader(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, _, uow, _, _, transactions := newAuthDatabase(t)
	at := time.Now().UTC().Truncate(time.Microsecond)
	transaction := testTransaction(t, at)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return transactions.Create(tx, transaction) }))

	lockCtx, cancelLock := context.WithCancel(ctx)
	defer cancelLock()
	acquired := make(chan struct{})
	released := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		result <- uow.Do(lockCtx, func(tx context.Context) error {
			_, err := transactions.FindByStateHashForUpdate(tx, transaction.StateHash())
			if err != nil {
				return err
			}
			close(acquired)
			<-released
			return nil
		})
	}()
	select {
	case <-acquired:
	case <-time.After(time.Second):
		require.Fail(t, "FOR UPDATE lock was not acquired")
	}

	blockedCtx, cancel := context.WithTimeout(ctx, 250*time.Millisecond)
	defer cancel()
	err := uow.Do(blockedCtx, func(tx context.Context) error {
		_, err := transactions.FindByStateHashForUpdate(tx, transaction.StateHash())
		return err
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "context deadline exceeded"))
	close(released)
	require.NoError(t, <-result)
}

func TestAuthRepositoriesDeleteExpiredBatch(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, _, uow, _, _, transactions := newAuthDatabase(t)
	now := time.Now().UTC().Truncate(time.Microsecond)
	tx := testTransaction(t, now.Add(-10*time.Minute))
	require.NoError(t, uow.Do(ctx, func(c context.Context) error { return transactions.Create(c, tx) }))
	require.NoError(t, uow.Do(ctx, func(c context.Context) error {
		n, err := transactions.DeleteExpiredBatch(c, now, 1)
		require.NoError(t, err)
		require.Equal(t, 1, n)
		return nil
	}))
}
