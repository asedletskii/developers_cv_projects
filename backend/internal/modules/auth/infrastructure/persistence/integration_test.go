//go:build integration

package persistence_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"crypto/sha256"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/pgtest"
)

func newAuthDatabase(t *testing.T) (context.Context, string, *database.UnitOfWork, *repository.IdentityRepository, *repository.SessionRepository, *repository.AuthTransactionRepository, *repository.UserDataCleanupRepository) {
	t.Helper()
	ctx, dsn := pgtest.NewDatabase(t)
	db, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(db)) })
	return ctx, dsn, database.NewUnitOfWork(db), repository.NewIdentityRepository(db), repository.NewSessionRepository(db), repository.NewAuthTransactionRepository(db), repository.NewUserDataCleanupRepository(db)
}

func insertUser(t *testing.T, ctx context.Context, dsn string, id sharedvalueobjects.ID) {
	t.Helper()
	cfg, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, cfg)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, "INSERT INTO users.users (id) VALUES ($1)", id.UUID())
	require.NoError(t, err)
}

func testIdentity(t *testing.T, userID sharedvalueobjects.ID, subject string, at time.Time) *domain.Identity {
	t.Helper()
	providerSubject, err := valueobjects.NewProviderSubject(subject)
	require.NoError(t, err)
	identity, err := domain.NewIdentity(domain.IdentityParams{UserID: userID, Provider: domain.ProviderTelegram, ProviderSubject: providerSubject, CreatedAt: at, UpdatedAt: at, LastLoginAt: at})
	require.NoError(t, err)
	return identity
}

func testSession(t *testing.T, userID sharedvalueobjects.ID, at time.Time) *domain.Session {
	t.Helper()
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	session, err := domain.NewSession(domain.SessionParams{UserID: userID, TokenHash: hash, CreatedAt: at, LastSeenAt: at, ExpiresAt: at.Add(time.Hour)})
	require.NoError(t, err)
	return session
}

func testSessionWithTokenByte(t *testing.T, userID sharedvalueobjects.ID, at time.Time, fill byte) *domain.Session {
	t.Helper()
	raw := make([]byte, valueobjects.HashSize)
	for i := range raw {
		raw[i] = fill
	}
	hash, err := valueobjects.NewSessionTokenHash(raw)
	require.NoError(t, err)
	session, err := domain.NewSession(domain.SessionParams{UserID: userID, TokenHash: hash, CreatedAt: at, LastSeenAt: at, ExpiresAt: at.Add(time.Hour)})
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
	tx, err := domain.NewAuthTransaction(domain.AuthTransactionParams{Provider: domain.ProviderTelegram, StateHash: state, NonceHash: nonce, BrowserBindingHash: binding, CodeVerifier: "verifier", RedirectURI: "https://example.test/callback", CreatedAt: at, ExpiresAt: at.Add(5 * time.Minute)})
	require.NoError(t, err)
	return tx
}

func TestAuthRepositoriesPersistAndCleanupUserData(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, dsn, uow, identities, sessions, transactions, cleanup := newAuthDatabase(t)
	userID := sharedvalueobjects.NewID()
	otherUserID := sharedvalueobjects.NewID()
	insertUser(t, ctx, dsn, userID)
	insertUser(t, ctx, dsn, otherUserID)
	at := time.Now().UTC().Truncate(time.Microsecond)

	identity := testIdentity(t, userID, "telegram-subject", at)
	session := testSession(t, userID, at)
	otherIdentity := testIdentity(t, otherUserID, "other-telegram-subject", at)
	otherSession := testSessionWithTokenByte(t, otherUserID, at, 1)
	transaction := testTransaction(t, at)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		require.NoError(t, identities.Create(tx, identity))
		require.NoError(t, sessions.Create(tx, session))
		require.NoError(t, identities.Create(tx, otherIdentity))
		require.NoError(t, sessions.Create(tx, otherSession))
		return transactions.Create(tx, transaction)
	}))

	loadedIdentity, err := identities.FindByProviderSubject(ctx, domain.ProviderTelegram, identity.ProviderSubject())
	require.NoError(t, err)
	require.Equal(t, identity.ID(), loadedIdentity.ID())
	updatedAt := at.Add(2 * time.Minute)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return identities.UpdateLastLogin(tx, identity.ID(), updatedAt) }))
	loadedIdentity, err = identities.FindByProviderSubject(ctx, domain.ProviderTelegram, identity.ProviderSubject())
	require.NoError(t, err)
	require.True(t, updatedAt.Equal(loadedIdentity.LastLoginAt()))

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

	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return sessions.Revoke(tx, session.ID(), at.Add(time.Minute)) }))
	_, err = sessions.FindActiveByTokenHash(ctx, session.TokenHash())
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return cleanup.DeleteByUserID(tx, userID) }))
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return cleanup.DeleteByUserID(tx, userID) }), "cleanup must be idempotent")

	db, err := database.NewGORM(dsn)
	require.NoError(t, err)
	defer database.Close(db)
	for _, table := range []string{"auth.user_identities", "auth.auth_sessions"} {
		var count int64
		require.NoError(t, db.Table(table).Where("user_id = ?", userID.UUID()).Count(&count).Error)
		require.Zero(t, count, table)
		require.NoError(t, db.Table(table).Where("user_id = ?", otherUserID.UUID()).Count(&count).Error)
		require.Equal(t, int64(1), count, table)
	}
}

func TestAuthTransactionForUpdateBlocksConcurrentReader(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx, _, uow, _, _, transactions, _ := newAuthDatabase(t)
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
	ctx, _, uow, _, _, transactions, _ := newAuthDatabase(t)
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
