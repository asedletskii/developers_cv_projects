package ports

import (
	"context"
	"time"
)

type AnalyticsRetentionRepository interface {
	DeleteExpiredClicks(context.Context, time.Time, int) (int, error)
	DeleteExpiredViews(context.Context, time.Time, int) (int, error)
}
