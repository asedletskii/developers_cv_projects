package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type SkillsInput struct {
	ID   identity.ID
	Hard []domain.HardSkill
	Soft []domain.SoftSkill
}

func (i SkillsInput) toDomain(resumeID identity.ID) domain.SkillsParams {
	return domain.SkillsParams{
		ID: i.ID, ResumeID: resumeID,
		Hard: append([]domain.HardSkill(nil), i.Hard...),
		Soft: append([]domain.SoftSkill(nil), i.Soft...),
	}
}

type ReplaceSkillsCommand struct {
	ResumeID identity.ID
	UserID   identity.ID
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

	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace skills", func(resume *domain.Resume) error {
		return resume.ReplaceSkills(c.Input.toDomain(c.ResumeID))
	})
}
