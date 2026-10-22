package commands

import (
	"context"
	"fmt"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume_link/application/ownership"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CreateResumeLinkCommand struct {
	UserID   sharedvalueobjects.ID
	ResumeID sharedvalueobjects.ID
	Name     string
}

type CreateResumeLinkHandler struct {
	repository ports.ResumeLinkRepository
	resumes    resumecontracts.OwnedResumeLocker
	uow        ports.UnitOfWork
}

func NewCreateResumeLinkHandler(repository ports.ResumeLinkRepository, resumes resumecontracts.OwnedResumeLocker, uow ports.UnitOfWork) *CreateResumeLinkHandler {
	return &CreateResumeLinkHandler{repository: repository, resumes: resumes, uow: uow}
}

func (h *CreateResumeLinkHandler) Handle(ctx context.Context, c CreateResumeLinkCommand) (sharedvalueobjects.ID, error) {
	if h == nil || h.repository == nil || h.resumes == nil || h.uow == nil {
		return sharedvalueobjects.ID{}, fmt.Errorf("create resume link handler dependencies are nil")
	}
	if c.UserID.IsZero() || c.ResumeID.IsZero() {
		return sharedvalueobjects.ID{}, fmt.Errorf("invalid create resume link command")
	}

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ResumeID: c.ResumeID,
		Name:     c.Name,
	})
	if err != nil {
		return sharedvalueobjects.ID{}, err
	}

	linkID := link.ID()
	err = h.uow.Do(ctx, func(tx context.Context) error {
		if err := ownership.EnsureResumeOwnedForUpdate(tx, h.resumes, c.ResumeID, c.UserID); err != nil {
			return err
		}
		count, err := h.repository.CountByResumeID(tx, c.ResumeID)
		if err != nil {
			return err
		}
		if err := domain.EnsureResumeLinkCapacity(count); err != nil {
			return err
		}
		return h.repository.Create(tx, link)
	})
	if err != nil {
		return sharedvalueobjects.ID{}, err
	}
	return linkID, nil
}
