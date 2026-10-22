package commands

import (
	"context"
	"fmt"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	application "developers-cv/backend/internal/modules/resume_link/application"
	"developers-cv/backend/internal/modules/resume_link/application/ownership"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type DeleteResumeLinkCommand struct {
	UserID sharedvalueobjects.ID
	LinkID sharedvalueobjects.ID
}

type DeleteResumeLinkHandler struct {
	repository ports.ResumeLinkRepository
	resumes    resumecontracts.OwnedResumeLocker
	cleanup    analyticscontracts.AnalyticsLinkCleanup
	uow        ports.UnitOfWork
}

func NewDeleteResumeLinkHandler(
	repository ports.ResumeLinkRepository,
	resumes resumecontracts.OwnedResumeLocker,
	cleanup analyticscontracts.AnalyticsLinkCleanup,
	uow ports.UnitOfWork,
) *DeleteResumeLinkHandler {
	return &DeleteResumeLinkHandler{repository: repository, resumes: resumes, cleanup: cleanup, uow: uow}
}

func (h *DeleteResumeLinkHandler) Handle(ctx context.Context, command DeleteResumeLinkCommand) error {
	if h == nil || h.repository == nil || h.resumes == nil || h.cleanup == nil || h.uow == nil {
		return fmt.Errorf("delete resume link handler dependencies are nil")
	}
	if command.UserID.IsZero() || command.LinkID.IsZero() {
		return fmt.Errorf("invalid delete resume link command")
	}

	return h.uow.Do(ctx, func(tx context.Context) error {
		link, err := h.repository.GetByID(tx, command.LinkID)
		if err != nil {
			return application.NormalizeLinkNotFoundError(err)
		}
		if link == nil {
			return fmt.Errorf("delete resume link: repository returned nil link")
		}
		if err := ownership.EnsureResumeOwnedForUpdate(tx, h.resumes, link.ResumeID(), command.UserID); err != nil {
			return err
		}
		locked, err := h.repository.GetByIDForUpdate(tx, command.LinkID)
		if err != nil {
			return application.NormalizeLinkNotFoundError(err)
		}
		if locked == nil {
			return fmt.Errorf("delete resume link: repository returned nil locked link")
		}
		if locked.ResumeID() != link.ResumeID() {
			return fmt.Errorf("delete resume link: link resume changed while locking")
		}
		if err := h.cleanup.DeleteByLinkIDs(tx, []sharedvalueobjects.ID{command.LinkID}); err != nil {
			return fmt.Errorf("delete resume link: clean analytics: %w", err)
		}
		return application.NormalizeLinkNotFoundError(h.repository.Delete(tx, command.LinkID))
	})
}
