package app

import (
	"net/http"
	"time"

	authcommands "developers-cv/backend/internal/modules/auth/application/commands"
	authports "developers-cv/backend/internal/modules/auth/application/ports"
	authoidctelegram "developers-cv/backend/internal/modules/auth/infrastructure/oidc/telegram"
	authhttp "developers-cv/backend/internal/modules/auth/transport/http"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/telemetry"
)

func newTelegramOIDC(cfg config.AuthConfig, client *http.Client, metrics *telemetry.Metrics) (authports.TelegramOIDC, error) {
	if cfg.Mode != "telegram" {
		return nil, nil
	}
	return authoidctelegram.New(authoidctelegram.Config{
		ClientID: cfg.Telegram.ClientID, ClientSecret: cfg.Telegram.ClientSecret,
		RedirectURL: cfg.Telegram.RedirectURL, IssuerURL: cfg.Telegram.IssuerURL,
	}, client, metrics)
}

func cookieConfigFromRuntime(cfg config.AuthConfig) authhttp.CookieConfig {
	return authhttp.CookieConfig{
		Name: cfg.Cookie.Name, Path: "/", Secure: cfg.Cookie.Secure,
		HTTPOnly: true, SameSite: http.SameSiteLaxMode,
		MaxAge: int(cfg.Session.AbsoluteTTL.Seconds()),
	}
}

func newTelegramHTTPClient() *http.Client { return &http.Client{Timeout: 10 * time.Second} }

func newCompleteTelegramLoginHandler(
	transactions authports.AuthTransactionRepository,
	identities authports.IdentityRepository,
	sessions authports.SessionRepository,
	oidc authports.TelegramOIDC,
	provisioner authports.UserProvisioner,
	tokenGenerator authports.TokenGenerator,
	uow authports.UnitOfWork,
	cfg config.AuthConfig,
	observer authports.AuthObserver,
) *authcommands.CompleteTelegramLoginHandler {
	return authcommands.NewCompleteTelegramLoginHandler(transactions, identities, sessions, oidc, provisioner, tokenGenerator, uow, cfg.Session.AbsoluteTTL, observer)
}
