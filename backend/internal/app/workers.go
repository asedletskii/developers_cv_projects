package app

import (
	analyticscleanup "developers-cv/backend/internal/modules/analytics/infrastructure/cleanup"
	authcleanup "developers-cv/backend/internal/modules/auth/infrastructure/cleanup"
)

func newBackgroundWorkers(authWorker *authcleanup.Worker, analyticsWorker *analyticscleanup.Worker) []backgroundWorker {
	workers := make([]backgroundWorker, 0, 2)
	if authWorker != nil {
		workers = append(workers, authWorker)
	}
	if analyticsWorker != nil {
		workers = append(workers, analyticsWorker)
	}
	return workers
}
