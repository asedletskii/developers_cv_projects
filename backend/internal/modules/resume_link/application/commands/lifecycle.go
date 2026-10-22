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

type ReplaceStatusResumeLinkCommand struct {
	UserID  sharedvalueobjects.ID
	LinkID  sharedvalueobjects.ID
	Active  bool
	Version int64
}

type ReplaceStatusResumeLinkHandler struct {
	repository ports.ResumeLinkRepository
	resumes    resumecontracts.OwnedResumeLocker
	uow        ports.UnitOfWork
}

func NewReplaceStatusResumeLinkHandler(repository ports.ResumeLinkRepository, resumes resumecontracts.OwnedResumeLocker, uow ports.UnitOfWork) *ReplaceStatusResumeLinkHandler {
	return &ReplaceStatusResumeLinkHandler{repository: repository, resumes: resumes, uow: uow}
}

func (h *ReplaceStatusResumeLinkHandler) Handle(ctx context.Context, c ReplaceStatusResumeLinkCommand) error {
	if h == nil || h.repository == nil || h.resumes == nil || h.uow == nil {
		return fmt.Errorf("replace resume link status handler dependencies are nil")
	}
	if c.UserID.IsZero() || c.LinkID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace resume link status command")
	}

	return h.uow.Do(ctx, func(tx context.Context) error {
		link, err := h.repository.GetByID(tx, c.LinkID)
		if err != nil {
			return application.NormalizeLinkNotFoundError(err)
		}
		if link == nil {
			return fmt.Errorf("replace resume link status: repository returned nil link")
		}
		if err := ownership.EnsureResumeOwnedForUpdate(tx, h.resumes, link.ResumeID(), c.UserID); err != nil {
			return err
		}
		if link.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}

		if link.IsActive() == c.Active {
			return nil
		}

		if c.Active {
			err = link.Activate()
		} else {
			err = link.Deactivate()
		}
		if err != nil {
			return err
		}
		return application.NormalizeLinkNotFoundError(h.repository.ReplaceStatus(tx, c.LinkID, c.Active, c.Version))
	})
}
