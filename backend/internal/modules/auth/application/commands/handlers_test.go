package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"crypto/sha256"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/application"
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
	userID sharedvalueobjects.ID
	calls  int
}

func (p *provisionerStub) CreateUser(context.Context) (sharedvalueobjects.ID, error) {
	p.calls++
	return p.userID, nil
}

type uowStub struct {
	calls int
	err   error
}

func (u *uowStub) Do(ctx context.Context, fn func(context.Context) error) error {
	u.calls++
	if u.err != nil {
		return u.err
	}
	return fn(ctx)
}

func testProvider(t *testing.T) domain.Provider {
	t.Helper()
	p, err := domain.NewProvider(string(domain.ProviderTelegram))
	require.NoError(t, err)
	return p
}
func testSubject(t *testing.T, value string) valueobjects.ProviderSubject {
	t.Helper()
	s, err := valueobjects.NewProviderSubject(value)
	require.NoError(t, err)
	return s
}
func testTransaction(t *testing.T, at time.Time) *domain.AuthTransaction {
	t.Helper()
	state := sha256.Sum256([]byte("state"))
	nonce := sha256.Sum256([]byte("nonce"))
	sh, err := valueobjects.NewStateHash(state[:])
	require.NoError(t, err)
	nh, err := valueobjects.NewNonceHash(nonce[:])
	require.NoError(t, err)
	binding := sha256.Sum256([]byte("browser-binding"))
	bb, err := valueobjects.NewBrowserBindingHash(binding[:])
	require.NoError(t, err)
	tx, err := domain.NewAuthTransaction(domain.AuthTransactionParams{Provider: testProvider(t), StateHash: sh, NonceHash: nh, BrowserBindingHash: bb, CodeVerifier: "verifier", RedirectURI: "https://example/callback", CreatedAt: at, ExpiresAt: at.Add(5 * time.Minute)})
	require.NoError(t, err)
	return tx
}
func testSession(t *testing.T, at time.Time, lastSeen time.Time) *domain.Session {
	t.Helper()
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	s, err := domain.NewSession(domain.SessionParams{UserID: sharedvalueobjects.NewID(), TokenHash: hash, CreatedAt: at, LastSeenAt: lastSeen, ExpiresAt: at.Add(time.Hour)})
	require.NoError(t, err)
	return s
}

func TestStartTelegramLoginCreatesTransactionAndURL(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	repo := &transactionRepositoryStub{}
	oidc := &oidcStub{url: "https://oauth.example/auth"}
	h := NewStartTelegramLoginHandler(repo, oidc, tokenGeneratorStub{state: []byte("state"), nonce: []byte("nonce"), verifier: []byte("verifier")}, &uowStub{}, "https://example/callback", 5*time.Minute)
	h.now = func() time.Time { return at }
	start, err := h.Handle(context.Background())
	require.NoError(t, err)
	require.Equal(t, oidc.url, start.AuthorizationURL)
	require.NotEmpty(t, start.BrowserBindingToken)
	require.Equal(t, 300, start.BrowserBindingMaxAge)
	require.NotNil(t, repo.transaction)
	require.Equal(t, "https://example/callback", repo.transaction.RedirectURI())
	require.False(t, repo.transaction.StateHash().IsZero())
	expectedBinding := sha256.Sum256([]byte("YnJvd3Nlci1iaW5kaW5n"))
	require.Equal(t, expectedBinding[:], repo.transaction.BrowserBindingHash().Bytes())
	require.False(t, repo.transaction.NonceHash().IsZero())
	require.Equal(t, "dmVyaWZpZXI", repo.transaction.CodeVerifier())
}

func TestStartTelegramLoginPropagatesRepositoryError(t *testing.T) {
	repo := &transactionRepositoryStub{createErr: errors.New("db failed")}
	h := NewStartTelegramLoginHandler(repo, &oidcStub{url: "url"}, tokenGeneratorStub{state: []byte("state"), nonce: []byte("nonce"), verifier: []byte("verifier")}, &uowStub{}, "callback", time.Minute)
	_, err := h.Handle(context.Background())
	require.ErrorContains(t, err, "db failed")
}

func TestCompleteTelegramLoginRejectsInvalidState(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	repo := &transactionRepositoryStub{findErr: domain.ErrEntityNotFound}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, application.ErrStateInvalid)
}

func TestCompleteTelegramLoginRejectsNilTransactionFromRepository(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	repo := &transactionRepositoryStub{returnNil: true}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.EqualError(t, err, "find auth transaction: repository returned nil transaction")
}

func TestCompleteTelegramLoginRejectsConsumedTransaction(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	consumedAt := at.Add(time.Minute)
	tx := testTransaction(t, at)
	require.NoError(t, tx.Consume(consumedAt))
	repo := &transactionRepositoryStub{transaction: tx}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at.Add(2 * time.Minute) }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, domain.ErrAuthTransactionConsumed)
}

func TestCompleteTelegramLoginRejectsExpiredTransaction(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx := testTransaction(t, at)
	repo := &transactionRepositoryStub{transaction: tx}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at.Add(6 * time.Minute) }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, domain.ErrAuthTransactionExpired)
}

func TestCompleteTelegramLoginPropagatesProviderExchangeFailure(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	exchangeErr := ports.ErrTokenExchangeFailed
	repo := &transactionRepositoryStub{transaction: testTransaction(t, at)}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{exchangeErr: exchangeErr}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, application.ErrIdentityVerificationFailed)
	require.ErrorIs(t, err, exchangeErr)
}

func TestCompleteTelegramLoginPropagatesIDTokenVerificationFailure(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	verificationErr := ports.ErrIDTokenInvalid
	repo := &transactionRepositoryStub{transaction: testTransaction(t, at)}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{exchangeErr: verificationErr}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, application.ErrIdentityVerificationFailed)
	require.ErrorIs(t, err, verificationErr)
}

func TestCompleteTelegramLoginPropagatesUnexpectedDatabaseError(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	dbErr := errors.New("database unavailable")
	uow := &uowStub{err: dbErr}
	repo := &transactionRepositoryStub{transaction: testTransaction(t, at)}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{identity: ports.VerifiedIdentity{Provider: testProvider(t), ProviderSubject: testSubject(t, "telegram-sub")}}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, uow, time.Hour)
	h.now = func() time.Time { return at }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, dbErr)
	require.Equal(t, 1, uow.calls)
}

func TestCompleteTelegramLoginCreatesUserIdentityAndSession(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx := testTransaction(t, at)
	repo := &transactionRepositoryStub{transaction: tx}
	userID := sharedvalueobjects.NewID()
	identityRepo := &identityRepositoryStub{findErr: domain.ErrEntityNotFound}
	sessions := &sessionRepositoryStub{}
	provisioner := &provisionerStub{userID: userID}
	oidc := &oidcStub{identity: ports.VerifiedIdentity{Provider: testProvider(t), ProviderSubject: testSubject(t, "telegram-sub")}}
	h := NewCompleteTelegramLoginHandler(repo, identityRepo, sessions, oidc, provisioner, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at }
	result, err := h.Handle(context.Background(), "code", "state", "browser-binding")
	require.NoError(t, err)
	require.Equal(t, userID, result.UserID)
	require.Equal(t, "c2Vzc2lvbi10b2tlbg", result.RawSessionToken)
	require.Equal(t, at.Add(time.Hour), result.ExpiresAt)
	require.Equal(t, 1, provisioner.calls)
	require.Equal(t, 1, identityRepo.createCalls)
	require.Equal(t, 1, repo.consumed)
	require.NotNil(t, sessions.session)
	require.NotEqual(t, []byte("session-token"), sessions.session.TokenHash().Bytes())
}

func TestCompleteTelegramLoginRejectsNilIdentityFromRepository(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	identityRepo := &identityRepositoryStub{returnNil: true}
	oidc := &oidcStub{identity: ports.VerifiedIdentity{Provider: testProvider(t), ProviderSubject: testSubject(t, "telegram-sub")}}
	uow := &uowStub{}
	h := NewCompleteTelegramLoginHandler(
		&transactionRepositoryStub{transaction: testTransaction(t, at)},
		identityRepo,
		&sessionRepositoryStub{},
		oidc,
		&provisionerStub{userID: sharedvalueobjects.NewID()},
		tokenGeneratorStub{session: []byte("session-token")},
		uow,
		time.Hour,
	)
	h.now = func() time.Time { return at }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.EqualError(t, err, "find auth identity: repository returned nil identity")
	require.Equal(t, 1, uow.calls)
}

func TestCompleteTelegramLoginReusesExistingIdentity(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	identityID := sharedvalueobjects.NewID()
	userID := sharedvalueobjects.NewID()
	existingIdentity, err := domain.RestoreIdentity(domain.IdentitySnapshot{ID: identityID, UserID: userID, Provider: testProvider(t), ProviderSubject: testSubject(t, "telegram-sub"), CreatedAt: at, UpdatedAt: at, LastLoginAt: at})
	require.NoError(t, err)
	repo := &transactionRepositoryStub{transaction: testTransaction(t, at)}
	identityRepo := &identityRepositoryStub{identity: existingIdentity}
	sessions := &sessionRepositoryStub{}
	oidc := &oidcStub{identity: ports.VerifiedIdentity{Provider: testProvider(t), ProviderSubject: testSubject(t, "telegram-sub")}}
	provisioner := &provisionerStub{userID: sharedvalueobjects.NewID()}
	h := NewCompleteTelegramLoginHandler(repo, identityRepo, sessions, oidc, provisioner, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at }
	result, err := h.Handle(context.Background(), "code", "state", "browser-binding")
	require.NoError(t, err)
	require.Equal(t, userID, result.UserID)
	require.Equal(t, 0, provisioner.calls)
	require.Equal(t, 1, identityRepo.updateCalls)
}

func TestCompleteTelegramLoginRetriesIdentityConflict(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx := testTransaction(t, at)
	userID := sharedvalueobjects.NewID()
	existingIdentity, err := domain.NewIdentity(domain.IdentityParams{UserID: userID, Provider: testProvider(t), ProviderSubject: testSubject(t, "telegram-sub"), CreatedAt: at, UpdatedAt: at, LastLoginAt: at})
	require.NoError(t, err)
	tracing := &raceIdentityRepository{identity: existingIdentity}
	oidc := &oidcStub{identity: ports.VerifiedIdentity{Provider: testProvider(t), ProviderSubject: testSubject(t, "telegram-sub")}}
	provisioner := &provisionerStub{userID: sharedvalueobjects.NewID()}
	h := NewCompleteTelegramLoginHandler(&transactionRepositoryStub{transaction: tx}, tracing, &sessionRepositoryStub{}, oidc, provisioner, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at }
	result, err := h.Handle(context.Background(), "code", "state", "browser-binding")
	require.NoError(t, err)
	require.Equal(t, userID, result.UserID)
	require.Equal(t, 1, tracing.createCalls)
	require.Equal(t, 2, tracing.findCalls)
	require.Equal(t, 1, provisioner.calls)
}

type raceIdentityRepository struct {
	identity               *domain.Identity
	createCalls, findCalls int
}

func (r *raceIdentityRepository) FindByProviderSubject(context.Context, domain.Provider, valueobjects.ProviderSubject) (*domain.Identity, error) {
	r.findCalls++
	if r.findCalls == 1 {
		return nil, domain.ErrEntityNotFound
	}
	return r.identity, nil
}
func (r *raceIdentityRepository) Create(context.Context, *domain.Identity) error {
	r.createCalls++
	return domain.ErrDuplicateEntity
}
func (r *raceIdentityRepository) UpdateLastLogin(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}

func TestCompleteTelegramLoginRejectsMissingBrowserBinding(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx := testTransaction(t, at)
	h := NewCompleteTelegramLoginHandler(&transactionRepositoryStub{transaction: tx}, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at }
	_, err := h.Handle(context.Background(), "code", "state", "")
	require.ErrorIs(t, err, application.ErrInvalidAuthFlow)
}

func TestCompleteTelegramLoginRejectsMismatchedBrowserBinding(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx := testTransaction(t, at)
	h := NewCompleteTelegramLoginHandler(&transactionRepositoryStub{transaction: tx}, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at }
	_, err := h.Handle(context.Background(), "code", "state", "wrong-binding")
	require.ErrorIs(t, err, application.ErrBrowserBindingInvalid)
}

func TestCompleteTelegramLoginAcceptsMatchingBrowserBinding(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx := testTransaction(t, at)
	provider := testProvider(t)
	subject := testSubject(t, "telegram-sub")
	h := NewCompleteTelegramLoginHandler(&transactionRepositoryStub{transaction: tx}, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{identity: ports.VerifiedIdentity{Provider: provider, ProviderSubject: subject}}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &uowStub{}, time.Hour)
	h.now = func() time.Time { return at }
	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")
	require.NoError(t, err)
}

func TestLogoutIsIdempotentForMissingSession(t *testing.T) {
	h := NewLogoutHandler(&sessionRepositoryStub{}, &uowStub{})
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	require.NoError(t, h.Handle(context.Background(), hash, time.Now().UTC()))
}

func TestLogoutRevokesSession(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	sessions := &sessionRepositoryStub{session: testSession(t, at, at)}
	uow := &uowStub{}
	h := NewLogoutHandler(sessions, uow)
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	require.NoError(t, h.Handle(context.Background(), hash, at.Add(time.Minute)))
	require.Equal(t, 1, sessions.revokeCalls)
	require.Equal(t, 1, uow.calls)
}

func TestLogoutRejectsNilSessionFromRepository(t *testing.T) {
	sessions := &sessionRepositoryStub{returnNil: true}
	h := NewLogoutHandler(sessions, &uowStub{})
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)

	err = h.Handle(context.Background(), hash, time.Now().UTC())

	require.EqualError(t, err, "logout: session repository returned nil session")
}

func TestLogoutRejectsZeroTimestamp(t *testing.T) {
	sessions := &sessionRepositoryStub{}
	h := NewLogoutHandler(sessions, &uowStub{})
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)

	err = h.Handle(context.Background(), hash, time.Time{})

	require.EqualError(t, err, "logout: timestamp must be set")
	require.Zero(t, sessions.revokeCalls)
}
