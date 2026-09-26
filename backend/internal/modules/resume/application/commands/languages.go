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

type LanguageInput struct {
	ID       valueobjects.ID
	Language enums.SpokenLanguage
	Grade    *enums.LanguageGrade
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

type ReorderLanguagesCommand struct {
	ResumeID    valueobjects.ID
	UserID      useridentity.ID
	LanguageIDs []valueobjects.ID
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
	a, err := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{ID: c.Input.ID, Language: c.Input.Language, Grade: c.Input.Grade, Position: 0})
	if err != nil {
		return valueobjects.ID{}, err
	}
	if err := h.u.Do(ctx, func(tx context.Context) error {
		return h.r.CreateLanguage(tx, ports.CreateLanguageRepositoryCommand{Language: a, ResumeID: c.ResumeID, UserID: c.UserID})
	}); err != nil {
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
	entity, err := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{ID: c.LanguageID, Language: c.Input.Language, Grade: c.Input.Grade, Position: 0})
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.u.UpdateLanguage(tx, ports.UpdateLanguageRepositoryCommand{Language: entity, ResumeID: c.ResumeID, UserID: c.UserID})
	})
}

type ReorderLanguagesHandler struct {
	r ports.ReorderLanguagesRepository
	u UnitOfWork
}

func NewReorderLanguagesHandler(r ports.ReorderLanguagesRepository, u UnitOfWork) *ReorderLanguagesHandler {
	return &ReorderLanguagesHandler{r: r, u: u}
}

func (h *ReorderLanguagesHandler) Handle(ctx context.Context, c ReorderLanguagesCommand) error {
	if h == nil || h.r == nil || h.u == nil {
		return fmt.Errorf("reorder languages dependencies are nil")
	}
	if len(c.LanguageIDs) == 0 {
		return fmt.Errorf("reorder languages ids are required: %w", domain.ErrInvalidEntity)
	}
	return h.u.Do(ctx, func(tx context.Context) error {
		return h.r.ReorderLanguages(tx, ports.ReorderLanguagesRepositoryCommand{ResumeID: c.ResumeID, UserID: c.UserID, LanguageIDs: c.LanguageIDs})
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
