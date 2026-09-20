package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/user/application/ports"
	"developers-cv/backend/internal/modules/user/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LocalizedProfileInput struct {
	Name       *string
	Surname    *string
	Patronymic *string
	Position   *string
	Location   *string
}

type ReplaceProfileCommand struct {
	UserID   sharedvalueobjects.ID
	RU       LocalizedProfileInput
	EN       LocalizedProfileInput
	PhotoURL *string
}

type ReplaceProfileHandler struct {
	repository ports.ProfileRepository
	uow        ports.UnitOfWork
}

func NewReplaceProfileHandler(repository ports.ProfileRepository, uow ports.UnitOfWork) *ReplaceProfileHandler {
	return &ReplaceProfileHandler{repository: repository, uow: uow}
}

func (h *ReplaceProfileHandler) Handle(ctx context.Context, c ReplaceProfileCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace user profile handler dependencies are nil")
	}
	ru, err := toDomain(c.RU)
	if err != nil {
		return err
	}
	en, err := toDomain(c.EN)
	if err != nil {
		return err
	}
	profile, err := domain.NewProfile(c.UserID, ru, en, c.PhotoURL)
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.ReplaceProfile(tx, profile)
	})
}

func toDomain(input LocalizedProfileInput) (domain.LocalizedProfile, error) {
	return domain.NewLocalizedProfile(input.Name, input.Surname, input.Patronymic, input.Position, input.Location)
}
