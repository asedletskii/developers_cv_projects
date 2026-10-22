package app

import (
	analyticscleanup "developers-cv/backend/internal/modules/analytics/infrastructure/cleanup"
	authcleanup "developers-cv/backend/internal/modules/auth/infrastructure/cleanup"
)

func newBackgroundWorkers(authWorker *authcleanup.Worker, analyticsWorker *analyticscleanup.Worker, orphanChecker *orphanChecker) []backgroundWorker {
	workers := make([]backgroundWorker, 0, 3)
	if authWorker != nil {
		workers = append(workers, authWorker)
	}
	if analyticsWorker != nil {
		workers = append(workers, analyticsWorker)
	}
	if orphanChecker != nil {
		workers = append(workers, orphanChecker)
	}
	return workers
}
