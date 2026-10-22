package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeLinkRepository interface {
	CountByResumeID(context.Context, sharedvalueobjects.ID) (int64, error)
	Create(context.Context, *domain.ResumeLink) error
	GetByID(context.Context, sharedvalueobjects.ID) (*domain.ResumeLink, error)
	GetByIDForUpdate(context.Context, sharedvalueobjects.ID) (*domain.ResumeLink, error)
	Delete(context.Context, sharedvalueobjects.ID) error
	Rename(context.Context, sharedvalueobjects.ID, string, int64) error
	ReplaceStatus(context.Context, sharedvalueobjects.ID, bool, int64) error
}

type ResumeLinkCleanupRepository interface {
	ListIDsForUpdateByResumeIDs(context.Context, []sharedvalueobjects.ID) ([]sharedvalueobjects.ID, error)
	DeleteByResumeIDs(context.Context, []sharedvalueobjects.ID) error
}
