package commands

import (
	"context"
	"fmt"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
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
		ID:       sharedvalueobjects.NewID(),
		ResumeID: c.ResumeID,
		Name:     c.Name,
	})
	if err != nil {
		return sharedvalueobjects.ID{}, err
	}

	linkID := link.ID()
	err = h.uow.Do(ctx, func(tx context.Context) error {
		if _, err := h.resumes.GetOwned(tx, resumequeries.GetResumeQuery{ResumeID: c.ResumeID, UserID: c.UserID}); err != nil {
			return err
		}
		return h.repository.Create(tx, link)
	})
	if err != nil {
		return sharedvalueobjects.ID{}, err
	}
	return linkID, nil
}
