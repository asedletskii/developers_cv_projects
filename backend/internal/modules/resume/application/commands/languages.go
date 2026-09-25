package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type LanguageInput struct {
	ID       valueobjects.ID
	Language enums.SpokenLanguage
	Grade    *enums.LanguageGrade
	Position int
}

type CreateLanguageCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Input        LanguageInput
}

type UpdateLanguageCommand struct {
	ResumeID, LanguageID valueobjects.ID
	UserID               useridentity.ID
	Input                LanguageInput
}

type DeleteLanguageCommand struct {
	ResumeID, LanguageID valueobjects.ID
	UserID               useridentity.ID
}

type CreateLanguageHandler struct {
	r ports.CreateLanguageRepository
	u UnitOfWork
}

func NewCreateLanguageHandler(r ports.CreateLanguageRepository, u UnitOfWork) *CreateLanguageHandler {
	return &CreateLanguageHandler{r: r, u: u}
}

func (h *CreateLanguageHandler) Handle(ctx context.Context, c CreateLanguageCommand) (valueobjects.ID, error) {
	if h == nil || h.r == nil || h.u == nil {
		return valueobjects.ID{}, fmt.Errorf("create language dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	c.Input.ID = c.ID
	a, err := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{ID: c.Input.ID, Language: c.Input.Language, Grade: c.Input.Grade, Position: c.Input.Position})
	if err != nil {
		return valueobjects.ID{}, err
	}
	err = h.u.Do(ctx, func(tx context.Context) error {
		return h.r.CreateLanguage(tx, ports.CreateLanguageRepositoryCommand{Language: a, ResumeID: c.ResumeID, UserID: c.UserID})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return a.ID(), nil
}

type UpdateLanguageHandler struct {
	u   ports.UpdateLanguageRepository
	uow UnitOfWork
}

func NewUpdateLanguageHandler(u ports.UpdateLanguageRepository, uow UnitOfWork) *UpdateLanguageHandler {
	return &UpdateLanguageHandler{u: u, uow: uow}
}

func (h *UpdateLanguageHandler) Handle(ctx context.Context, c UpdateLanguageCommand) error {
	if h == nil || h.u == nil || h.uow == nil {
		return fmt.Errorf("update language dependencies are nil")
	}
	entity, err := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{ID: c.LanguageID, Language: c.Input.Language, Grade: c.Input.Grade, Position: c.Input.Position})
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.u.UpdateLanguage(tx, ports.UpdateLanguageRepositoryCommand{Language: entity, ResumeID: c.ResumeID, UserID: c.UserID})
	})
}

type DeleteLanguageHandler struct {
	r ports.DeleteLanguageRepository
	u UnitOfWork
}

func NewDeleteLanguageHandler(r ports.DeleteLanguageRepository, u UnitOfWork) *DeleteLanguageHandler {
	return &DeleteLanguageHandler{r: r, u: u}
}

func (h *DeleteLanguageHandler) Handle(ctx context.Context, c DeleteLanguageCommand) error {
	if h == nil || h.r == nil || h.u == nil {
		return fmt.Errorf("delete language dependencies are nil")
	}
	return h.u.Do(ctx, func(tx context.Context) error {
		return h.r.DeleteLanguage(tx, ports.DeleteLanguageRepositoryCommand{ResumeID: c.ResumeID, LanguageID: c.LanguageID, UserID: c.UserID})
	})
}
