package commands

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/auth/application"
	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type AuthenticatedSession struct {
	UserID          sharedvalueobjects.ID
	RawSessionToken string
	ExpiresAt       time.Time
}

type CompleteTelegramLoginHandler struct {
	transactions   ports.AuthTransactionRepository
	identities     ports.IdentityRepository
	sessions       ports.SessionRepository
	oidc           ports.TelegramOIDC
	provisioner    ports.UserProvisioner
	tokenGenerator ports.TokenGenerator
	uow            ports.UnitOfWork
	observer       ports.AuthObserver
	sessionTTL     time.Duration
	now            func() time.Time
}

func NewCompleteTelegramLoginHandler(
	transactions ports.AuthTransactionRepository,
	identities ports.IdentityRepository,
	sessions ports.SessionRepository,
	oidc ports.TelegramOIDC,
	provisioner ports.UserProvisioner,
	tokenGenerator ports.TokenGenerator,
	uow ports.UnitOfWork,
	sessionTTL time.Duration,
	observers ...ports.AuthObserver,
) *CompleteTelegramLoginHandler {
	var observer ports.AuthObserver
	if len(observers) > 0 {
		observer = observers[0]
	}
	return &CompleteTelegramLoginHandler{
		transactions: transactions, identities: identities, sessions: sessions, oidc: oidc,
		provisioner: provisioner, tokenGenerator: tokenGenerator, uow: uow, observer: observer, sessionTTL: sessionTTL,
		now: func() time.Time { return time.Now().UTC() },
	}
}

func (h *CompleteTelegramLoginHandler) Handle(
	ctx context.Context,
	code string,
	state string,
	browserBindingToken string,
) (AuthenticatedSession, error) {
	if h == nil ||
		h.transactions == nil ||
		h.identities == nil ||
		h.sessions == nil ||
		h.oidc == nil ||
		h.provisioner == nil ||
		h.tokenGenerator == nil ||
		h.uow == nil {
		return AuthenticatedSession{}, fmt.Errorf("complete telegram login handler dependencies are nil")
	}
	if code == "" || state == "" || browserBindingToken == "" || h.sessionTTL <= 0 {
		return AuthenticatedSession{}, fmt.Errorf("%w: invalid complete telegram login request", application.ErrInvalidAuthFlow)
	}

	stateHash, err := valueobjects.NewStateHash(h.tokenGenerator.Hash([]byte(state)))
	if err != nil {
		return AuthenticatedSession{}, fmt.Errorf("create state hash: %w", err)
	}
	transaction, err := h.transactions.FindByStateHash(ctx, stateHash)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return AuthenticatedSession{}, application.ErrStateInvalid
		}
		return AuthenticatedSession{}, err
	}
	if transaction == nil {
		return AuthenticatedSession{}, fmt.Errorf("find auth transaction: repository returned nil transaction")
	}
	initialNow := h.now().UTC()
	if transaction.IsConsumed() {
		return AuthenticatedSession{}, domain.ErrAuthTransactionConsumed
	}
	if transaction.IsExpired(initialNow) {
		return AuthenticatedSession{}, domain.ErrAuthTransactionExpired
	}
	browserBindingHash, err := valueobjects.NewBrowserBindingHash(h.tokenGenerator.Hash([]byte(browserBindingToken)))
	if err != nil {
		return AuthenticatedSession{}, fmt.Errorf("create browser binding hash: %w", err)
	}
	if subtle.ConstantTimeCompare(browserBindingHash.Bytes(), transaction.BrowserBindingHash().Bytes()) != 1 {
		return AuthenticatedSession{}, application.ErrBrowserBindingInvalid
	}

	identity, err := h.oidc.Exchange(ctx, code, transaction.RedirectURI(), transaction.CodeVerifier(), transaction.NonceHash())
	if err != nil {
		if errors.Is(err, ports.ErrProviderUnavailable) {
			return AuthenticatedSession{}, fmt.Errorf("%w: %w", application.ErrProviderUnavailable, err)
		}
		return AuthenticatedSession{}, application.WrapIdentityVerificationError(err)
	}
	now := h.now().UTC()
	if transaction.IsExpired(now) {
		return AuthenticatedSession{}, domain.ErrAuthTransactionExpired
	}
	if identity.Provider != transaction.Provider() {
		return AuthenticatedSession{}, fmt.Errorf("telegram identity provider does not match transaction")
	}

	for attempt := 0; attempt < 2; attempt++ {
		result, err := h.completeAttempt(ctx, transaction.ID(), stateHash, identity)
		if err == nil {
			return result, nil
		}
		if !errors.Is(err, errIdentityCreationConflict) || attempt == 1 {
			if errors.Is(err, errIdentityCreationConflict) {
				return AuthenticatedSession{}, application.ErrConcurrentIdentityCreation
			}
			return AuthenticatedSession{}, err
		}
		transaction, err = h.transactions.FindByStateHash(ctx, stateHash)
		if err != nil {
			return AuthenticatedSession{}, err
		}
		if transaction == nil {
			return AuthenticatedSession{}, fmt.Errorf("find auth transaction for retry: repository returned nil transaction")
		}
		if transaction.IsConsumed() {
			return AuthenticatedSession{}, domain.ErrAuthTransactionConsumed
		}
		retryNow := h.now().UTC()
		if transaction.IsExpired(retryNow) {
			return AuthenticatedSession{}, domain.ErrAuthTransactionExpired
		}
	}
	return AuthenticatedSession{}, application.ErrConcurrentIdentityCreation
}

var errIdentityCreationConflict = errors.New("identity creation conflict")

func (h *CompleteTelegramLoginHandler) completeAttempt(
	ctx context.Context,
	transactionID sharedvalueobjects.ID,
	stateHash valueobjects.StateHash,
	verified ports.VerifiedIdentity,
) (AuthenticatedSession, error) {
	rawToken, err := h.tokenGenerator.GenerateSessionToken()
	if err != nil {
		return AuthenticatedSession{}, fmt.Errorf("generate session token: %w", err)
	}
	tokenHash, err := valueobjects.NewSessionTokenHash(h.tokenGenerator.Hash(rawToken))
	if err != nil {
		return AuthenticatedSession{}, fmt.Errorf("create session token hash: %w", err)
	}

	var result AuthenticatedSession
	err = h.uow.Do(ctx, func(tx context.Context) error {
		transaction, err := h.transactions.FindByStateHashForUpdate(tx, stateHash)
		if err != nil {
			return err
		}
		if transaction == nil {
			return fmt.Errorf("lock auth transaction: repository returned nil transaction")
		}
		if transaction.ID() != transactionID {
			return domain.ErrAuthTransactionConsumed
		}
		identityCtx, finishIdentitySpan := tx, func(error) {}
		if h.observer != nil {
			identityCtx, finishIdentitySpan = h.observer.StartSpan(tx, "auth.resolve_identity", verified.Provider)
		}
		authIdentity, err := h.identities.FindByProviderSubject(identityCtx, verified.Provider, verified.ProviderSubject)
		if err != nil && !errors.Is(err, domain.ErrEntityNotFound) {
			finishIdentitySpan(err)
			return err
		}
		if err == nil && authIdentity == nil {
			err = fmt.Errorf("find auth identity: repository returned nil identity")
			finishIdentitySpan(err)
			return err
		}
		finishIdentitySpan(nil)

		loginAt := h.now().UTC()
		var userID sharedvalueobjects.ID
		if errors.Is(err, domain.ErrEntityNotFound) {
			userID, err = h.provisioner.CreateUser(tx)
			if err != nil {
				return err
			}
			authIdentity, err = domain.NewIdentity(domain.IdentityParams{
				UserID: userID, Provider: verified.Provider,
				ProviderSubject: verified.ProviderSubject, CreatedAt: loginAt, UpdatedAt: loginAt, LastLoginAt: loginAt,
			})
			if err != nil {
				return err
			}
			if err := h.identities.Create(tx, authIdentity); err != nil {
				if errors.Is(err, domain.ErrDuplicateEntity) {
					return fmt.Errorf("%w: %w", errIdentityCreationConflict, err)
				}
				return err
			}
		} else {
			userID = authIdentity.UserID()
			if err := h.identities.UpdateLastLogin(tx, authIdentity.ID(), loginAt); err != nil {
				return err
			}
		}

		sessionCtx, finishSessionSpan := tx, func(error) {}
		if h.observer != nil {
			sessionCtx, finishSessionSpan = h.observer.StartSpan(tx, "auth.create_session", verified.Provider)
		}
		var sessionSpanErr error
		defer func() { finishSessionSpan(sessionSpanErr) }()
		failSession := func(err error) error {
			sessionSpanErr = err
			return err
		}
		expiresAt := loginAt.Add(h.sessionTTL)
		session, err := domain.NewSession(domain.SessionParams{
			UserID: userID, TokenHash: tokenHash,
			CreatedAt: loginAt, LastSeenAt: loginAt, ExpiresAt: expiresAt,
		})
		if err != nil {
			return failSession(err)
		}
		if err := h.sessions.Create(sessionCtx, session); err != nil {
			return failSession(err)
		}
		if err := transaction.Consume(loginAt); err != nil {
			return failSession(err)
		}
		if err := h.transactions.Consume(tx, transaction.ID(), loginAt); err != nil {
			return failSession(err)
		}
		result = AuthenticatedSession{UserID: userID, RawSessionToken: encodeOpaque(rawToken), ExpiresAt: expiresAt}
		return nil
	})
	return result, err
}
