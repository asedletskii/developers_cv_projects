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
type GetSkillsQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type UpdateSkillsCommand struct {
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
	Hard, Soft []string
}
type SkillsRepository interface {
	CreateSkills(context.Context, *domain.Skills, useridentity.ID) error
	GetSkills(context.Context, GetSkillsQuery) (*domain.Skills, error)
	UpdateSkills(context.Context, *domain.Skills, useridentity.ID) error
}

type CreateSkillsHandler struct {
	r SkillsRepository
	u UnitOfWork
}

func NewCreateSkillsHandler(r SkillsRepository, u UnitOfWork) *CreateSkillsHandler {
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
		return h.r.CreateSkills(tx, a, c.UserID)
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return a.ID(), nil
}

type GetSkillsHandler struct{ r SkillsRepository }

func NewGetSkillsHandler(r SkillsRepository) *GetSkillsHandler { return &GetSkillsHandler{r} }
func (h *GetSkillsHandler) Handle(ctx context.Context, q GetSkillsQuery) (*domain.Skills, error) {
	if h == nil || h.r == nil {
		return nil, fmt.Errorf("get skills dependencies are nil")
	}
	return h.r.GetSkills(ctx, q)
}

type UpdateSkillsHandler struct {
	r SkillsRepository
	u UnitOfWork
}

func NewUpdateSkillsHandler(r SkillsRepository, u UnitOfWork) *UpdateSkillsHandler {
	return &UpdateSkillsHandler{r, u}
}
func (h *UpdateSkillsHandler) Handle(ctx context.Context, c UpdateSkillsCommand) error {
	if h == nil || h.r == nil || h.u == nil {
		return fmt.Errorf("update skills dependencies are nil")
	}
	return h.u.Do(ctx, func(tx context.Context) error {
		cur, e := h.r.GetSkills(tx, GetSkillsQuery{ResumeID: c.ResumeID, UserID: c.UserID})
		if e != nil {
			return e
		}
		a, e := domain.NewSkills(domain.SkillsParams{ID: cur.ID(), ResumeID: c.ResumeID, Hard: toHard(c.Hard), Soft: toSoft(c.Soft)})
		if e != nil {
			return e
		}
		return h.r.UpdateSkills(tx, a, c.UserID)
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
