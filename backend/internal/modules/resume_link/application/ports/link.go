package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/domain/valueobjects"
)

type ResumeLinkRepository interface {
	Create(context.Context, *domain.ResumeLink) error
	GetByID(context.Context, valueobjects.ID) (*domain.ResumeLink, error)
	Delete(context.Context, valueobjects.ID) error
	ReplaceStatus(context.Context, valueobjects.ID, bool, int64) error
	RecordEvent(context.Context, *domain.ResumeLinkEvent) error
}

type LinkEventRecorder interface {
	Record(context.Context, *domain.ResumeLinkEvent) error
}
