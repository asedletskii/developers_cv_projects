package commands

import (
	"context"
	"fmt"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type DeleteResumeLinkCommand struct {
	UserID sharedvalueobjects.ID
	LinkID sharedvalueobjects.ID
}

type DeleteResumeLinkHandler struct {
	repository ports.ResumeLinkRepository
	resumes    ports.ResumeOwnershipReader
	uow        ports.UnitOfWork
}

func NewDeleteResumeLinkHandler(repository ports.ResumeLinkRepository, resumes ports.ResumeOwnershipReader, uow ports.UnitOfWork) *DeleteResumeLinkHandler {
	return &DeleteResumeLinkHandler{repository: repository, resumes: resumes, uow: uow}
}

func (h *DeleteResumeLinkHandler) Handle(ctx context.Context, command DeleteResumeLinkCommand) error {
	if h == nil || h.repository == nil || h.resumes == nil || h.uow == nil {
		return fmt.Errorf("delete resume link handler dependencies are nil")
	}
	if command.UserID.IsZero() || command.LinkID.IsZero() {
		return fmt.Errorf("invalid delete resume link command")
	}

	return h.uow.Do(ctx, func(tx context.Context) error {
		link, err := h.repository.GetByID(tx, command.LinkID)
		if err != nil {
			return err
		}
		if _, err := h.resumes.GetOwned(tx, resumequeries.GetResumeQuery{
			ResumeID: link.ResumeID(),
			UserID:   command.UserID,
		}); err != nil {
			return err
		}
		return h.repository.Delete(tx, command.LinkID)
	})
}
