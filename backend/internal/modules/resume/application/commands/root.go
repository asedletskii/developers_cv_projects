package commands

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateResumeRootCommand struct {
	ID           valueobjects.ID
	UserID       useridentity.ID
	Title        string
	Presentation domain.ResumePresentation
}

type UpdateResumeRootCommand struct {
	ResumeID     valueobjects.ID
	UserID       useridentity.ID
	Title        string
	Presentation domain.ResumePresentation
}

type GetResumeRootRepositoryQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type CreateResumeRootRepositoryCommand struct{ Resume *domain.ResumeRoot }
type UpdateResumeRootRepositoryCommand struct{ Resume *domain.ResumeRoot }

type CreateResumeRootRepository interface {
	CreateRoot(context.Context, CreateResumeRootRepositoryCommand) error
}
type GetResumeRootRepository interface {
	GetRoot(context.Context, GetResumeRootRepositoryQuery) (*domain.ResumeRoot, error)
}
type UpdateResumeRootRepository interface {
	UpdateRoot(context.Context, UpdateResumeRootRepositoryCommand) error
}
type DeleteResumeRootRepository interface {
	DeleteRoot(context.Context, GetResumeRootRepositoryQuery) error
}

type CreateResumeRootHandler struct {
	repository CreateResumeRootRepository
	uow        UnitOfWork
}

func NewCreateResumeRootHandler(repository CreateResumeRootRepository, uow UnitOfWork) *CreateResumeRootHandler {
	return &CreateResumeRootHandler{repository: repository, uow: uow}
}
func (h *CreateResumeRootHandler) Handle(ctx context.Context, c CreateResumeRootCommand) (valueobjects.ID, error) {
	if h == nil || h.repository == nil || h.uow == nil {
		return valueobjects.ID{}, fmt.Errorf("create resume root handler dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	now := time.Now().UTC()
	err := h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := domain.NewResumeRoot(domain.ResumeRootParams{ID: c.ID, UserID: c.UserID, Title: c.Title, CreatedAt: now, UpdatedAt: now, Presentation: c.Presentation})
		if err != nil {
			return err
		}
		return h.repository.CreateRoot(tx, CreateResumeRootRepositoryCommand{Resume: resume})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return c.ID, nil
}

type GetResumeRootHandler struct{ repository GetResumeRootRepository }

func NewGetResumeRootHandler(repository GetResumeRootRepository) *GetResumeRootHandler {
	return &GetResumeRootHandler{repository: repository}
}
func (h *GetResumeRootHandler) Handle(ctx context.Context, q GetResumeRootRepositoryQuery) (*domain.ResumeRoot, error) {
	if h == nil || h.repository == nil {
		return nil, fmt.Errorf("get resume root handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("invalid get resume root query")
	}
	return h.repository.GetRoot(ctx, q)
}

type UpdateResumeRootHandler struct {
	repository GetResumeRootRepository
	updater    UpdateResumeRootRepository
	uow        UnitOfWork
}

func NewUpdateResumeRootHandler(repository GetResumeRootRepository, updater UpdateResumeRootRepository, uow UnitOfWork) *UpdateResumeRootHandler {
	return &UpdateResumeRootHandler{repository: repository, updater: updater, uow: uow}
}
func (h *UpdateResumeRootHandler) Handle(ctx context.Context, c UpdateResumeRootCommand) error {
	if h == nil || h.repository == nil || h.updater == nil || h.uow == nil {
		return fmt.Errorf("update resume root handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() {
		return fmt.Errorf("invalid update resume root command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		current, err := h.repository.GetRoot(tx, GetResumeRootRepositoryQuery{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		updated, err := domain.NewResumeRoot(domain.ResumeRootParams{ID: current.ID(), UserID: current.UserID(), Title: c.Title, CreatedAt: current.CreatedAt(), UpdatedAt: time.Now().UTC(), Presentation: c.Presentation})
		if err != nil {
			return err
		}
		return h.updater.UpdateRoot(tx, UpdateResumeRootRepositoryCommand{Resume: updated})
	})
}

type DeleteResumeRootHandler struct {
	repository DeleteResumeRootRepository
	uow        UnitOfWork
}

func NewDeleteResumeRootHandler(repository DeleteResumeRootRepository, uow UnitOfWork) *DeleteResumeRootHandler {
	return &DeleteResumeRootHandler{repository: repository, uow: uow}
}
func (h *DeleteResumeRootHandler) Handle(ctx context.Context, q GetResumeRootRepositoryQuery) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("delete resume root handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return fmt.Errorf("invalid delete resume root query")
	}
	return h.uow.Do(ctx, func(tx context.Context) error { return h.repository.DeleteRoot(tx, q) })
}
