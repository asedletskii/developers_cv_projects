package app

import (
	authcleanup "developers-cv/backend/internal/modules/auth/infrastructure/cleanup"
	resumelinkcleanup "developers-cv/backend/internal/modules/resume_link/infrastructure/cleanup"
)

func newBackgroundWorkers(authWorker *authcleanup.Worker, analyticsWorker *resumelinkcleanup.Worker) []backgroundWorker {
	workers := make([]backgroundWorker, 0, 2)
	if authWorker != nil {
		workers = append(workers, authWorker)
	}
	if analyticsWorker != nil {
		workers = append(workers, analyticsWorker)
	}
	return workers
}
