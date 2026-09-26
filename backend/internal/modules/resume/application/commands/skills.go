package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CreateSkillsCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Hard         []enums.HardSkill
	Soft         []enums.SoftSkill
}
type UpdateSkillsCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Hard     []enums.HardSkill
	Soft     []enums.SoftSkill
}

type CreateSkillsHandler struct {
	r ports.CreateSkillsRepository
	u UnitOfWork
}

func NewCreateSkillsHandler(r ports.CreateSkillsRepository, u UnitOfWork) *CreateSkillsHandler {
	return &CreateSkillsHandler{r, u}
}
func (h *CreateSkillsHandler) Handle(ctx context.Context, c CreateSkillsCommand) (valueobjects.ID, error) {
	if h == nil || h.r == nil || h.u == nil {
		return valueobjects.ID{}, fmt.Errorf("create skills dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	var a *domain.Skills
	err := h.u.Do(ctx, func(tx context.Context) error {
		var e error
		a, e = domain.NewSkills(domain.SkillsParams{ID: c.ID, ResumeID: c.ResumeID, Hard: c.Hard, Soft: c.Soft})
		if e != nil {
			return e
		}
		return h.r.CreateSkills(tx, ports.CreateSkillsRepositoryCommand{Skills: a, UserID: c.UserID})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return a.ID(), nil
}

type UpdateSkillsHandler struct {
	u   ports.UpdateSkillsRepository
	uow UnitOfWork
}

func NewUpdateSkillsHandler(u ports.UpdateSkillsRepository, uow UnitOfWork) *UpdateSkillsHandler {
	return &UpdateSkillsHandler{u: u, uow: uow}
}
func (h *UpdateSkillsHandler) Handle(ctx context.Context, c UpdateSkillsCommand) error {
	if h == nil || h.u == nil || h.uow == nil {
		return fmt.Errorf("update skills dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() {
		return fmt.Errorf("invalid update skills command")
	}
	if _, err := domain.NewSkills(domain.SkillsParams{
		ID: valueobjects.NewID(), ResumeID: c.ResumeID, Hard: c.Hard, Soft: c.Soft,
	}); err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.u.UpdateSkills(tx, ports.UpdateSkillsRepositoryCommand{
			ResumeID: c.ResumeID, UserID: c.UserID, Hard: c.Hard, Soft: c.Soft,
		})
	})
}
