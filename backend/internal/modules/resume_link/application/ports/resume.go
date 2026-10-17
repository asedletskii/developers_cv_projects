package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type ResumeOwnershipQuery = queries.GetResumeQuery

type ResumeOwnershipReader interface {
	OwnsResume(context.Context, ResumeOwnershipQuery) (bool, error)
}
