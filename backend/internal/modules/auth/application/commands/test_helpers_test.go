package commands

import (
	"context"
	"crypto/sha256"
	"testing"
	"time"

	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type tokenGeneratorStub struct {
	state, nonce, verifier, session []byte
}

func (g tokenGeneratorStub) GenerateState() ([]byte, error) { return g.state, nil }
func (g tokenGeneratorStub) GenerateBrowserBinding() ([]byte, error) {
	return []byte("browser-binding"), nil
}
func (g tokenGeneratorStub) GenerateNonce() ([]byte, error)        { return g.nonce, nil }
func (g tokenGeneratorStub) GeneratePKCEVerifier() ([]byte, error) { return g.verifier, nil }
func (g tokenGeneratorStub) GenerateSessionToken() ([]byte, error) { return g.session, nil }
func (g tokenGeneratorStub) Hash(value []byte) []byte {
	hash := sha256.Sum256(value)
	return hash[:]
}

type transactionRepositoryStub struct {
	transaction *domain.AuthTransaction
	createErr   error
	findErr     error
	returnNil   bool
	consumed    int
}

func (r *transactionRepositoryStub) Create(_ context.Context, transaction *domain.AuthTransaction) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.transaction = transaction
	return nil
}
func (r *transactionRepositoryStub) FindByStateHash(_ context.Context, _ valueobjects.StateHash) (*domain.AuthTransaction, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	if r.returnNil {
		return nil, nil
	}
	return r.transaction, nil
}
func (r *transactionRepositoryStub) FindByStateHashForUpdate(ctx context.Context, _ valueobjects.StateHash) (*domain.AuthTransaction, error) {
	return r.FindByStateHash(ctx, valueobjects.StateHash{})
}
func (r *transactionRepositoryStub) Consume(_ context.Context, _ sharedvalueobjects.ID, _ time.Time) error {
	r.consumed++
	return nil
}
func (r *transactionRepositoryStub) DeleteExpiredBatch(context.Context, time.Time, int) (int, error) {
	return 0, nil
}

type oidcStub struct {
	url          string
	identity     ports.VerifiedIdentity
	exchangeErr  error
	gotNonceHash valueobjects.NonceHash
}

func (o *oidcStub) AuthorizationURL(ctx context.Context, params ports.TelegramAuthorizationParams) (string, error) {
	return o.url, nil
}
func (o *oidcStub) Exchange(_ context.Context, _ string, _ string, _ string, nonceHash valueobjects.NonceHash) (ports.VerifiedIdentity, error) {
	o.gotNonceHash = nonceHash
	if o.exchangeErr != nil {
		return ports.VerifiedIdentity{}, o.exchangeErr
	}
	return o.identity, nil
}

type identityRepositoryStub struct {
	identity    *domain.Identity
	findErr     error
	returnNil   bool
	createErr   error
	updateCalls int
	createCalls int
}

func (r *identityRepositoryStub) FindByProviderSubject(context.Context, domain.Provider, valueobjects.ProviderSubject) (*domain.Identity, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	if r.returnNil {
		return nil, nil
	}
	if r.identity == nil {
		return nil, domain.ErrEntityNotFound
	}
	return r.identity, nil
}
func (r *identityRepositoryStub) Create(_ context.Context, identity *domain.Identity) error {
	r.createCalls++
	if r.createErr != nil {
		return r.createErr
	}
	r.identity = identity
	return nil
}
func (r *identityRepositoryStub) UpdateLastLogin(context.Context, sharedvalueobjects.ID, time.Time) error {
	r.updateCalls++
	return nil
}

type sessionRepositoryStub struct {
	session     *domain.Session
	returnNil   bool
	createErr   error
	touchCalls  int
	revokeCalls int
}

func (r *sessionRepositoryStub) Create(_ context.Context, session *domain.Session) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.session = session
	return nil
}
func (r *sessionRepositoryStub) FindActiveByTokenHash(context.Context, valueobjects.SessionTokenHash) (*domain.Session, error) {
	if r.returnNil {
		return nil, nil
	}
	if r.session == nil {
		return nil, domain.ErrEntityNotFound
	}
	return r.session, nil
}
func (r *sessionRepositoryStub) Revoke(context.Context, sharedvalueobjects.ID, time.Time) error {
	r.revokeCalls++
	return nil
}
func (r *sessionRepositoryStub) RevokeAllByUserID(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (r *sessionRepositoryStub) Touch(_ context.Context, _ sharedvalueobjects.ID, at time.Time) error {
	r.touchCalls++
	if r.session != nil {
		return r.session.Touch(at)
	}
	return nil
}
func (r *sessionRepositoryStub) DeleteExpiredBatch(context.Context, time.Time, int) (int, error) {
	return 0, nil
}

type provisionerStub struct {
	userID        sharedvalueobjects.ID
	calls         int
	ensureCalls   int
	ensureUserErr error
}

func (p *provisionerStub) CreateUser(context.Context) (sharedvalueobjects.ID, error) {
	p.calls++
	return p.userID, nil
}

func (p *provisionerStub) EnsureUserExists(context.Context, sharedvalueobjects.ID) error {
	p.ensureCalls++
	return p.ensureUserErr
}

type unitOfWorkStub struct {
	calls int
	err   error
}

func (u *unitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	u.calls++
	if u.err != nil {
		return u.err
	}
	return fn(ctx)
}

type cleanupTransactionRepositoryStub struct {
	deleted int
	err     error
}

func (r *cleanupTransactionRepositoryStub) Create(context.Context, *domain.AuthTransaction) error {
	return nil
}
func (r *cleanupTransactionRepositoryStub) FindByStateHash(context.Context, valueobjects.StateHash) (*domain.AuthTransaction, error) {
	return nil, domain.ErrEntityNotFound
}
func (r *cleanupTransactionRepositoryStub) FindByStateHashForUpdate(context.Context, valueobjects.StateHash) (*domain.AuthTransaction, error) {
	return nil, domain.ErrEntityNotFound
}
func (r *cleanupTransactionRepositoryStub) Consume(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (r *cleanupTransactionRepositoryStub) DeleteExpiredBatch(context.Context, time.Time, int) (int, error) {
	return r.deleted, r.err
}

type cleanupSessionRepositoryStub struct {
	deleted int
	err     error
}

func (r *cleanupSessionRepositoryStub) Create(context.Context, *domain.Session) error { return nil }
func (r *cleanupSessionRepositoryStub) FindActiveByTokenHash(context.Context, valueobjects.SessionTokenHash) (*domain.Session, error) {
	return nil, domain.ErrEntityNotFound
}
func (r *cleanupSessionRepositoryStub) Revoke(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (r *cleanupSessionRepositoryStub) RevokeAllByUserID(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (r *cleanupSessionRepositoryStub) Touch(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (r *cleanupSessionRepositoryStub) DeleteExpiredBatch(context.Context, time.Time, int) (int, error) {
	return r.deleted, r.err
}

type raceIdentityRepositoryStub struct {
	identity               *domain.Identity
	createCalls, findCalls int
}

func (r *raceIdentityRepositoryStub) FindByProviderSubject(context.Context, domain.Provider, valueobjects.ProviderSubject) (*domain.Identity, error) {
	r.findCalls++
	if r.findCalls == 1 {
		return nil, domain.ErrEntityNotFound
	}
	return r.identity, nil
}
func (r *raceIdentityRepositoryStub) Create(context.Context, *domain.Identity) error {
	r.createCalls++
	return domain.ErrDuplicateEntity
}
func (r *raceIdentityRepositoryStub) UpdateLastLogin(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}

func testProvider(t *testing.T) domain.Provider {
	t.Helper()
	provider, err := domain.NewProvider(string(domain.ProviderTelegram))
	if err != nil {
		t.Fatal(err)
	}
	return provider
}

func testSubject(t *testing.T, value string) valueobjects.ProviderSubject {
	t.Helper()
	subject, err := valueobjects.NewProviderSubject(value)
	if err != nil {
		t.Fatal(err)
	}
	return subject
}

func testTransaction(t *testing.T, at time.Time) *domain.AuthTransaction {
	t.Helper()
	state := sha256.Sum256([]byte("state"))
	nonce := sha256.Sum256([]byte("nonce"))
	stateHash, err := valueobjects.NewStateHash(state[:])
	if err != nil {
		t.Fatal(err)
	}
	nonceHash, err := valueobjects.NewNonceHash(nonce[:])
	if err != nil {
		t.Fatal(err)
	}
	binding := sha256.Sum256([]byte("browser-binding"))
	bindingHash, err := valueobjects.NewBrowserBindingHash(binding[:])
	if err != nil {
		t.Fatal(err)
	}
	transaction, err := domain.NewAuthTransaction(domain.AuthTransactionParams{
		Provider: testProvider(t), StateHash: stateHash, NonceHash: nonceHash, BrowserBindingHash: bindingHash,
		CodeVerifier: "verifier", RedirectURI: "https://example/callback", CreatedAt: at, ExpiresAt: at.Add(5 * time.Minute),
	})
	if err != nil {
		t.Fatal(err)
	}
	return transaction
}

func testSession(t *testing.T, at time.Time, lastSeen time.Time) *domain.Session {
	t.Helper()
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	if err != nil {
		t.Fatal(err)
	}
	session, err := domain.NewSession(domain.SessionParams{
		UserID: sharedvalueobjects.NewID(), TokenHash: hash, CreatedAt: at,
		LastSeenAt: lastSeen, ExpiresAt: at.Add(time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}
	return session
}
