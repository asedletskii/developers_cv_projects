package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ReplaceProjectsCommand struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Version  int64
	Input    ProjectsInput
}

type ReplaceProjectsHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewReplaceProjectsHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceProjectsHandler {
	return &ReplaceProjectsHandler{repository: repository, uow: uow}
}

func (h *ReplaceProjectsHandler) Handle(ctx context.Context, c ReplaceProjectsCommand) error {
	if h == nil {
		return fmt.Errorf("replace projects handler dependencies are nil")
	}
	params, err := c.Input.toDomain()
	if err != nil {
		return err
	}
	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace projects", func(resume *domain.Resume) error {
		return resume.ReplaceProjects(params)
	})
}
