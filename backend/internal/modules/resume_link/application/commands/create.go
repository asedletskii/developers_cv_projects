package commands

import (
	"context"
	"fmt"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/shared/domain/valueobjects"
)

type CreateResumeLinkCommand struct {
	UserID   valueobjects.ID
	ResumeID valueobjects.ID
	Name     string
}

type UnitOfWork interface {
	Do(context.Context, func(context.Context) error) error
}

type CreateResumeLinkHandler struct {
	repository ports.ResumeLinkRepository
	resumes    ports.ResumeOwnershipReader
	uow        UnitOfWork
}

func NewCreateResumeLinkHandler(repository ports.ResumeLinkRepository, resumes ports.ResumeOwnershipReader, uow UnitOfWork) *CreateResumeLinkHandler {
	return &CreateResumeLinkHandler{repository: repository, resumes: resumes, uow: uow}
}

func (h *CreateResumeLinkHandler) Handle(ctx context.Context, c CreateResumeLinkCommand) (valueobjects.ID, error) {
	if h == nil || h.repository == nil || h.resumes == nil || h.uow == nil {
		return valueobjects.ID{}, fmt.Errorf("create resume link handler dependencies are nil")
	}
	if c.UserID.IsZero() || c.ResumeID.IsZero() {
		return valueobjects.ID{}, fmt.Errorf("invalid create resume link command")
	}

	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ID:       valueobjects.NewID(),
		ResumeID: c.ResumeID,
		Name:     c.Name,
	})
	if err != nil {
		return valueobjects.ID{}, err
	}

	linkID := link.ID()
	err = h.uow.Do(ctx, func(tx context.Context) error {
		if _, err := h.resumes.GetOwned(tx, resumequeries.GetResumeQuery{ResumeID: c.ResumeID, UserID: c.UserID}); err != nil {
			return err
		}
		return h.repository.Create(tx, link)
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return linkID, nil
}
