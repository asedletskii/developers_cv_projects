package ports

import (
	"context"
	"time"
)

type RetentionRepository interface {
	DeleteExpiredClicks(context.Context, time.Time, int) (int, error)
	DeleteExpiredViews(context.Context, time.Time, int) (int, error)
}
