package commands

import (
	"context"
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

func (h *CompleteTelegramLoginHandler) Handle(ctx context.Context, code string, state string) (AuthenticatedSession, error) {
	if h == nil || h.transactions == nil || h.identities == nil || h.sessions == nil || h.oidc == nil || h.provisioner == nil || h.tokenGenerator == nil || h.uow == nil {
		return AuthenticatedSession{}, fmt.Errorf("complete telegram login handler dependencies are nil")
	}
	if code == "" || state == "" || h.sessionTTL <= 0 {
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
	initialNow := h.now().UTC()
	if transaction.IsConsumed() {
		return AuthenticatedSession{}, domain.ErrAuthTransactionConsumed
	}
	if transaction.IsExpired(initialNow) {
		return AuthenticatedSession{}, domain.ErrAuthTransactionExpired
	}

	identity, err := h.oidc.Exchange(ctx, code, transaction.RedirectURI(), transaction.CodeVerifier(), transaction.NonceHash())
	if err != nil {
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
		if transaction.IsConsumed() {
			return AuthenticatedSession{}, domain.ErrAuthTransactionConsumed
		}
		if transaction.IsExpired(now) {
			return AuthenticatedSession{}, domain.ErrAuthTransactionExpired
		}
	}
	return AuthenticatedSession{}, application.ErrConcurrentIdentityCreation
}

var errIdentityCreationConflict = errors.New("identity creation conflict")

func (h *CompleteTelegramLoginHandler) completeAttempt(ctx context.Context, transactionID sharedvalueobjects.ID, stateHash valueobjects.StateHash, verified ports.VerifiedIdentity) (AuthenticatedSession, error) {
	var result AuthenticatedSession
	err := h.uow.Do(ctx, func(tx context.Context) error {
		transaction, err := h.transactions.FindByStateHashForUpdate(tx, stateHash)
		if err != nil {
			return err
		}
		if transaction.ID() != transactionID {
			return domain.ErrAuthTransactionConsumed
		}
		identityCtx, finishIdentitySpan := tx, func(error) {}
		if h.observer != nil {
			identityCtx, finishIdentitySpan = h.observer.StartSpan(tx, "auth.resolve_identity", verified.Provider)
		}
		identity, err := h.identities.FindByProviderSubject(identityCtx, verified.Provider, verified.ProviderSubject)
		if err != nil && !errors.Is(err, domain.ErrEntityNotFound) {
			finishIdentitySpan(err)
			return err
		}
		finishIdentitySpan(nil)

		// The login timestamp must be sampled after identity resolution. Under concurrent
		// first-login attempts, the identity may have been created by the other
		// transaction while this transaction was waiting on its locks. Sampling now
		// before that wait can produce a timestamp earlier than identity.created_at,
		// which is rejected by user_identities_timestamps_check.
		loginAt := h.now().UTC()
		var userID sharedvalueobjects.ID
		if errors.Is(err, domain.ErrEntityNotFound) {
			userID, err = h.provisioner.CreateUser(tx)
			if err != nil {
				return err
			}
			identity, err = domain.NewIdentity(domain.IdentityParams{
				ID: sharedvalueobjects.NewID(), UserID: userID, Provider: verified.Provider,
				ProviderSubject: verified.ProviderSubject, CreatedAt: loginAt, UpdatedAt: loginAt, LastLoginAt: loginAt,
			})
			if err != nil {
				return err
			}
			if err := h.identities.Create(tx, identity); err != nil {
				if errors.Is(err, domain.ErrDuplicateEntity) {
					return fmt.Errorf("%w: %w", errIdentityCreationConflict, err)
				}
				return err
			}
		} else {
			userID = identity.UserID()
			if err := h.identities.UpdateLastLogin(tx, identity.ID(), loginAt); err != nil {
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
		rawToken, err := h.tokenGenerator.GenerateSessionToken()
		if err != nil {
			return failSession(fmt.Errorf("generate session token: %w", err))
		}
		tokenHash, err := valueobjects.NewSessionTokenHash(h.tokenGenerator.Hash(rawToken))
		if err != nil {
			return failSession(fmt.Errorf("create session token hash: %w", err))
		}
		expiresAt := loginAt.Add(h.sessionTTL)
		session, err := domain.NewSession(domain.SessionParams{
			ID: sharedvalueobjects.NewID(), UserID: userID, TokenHash: tokenHash,
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
