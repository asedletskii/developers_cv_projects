package ports

import (
	"context"
	"time"

	"developers-cv/backend/internal/modules/auth/domain"
	valueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type AuthObserver interface {
	LoginStarted(context.Context, domain.Provider)
	LoginCallback(context.Context, domain.Provider)
	LoginSucceeded(context.Context, domain.Provider, valueobjects.ID, time.Duration)
	LoginFailed(context.Context, domain.Provider, string, time.Duration)
	Logout(context.Context, valueobjects.ID, bool)
	SessionRejected(context.Context, string)
	StartSpan(context.Context, string, domain.Provider) (context.Context, func(error))
}
