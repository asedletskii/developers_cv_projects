package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type LanguageInput struct {
	ID       valueobjects.ID
	Language enums.LanguageType
	Grade    *enums.LanguageGrade
	Position int
}
type CreateLanguageCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Input        LanguageInput
}
type ListLanguagesQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetLanguageQuery struct {
	ResumeID, LanguageID valueobjects.ID
	UserID               useridentity.ID
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
type LanguagesRepository interface {
	CreateLanguage(context.Context, LanguageInput, valueobjects.ID, useridentity.ID) error
	ListLanguages(context.Context, ListLanguagesQuery) ([]domain.ResumeLanguageEntry, error)
	GetLanguage(context.Context, GetLanguageQuery) (domain.ResumeLanguageEntry, error)
	UpdateLanguage(context.Context, UpdateLanguageCommand) error
	DeleteLanguage(context.Context, DeleteLanguageCommand) error
}
type CreateLanguageHandler struct {
	r LanguagesRepository
	u UnitOfWork
}

func NewCreateLanguageHandler(r LanguagesRepository, u UnitOfWork) *CreateLanguageHandler {
	return &CreateLanguageHandler{r, u}
}
func (h *CreateLanguageHandler) Handle(ctx context.Context, c CreateLanguageCommand) (valueobjects.ID, error) {
	if h == nil || h.r == nil || h.u == nil {
		return valueobjects.ID{}, fmt.Errorf("create language dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	c.Input.ID = c.ID
	if _, err := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{
		ID: c.Input.ID, Language: c.Input.Language, Grade: c.Input.Grade, Position: c.Input.Position,
	}); err != nil {
		return valueobjects.ID{}, err
	}
	return c.ID, h.u.Do(ctx, func(tx context.Context) error { return h.r.CreateLanguage(tx, c.Input, c.ResumeID, c.UserID) })
}

type ListLanguagesHandler struct{ r LanguagesRepository }

func NewListLanguagesHandler(r LanguagesRepository) *ListLanguagesHandler {
	return &ListLanguagesHandler{r}
}
func (h *ListLanguagesHandler) Handle(ctx context.Context, q ListLanguagesQuery) ([]domain.ResumeLanguageEntry, error) {
	return h.r.ListLanguages(ctx, q)
}

type GetLanguageHandler struct{ r LanguagesRepository }

func NewGetLanguageHandler(r LanguagesRepository) *GetLanguageHandler { return &GetLanguageHandler{r} }
func (h *GetLanguageHandler) Handle(ctx context.Context, q GetLanguageQuery) (domain.ResumeLanguageEntry, error) {
	return h.r.GetLanguage(ctx, q)
}

type UpdateLanguageHandler struct {
	r LanguagesRepository
	u UnitOfWork
}

func NewUpdateLanguageHandler(r LanguagesRepository, u UnitOfWork) *UpdateLanguageHandler {
	return &UpdateLanguageHandler{r, u}
}
func (h *UpdateLanguageHandler) Handle(ctx context.Context, c UpdateLanguageCommand) error {
	c.Input.ID = c.LanguageID
	if _, err := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{
		ID: c.Input.ID, Language: c.Input.Language, Grade: c.Input.Grade, Position: c.Input.Position,
	}); err != nil {
		return err
	}
	return h.u.Do(ctx, func(tx context.Context) error { return h.r.UpdateLanguage(tx, c) })
}

type DeleteLanguageHandler struct {
	r LanguagesRepository
	u UnitOfWork
}

func NewDeleteLanguageHandler(r LanguagesRepository, u UnitOfWork) *DeleteLanguageHandler {
	return &DeleteLanguageHandler{r, u}
}
func (h *DeleteLanguageHandler) Handle(ctx context.Context, c DeleteLanguageCommand) error {
	return h.u.Do(ctx, func(tx context.Context) error { return h.r.DeleteLanguage(tx, c) })
}
