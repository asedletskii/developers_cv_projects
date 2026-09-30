//go:build integration

package commands

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/repository"
	usercommands "developers-cv/backend/internal/modules/user/application/commands"
	userrepo "developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type integrationTokenGenerator struct{ counter atomic.Uint64 }

func (g *integrationTokenGenerator) GenerateState() ([]byte, error) { return []byte("unused"), nil }
func (g *integrationTokenGenerator) GenerateNonce() ([]byte, error) { return []byte("unused"), nil }
func (g *integrationTokenGenerator) GeneratePKCEVerifier() ([]byte, error) {
	return []byte("verifier"), nil
}
func (g *integrationTokenGenerator) GenerateSessionToken() ([]byte, error) {
	n := g.counter.Add(1)
	return []byte(fmt.Sprintf("session-token-%d", n)), nil
}
func (g *integrationTokenGenerator) Hash(value []byte) []byte {
	h := sha256.Sum256(value)
	return h[:]
}

type integrationOIDCStub struct{ subject valueobjects.ProviderSubject }

func (o integrationOIDCStub) AuthorizationURL(params ports.TelegramAuthorizationParams) string {
	return ""
}
func (o integrationOIDCStub) Exchange(context.Context, string, string, string, valueobjects.NonceHash) (ports.VerifiedIdentity, error) {
	return ports.VerifiedIdentity{Provider: testIntegrationProvider(), ProviderSubject: o.subject}, nil
}

type blockingIdentityRepository struct {
	inner       ports.IdentityRepository
	barrier     chan struct{}
	entered     atomic.Int32
	barrierOnce sync.Once
}

func (r *blockingIdentityRepository) FindByProviderSubject(ctx context.Context, provider domain.Provider, subject valueobjects.ProviderSubject) (*domain.Identity, error) {
	if r.entered.Add(1) <= 2 {
		if r.entered.Load() == 2 {
			r.barrierOnce.Do(func() { close(r.barrier) })
		}
		<-r.barrier
	}
	return r.inner.FindByProviderSubject(ctx, provider, subject)
}
func (r *blockingIdentityRepository) Create(ctx context.Context, identity *domain.Identity) error {
	return r.inner.Create(ctx, identity)
}
func (r *blockingIdentityRepository) UpdateLastLogin(ctx context.Context, identityID sharedvalueobjects.ID, at time.Time) error {
	return r.inner.UpdateLastLogin(ctx, identityID, at)
}

func TestCompleteTelegramLoginConcurrentFirstLoginUsesOneIdentityAndLeavesNoOrphanUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, container.Terminate(ctx)) })
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyAuthIntegrationMigrations(t, ctx, dsn)

	db, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(db)) })

	uow := database.NewUnitOfWork(db)
	identities := repository.NewIdentityRepository(db)
	sessions := repository.NewSessionRepository(db)
	transactions := repository.NewAuthTransactionRepository(db)
	users := userrepo.NewUserRepository(db)
	provisioner := usercommands.NewCreateUserHandler(users)
	identity := func() valueobjects.ProviderSubject {
		subject, err := valueobjects.NewProviderSubject("same-telegram-subject")
		require.NoError(t, err)
		return subject
	}()
	blockedIdentities := &blockingIdentityRepository{inner: identities, barrier: make(chan struct{})}
	generator := &integrationTokenGenerator{}
	oidc := integrationOIDCStub{subject: identity}

	handler := NewCompleteTelegramLoginHandler(transactions, blockedIdentities, sessions, oidc, provisioner, generator, uow, time.Hour)
	now := time.Now().UTC().Truncate(time.Microsecond)
	var baselineUsersCount int64
	require.NoError(t, db.Table("users").Count(&baselineUsersCount).Error)
	createTransaction := func(state string) {
		hash := generator.Hash([]byte(state))
		stateHash, err := valueobjects.NewStateHash(hash)
		require.NoError(t, err)
		nonceHash, err := valueobjects.NewNonceHash(generator.Hash([]byte("nonce-" + state)))
		require.NoError(t, err)
		tx, err := domain.NewAuthTransaction(domain.AuthTransactionParams{
			ID: sharedvalueobjects.NewID(), Provider: testIntegrationProvider(), StateHash: stateHash, NonceHash: nonceHash,
			CodeVerifier: "verifier", RedirectURI: "https://example.test/callback", CreatedAt: now, ExpiresAt: now.Add(5 * time.Minute),
		})
		require.NoError(t, err)
		require.NoError(t, uow.Do(ctx, func(txCtx context.Context) error { return transactions.Create(txCtx, tx) }))
	}
	createTransaction("state-a")
	createTransaction("state-b")

	type result struct {
		session AuthenticatedSession
		err     error
	}
	results := make(chan result, 2)
	var wg sync.WaitGroup
	for _, state := range []string{"state-a", "state-b"} {
		wg.Add(1)
		go func(state string) {
			defer wg.Done()
			session, err := handler.Handle(ctx, "code-"+state, state)
			results <- result{session: session, err: err}
		}(state)
	}
	wg.Wait()
	close(results)

	var successes int
	for result := range results {
		if result.err == nil {
			successes++
			require.NotEmpty(t, result.session.RawSessionToken)
		} else {
			require.NoError(t, result.err)
		}
	}
	require.Equal(t, 2, successes)

	var usersCount, identitiesCount, sessionsCount int64
	require.NoError(t, db.Table("users").Count(&usersCount).Error)
	require.NoError(t, db.Table("user_identities").Count(&identitiesCount).Error)
	require.NoError(t, db.Table("auth_sessions").Count(&sessionsCount).Error)
	require.Equal(t, baselineUsersCount+1, usersCount)
	require.Equal(t, int64(1), identitiesCount)
	require.Equal(t, int64(2), sessionsCount)
}

func applyAuthIntegrationMigrations(t *testing.T, ctx context.Context, dsn string) {
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

func testIntegrationProvider() domain.Provider {
	provider, err := domain.NewProvider(string(domain.ProviderTelegram))
	if err != nil {
		panic(err)
	}
	return provider
}
