package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type SkillsInput struct {
	ID   sharedvalueobjects.ID
	Hard []string
	Soft []string
}

func (i SkillsInput) toDomain() domain.SkillsParams {
	return domain.SkillsParams{
		ID:   i.ID,
		Hard: append([]string(nil), i.Hard...),
		Soft: append([]string(nil), i.Soft...),
	}
}

type ReplaceSkillsCommand struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Version  int64
	Input    SkillsInput
}

type ReplaceSkillsHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewReplaceSkillsHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceSkillsHandler {
	return &ReplaceSkillsHandler{repository: repository, uow: uow}
}

func (h *ReplaceSkillsHandler) Handle(ctx context.Context, c ReplaceSkillsCommand) error {
	if h == nil {
		return fmt.Errorf("replace skills handler dependencies are nil")
	}

	skills := c.Input.toDomain()
	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace skills", ports.ResumeMutationSkills, func(resume *domain.Resume) error {
		return resume.ReplaceSkills(skills)
	})
}
