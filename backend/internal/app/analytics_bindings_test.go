package app

import (
	analyticscommands "developers-cv/backend/internal/modules/analytics/application/commands"
	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/analyticsaccess"
)

var (
	_ analyticscontracts.PublicVisitCreator             = (*analyticscommands.CreatePublicVisitHandler)(nil)
	_ analyticscontracts.PublicInteractionEventRecorder = (*analyticscommands.RecordPublicInteractionEventHandler)(nil)
	_ analyticscontracts.AnalyticsAccessReader          = (*analyticsaccess.Reader)(nil)
)
