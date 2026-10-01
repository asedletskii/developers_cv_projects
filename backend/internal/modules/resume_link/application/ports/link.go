package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeLinkRepository interface {
	Create(context.Context, *domain.ResumeLink) error
	GetByID(context.Context, sharedvalueobjects.ID) (*domain.ResumeLink, error)
	Delete(context.Context, sharedvalueobjects.ID) error
	ReplaceStatus(context.Context, sharedvalueobjects.ID, bool, int64) error
	RecordEvent(context.Context, *domain.ResumeLinkEvent) error
}

type LinkEventRecorder interface {
	Record(context.Context, *domain.ResumeLinkEvent) error
}
