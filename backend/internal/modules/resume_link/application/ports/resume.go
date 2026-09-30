package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type ResumeOwnershipReader interface {
	GetOwned(context.Context, queries.GetResumeQuery) (queries.ResumeView, error)
}
