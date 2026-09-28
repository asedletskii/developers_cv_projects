package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type LocalizedProfileInput struct {
	Name       *string
	Surname    *string
	Patronymic *string
	Position   *string
	Location   *string
}

type ReplaceProfileCommand struct {
	UserID    valueobjects.ID
	RU        LocalizedProfileInput
	EN        LocalizedProfileInput
	PhotoPath *string
}

type ProfileRepository interface {
	ReplaceProfile(ctx context.Context, profile *domain.Profile) error
}

type ReplaceProfileHandler struct {
	repository ProfileRepository
	uow        UnitOfWork
}

func NewReplaceProfileHandler(repository ProfileRepository, uow UnitOfWork) *ReplaceProfileHandler {
	return &ReplaceProfileHandler{repository: repository, uow: uow}
}

func (h *ReplaceProfileHandler) Handle(ctx context.Context, c ReplaceProfileCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace user profile handler dependencies are nil")
	}
	if c.UserID.IsZero() {
		return fmt.Errorf("invalid replace user profile command")
	}

	profile, err := domain.NewProfile(c.UserID, toDomain(c.RU), toDomain(c.EN), c.PhotoPath)
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.ReplaceProfile(tx, profile)
	})
}

func toDomain(input LocalizedProfileInput) domain.LocalizedProfile {
	return domain.LocalizedProfile{Name: input.Name, Surname: input.Surname, Patronymic: input.Patronymic, Position: input.Position, Location: input.Location}
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}
