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

func TestStartTelegramLoginCreatesTransactionAndURL(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	repo := &transactionRepositoryStub{}
	oidc := &oidcStub{url: "https://oauth.example/auth"}
	h := NewStartTelegramLoginHandler(repo, oidc, tokenGeneratorStub{state: []byte("state"), nonce: []byte("nonce"), verifier: []byte("verifier")}, &unitOfWorkStub{}, "https://example/callback", 5*time.Minute)
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
	h := NewStartTelegramLoginHandler(repo, &oidcStub{url: "url"}, tokenGeneratorStub{state: []byte("state"), nonce: []byte("nonce"), verifier: []byte("verifier")}, &unitOfWorkStub{}, "callback", time.Minute)
	_, err := h.Handle(context.Background())
	require.ErrorContains(t, err, "db failed")
}

func TestCompleteTelegramLoginRejectsInvalidState(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	repo := &transactionRepositoryStub{findErr: domain.ErrEntityNotFound}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
	h.now = func() time.Time { return at }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, application.ErrStateInvalid)
}

func TestCompleteTelegramLoginRejectsNilTransactionFromRepository(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	repo := &transactionRepositoryStub{returnNil: true}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
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
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
	h.now = func() time.Time { return at.Add(2 * time.Minute) }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, domain.ErrAuthTransactionConsumed)
}

func TestCompleteTelegramLoginRejectsExpiredTransaction(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx := testTransaction(t, at)
	repo := &transactionRepositoryStub{transaction: tx}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
	h.now = func() time.Time { return at.Add(6 * time.Minute) }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, domain.ErrAuthTransactionExpired)
}

func TestCompleteTelegramLoginPropagatesProviderExchangeFailure(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	exchangeErr := ports.ErrTokenExchangeFailed
	repo := &transactionRepositoryStub{transaction: testTransaction(t, at)}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{exchangeErr: exchangeErr}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
	h.now = func() time.Time { return at }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, application.ErrIdentityVerificationFailed)
	require.ErrorIs(t, err, exchangeErr)
}

func TestCompleteTelegramLoginPropagatesIDTokenVerificationFailure(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	verificationErr := ports.ErrIDTokenInvalid
	repo := &transactionRepositoryStub{transaction: testTransaction(t, at)}
	h := NewCompleteTelegramLoginHandler(repo, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{exchangeErr: verificationErr}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
	h.now = func() time.Time { return at }

	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, application.ErrIdentityVerificationFailed)
	require.ErrorIs(t, err, verificationErr)
}

func TestCompleteTelegramLoginPropagatesUnexpectedDatabaseError(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	dbErr := errors.New("database unavailable")
	uow := &unitOfWorkStub{err: dbErr}
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
	h := NewCompleteTelegramLoginHandler(repo, identityRepo, sessions, oidc, provisioner, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
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
	uow := &unitOfWorkStub{}
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
	h := NewCompleteTelegramLoginHandler(repo, identityRepo, sessions, oidc, provisioner, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
	h.now = func() time.Time { return at }
	result, err := h.Handle(context.Background(), "code", "state", "browser-binding")
	require.NoError(t, err)
	require.Equal(t, userID, result.UserID)
	require.Equal(t, 0, provisioner.calls)
	require.Equal(t, 1, provisioner.ensureCalls)
	require.Equal(t, 1, identityRepo.updateCalls)
}

func TestCompleteTelegramLoginRejectsMissingUserForExistingIdentity(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	identityID := sharedvalueobjects.NewID()
	userID := sharedvalueobjects.NewID()
	existingIdentity, err := domain.RestoreIdentity(domain.IdentitySnapshot{ID: identityID, UserID: userID, Provider: testProvider(t), ProviderSubject: testSubject(t, "telegram-sub"), CreatedAt: at, UpdatedAt: at, LastLoginAt: at})
	require.NoError(t, err)
	identityRepo := &identityRepositoryStub{identity: existingIdentity}
	sessions := &sessionRepositoryStub{}
	userErr := errors.New("user not found")
	provisioner := &provisionerStub{userID: sharedvalueobjects.NewID(), ensureUserErr: userErr}
	h := NewCompleteTelegramLoginHandler(
		&transactionRepositoryStub{transaction: testTransaction(t, at)},
		identityRepo,
		sessions,
		&oidcStub{identity: ports.VerifiedIdentity{Provider: testProvider(t), ProviderSubject: testSubject(t, "telegram-sub")}},
		provisioner,
		tokenGeneratorStub{session: []byte("session-token")},
		&unitOfWorkStub{},
		time.Hour,
	)
	h.now = func() time.Time { return at }

	_, err = h.Handle(context.Background(), "code", "state", "browser-binding")

	require.ErrorIs(t, err, userErr)
	require.Equal(t, 1, provisioner.ensureCalls)
	require.Equal(t, 0, identityRepo.updateCalls)
	require.Nil(t, sessions.session)
}

func TestCompleteTelegramLoginRetriesIdentityConflict(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx := testTransaction(t, at)
	userID := sharedvalueobjects.NewID()
	existingIdentity, err := domain.NewIdentity(domain.IdentityParams{UserID: userID, Provider: testProvider(t), ProviderSubject: testSubject(t, "telegram-sub"), CreatedAt: at, UpdatedAt: at, LastLoginAt: at})
	require.NoError(t, err)
	tracing := &raceIdentityRepositoryStub{identity: existingIdentity}
	oidc := &oidcStub{identity: ports.VerifiedIdentity{Provider: testProvider(t), ProviderSubject: testSubject(t, "telegram-sub")}}
	provisioner := &provisionerStub{userID: sharedvalueobjects.NewID()}
	h := NewCompleteTelegramLoginHandler(&transactionRepositoryStub{transaction: tx}, tracing, &sessionRepositoryStub{}, oidc, provisioner, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
	h.now = func() time.Time { return at }
	result, err := h.Handle(context.Background(), "code", "state", "browser-binding")
	require.NoError(t, err)
	require.Equal(t, userID, result.UserID)
	require.Equal(t, 1, tracing.createCalls)
	require.Equal(t, 2, tracing.findCalls)
	require.Equal(t, 1, provisioner.calls)
}

func TestCompleteTelegramLoginRejectsMissingBrowserBinding(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx := testTransaction(t, at)
	h := NewCompleteTelegramLoginHandler(&transactionRepositoryStub{transaction: tx}, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
	h.now = func() time.Time { return at }
	_, err := h.Handle(context.Background(), "code", "state", "")
	require.ErrorIs(t, err, application.ErrInvalidAuthFlow)
}

func TestCompleteTelegramLoginRejectsMismatchedBrowserBinding(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx := testTransaction(t, at)
	h := NewCompleteTelegramLoginHandler(&transactionRepositoryStub{transaction: tx}, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
	h.now = func() time.Time { return at }
	_, err := h.Handle(context.Background(), "code", "state", "wrong-binding")
	require.ErrorIs(t, err, application.ErrBrowserBindingInvalid)
}

func TestCompleteTelegramLoginAcceptsMatchingBrowserBinding(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx := testTransaction(t, at)
	provider := testProvider(t)
	subject := testSubject(t, "telegram-sub")
	h := NewCompleteTelegramLoginHandler(&transactionRepositoryStub{transaction: tx}, &identityRepositoryStub{}, &sessionRepositoryStub{}, &oidcStub{identity: ports.VerifiedIdentity{Provider: provider, ProviderSubject: subject}}, &provisionerStub{userID: sharedvalueobjects.NewID()}, tokenGeneratorStub{session: []byte("session-token")}, &unitOfWorkStub{}, time.Hour)
	h.now = func() time.Time { return at }
	_, err := h.Handle(context.Background(), "code", "state", "browser-binding")
	require.NoError(t, err)
}

func TestLogoutIsIdempotentForMissingSession(t *testing.T) {
	h := NewLogoutHandler(&sessionRepositoryStub{}, &unitOfWorkStub{})
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	require.NoError(t, h.Handle(context.Background(), hash, time.Now().UTC()))
}

func TestLogoutRevokesSession(t *testing.T) {
	at := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	sessions := &sessionRepositoryStub{session: testSession(t, at, at)}
	uow := &unitOfWorkStub{}
	h := NewLogoutHandler(sessions, uow)
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	require.NoError(t, h.Handle(context.Background(), hash, at.Add(time.Minute)))
	require.Equal(t, 1, sessions.revokeCalls)
	require.Equal(t, 1, uow.calls)
}

func TestLogoutRejectsNilSessionFromRepository(t *testing.T) {
	sessions := &sessionRepositoryStub{returnNil: true}
	h := NewLogoutHandler(sessions, &unitOfWorkStub{})
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)

	err = h.Handle(context.Background(), hash, time.Now().UTC())

	require.EqualError(t, err, "logout: session repository returned nil session")
}

func TestLogoutRejectsZeroTimestamp(t *testing.T) {
	sessions := &sessionRepositoryStub{}
	h := NewLogoutHandler(sessions, &unitOfWorkStub{})
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)

	err = h.Handle(context.Background(), hash, time.Time{})

	require.EqualError(t, err, "logout: timestamp must be set")
	require.Zero(t, sessions.revokeCalls)
}
