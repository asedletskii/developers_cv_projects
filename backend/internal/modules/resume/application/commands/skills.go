package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateSkillsCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Hard, Soft   []string
}
type GetSkillsCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type UpdateSkillsCommand struct {
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
	Hard, Soft []string
}
type CreateSkillsRepository interface {
	CreateSkills(context.Context, CreateSkillsRepositoryCommand) error
}
type GetSkillsRepository interface {
	GetSkills(context.Context, GetSkillsCommand) (*domain.Skills, error)
}
type UpdateSkillsRepository interface {
	UpdateSkills(context.Context, UpdateSkillsRepositoryCommand) error
}
type CreateSkillsRepositoryCommand struct {
	Skills *domain.Skills
	UserID useridentity.ID
}
type UpdateSkillsRepositoryCommand struct {
	Skills *domain.Skills
	UserID useridentity.ID
}

type CreateSkillsHandler struct {
	r CreateSkillsRepository
	u UnitOfWork
}

func NewCreateSkillsHandler(r CreateSkillsRepository, u UnitOfWork) *CreateSkillsHandler {
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
		a, e = domain.NewSkills(domain.SkillsParams{ID: c.ID, ResumeID: c.ResumeID, Hard: toHard(c.Hard), Soft: toSoft(c.Soft)})
		if e != nil {
			return e
		}
		return h.r.CreateSkills(tx, CreateSkillsRepositoryCommand{Skills: a, UserID: c.UserID})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return a.ID(), nil
}

type UpdateSkillsHandler struct {
	r   GetSkillsRepository
	u   UpdateSkillsRepository
	uow UnitOfWork
}

func NewUpdateSkillsHandler(r GetSkillsRepository, u UpdateSkillsRepository, uow UnitOfWork) *UpdateSkillsHandler {
	return &UpdateSkillsHandler{r: r, u: u, uow: uow}
}
func (h *UpdateSkillsHandler) Handle(ctx context.Context, c UpdateSkillsCommand) error {
	if h == nil || h.r == nil || h.u == nil || h.uow == nil {
		return fmt.Errorf("update skills dependencies are nil")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		cur, e := h.r.GetSkills(tx, GetSkillsCommand{ResumeID: c.ResumeID, UserID: c.UserID})
		if e != nil {
			return e
		}
		if e := cur.Replace(domain.SkillsParams{ID: cur.ID(), ResumeID: c.ResumeID, Hard: toHard(c.Hard), Soft: toSoft(c.Soft)}); e != nil {
			return e
		}
		return h.u.UpdateSkills(tx, UpdateSkillsRepositoryCommand{Skills: cur, UserID: c.UserID})
	})
}
func toHard(v []string) []enums.HardSkill {
	out := make([]enums.HardSkill, len(v))
	for i, x := range v {
		out[i] = enums.HardSkill(x)
	}
	return out
}
func toSoft(v []string) []enums.SoftSkill {
	out := make([]enums.SoftSkill, len(v))
	for i, x := range v {
		out[i] = enums.SoftSkill(x)
	}
	return out
}
