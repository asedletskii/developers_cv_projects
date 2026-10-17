package commands

import (
	"context"
	"fmt"

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
	resumes    ports.ResumeOwnershipReader
	uow        ports.UnitOfWork
}

func NewCreateResumeLinkHandler(repository ports.ResumeLinkRepository, resumes ports.ResumeOwnershipReader, uow ports.UnitOfWork) *CreateResumeLinkHandler {
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
		if err := ownership.EnsureResumeOwned(tx, h.resumes, ports.ResumeOwnershipQuery{ResumeID: c.ResumeID, UserID: c.UserID}); err != nil {
			return err
		}
		if err := h.repository.LockResumeForLinkCreation(tx, c.ResumeID); err != nil {
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
