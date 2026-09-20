package commands

import (
	"context"
	"fmt"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	application "developers-cv/backend/internal/modules/resume_link/application"
	"developers-cv/backend/internal/modules/resume_link/application/ownership"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type RenameResumeLinkCommand struct {
	UserID  sharedvalueobjects.ID
	LinkID  sharedvalueobjects.ID
	Name    string
	Version int64
}

type RenameResumeLinkHandler struct {
	repository ports.ResumeLinkRepository
	resumes    resumecontracts.OwnedResumeLocker
	uow        ports.UnitOfWork
}

func NewRenameResumeLinkHandler(
	repository ports.ResumeLinkRepository,
	resumes resumecontracts.OwnedResumeLocker,
	uow ports.UnitOfWork,
) *RenameResumeLinkHandler {
	return &RenameResumeLinkHandler{repository: repository, resumes: resumes, uow: uow}
}

func (h *RenameResumeLinkHandler) Handle(ctx context.Context, c RenameResumeLinkCommand) error {
	if h == nil || h.repository == nil || h.resumes == nil || h.uow == nil {
		return fmt.Errorf("rename resume link handler dependencies are nil")
	}
	if c.UserID.IsZero() || c.LinkID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid rename resume link command")
	}

	return h.uow.Do(ctx, func(tx context.Context) error {
		link, err := h.repository.GetByID(tx, c.LinkID)
		if err != nil {
			return application.NormalizeLinkNotFoundError(err)
		}
		if link == nil {
			return fmt.Errorf("rename resume link: repository returned nil link")
		}
		if err := ownership.EnsureResumeOwnedForUpdate(tx, h.resumes, link.ResumeID(), c.UserID); err != nil {
			return err
		}
		if link.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := link.Rename(c.Name); err != nil {
			return err
		}
		return application.NormalizeLinkNotFoundError(h.repository.Rename(tx, c.LinkID, link.Name(), c.Version))
	})
}
