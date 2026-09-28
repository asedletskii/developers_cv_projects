package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type SkillsInput struct {
	ID   valueobjects.ID
	Hard []enums.HardSkill
	Soft []enums.SoftSkill
}

func (i SkillsInput) toDomain(resumeID valueobjects.ID) domain.SkillsParams {
	return domain.SkillsParams{
		ID: i.ID, ResumeID: resumeID,
		Hard: append([]enums.HardSkill(nil), i.Hard...),
		Soft: append([]enums.SoftSkill(nil), i.Soft...),
	}
}

type ReplaceSkillsCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Input    SkillsInput
}

type ReplaceSkillsHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceSkillsHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceSkillsHandler {
	return &ReplaceSkillsHandler{repository: repository, uow: uow}
}

func (h *ReplaceSkillsHandler) Handle(ctx context.Context, c ReplaceSkillsCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace skills handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace skills command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if resume.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := resume.ReplaceSkills(c.Input.toDomain(c.ResumeID)); err != nil {
			return err
		}
		return h.repository.Save(tx, resume)
	})
}
